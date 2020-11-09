package onetimetoken

import (
	"fmt"
	"net/url"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/pkg/errors"
)

type converter struct {
	legacyConnectorURL string
}

func NewConverter(legacyConnectorURL string) *converter {
	return &converter{legacyConnectorURL}
}

func (c converter) ToGraphQLForRuntime(model model.OneTimeToken) externalschema.OneTimeTokenForRuntime {
	return externalschema.OneTimeTokenForRuntime{
		TokenWithURL: externalschema.TokenWithURL{
			Token:        model.Token,
			ConnectorURL: model.ConnectorURL,
		},
	}
}

func (c converter) ToGraphQLForApplication(model model.OneTimeToken) (externalschema.OneTimeTokenForApplication, error) {
	legacyConnectorURL, err := url.Parse(c.legacyConnectorURL)
	if err != nil {
		return externalschema.OneTimeTokenForApplication{}, errors.Wrapf(err, "while parsing string (%s) as the URL", c.legacyConnectorURL)
	}

	if legacyConnectorURL.RawQuery != "" {
		legacyConnectorURL.RawQuery += "&"
	}
	legacyConnectorURL.RawQuery += fmt.Sprintf("token=%s", model.Token)

	return externalschema.OneTimeTokenForApplication{
		TokenWithURL: externalschema.TokenWithURL{
			Token:        model.Token,
			ConnectorURL: model.ConnectorURL,
		},
		LegacyConnectorURL: legacyConnectorURL.String(),
	}, nil
}
