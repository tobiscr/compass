#!/bin/bash

# set -e

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  RESOURCE_PATH=$1
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/runtime.yaml

export KUBECONFIG=$admiral_cluster

kubectl delete -f $RESOURCE_PATH