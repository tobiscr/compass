// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
	mock "github.com/stretchr/testify/mock"
)

// APIRepository is an autogenerated mock type for the APIRepository type
type APIRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *APIRepository) Create(ctx context.Context, item *model.APIDefinition) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.APIDefinition) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllByApplicationID provides a mock function with given fields: ctx, tenant, id
func (_m *APIRepository) DeleteAllByApplicationID(ctx context.Context, tenant string, id string) error {
	ret := _m.Called(ctx, tenant, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListForApplication provides a mock function with given fields: ctx, tenant, applicationID, pageSize, cursor
func (_m *APIRepository) ListForApplication(ctx context.Context, tenant string, applicationID string, pageSize int, cursor string) (*model.APIDefinitionPage, error) {
	ret := _m.Called(ctx, tenant, applicationID, pageSize, cursor)

	var r0 *model.APIDefinitionPage
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) *model.APIDefinitionPage); ok {
		r0 = rf(ctx, tenant, applicationID, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinitionPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, tenant, applicationID, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}