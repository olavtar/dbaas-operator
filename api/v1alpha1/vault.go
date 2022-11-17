/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"fmt"
	"log"
	"strings"

	auth "github.com/hashicorp/vault/api/auth/kubernetes"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	vault "github.com/hashicorp/vault/api"
)

//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list

// Fetches a key-value secret (kv-v2) after authenticating to Vault with a Kubernetes service account.
func GetSecretFromVault(k8sClient k8s.Client, role, path, sa, namespace string) (map[string]interface{}, error) {
	// If set, the VAULT_ADDR environment variable will be the address that
	// your pod uses to communicate with Vault.
	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = "http://vault-service-vault-infra.apps.rhoda-411-lab.51ty.p1.openshiftapps.com"
	client, err := vault.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	var serviceAccount v1.ServiceAccount
	if err := k8sClient.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: sa}, &serviceAccount); err != nil {
		fmt.Printf("GetSecretFromVault serviceAccount:%v", err)
		return nil, err
	}
	var saSecret v1.Secret
	var saName string
	fmt.Printf("GetSecretFromVault serviceAccount:%v", serviceAccount)
	fmt.Printf("GetSecretFromVault serviceAccountName:%v", serviceAccount.Name)
	fmt.Printf("GetSecretFromVault serviceAccountNameSpace:%v", serviceAccount.Namespace)
	fmt.Printf("GetSecretFromVault serviceAccountSecret:%v", serviceAccount.Secrets)

	for _, s := range serviceAccount.Secrets {
		fmt.Printf("GetSecretFromVault secrets data:%v", s.Name)
		if strings.Contains(s.Name, "-token-") {
			saName = s.Name
			fmt.Printf("GetSecretFromVault saName:%v", saName)
			break
		}
	}

	if err := k8sClient.Get(context.Background(),
		types.NamespacedName{Name: saName, Namespace: namespace},
		&saSecret); err != nil {
		fmt.Printf("GetSecretFromVault saSecret error:%v", err)
		return nil, err
	}

	token, ok := saSecret.Data["token"]
	if !ok {
		return nil, fmt.Errorf("invalid token in secret")
	}
	k8sAuth, err := auth.NewKubernetesAuth(
		role,
		auth.WithServiceAccountToken(string(token)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	secret, err := client.Logical().Read(path)
	fmt.Printf("GetSecretFromVault path:%v", path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret from path %v: %w", path, err)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	for _, dataValues := range data {
		fmt.Printf("GetSecretFromVault dataValue:%v", dataValues)
	}
	if !ok {
		return nil, fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}
	return data, nil
}

func StoreSecrets(k8sClient k8s.Client, role, sa, namespace string) (bool, error) {
	fmt.Println("\nStoreSecrets()")

	config := vault.DefaultConfig()

	config.Address = "http://vault-service-vault-infra.apps.rhoda-411-lab.51ty.p1.openshiftapps.com"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("hvs.8EQ0V3S2ORN0jwnEdwRULq04")

	secretData := map[string]interface{}{
		"password": "test12345",
	}

	// Write a secret
	_, err = client.KVv2("dbaas/").Put(context.Background(), "test", secretData)

	//output, err := client.Logical().Write("secret/data/dbaas", secretData)
	//fmt.Println(output)
	if err != nil {
		fmt.Printf("unable to write secret: %v", err)
		return false, err
	} else {
		fmt.Println("\nSecret written successfully.")
		return true, nil

	}
}
