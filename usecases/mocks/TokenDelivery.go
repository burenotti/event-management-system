// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/burenotti/rtu-it-lab-recruit/model"
	mock "github.com/stretchr/testify/mock"
)

// TokenDelivery is an autogenerated mock type for the TokenDelivery type
type TokenDelivery struct {
	mock.Mock
}

// SendActivationToken provides a mock function with given fields: ctx, user, token
func (_m *TokenDelivery) SendActivationToken(ctx context.Context, user *model.User, token string) error {
	ret := _m.Called(ctx, user, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, string) error); ok {
		r0 = rf(ctx, user, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTokenDelivery interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenDelivery creates a new instance of TokenDelivery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenDelivery(t mockConstructorTestingTNewTokenDelivery) *TokenDelivery {
	mock := &TokenDelivery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}