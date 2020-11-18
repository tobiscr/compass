#!/bin/bash

set -e

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -gt "0" ]; then
  dep_name=$1
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml

export KUBECONFIG=$admiral_cluster

kubectl get dependency $1 --ignore-not-found