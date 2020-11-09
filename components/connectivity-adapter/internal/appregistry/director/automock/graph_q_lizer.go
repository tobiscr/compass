// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	externalschema "github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	mock "github.com/stretchr/testify/mock"
)

// GraphQLizer is an autogenerated mock type for the GraphQLizer type
type GraphQLizer struct {
	mock.Mock
}

// APIDefinitionInputToGQL provides a mock function with given fields: in
func (_m *GraphQLizer) APIDefinitionInputToGQL(in externalschema.APIDefinitionInput) (string, error) {
	ret := _m.Called(in)

	var r0 string
	if rf, ok := ret.Get(0).(func(externalschema.APIDefinitionInput) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(externalschema.APIDefinitionInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DocumentInputToGQL provides a mock function with given fields: in
func (_m *GraphQLizer) DocumentInputToGQL(in *externalschema.DocumentInput) (string, error) {
	ret := _m.Called(in)

	var r0 string
	if rf, ok := ret.Get(0).(func(*externalschema.DocumentInput) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*externalschema.DocumentInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventDefinitionInputToGQL provides a mock function with given fields: in
func (_m *GraphQLizer) EventDefinitionInputToGQL(in externalschema.EventDefinitionInput) (string, error) {
	ret := _m.Called(in)

	var r0 string
	if rf, ok := ret.Get(0).(func(externalschema.EventDefinitionInput) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(externalschema.EventDefinitionInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageCreateInputToGQL provides a mock function with given fields: in
func (_m *GraphQLizer) PackageCreateInputToGQL(in externalschema.PackageCreateInput) (string, error) {
	ret := _m.Called(in)

	var r0 string
	if rf, ok := ret.Get(0).(func(externalschema.PackageCreateInput) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(externalschema.PackageCreateInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageUpdateInputToGQL provides a mock function with given fields: in
func (_m *GraphQLizer) PackageUpdateInputToGQL(in externalschema.PackageUpdateInput) (string, error) {
	ret := _m.Called(in)

	var r0 string
	if rf, ok := ret.Get(0).(func(externalschema.PackageUpdateInput) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(externalschema.PackageUpdateInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
