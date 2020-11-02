#!/bin/bash

set -e

function cleanup() {
  rm -rf dep.yaml
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  dep=$1
fi

if [ "$#" -gt "1" ]; then
  remote=$1
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml
remote_cluster=${ROOT_PATH}/kubeconfigs/${remote}

export KUBECONFIG=$admiral_cluster

kubectl delete dependency ${dep}
kubectl rollout restart deployment/admiral -n admiral

export KUBECONFIG=$remote_cluster

kubectl delete -all serviceentries -n admiral-sync
