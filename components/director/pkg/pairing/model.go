package pairing

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

type RequestData struct {
	Application externalschema.Application
	Tenant      string
}

type ResponseData struct {
	Token string
}
