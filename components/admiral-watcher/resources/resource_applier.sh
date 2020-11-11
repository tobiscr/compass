#!/bin/bash

set -e

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  RESOURCE_PATH=$1
fi

runtime_cluster=${ROOT_PATH}/kubeconfigs/runtime.yaml

export KUBECONFIG=$runtime_cluster

kubectl apply -f $RESOURCE_PATH