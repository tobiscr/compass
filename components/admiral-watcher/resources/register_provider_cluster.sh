#!/bin/bash

set -e

function cleanup() {
  export KUBECONFIG=
}

trap cleanup EXIT

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ADMIRAL_HOME=${ROOT_PATH}/admiral-install-v1.0

if [ "$#" -gt "0" ]; then
  remote=$1
fi

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml
remote_cluster=${ROOT_PATH}/kubeconfigs/${remote}

export KUBECONFIG=$remote_cluster
kubectl apply -f $ADMIRAL_HOME/yaml/remotecluster.yaml

export KUBECONFIG=$remote_cluster
$ADMIRAL_HOME/scripts/cluster-secret.sh $admiral_cluster $remote_cluster admiral