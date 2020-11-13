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

admiral_cluster=${ROOT_PATH}/kubeconfigs/admiral.yaml

export KUBECONFIG=$admiral_cluster

echo $dep

cat <<EOF | kubectl apply -f -
$dep
EOF

#kubectl apply -f dep.yaml