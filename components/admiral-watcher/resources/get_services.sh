#!/bin/bash

set -e

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

admiral_cluster=${ROOT_PATH}/kubeconfigs/runtime.yaml

export KUBECONFIG=$admiral_cluster

kubectl get svc -l 'app=admiral-watcher' -o=jsonpath='{.items[*].metadata.name}'