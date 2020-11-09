package externalschema_test

import (
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation/inputvalidationtest"
	"github.com/stretchr/testify/require"
)

func TestWebhookInput_Validate_Type(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         externalschema.ApplicationWebhookType
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         externalschema.ApplicationWebhookTypeConfigurationChanged,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid - Empty",
			Value:         inputvalidationtest.EmptyString,
			ExpectedValid: false,
		},
		{
			Name:          "Invalid - Not enum",
			Value:         "invalid",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := fixValidWebhookInput()
			sut.Type = testCase.Value
			//WHEN
			err := sut.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestWebhookInput_Validate_URL(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         inputvalidationtest.ValidURL,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid - Empty string",
			Value:         inputvalidationtest.EmptyString,
			ExpectedValid: false,
		},
		{
			Name:          "Invalid - Invalid URL",
			Value:         inputvalidationtest.InvalidURL,
			ExpectedValid: false,
		},
		{
			Name:          "Invalid - Too long",
			Value:         inputvalidationtest.URL257Long,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := fixValidWebhookInput()
			sut.URL = testCase.Value
			//WHEN
			err := sut.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestWebhookInput_Validate_Auth(t *testing.T) {
	auth := fixValidAuthInput()
	testCases := []struct {
		Name          string
		Value         *externalschema.AuthInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         &auth,
			ExpectedValid: true,
		},
		{
			Name:          "ExpectedValid - nil",
			Value:         nil,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid - Nested validation error",
			Value:         &externalschema.AuthInput{Credential: &externalschema.CredentialDataInput{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := fixValidWebhookInput()
			sut.Auth = testCase.Value
			//WHEN
			err := sut.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func fixValidWebhookInput() externalschema.WebhookInput {
	return externalschema.WebhookInput{
		Type: externalschema.ApplicationWebhookTypeConfigurationChanged,
		URL:  inputvalidationtest.ValidURL,
	}
}
