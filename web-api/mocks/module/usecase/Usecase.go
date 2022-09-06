// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	web "go-unit-test/web-api/module/model/web"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Register provides a mock function with given fields: ctx, request
func (_m *Usecase) Register(ctx context.Context, request web.RegisterRequest) error {
	ret := _m.Called(ctx, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, web.RegisterRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyRegister provides a mock function with given fields: ctx, request
func (_m *Usecase) VerifyRegister(ctx context.Context, request web.VerifyRegisterRequest) (*web.VerifyRegisterResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *web.VerifyRegisterResponse
	if rf, ok := ret.Get(0).(func(context.Context, web.VerifyRegisterRequest) *web.VerifyRegisterResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*web.VerifyRegisterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, web.VerifyRegisterRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUsecase creates a new instance of Usecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUsecase(t mockConstructorTestingTNewUsecase) *Usecase {
	mock := &Usecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}