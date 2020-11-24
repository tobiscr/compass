module github.com/kyma-incubator/compass/tests

go 1.14

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/kyma-incubator/compass v0.0.0-20201123100504-8408fe755023
	github.com/kyma-incubator/compass/components/director v0.0.0-20201123100504-8408fe755023
	github.com/kyma-incubator/compass/components/gateway v0.0.0-20200429083609-7d80a85180c6
	github.com/kyma-incubator/compass/tests/director v0.0.0-20201123100504-8408fe755023
	github.com/machinebox/graphql v0.2.3-0.20181106130121-3a9253180225
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/vrischmann/envconfig v1.3.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)
