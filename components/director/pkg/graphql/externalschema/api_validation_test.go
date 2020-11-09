package externalschema_test

import (
	"strings"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation/inputvalidationtest"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
	"github.com/stretchr/testify/require"
)

func TestAPIDefinitionInput_Validate_Name(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         "name-123.com",
			ExpectedValid: true,
		},
		{
			Name:          "Valid Printable ASCII",
			Value:         "V1 +=_-)(*&^%$#@!?/>.<,|\\\"':;}{][",
			ExpectedValid: true,
		},
		{
			Name:          "Empty string",
			Value:         inputvalidationtest.EmptyString,
			ExpectedValid: false,
		},
		{
			Name:          "String longer than 100 chars",
			Value:         inputvalidationtest.String129Long,
			ExpectedValid: false,
		},
		{
			Name:          "String contains invalid ASCII",
			Value:         "ąćńłóęǖǘǚǜ",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPIDefinitionInput()
			obj.Name = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPIDefinitionInput_Validate_Description(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         *string
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         str.Ptr("this is a valid description"),
			ExpectedValid: true,
		},
		{
			Name:          "Nil pointer",
			Value:         nil,
			ExpectedValid: true,
		},
		{
			Name:          "Empty string",
			Value:         str.Ptr(inputvalidationtest.EmptyString),
			ExpectedValid: true,
		},
		{
			Name:          "String longer than 2000 chars",
			Value:         str.Ptr(inputvalidationtest.String2001Long),
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPIDefinitionInput()
			obj.Description = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPIDefinitionInput_Validate_TargetURL(t *testing.T) {
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
			Name:          "URL longer than 256",
			Value:         "kyma-project.io/" + strings.Repeat("a", 241),
			ExpectedValid: false,
		},
		{
			Name:          "Invalid, space in URL",
			Value:         "https://kyma test project.io",
			ExpectedValid: false,
		},
		{
			Name:          "Invalid, no protocol",
			Value:         "kyma-project.io",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			app := fixValidAPIDefinitionInput()
			app.TargetURL = testCase.Value
			//WHEN
			err := app.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPIDefinitionInput_Validate_Group(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         *string
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         str.Ptr("this is a valid description"),
			ExpectedValid: true,
		},
		{
			Name:          "Nil pointer",
			Value:         nil,
			ExpectedValid: true,
		},
		{
			Name:          "Empty string",
			Value:         str.Ptr(inputvalidationtest.EmptyString),
			ExpectedValid: true,
		},
		{
			Name:          "String longer than 36 chars",
			Value:         str.Ptr(inputvalidationtest.String37Long),
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPIDefinitionInput()
			obj.Group = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPIDefinitionInput_Validate_APISpecInput(t *testing.T) {
	validObj := fixValidAPISpecInput()
	emptyObj := externalschema.APISpecInput{}

	testCases := []struct {
		Name          string
		Value         *externalschema.APISpecInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid obj",
			Value:         &validObj,
			ExpectedValid: true,
		},
		{
			Name:          "Nil object",
			Value:         nil,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid object",
			Value:         &emptyObj,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPIDefinitionInput()
			obj.Spec = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPIDefinitionInput_Validate_Version(t *testing.T) {
	validObj := fixValidVersionInput()
	emptyObj := externalschema.VersionInput{}

	testCases := []struct {
		Name          string
		Value         *externalschema.VersionInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid obj",
			Value:         &validObj,
			ExpectedValid: true,
		},
		{
			Name:          "Nil object",
			Value:         nil,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid object",
			Value:         &emptyObj,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPIDefinitionInput()
			obj.Version = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPISpecInput_Validate_Type(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         externalschema.APISpecType
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         externalschema.APISpecTypeOpenAPI,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid object",
			Value:         externalschema.APISpecType("INVALID"),
			ExpectedValid: false,
		},
		{
			Name:          "Invalid default value",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPISpecInput()
			obj.Type = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPISpecInput_Validate_Format(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         externalschema.SpecFormat
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid JSON",
			Value:         externalschema.SpecFormatJSON,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid object",
			Value:         externalschema.SpecFormat("INVALID"),
			ExpectedValid: false,
		},
		{
			Name:          "Invalid default value",
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPISpecInput()
			obj.Format = testCase.Value
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPISpecInput_Validate_TypeODataWithFormat(t *testing.T) {
	testCases := []struct {
		Name          string
		InputType     externalschema.APISpecType
		InputFormat   externalschema.SpecFormat
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid ODATA with XML",
			InputType:     externalschema.APISpecTypeOdata,
			InputFormat:   externalschema.SpecFormatXML,
			ExpectedValid: true,
		},
		{
			Name:          "ExpectedValid ODATA with JSON",
			InputType:     externalschema.APISpecTypeOdata,
			InputFormat:   externalschema.SpecFormatJSON,
			ExpectedValid: true,
		},
		{
			Name:          "Invalid ODATA with YAML",
			InputType:     externalschema.APISpecTypeOdata,
			InputFormat:   externalschema.SpecFormatYaml,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPISpecInput()
			obj.Type = testCase.InputType
			obj.Format = testCase.InputFormat
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPISpecInput_Validate_TypeOpenAPIWithFormat(t *testing.T) {
	testCases := []struct {
		Name          string
		InputType     externalschema.APISpecType
		InputFormat   externalschema.SpecFormat
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid OpenAPI with JSON",
			InputType:     externalschema.APISpecTypeOpenAPI,
			InputFormat:   externalschema.SpecFormatJSON,
			ExpectedValid: true,
		},
		{
			Name:          "ExpectedValid OpenAPI with YAML",
			InputType:     externalschema.APISpecTypeOpenAPI,
			InputFormat:   externalschema.SpecFormatYaml,
			ExpectedValid: true,
		},
		{
			Name:          "invalid OpenAPI with XML",
			InputType:     externalschema.APISpecTypeOpenAPI,
			InputFormat:   externalschema.SpecFormatXML,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPISpecInput()
			obj.Type = testCase.InputType
			obj.Format = testCase.InputFormat
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAPISpecInput_Validate_FetchRequest(t *testing.T) {
	validObj := fixValidFetchRequestInput()
	emptyObj := externalschema.FetchRequestInput{}

	testCases := []struct {
		Name          string
		Value         *externalschema.FetchRequestInput
		DataClob      *externalschema.CLOB
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid obj",
			Value:         &validObj,
			ExpectedValid: true,
		},
		{
			Name:          "Nil object",
			Value:         nil,
			DataClob:      fixCLOB("data"),
			ExpectedValid: true,
		},
		{
			Name:          "Invalid object",
			Value:         &emptyObj,
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidAPISpecInput()
			obj.FetchRequest = testCase.Value
			obj.Data = testCase.DataClob
			//WHEN
			err := obj.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func fixValidAPISpecInput() externalschema.APISpecInput {
	return externalschema.APISpecInput{
		Type:   externalschema.APISpecTypeOpenAPI,
		Format: externalschema.SpecFormatJSON,
		Data:   fixCLOB("data"),
	}
}

func fixValidAPIDefinitionInput() externalschema.APIDefinitionInput {
	return externalschema.APIDefinitionInput{
		Name:      inputvalidationtest.ValidName,
		TargetURL: inputvalidationtest.ValidURL,
	}
}
