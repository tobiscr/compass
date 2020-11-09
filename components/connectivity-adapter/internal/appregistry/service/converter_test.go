package service_test

import (
	"encoding/json"
	"testing"

	"github.com/kyma-incubator/compass/components/connectivity-adapter/internal/appregistry/model"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/internal/appregistry/service"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_DetailsToGraphQLCreateInput(t *testing.T) {
	additionalQueryParamsSerialized := externalschema.QueryParamsSerialized(`{"q1":["a","b"],"q2":["c","d"]}`)
	additionalHeadersSerialized := externalschema.HttpHeadersSerialized(`{"h1":["e","f"],"h2":["g","h"]}`)

	type testCase struct {
		given    model.ServiceDetails
		expected externalschema.PackageCreateInput
	}

	conv := service.NewConverter()

	for name, tc := range map[string]testCase{
		"name and description propagated to api": {
			given: model.ServiceDetails{Name: "name", Description: "description", Api: &model.API{}},
			expected: externalschema.PackageCreateInput{
				Name:                "name",
				Description:         ptrString("description"),
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Name:        "name",
						Description: ptrString("description"),
					},
				},
			},
		},
		"API with only URL provided": {
			given: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						TargetURL: "http://target.url",
					},
				},
			},
		},
		"API with empty credentials": {
			given: model.ServiceDetails{
				Api: &model.API{
					TargetUrl:   "http://target.url",
					Credentials: &model.CredentialsWithCSRF{},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						TargetURL: "http://target.url",
					},
				},
			},
		},
		"ODATA API provided": {
			given: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
					ApiType:   "ODATA",
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						TargetURL: "http://target.url",
						Spec: &externalschema.APISpecInput{
							Type:   externalschema.APISpecTypeOdata,
							Format: externalschema.SpecFormatXML,
							FetchRequest: &externalschema.FetchRequestInput{
								URL: "http://target.url/$metadata",
							},
						},
					},
				},
			},
		},

		"API other than ODATA provided": {
			given: model.ServiceDetails{
				Api: &model.API{
					ApiType: "anything else",
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatYaml,
						},
					},
				},
			},
		},

		"API with directly spec provided in YAML": {
			given: model.ServiceDetails{
				Api: &model.API{
					Spec: json.RawMessage(`openapi: "3.0.0"`),
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							Data:   ptrClob(externalschema.CLOB(`openapi: "3.0.0"`)),
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatYaml,
						},
					},
				},
			},
		},

		"API with directly spec provided in JSON": {
			given: model.ServiceDetails{
				Api: &model.API{
					Spec: json.RawMessage(`{"spec":"v0.0.1"}`),
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							Data:   ptrClob(externalschema.CLOB(`{"spec":"v0.0.1"}`)),
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatJSON,
						},
					},
				},
			},
		},

		"API with directly spec provided in XML": {
			given: model.ServiceDetails{
				Api: &model.API{
					Spec: json.RawMessage(`<spec></spec>"`),
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							Data:   ptrClob(externalschema.CLOB(`<spec></spec>"`)),
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatXML,
						},
					},
				},
			},
		},

		"API with query params and headers stored in old fields": {
			given: model.ServiceDetails{
				Name: "foo",
				Api: &model.API{
					QueryParameters: &map[string][]string{
						"q1": {"a", "b"},
						"q2": {"c", "d"},
					},
					Headers: &map[string][]string{
						"h1": {"e", "f"},
						"h2": {"g", "h"},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				Name: "foo",
				DefaultInstanceAuth: &externalschema.AuthInput{
					AdditionalQueryParamsSerialized: &additionalQueryParamsSerialized,
					AdditionalHeadersSerialized:     &additionalHeadersSerialized,
				},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{Name: "foo"},
				},
			},
		},
		"API with query params and headers stored in the new fields": {
			given: model.ServiceDetails{
				Api: &model.API{
					RequestParameters: &model.RequestParameters{
						QueryParameters: &map[string][]string{
							"q1": {"a", "b"},
							"q2": {"c", "d"},
						},
						Headers: &map[string][]string{
							"h1": {"e", "f"},
							"h2": {"g", "h"},
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{
					AdditionalQueryParamsSerialized: &additionalQueryParamsSerialized,
					AdditionalHeadersSerialized:     &additionalHeadersSerialized,
				},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{},
				},
			},
		},
		"API with query params and headers stored in old and new fields": {
			given: model.ServiceDetails{
				Api: &model.API{
					RequestParameters: &model.RequestParameters{
						QueryParameters: &map[string][]string{
							"q1": {"a", "b"},
							"q2": {"c", "d"},
						},
						Headers: &map[string][]string{
							"h1": {"e", "f"},
							"h2": {"g", "h"},
						}},
					QueryParameters: &map[string][]string{
						"old": {"old"},
					},
					Headers: &map[string][]string{
						"old": {"old"},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{
					AdditionalQueryParamsSerialized: &additionalQueryParamsSerialized,
					AdditionalHeadersSerialized:     &additionalHeadersSerialized,
				},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{},
				},
			},
		},
		"API protected with basic": {
			given: model.ServiceDetails{
				Api: &model.API{
					Credentials: &model.CredentialsWithCSRF{
						BasicWithCSRF: &model.BasicAuthWithCSRF{
							BasicAuth: model.BasicAuth{
								Username: "user",
								Password: "password",
							},
							CSRFInfo: &model.CSRFInfo{TokenEndpointURL: "foo.bar"},
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{
					Credential: &externalschema.CredentialDataInput{
						Basic: &externalschema.BasicCredentialDataInput{
							Username: "user",
							Password: "password",
						},
					},
					RequestAuth: &externalschema.CredentialRequestAuthInput{
						Csrf: &externalschema.CSRFTokenCredentialRequestAuthInput{TokenEndpointURL: "foo.bar"},
					},
				},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{},
				},
			},
		},
		"API protected with oauth": {
			given: model.ServiceDetails{
				Api: &model.API{
					Credentials: &model.CredentialsWithCSRF{
						OauthWithCSRF: &model.OauthWithCSRF{
							Oauth: model.Oauth{
								ClientID:     "client_id",
								ClientSecret: "client_secret",
								URL:          "http://oauth.url",
								RequestParameters: &model.RequestParameters{ // TODO this field is not mapped at all
									QueryParameters: &map[string][]string{
										"q1": {"a", "b"},
										"q2": {"c", "d"},
									},
									Headers: &map[string][]string{
										"h1": {"e", "f"},
										"h2": {"g", "h"},
									},
								},
							},
							CSRFInfo: &model.CSRFInfo{TokenEndpointURL: "foo.bar"},
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{
					Credential: &externalschema.CredentialDataInput{
						Oauth: &externalschema.OAuthCredentialDataInput{
							URL:          "http://oauth.url",
							ClientID:     "client_id",
							ClientSecret: "client_secret",
						},
					},
					RequestAuth: &externalschema.CredentialRequestAuthInput{
						Csrf: &externalschema.CSRFTokenCredentialRequestAuthInput{TokenEndpointURL: "foo.bar"},
					},
				},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{},
				},
			},
		},
		"API specification mapped to fetch request": {
			given: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://specification.url",
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							FetchRequest: &externalschema.FetchRequestInput{
								URL: "http://specification.url",
							},
							Format: externalschema.SpecFormatJSON,
							Type:   externalschema.APISpecTypeOpenAPI,
						},
					},
				},
			},
		},
		"API specification with basic auth converted to fetch request": {
			given: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://specification.url",
					SpecificationCredentials: &model.Credentials{
						Basic: &model.BasicAuth{
							Username: "username",
							Password: "password",
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							FetchRequest: &externalschema.FetchRequestInput{
								URL: "http://specification.url",
								Auth: &externalschema.AuthInput{
									Credential: &externalschema.CredentialDataInput{
										Basic: &externalschema.BasicCredentialDataInput{
											Username: "username",
											Password: "password",
										},
									},
								},
							},
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatJSON,
						},
					},
				},
			},
		},
		"API specification with oauth converted to fetch request": {
			given: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://specification.url",
					SpecificationCredentials: &model.Credentials{
						Oauth: &model.Oauth{
							URL:               "http://oauth.url",
							ClientID:          "client_id",
							ClientSecret:      "client_secret",
							RequestParameters: nil, // TODO not supported
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							FetchRequest: &externalschema.FetchRequestInput{
								URL: "http://specification.url",
								Auth: &externalschema.AuthInput{
									Credential: &externalschema.CredentialDataInput{
										Oauth: &externalschema.OAuthCredentialDataInput{
											URL:          "http://oauth.url",
											ClientID:     "client_id",
											ClientSecret: "client_secret",
										},
									},
								},
							},
							Type:   externalschema.APISpecTypeOpenAPI,
							Format: externalschema.SpecFormatJSON,
						},
					},
				},
			},
		},
		"API specification with request parameters converted to fetch request": {
			given: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://specification.url",
					SpecificationRequestParameters: &model.RequestParameters{
						QueryParameters: &map[string][]string{
							"q1": {"a", "b"},
							"q2": {"c", "d"},
						},
						Headers: &map[string][]string{
							"h1": {"e", "f"},
							"h2": {"g", "h"},
						},
					},
				},
			},
			expected: externalschema.PackageCreateInput{
				DefaultInstanceAuth: &externalschema.AuthInput{},
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{
						Spec: &externalschema.APISpecInput{
							FetchRequest: &externalschema.FetchRequestInput{
								URL: "http://specification.url",
								Auth: &externalschema.AuthInput{
									AdditionalQueryParamsSerialized: &additionalQueryParamsSerialized,
									AdditionalHeadersSerialized:     &additionalHeadersSerialized,
								},
							},
							Format: externalschema.SpecFormatJSON,
							Type:   externalschema.APISpecTypeOpenAPI,
						},
					},
				},
			},
		},
		"Event": {
			given: model.ServiceDetails{
				Name: "foo",
				Events: &model.Events{
					Spec: json.RawMessage(`asyncapi: "1.2.0"`),
				},
			},
			expected: externalschema.PackageCreateInput{
				Name:                "foo",
				DefaultInstanceAuth: &externalschema.AuthInput{},
				EventDefinitions: []*externalschema.EventDefinitionInput{
					{
						Name: "foo",
						Spec: &externalschema.EventSpecInput{
							Data:   ptrClob(`asyncapi: "1.2.0"`),
							Type:   externalschema.EventSpecTypeAsyncAPI,
							Format: externalschema.SpecFormatYaml,
						},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// WHEN
			actual, err := conv.DetailsToGraphQLCreateInput(tc.given)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestConverter_GraphQLToServiceDetails(t *testing.T) {
	type testCase struct {
		given    externalschema.PackageExt
		expected model.ServiceDetails
	}
	conv := service.NewConverter()

	testSvcRef := service.LegacyServiceReference{
		ID:         "foo",
		Identifier: "",
	}

	for name, tc := range map[string]testCase{
		"name and description is loaded from Package": {
			given: externalschema.PackageExt{
				Package: externalschema.Package{Name: "foo", Description: ptrString("description")},
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							APIDefinition: externalschema.APIDefinition{},
						},
					},
				},
			},
			expected: model.ServiceDetails{
				Name:        "foo",
				Description: "description",
				Api:         &model.API{},
				Labels:      emptyLabels(),
			},
		},
		"simple API": {
			given: externalschema.PackageExt{
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							APIDefinition: externalschema.APIDefinition{
								TargetURL: "http://target.url",
							},
						},
					},
				},
			},
			expected: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with additional headers and query params": {
			given: externalschema.PackageExt{
				Package: externalschema.Package{
					DefaultInstanceAuth: &externalschema.Auth{
						AdditionalQueryParams: &externalschema.QueryParams{
							"q1": []string{"a", "b"},
							"q2": []string{"c", "d"},
						},
						AdditionalHeaders: &externalschema.HttpHeaders{
							"h1": []string{"e", "f"},
							"h2": []string{"g", "h"},
						},
					},
				},
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							APIDefinition: externalschema.APIDefinition{
								TargetURL: "http://target.url",
							},
						},
					},
				},
			},
			expected: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
					Headers: &map[string][]string{
						"h1": {"e", "f"},
						"h2": {"g", "h"}},
					QueryParameters: &map[string][]string{
						"q1": {"a", "b"},
						"q2": {"c", "d"},
					},
					RequestParameters: &model.RequestParameters{
						Headers: &map[string][]string{
							"h1": {"e", "f"},
							"h2": {"g", "h"}},
						QueryParameters: &map[string][]string{
							"q1": {"a", "b"},
							"q2": {"c", "d"},
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with Basic Auth": {
			given: externalschema.PackageExt{
				Package: externalschema.Package{
					DefaultInstanceAuth: &externalschema.Auth{
						Credential: &externalschema.BasicCredentialData{
							Username: "username",
							Password: "password",
						},
					},
				},
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							APIDefinition: externalschema.APIDefinition{
								TargetURL: "http://target.url",
							},
						},
					},
				},
			},
			expected: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
					Credentials: &model.CredentialsWithCSRF{
						BasicWithCSRF: &model.BasicAuthWithCSRF{
							BasicAuth: model.BasicAuth{
								Username: "username",
								Password: "password",
							},
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with Oauth": {
			given: externalschema.PackageExt{
				Package: externalschema.Package{
					DefaultInstanceAuth: &externalschema.Auth{
						Credential: &externalschema.OAuthCredentialData{
							URL:          "http://oauth.url",
							ClientID:     "client_id",
							ClientSecret: "client_secret",
						},
					},
				},
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							APIDefinition: externalschema.APIDefinition{
								TargetURL: "http://target.url",
							},
						},
					},
				},
			},
			expected: model.ServiceDetails{
				Api: &model.API{
					TargetUrl: "http://target.url",
					Credentials: &model.CredentialsWithCSRF{
						OauthWithCSRF: &model.OauthWithCSRF{
							Oauth: model.Oauth{
								URL:          "http://oauth.url",
								ClientID:     "client_id",
								ClientSecret: "client_secret",
							},
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with FetchRequest (query params and headers)": {
			given: externalschema.PackageExt{
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							Spec: &externalschema.APISpecExt{
								FetchRequest: &externalschema.FetchRequest{
									URL: "http://apispec.url",
									Auth: &externalschema.Auth{
										AdditionalQueryParams: &externalschema.QueryParams{
											"q1": {"a", "b"},
											"q2": {"c", "d"},
										},
										AdditionalHeaders: &externalschema.HttpHeaders{
											"h1": {"e", "f"},
											"h2": {"g", "h"},
										},
									},
								}}},
					},
				}},
			expected: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://apispec.url",
					SpecificationRequestParameters: &model.RequestParameters{
						Headers: &map[string][]string{
							"h1": {"e", "f"},
							"h2": {"g", "h"}},
						QueryParameters: &map[string][]string{
							"q1": {"a", "b"},
							"q2": {"c", "d"},
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with Fetch Request protected with Basic Auth": {
			given: externalschema.PackageExt{
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							Spec: &externalschema.APISpecExt{
								FetchRequest: &externalschema.FetchRequest{
									URL: "http://apispec.url",
									Auth: &externalschema.Auth{
										Credential: &externalschema.BasicCredentialData{
											Username: "username",
											Password: "password",
										},
									},
								}}},
					},
				}},
			expected: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://apispec.url",
					SpecificationCredentials: &model.Credentials{
						Basic: &model.BasicAuth{
							Username: "username",
							Password: "password",
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"simple API with Fetch Request protected with Oauth": {
			given: externalschema.PackageExt{
				APIDefinitions: externalschema.APIDefinitionPageExt{
					Data: []*externalschema.APIDefinitionExt{
						{
							Spec: &externalschema.APISpecExt{
								FetchRequest: &externalschema.FetchRequest{
									URL: "http://apispec.url",
									Auth: &externalschema.Auth{
										Credential: &externalschema.OAuthCredentialData{
											URL:          "http://oauth.url",
											ClientID:     "client_id",
											ClientSecret: "client_secret",
										},
									},
								}}},
					},
				}},
			expected: model.ServiceDetails{
				Api: &model.API{
					SpecificationUrl: "http://apispec.url",
					SpecificationCredentials: &model.Credentials{
						Oauth: &model.Oauth{
							URL:          "http://oauth.url",
							ClientID:     "client_id",
							ClientSecret: "client_secret",
						},
					},
				},
				Labels: emptyLabels(),
			},
		},
		"events": {
			given: externalschema.PackageExt{
				EventDefinitions: externalschema.EventAPIDefinitionPageExt{
					Data: []*externalschema.EventAPIDefinitionExt{{
						Spec: &externalschema.EventAPISpecExt{
							EventSpec: externalschema.EventSpec{
								Data: ptrClob(`asyncapi: "1.2.0"`),
							},
						}},
					},
				},
			},
			expected: model.ServiceDetails{
				Events: &model.Events{
					Spec: json.RawMessage(`asyncapi: "1.2.0"`),
				},
				Labels: emptyLabels(),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// WHEN
			actual, err := conv.GraphQLToServiceDetails(tc.given, testSvcRef)
			// THEN
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}

	t.Run("identifier provided", func(t *testing.T) {
		in := externalschema.PackageExt{
			APIDefinitions: externalschema.APIDefinitionPageExt{
				Data: []*externalschema.APIDefinitionExt{
					{
						APIDefinition: externalschema.APIDefinition{
							TargetURL: "http://target.url",
						},
					},
				},
			},
		}
		inSvcRef := service.LegacyServiceReference{
			ID:         "foo",
			Identifier: "test",
		}
		expected := model.ServiceDetails{
			Identifier: "test",
			Api: &model.API{
				TargetUrl: "http://target.url",
			},
			Labels: emptyLabels(),
		}
		// WHEN
		actual, err := conv.GraphQLToServiceDetails(in, inSvcRef)
		// THEN
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestConverter_ServiceDetailsToService(t *testing.T) {
	//GIVEN
	input := model.ServiceDetails{
		Provider:         "provider",
		Name:             "name",
		Description:      "description",
		ShortDescription: "short description",
		Identifier:       "identifie",
		Labels:           &map[string]string{"blalb": "blalba"},
	}
	id := "id"

	//WHEN
	conv := service.NewConverter()
	output, err := conv.ServiceDetailsToService(input, id)

	//THEN
	require.NoError(t, err)
	assert.Equal(t, input.Provider, output.Provider)
	assert.Equal(t, input.Name, output.Name)
	assert.Equal(t, input.Description, output.Description)
	assert.Equal(t, input.Identifier, output.Identifier)
	assert.Equal(t, input.Labels, output.Labels)
}

func TestConverter_GraphQLCreateInputToUpdateInput(t *testing.T) {
	desc := "Desc"
	schema := externalschema.JSONSchema("foo")
	auth := externalschema.AuthInput{Credential: &externalschema.CredentialDataInput{Basic: &externalschema.BasicCredentialDataInput{
		Username: "foo",
		Password: "bar",
	}}}
	in := externalschema.PackageCreateInput{
		Name:                           "foo",
		Description:                    &desc,
		InstanceAuthRequestInputSchema: &schema,
		DefaultInstanceAuth:            &auth,
	}
	expected := externalschema.PackageUpdateInput{
		Name:                           "foo",
		Description:                    &desc,
		InstanceAuthRequestInputSchema: &schema,
		DefaultInstanceAuth:            &auth,
	}

	conv := service.NewConverter()

	res := conv.GraphQLCreateInputToUpdateInput(in)

	assert.Equal(t, expected, res)
}

func TestConverter_DetailsToGraphQLInput_TestSpecsRecognition(t *testing.T) {
	// GIVEN
	conv := service.NewConverter()

	// API
	apiCases := []struct {
		Name           string
		InputAPI       model.API
		ExpectedType   externalschema.APISpecType
		ExpectedFormat externalschema.SpecFormat
	}{
		{
			Name:           "OpenAPI + YAML",
			InputAPI:       fixAPIOpenAPIYAML(),
			ExpectedType:   externalschema.APISpecTypeOpenAPI,
			ExpectedFormat: externalschema.SpecFormatYaml,
		},
		{
			Name:           "OpenAPI + JSON",
			InputAPI:       fixAPIOpenAPIJSON(),
			ExpectedType:   externalschema.APISpecTypeOpenAPI,
			ExpectedFormat: externalschema.SpecFormatJSON,
		},
		{
			Name:           "OData + XML",
			InputAPI:       fixAPIODataXML(),
			ExpectedType:   externalschema.APISpecTypeOdata,
			ExpectedFormat: externalschema.SpecFormatXML,
		},
	}

	for _, testCase := range apiCases {
		t.Run(testCase.Name, func(t *testing.T) {
			in := model.ServiceDetails{Api: &testCase.InputAPI}

			// WHEN
			out, err := conv.DetailsToGraphQLCreateInput(in)

			// THEN
			require.NoError(t, err)
			require.Len(t, out.APIDefinitions, 1)
			require.NotNil(t, out.APIDefinitions[0].Spec)
			assert.Equal(t, testCase.ExpectedType, out.APIDefinitions[0].Spec.Type)
			assert.Equal(t, testCase.ExpectedFormat, out.APIDefinitions[0].Spec.Format)
		})
	}

	// Events
	eventsCases := []struct {
		Name           string
		InputEvents    model.Events
		ExpectedType   externalschema.EventSpecType
		ExpectedFormat externalschema.SpecFormat
	}{
		{
			Name:           "Async API + JSON",
			InputEvents:    fixEventsAsyncAPIJSON(),
			ExpectedType:   externalschema.EventSpecTypeAsyncAPI,
			ExpectedFormat: externalschema.SpecFormatJSON,
		},
		{
			Name:           "Async API + YAML",
			InputEvents:    fixEventsAsyncAPIYAML(),
			ExpectedType:   externalschema.EventSpecTypeAsyncAPI,
			ExpectedFormat: externalschema.SpecFormatYaml,
		},
	}

	for _, testCase := range eventsCases {
		t.Run(testCase.Name, func(t *testing.T) {
			in := model.ServiceDetails{Events: &testCase.InputEvents}

			// WHEN
			out, err := conv.DetailsToGraphQLCreateInput(in)

			// THEN
			require.NoError(t, err)
			require.Len(t, out.EventDefinitions, 1)
			require.NotNil(t, out.EventDefinitions[0].Spec)
			assert.Equal(t, testCase.ExpectedType, out.EventDefinitions[0].Spec.Type)
			assert.Equal(t, testCase.ExpectedFormat, out.EventDefinitions[0].Spec.Format)
		})
	}
}

func emptyLabels() *map[string]string {
	return &map[string]string{}
}

func ptrString(in string) *string {
	return &in
}

func ptrClob(in externalschema.CLOB) *externalschema.CLOB {
	return &in
}
