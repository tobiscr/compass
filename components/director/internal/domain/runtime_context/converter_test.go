package runtime_context_test

import (
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime_context"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestConverter_ToGraphQL(t *testing.T) {
	id := "test_id"
	runtimeID := "test_runtime_id"
	tenant := "test_tenant"
	key := "key"
	val := "val"
	// given
	testCases := []struct {
		Name     string
		Input    *model.RuntimeContext
		Expected *externalschema.RuntimeContext
	}{
		{
			Name: "All properties given",
			Input: &model.RuntimeContext{
				ID:        id,
				RuntimeID: runtimeID,
				Tenant:    tenant,
				Key:       key,
				Value:     val,
			},
			Expected: &externalschema.RuntimeContext{
				ID:    id,
				Key:   key,
				Value: val,
			},
		},
		{
			Name:     "Empty",
			Input:    &model.RuntimeContext{},
			Expected: &externalschema.RuntimeContext{},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// when
			converter := runtime_context.NewConverter()
			res := converter.ToGraphQL(testCase.Input)

			// then
			assert.Equal(t, testCase.Expected, res)
		})
	}
}

func TestConverter_MultipleToGraphQL(t *testing.T) {
	id := "test_id"
	runtimeID := "test_runtime_id"
	tenant := "test_tenant"
	key := "key"
	val := "val"

	// given
	input := []*model.RuntimeContext{
		{
			ID:        id,
			RuntimeID: runtimeID,
			Tenant:    tenant,
			Key:       key,
			Value:     val,
		},
		{
			ID:        id + "2",
			RuntimeID: runtimeID + "2",
			Tenant:    tenant + "2",
			Key:       key + "2",
			Value:     val + "2",
		},
		nil,
	}
	expected := []*externalschema.RuntimeContext{
		{
			ID:    id,
			Key:   key,
			Value: val,
		},
		{
			ID:    id + "2",
			Key:   key + "2",
			Value: val + "2",
		},
	}

	// when
	converter := runtime_context.NewConverter()
	res := converter.MultipleToGraphQL(input)

	// then
	assert.Equal(t, expected, res)
}

func TestConverter_InputFromGraphQL(t *testing.T) {
	key := "key"
	val := "val"
	labels := externalschema.Labels(map[string]interface{}{
		"test": "test",
	})
	runtimeID := "runtime_id"

	// given
	testCases := []struct {
		Name     string
		Input    externalschema.RuntimeContextInput
		Expected model.RuntimeContextInput
	}{
		{
			Name: "All properties given",
			Input: externalschema.RuntimeContextInput{
				Key:    key,
				Value:  val,
				Labels: &labels,
			},
			Expected: model.RuntimeContextInput{
				Key:       key,
				Value:     val,
				RuntimeID: runtimeID,
				Labels:    labels,
			},
		},
		{
			Name:  "Empty",
			Input: externalschema.RuntimeContextInput{},
			Expected: model.RuntimeContextInput{
				RuntimeID: runtimeID,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// when
			converter := runtime_context.NewConverter()
			res := converter.InputFromGraphQL(testCase.Input, runtimeID)

			// then
			assert.Equal(t, testCase.Expected, res)
		})
	}
}
