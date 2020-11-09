// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	http "net/http"

	director "github.com/kyma-incubator/compass/components/connectivity-adapter/internal/connectorservice/director"

	mock "github.com/stretchr/testify/mock"
)

// ClientProvider is an autogenerated mock type for the ClientProvider type
type ClientProvider struct {
	mock.Mock
}

// Client provides a mock function with given fields: r
func (_m *ClientProvider) Client(r *http.Request) director.Client {
	ret := _m.Called(r)

	var r0 director.Client
	if rf, ok := ret.Get(0).(func(*http.Request) director.Client); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(director.Client)
		}
	}

	return r0
}
