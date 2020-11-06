#!/bin/bash

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  remote_provider=$1
fi

if [ "$#" -gt "1" ]; then
  remote_consumer=$2
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml
remote_provider_cluster=${ROOT_PATH}/kubeconfigs/${remote_provider}
remote_consumer_cluster=${ROOT_PATH}/kubeconfigs/${remote_consumer}

export KUBECONFIG=$remote_provider_cluster
kubectl delete namespace admiral-sync
kubectl delete namespace sample
CLUSTER_NAME=$(kubectl config view --minify=true -o "jsonpath={.clusters[].name}")

export KUBECONFIG=$admiral_cluster
kubectl delete secret ${CLUSTER_NAME} -n admiral

export KUBECONFIG=$remote_consumer_cluster
kubectl delete se --all -n admiral-sync
kubectl delete dr --all -n admiral-sync

export KUBECONFIG=$admiral_cluster
kubectl rollout restart deployment/admiral -n admiral

export KUBECONFIG=$remote_consumer_cluster
kubectl delete se --all -n admiral-sync
kubectl delete dr --all -n admiral-sync

