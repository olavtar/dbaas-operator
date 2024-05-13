ARG VERSION=23.0.1

#FROM scratch

FROM quay.io/fedora/fedora:36

FROM quay.io/jitesoft/alpine:3.18.0

FROM debian:12-slim AS builder
#
#FROM builder
#
#FROM mikefarah/yq:4 AS yq
#
#FROM debian:12-slim AS base

#FROM test