// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
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

// CreateMany provides a mock function with given fields: ctx, item
func (_m *APIRepository) CreateMany(ctx context.Context, item []*model.APIDefinition) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.APIDefinition) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenantID, id
func (_m *APIRepository) Delete(ctx context.Context, tenantID string, id string) error {
	ret := _m.Called(ctx, tenantID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, tenant, id
func (_m *APIRepository) Exists(ctx context.Context, tenant string, id string) (bool, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, tenantID, id
func (_m *APIRepository) GetByID(ctx context.Context, tenantID string, id string) (*model.APIDefinition, error) {
	ret := _m.Called(ctx, tenantID, id)

	var r0 *model.APIDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.APIDefinition); ok {
		r0 = rf(ctx, tenantID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForPackage provides a mock function with given fields: ctx, tenant, id, packageID
func (_m *APIRepository) GetForPackage(ctx context.Context, tenant string, id string, packageID string) (*model.APIDefinition, error) {
	ret := _m.Called(ctx, tenant, id, packageID)

	var r0 *model.APIDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *model.APIDefinition); ok {
		r0 = rf(ctx, tenant, id, packageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, tenant, id, packageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForPackage provides a mock function with given fields: ctx, tenantID, packageID, pageSize, cursor
func (_m *APIRepository) ListForPackage(ctx context.Context, tenantID string, packageID string, pageSize int, cursor string) (*model.APIDefinitionPage, error) {
	ret := _m.Called(ctx, tenantID, packageID, pageSize, cursor)

	var r0 *model.APIDefinitionPage
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) *model.APIDefinitionPage); ok {
		r0 = rf(ctx, tenantID, packageID, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinitionPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, tenantID, packageID, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, item
func (_m *APIRepository) Update(ctx context.Context, item *model.APIDefinition) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.APIDefinition) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}