module github.com/kyma-incubator/compass/tests-original

go 1.13

require (
	github.com/99designs/gqlgen v0.10.2 // indirect
	github.com/avast/retry-go v2.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/kyma-incubator/compass v0.0.0-20201123100504-8408fe755023 // indirect
	github.com/kyma-incubator/compass/components/director v0.0.0-20201109133626-4876e6d3caae
	github.com/kyma-incubator/compass/components/gateway v0.0.0-20201123100504-8408fe755023 // indirect
	github.com/kyma-incubator/compass/tests/connectivity-adapter v0.0.0-20201123100504-8408fe755023
	github.com/machinebox/graphql v0.2.3-0.20181106130121-3a9253180225
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/vrischmann/envconfig v1.2.0
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	k8s.io/api v0.19.4 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8
)
