#!/usr/local/bin/zsh

TEST_DEF_NAME=$1
TEST_DEF_NAMESPACE=${2:-compass-system}

kubectl delete clustertestsuite $TEST_DEF_NAME -n TEST_DEF_NAMESPACE

echo "NAME: $TEST_DEF_NAME"
echo "NAMESPACE: $TEST_DEF_NAMESPACE"

cat <<EOF | kubectl apply -f -
apiVersion: testing.kyma-project.io/v1alpha1
kind: ClusterTestSuite
metadata:
  name: $TEST_DEF_NAME
spec:
  concurrency: 10
  maxRetries: 10
  selectors:
    matchNames:
    - name: $TEST_DEF_NAME
      namespace: $TEST_DEF_NAMESPACE
EOF
