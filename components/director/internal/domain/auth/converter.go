package auth

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	"github.com/pkg/errors"
)

type converter struct {
}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGraphQL(in *model.Auth) (*externalschema.Auth, error) {
	if in == nil {
		return nil, nil
	}

	var headers *externalschema.HttpHeaders
	var headersSerialized *externalschema.HttpHeadersSerialized
	if len(in.AdditionalHeaders) != 0 {
		var value externalschema.HttpHeaders = in.AdditionalHeaders
		headers = &value

		serialized, err := externalschema.NewHttpHeadersSerialized(in.AdditionalHeaders)
		if err != nil {
			return nil, errors.Wrap(err, "while marshaling AdditionalHeaders")
		}
		headersSerialized = &serialized
	}

	var params *externalschema.QueryParams
	var paramsSerialized *externalschema.QueryParamsSerialized
	if len(in.AdditionalQueryParams) != 0 {
		var value externalschema.QueryParams = in.AdditionalQueryParams
		params = &value

		serialized, err := externalschema.NewQueryParamsSerialized(in.AdditionalQueryParams)
		if err != nil {
			return nil, errors.Wrap(err, "while marshaling AdditionalQueryParams")
		}
		paramsSerialized = &serialized
	}

	return &externalschema.Auth{
		Credential:                      c.credentialToGraphQL(in.Credential),
		AdditionalHeaders:               headers,
		AdditionalHeadersSerialized:     headersSerialized,
		AdditionalQueryParams:           params,
		AdditionalQueryParamsSerialized: paramsSerialized,
		RequestAuth:                     c.requestAuthToGraphQL(in.RequestAuth),
	}, nil
}

func (c *converter) InputFromGraphQL(in *externalschema.AuthInput) (*model.AuthInput, error) {
	if in == nil {
		return nil, nil
	}

	credential := c.credentialInputFromGraphQL(in.Credential)

	additionalHeaders, err := c.headersFromGraphQL(in.AdditionalHeaders, in.AdditionalHeadersSerialized)
	if err != nil {
		return nil, errors.Wrap(err, "while converting AdditionalHeaders from GraphQL input")
	}

	additionalQueryParams, err := c.queryParamsFromGraphQL(in.AdditionalQueryParams, in.AdditionalQueryParamsSerialized)
	if err != nil {
		return nil, errors.Wrap(err, "while converting AdditionalQueryParams from GraphQL input")
	}

	reqAuth, err := c.requestAuthInputFromGraphQL(in.RequestAuth)
	if err != nil {
		return nil, err
	}

	return &model.AuthInput{
		Credential:            credential,
		AdditionalHeaders:     additionalHeaders,
		AdditionalQueryParams: additionalQueryParams,
		RequestAuth:           reqAuth,
	}, nil
}

func (c *converter) requestAuthToGraphQL(in *model.CredentialRequestAuth) *externalschema.CredentialRequestAuth {
	if in == nil {
		return nil
	}

	var csrf *externalschema.CSRFTokenCredentialRequestAuth
	if in.Csrf != nil {
		var headers *externalschema.HttpHeaders
		if len(in.Csrf.AdditionalHeaders) != 0 {
			var value externalschema.HttpHeaders = in.Csrf.AdditionalHeaders
			headers = &value
		}

		var params *externalschema.QueryParams
		if len(in.Csrf.AdditionalQueryParams) != 0 {
			var value externalschema.QueryParams = in.Csrf.AdditionalQueryParams
			params = &value
		}

		csrf = &externalschema.CSRFTokenCredentialRequestAuth{
			TokenEndpointURL:      in.Csrf.TokenEndpointURL,
			AdditionalQueryParams: params,
			AdditionalHeaders:     headers,
			Credential:            c.credentialToGraphQL(in.Csrf.Credential),
		}
	}

	return &externalschema.CredentialRequestAuth{
		Csrf: csrf,
	}
}

func (c *converter) requestAuthInputFromGraphQL(in *externalschema.CredentialRequestAuthInput) (*model.CredentialRequestAuthInput, error) {
	if in == nil {
		return nil, nil
	}

	var csrf *model.CSRFTokenCredentialRequestAuthInput
	if in.Csrf != nil {
		additionalHeaders, err := c.headersFromGraphQL(in.Csrf.AdditionalHeaders, in.Csrf.AdditionalHeadersSerialized)
		if err != nil {
			return nil, errors.Wrap(err, "while converting CSRF AdditionalHeaders from GraphQL input")
		}

		additionalQueryParams, err := c.queryParamsFromGraphQL(in.Csrf.AdditionalQueryParams, in.Csrf.AdditionalQueryParamsSerialized)
		if err != nil {
			return nil, errors.Wrap(err, "while converting CSRF AdditionalQueryParams from GraphQL input")
		}

		csrf = &model.CSRFTokenCredentialRequestAuthInput{
			TokenEndpointURL:      in.Csrf.TokenEndpointURL,
			AdditionalQueryParams: additionalQueryParams,
			AdditionalHeaders:     additionalHeaders,
			Credential:            c.credentialInputFromGraphQL(in.Csrf.Credential),
		}
	}

	return &model.CredentialRequestAuthInput{
		Csrf: csrf,
	}, nil
}

func (c *converter) headersFromGraphQL(headers *externalschema.HttpHeaders, headersSerialized *externalschema.HttpHeadersSerialized) (map[string][]string, error) {
	var h map[string][]string

	if headersSerialized != nil {
		return headersSerialized.Unmarshal()
	} else if headers != nil {
		h = *headers
	}

	return h, nil
}

func (c *converter) queryParamsFromGraphQL(params *externalschema.QueryParams, paramsSerialized *externalschema.QueryParamsSerialized) (map[string][]string, error) {
	var p map[string][]string

	if paramsSerialized != nil {
		return paramsSerialized.Unmarshal()
	} else if params != nil {
		p = *params
	}

	return p, nil
}

func (c *converter) credentialInputFromGraphQL(in *externalschema.CredentialDataInput) *model.CredentialDataInput {
	if in == nil {
		return nil
	}

	var basic *model.BasicCredentialDataInput
	var oauth *model.OAuthCredentialDataInput

	if in.Basic != nil {
		basic = &model.BasicCredentialDataInput{
			Username: in.Basic.Username,
			Password: in.Basic.Password,
		}
	} else if in.Oauth != nil {
		oauth = &model.OAuthCredentialDataInput{
			URL:          in.Oauth.URL,
			ClientID:     in.Oauth.ClientID,
			ClientSecret: in.Oauth.ClientSecret,
		}
	}

	return &model.CredentialDataInput{
		Basic: basic,
		Oauth: oauth,
	}
}

func (c *converter) credentialToGraphQL(in model.CredentialData) externalschema.CredentialData {
	var credential externalschema.CredentialData
	if in.Basic != nil {
		credential = externalschema.BasicCredentialData{
			Username: in.Basic.Username,
			Password: in.Basic.Password,
		}
	} else if in.Oauth != nil {
		credential = externalschema.OAuthCredentialData{
			URL:          in.Oauth.URL,
			ClientID:     in.Oauth.ClientID,
			ClientSecret: in.Oauth.ClientSecret,
		}
	}

	return credential
}
