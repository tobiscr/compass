package externalschema_test

import (
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation/inputvalidationtest"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
	"github.com/stretchr/testify/require"
)

func TestPackageCreateInput_Validate_Name(t *testing.T) {
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
			obj := fixValidPackageCreateInput()
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

func TestPackageCreateInput_Validate_Description(t *testing.T) {
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
			obj := fixValidPackageCreateInput()
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

func TestPackageCreateInput_Validate_DefaultInstanceAuth(t *testing.T) {
	validObj := fixValidAuthInput()

	testCases := []struct {
		Name          string
		Value         *externalschema.AuthInput
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
			Name:          "Invalid - Nested validation error",
			Value:         &externalschema.AuthInput{Credential: &externalschema.CredentialDataInput{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidPackageCreateInput()
			obj.DefaultInstanceAuth = testCase.Value
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

func TestPackageCreateInput_Validate_InstanceAuthRequestInputSchema(t *testing.T) {
	schema := externalschema.JSONSchema("Test")
	emptySchema := externalschema.JSONSchema("")
	testCases := []struct {
		Name          string
		Value         *externalschema.JSONSchema
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         &schema,
			ExpectedValid: true,
		},
		{
			Name:          "Empty schema",
			Value:         &emptySchema,
			ExpectedValid: false,
		},
		{
			Name:          "Nil pointer",
			Value:         nil,
			ExpectedValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidPackageCreateInput()
			obj.InstanceAuthRequestInputSchema = testCase.Value
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

func TestPackageCreateInput_Validate_APIs(t *testing.T) {
	validObj := fixValidAPIDefinitionInput()

	testCases := []struct {
		Name          string
		Value         []*externalschema.APIDefinitionInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid array",
			Value:         []*externalschema.APIDefinitionInput{&validObj},
			ExpectedValid: true,
		},
		{
			Name:          "Empty array",
			Value:         []*externalschema.APIDefinitionInput{},
			ExpectedValid: true,
		},
		{
			Name:          "Array with invalid object",
			Value:         []*externalschema.APIDefinitionInput{{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			app := fixValidPackageCreateInput()
			app.APIDefinitions = testCase.Value
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

func TestPackageCreateInput_Validate_EventAPIs(t *testing.T) {
	validObj := fixValidEventAPIDefinitionInput()

	testCases := []struct {
		Name          string
		Value         []*externalschema.EventDefinitionInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid array",
			Value:         []*externalschema.EventDefinitionInput{&validObj},
			ExpectedValid: true,
		},
		{
			Name:          "Empty array",
			Value:         []*externalschema.EventDefinitionInput{},
			ExpectedValid: true,
		},
		{
			Name:          "Array with invalid object",
			Value:         []*externalschema.EventDefinitionInput{{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			app := fixValidPackageCreateInput()
			app.EventDefinitions = testCase.Value
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

func TestPackageCreateInput_Validate_Documents(t *testing.T) {
	validDoc := fixValidDocument()

	testCases := []struct {
		Name          string
		Value         []*externalschema.DocumentInput
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid array",
			Value:         []*externalschema.DocumentInput{&validDoc},
			ExpectedValid: true,
		},
		{
			Name:          "Empty array",
			Value:         []*externalschema.DocumentInput{},
			ExpectedValid: true,
		},
		{
			Name:          "Array with invalid object",
			Value:         []*externalschema.DocumentInput{{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			app := fixValidPackageCreateInput()
			app.Documents = testCase.Value
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

func TestPackageUpdateInput_Validate_Name(t *testing.T) {
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
			obj := fixValidPackageUpdateInput()
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

func TestPackageUpdateInput_Validate_Description(t *testing.T) {
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
			obj := fixValidPackageUpdateInput()
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

func TestPackageUpdateInput_Validate_DefaultInstanceAuth(t *testing.T) {
	validObj := fixValidAuthInput()

	testCases := []struct {
		Name          string
		Value         *externalschema.AuthInput
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
			Name:          "Invalid - Nested validation error",
			Value:         &externalschema.AuthInput{Credential: &externalschema.CredentialDataInput{}},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidPackageUpdateInput()
			obj.DefaultInstanceAuth = testCase.Value
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

func TestPackageUpdateInput_Validate_InstanceAuthRequestInputSchema(t *testing.T) {
	schema := externalschema.JSONSchema("Test")
	emptySchema := externalschema.JSONSchema("")
	testCases := []struct {
		Name          string
		Value         *externalschema.JSONSchema
		ExpectedValid bool
	}{
		{
			Name:          "ExpectedValid",
			Value:         &schema,
			ExpectedValid: true,
		},
		{
			Name:          "Empty schema",
			Value:         &emptySchema,
			ExpectedValid: false,
		},
		{
			Name:          "Nil pointer",
			Value:         nil,
			ExpectedValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			obj := fixValidPackageUpdateInput()
			obj.InstanceAuthRequestInputSchema = testCase.Value
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

func TestPackageInstanceAuthRequestInput_Validate(t *testing.T) {
	//GIVEN
	val := externalschema.JSON("{\"foo\": \"bar\"}")
	testCases := []struct {
		Name          string
		Value         externalschema.PackageInstanceAuthRequestInput
		ExpectedValid bool
	}{
		{
			Name:          "Empty",
			Value:         externalschema.PackageInstanceAuthRequestInput{},
			ExpectedValid: true,
		},
		{
			Name: "InputParams and Context set",
			Value: externalschema.PackageInstanceAuthRequestInput{
				Context:     &val,
				InputParams: &val,
			},
			ExpectedValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//WHEN
			err := testCase.Value.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestPackageInstanceAuthSetInput_Validate(t *testing.T) {
	//GIVEN
	authInput := fixValidAuthInput()
	str := "foo"
	testCases := []struct {
		Name          string
		Value         externalschema.PackageInstanceAuthSetInput
		ExpectedValid bool
	}{
		{
			Name: "Auth",
			Value: externalschema.PackageInstanceAuthSetInput{
				Auth: &authInput,
			},
			ExpectedValid: true,
		},
		{
			Name: "Failed Status",
			Value: externalschema.PackageInstanceAuthSetInput{
				Status: &externalschema.PackageInstanceAuthStatusInput{
					Condition: externalschema.PackageInstanceAuthSetStatusConditionInputFailed,
					Reason:    str,
					Message:   str,
				},
			},
			ExpectedValid: true,
		},
		{
			Name: "Success Status",
			Value: externalschema.PackageInstanceAuthSetInput{
				Status: &externalschema.PackageInstanceAuthStatusInput{
					Condition: externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded,
				},
			},
			ExpectedValid: false,
		},
		{
			Name: "Auth and Success Status",
			Value: externalschema.PackageInstanceAuthSetInput{
				Auth: &authInput,
				Status: &externalschema.PackageInstanceAuthStatusInput{
					Condition: externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded,
					Message:   str,
					Reason:    str,
				},
			},
			ExpectedValid: true,
		},
		{
			Name: "Auth and Failure Status",
			Value: externalschema.PackageInstanceAuthSetInput{
				Auth: &authInput,
				Status: &externalschema.PackageInstanceAuthStatusInput{
					Condition: externalschema.PackageInstanceAuthSetStatusConditionInputFailed,
				},
			},
			ExpectedValid: false,
		},
		{
			Name: "Empty objects",
			Value: externalschema.PackageInstanceAuthSetInput{
				Auth:   &externalschema.AuthInput{},
				Status: &externalschema.PackageInstanceAuthStatusInput{},
			},
			ExpectedValid: false,
		},
		{
			Name:          "Empty",
			Value:         externalschema.PackageInstanceAuthSetInput{},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//WHEN
			err := testCase.Value.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestPackageInstanceAuthStatusInput_Validate(t *testing.T) {
	//GIVEN
	str := "foo"
	testCases := []struct {
		Name          string
		Value         externalschema.PackageInstanceAuthStatusInput
		ExpectedValid bool
	}{
		{
			Name: "Success",
			Value: externalschema.PackageInstanceAuthStatusInput{
				Condition: externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded,
				Message:   str,
				Reason:    str,
			},
			ExpectedValid: true,
		},
		{
			Name: "No reason provided",
			Value: externalschema.PackageInstanceAuthStatusInput{
				Condition: externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded,
				Message:   str,
			},
			ExpectedValid: false,
		},
		{
			Name: "No message provided",
			Value: externalschema.PackageInstanceAuthStatusInput{
				Condition: externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded,
				Reason:    str,
			},
			ExpectedValid: false,
		},
		{
			Name: "No condition provided",
			Value: externalschema.PackageInstanceAuthStatusInput{
				Message: str,
				Reason:  str,
			},
			ExpectedValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//WHEN
			err := testCase.Value.Validate()
			//THEN
			if testCase.ExpectedValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func fixValidPackageCreateInput() externalschema.PackageCreateInput {
	return externalschema.PackageCreateInput{
		Name: inputvalidationtest.ValidName,
	}
}

func fixValidPackageUpdateInput() externalschema.PackageUpdateInput {
	return externalschema.PackageUpdateInput{
		Name: inputvalidationtest.ValidName,
	}
}
