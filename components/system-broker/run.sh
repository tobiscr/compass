#!/usr/bin/env bash

# This script is responsible for running System Broker.

MINIKUBE_IP=$(minikube ip)
KUBECONFIG=$HOME/.kube/config
go run cmd/main.go --kubeconfig $KUBECONFIG --master $MINIKUBE_IP
