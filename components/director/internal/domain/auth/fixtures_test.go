package auth_test

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

var (
	authUsername = "user"
	authPassword = "password"
	authEndpoint = "url"
	authMap      = map[string][]string{
		"foo": {"bar", "baz"},
		"too": {"tar", "taz"},
	}
	authMapSerialized     = "{\"foo\":[\"bar\",\"baz\"],\"too\":[\"tar\",\"taz\"]}"
	authHeaders           = externalschema.HttpHeaders(authMap)
	authHeadersSerialized = externalschema.HttpHeadersSerialized(authMapSerialized)
	authParams            = externalschema.QueryParams(authMap)
	authParamsSerialized  = externalschema.QueryParamsSerialized(authMapSerialized)
)

func fixDetailedAuth() *model.Auth {
	return &model.Auth{
		Credential: model.CredentialData{
			Basic: &model.BasicCredentialData{
				Username: authUsername,
				Password: authPassword,
			},
			Oauth: nil,
		},
		AdditionalHeaders:     authMap,
		AdditionalQueryParams: authMap,
		RequestAuth: &model.CredentialRequestAuth{
			Csrf: &model.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: authEndpoint,
				Credential: model.CredentialData{
					Basic: &model.BasicCredentialData{
						Username: authUsername,
						Password: authPassword,
					},
					Oauth: nil,
				},
				AdditionalHeaders:     authMap,
				AdditionalQueryParams: authMap,
			},
		},
	}
}

func fixDetailedGQLAuth() *externalschema.Auth {
	return &externalschema.Auth{
		Credential: externalschema.BasicCredentialData{
			Username: authUsername,
			Password: authPassword,
		},
		AdditionalHeaders:               &authHeaders,
		AdditionalHeadersSerialized:     &authHeadersSerialized,
		AdditionalQueryParams:           &authParams,
		AdditionalQueryParamsSerialized: &authParamsSerialized,
		RequestAuth: &externalschema.CredentialRequestAuth{
			Csrf: &externalschema.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: authEndpoint,
				Credential: externalschema.BasicCredentialData{
					Username: authUsername,
					Password: authPassword,
				},
				AdditionalHeaders:     &authHeaders,
				AdditionalQueryParams: &authParams,
			},
		},
	}
}

func fixDetailedAuthInput() *model.AuthInput {
	return &model.AuthInput{
		Credential: &model.CredentialDataInput{
			Basic: &model.BasicCredentialDataInput{
				Username: authUsername,
				Password: authPassword,
			},
			Oauth: nil,
		},
		AdditionalHeaders:     authMap,
		AdditionalQueryParams: authMap,
		RequestAuth: &model.CredentialRequestAuthInput{
			Csrf: &model.CSRFTokenCredentialRequestAuthInput{
				TokenEndpointURL: authEndpoint,
				Credential: &model.CredentialDataInput{
					Basic: &model.BasicCredentialDataInput{
						Username: authUsername,
						Password: authPassword,
					},
					Oauth: nil,
				},
				AdditionalHeaders:     authMap,
				AdditionalQueryParams: authMap,
			},
		},
	}
}

func fixDetailedGQLAuthInput() *externalschema.AuthInput {
	return &externalschema.AuthInput{
		Credential: &externalschema.CredentialDataInput{
			Basic: &externalschema.BasicCredentialDataInput{
				Username: authUsername,
				Password: authPassword,
			},
			Oauth: nil,
		},
		AdditionalHeadersSerialized:     &authHeadersSerialized,
		AdditionalQueryParamsSerialized: &authParamsSerialized,
		RequestAuth: &externalschema.CredentialRequestAuthInput{
			Csrf: &externalschema.CSRFTokenCredentialRequestAuthInput{
				TokenEndpointURL: authEndpoint,
				Credential: &externalschema.CredentialDataInput{
					Basic: &externalschema.BasicCredentialDataInput{
						Username: authUsername,
						Password: authPassword,
					},
					Oauth: nil,
				},
				AdditionalHeadersSerialized:     &authHeadersSerialized,
				AdditionalQueryParamsSerialized: &authParamsSerialized,
			},
		},
	}
}

func fixDetailedGQLAuthInputDeprecated() *externalschema.AuthInput {
	return &externalschema.AuthInput{
		Credential: &externalschema.CredentialDataInput{
			Basic: &externalschema.BasicCredentialDataInput{
				Username: authUsername,
				Password: authPassword,
			},
			Oauth: nil,
		},
		AdditionalHeaders:     &authHeaders,
		AdditionalQueryParams: &authParams,
		RequestAuth: &externalschema.CredentialRequestAuthInput{
			Csrf: &externalschema.CSRFTokenCredentialRequestAuthInput{
				TokenEndpointURL: authEndpoint,
				Credential: &externalschema.CredentialDataInput{
					Basic: &externalschema.BasicCredentialDataInput{
						Username: authUsername,
						Password: authPassword,
					},
					Oauth: nil,
				},
				AdditionalHeaders:     &authHeaders,
				AdditionalQueryParams: &authParams,
			},
		},
	}
}
