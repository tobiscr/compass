package externalschema_test

import (
	"encoding/json"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation/inputvalidationtest"
	"github.com/stretchr/testify/require"
)

var (
	validSchema   = externalschema.JSONSchema(`{"type": "string"}`)
	invalidSchema = externalschema.JSONSchema(`{invalid`)
)

func TestLabelDefinitionInput_Validate(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         externalschema.LabelDefinitionInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         fixValidLabelDefinitionInput(),
			ExpectedValid: true,
		},
		{
			Name: "ExpectedValid - Schema provided",
			Value: externalschema.LabelDefinitionInput{
				Key:    "ok",
				Schema: &validSchema,
			},
			ExpectedValid: true,
		},
		{
			Name: "ExpectedValid - Scenarios schema",
			Value: externalschema.LabelDefinitionInput{
				Key:    model.ScenariosKey,
				Schema: fixScenariosSchema(t),
			},
			ExpectedValid: true,
		},
		{
			Name: "Invalid - Invalid schema format",
			Value: externalschema.LabelDefinitionInput{
				Key:    "ok",
				Schema: &invalidSchema,
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - Scenarios schema invalid format",
			Value: externalschema.LabelDefinitionInput{
				Key:    model.ScenariosKey,
				Schema: &invalidSchema,
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - Scenarios schema invalid",
			Value: externalschema.LabelDefinitionInput{
				Key:    model.ScenariosKey,
				Schema: &validSchema,
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - Scenarios schema nil",
			Value: externalschema.LabelDefinitionInput{
				Key:    model.ScenariosKey,
				Schema: nil,
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - Scenarios schema with enum value which does not meet the regex - enum value contains invalid character",
			Value: externalschema.LabelDefinitionInput{
				Key: model.ScenariosKey,
				Schema: jsonSchemaPtr(`{
					"type":        "array",
					"minItems":    1,
					"uniqueItems": true,
					"items": {
						"type": "string",
						"enum": ["DEFAULT", "Abc&Cde"]
					}
				}`),
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - Scenarios schema with enum value which does not meet the regex - enum value too long",
			Value: externalschema.LabelDefinitionInput{
				Key: model.ScenariosKey,
				Schema: jsonSchemaPtr(`{
					"type":        "array",
					"minItems":    1,
					"uniqueItems": true,
					"items": {
						"type": "string",
						"enum": ["DEFAULT", "Abcdefghijklmnopqrstuvwxyz1234567890Abcdefghijklmnopqrstuvwxyz1234567890Abcdefghijklmnopqrstuvwxyz1234567890Abcdefghijklmnopqrstuvwxyz1234567890"]
					}	
				}`),
			},
			ExpectedValid: false,
		},
		{
			Name: "Invalid - scenarios schema without DEFAULT enum value",
			Value: externalschema.LabelDefinitionInput{
				Key: model.ScenariosKey,
				Schema: jsonSchemaPtr(`{
					"type":        "array",
					"minItems":    1,
					"uniqueItems": true,
					"items": {
						"type": "string",
						"enum": ["Abc"]
					}
				}`),
			},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := testCase.Value
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

func TestLabelDefinitionInput_Validate_Key(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         "valid",
			ExpectedValid: true,
		},
		{
			Name:          "Invalid - Empty",
			Value:         inputvalidationtest.EmptyString,
			ExpectedValid: false,
		},
		{
			Name:          "Invalid - Too long",
			Value:         inputvalidationtest.String257Long,
			ExpectedValid: false,
		},
		{
			Name:          "Invalid - Unsupported characters in key",
			Value:         "not/valid",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := fixValidLabelDefinitionInput()
			sut.Key = testCase.Value
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

func TestLabelDefinitionInput_Validate_Schema(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         *externalschema.JSONSchema
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         &validSchema,
			ExpectedValid: true,
		},
		{
			Name:          "ExpectedValid - Nil",
			Value:         (*externalschema.JSONSchema)(nil),
			ExpectedValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			sut := fixValidLabelDefinitionInput()
			sut.Schema = testCase.Value
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

func fixValidLabelDefinitionInput() externalschema.LabelDefinitionInput {
	return externalschema.LabelDefinitionInput{
		Key:    "valid",
		Schema: nil,
	}
}

func jsonSchemaPtr(schema string) *externalschema.JSONSchema {
	s := externalschema.JSONSchema(schema)
	return &s
}

func fixScenariosSchema(t *testing.T) *externalschema.JSONSchema {
	marshalled, err := json.Marshal(model.ScenariosSchema)
	require.NoError(t, err)
	return jsonSchemaPtr(string(marshalled))
}
