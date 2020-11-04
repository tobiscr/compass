#!/bin/bash

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  remote=$1
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml
remote_cluster=${ROOT_PATH}/kubeconfigs/${remote}

export KUBECONFIG=$remote_cluster

kubectl delete namespace admiral-sync
kubectl delete namespace sample
CLUSTER_NAME=$(kubectl config view --minify=true -o "jsonpath={.clusters[].name}")

export KUBECONFIG=$admiral_cluster

kubectl delete secret ${CLUSTER_NAME} -n admiral

