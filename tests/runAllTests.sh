#!/usr/local/bin/zsh

SUITE_NAME="compass-e2e-tests-suite"

kubectl delete clustertestsuite $SUITE_NAME

cat <<EOF | kubectl apply -f -
apiVersion: testing.kyma-project.io/v1alpha1
kind: ClusterTestSuite
metadata:
  name: $SUITE_NAME
spec:
  concurrency: 10
  maxRetries: 10
EOF
