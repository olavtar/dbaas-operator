/*
Copyright 2023 The OpenShift Database Access Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DBaaSProviderStatus defines the observed state of a DBaaSProvider object.
type DBaaSProviderStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// DBaaSProvider defines the schema for the DBaaSProvider API.
type DBaaSProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBaaSProviderSpec   `json:"spec,omitempty"`
	Status DBaaSProviderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DBaaSProviderList contains a list of DBaaSProviders.
type DBaaSProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBaaSProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBaaSProvider{}, &DBaaSProviderList{})
}
