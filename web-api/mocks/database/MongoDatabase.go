// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	database "go-unit-test/web-api/database"

	mock "github.com/stretchr/testify/mock"
)

// MongoDatabase is an autogenerated mock type for the MongoDatabase type
type MongoDatabase struct {
	mock.Mock
}

// FindOne provides a mock function with given fields: ctx, payload
func (_m *MongoDatabase) FindOne(ctx context.Context, payload database.FindOne) error {
	ret := _m.Called(ctx, payload)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, database.FindOne) error); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertOne provides a mock function with given fields: ctx, payload
func (_m *MongoDatabase) InsertOne(ctx context.Context, payload database.InsertOne) (*string, error) {
	ret := _m.Called(ctx, payload)

	var r0 *string
	if rf, ok := ret.Get(0).(func(context.Context, database.InsertOne) *string); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, database.InsertOne) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMongoDatabase interface {
	mock.TestingT
	Cleanup(func())
}

// NewMongoDatabase creates a new instance of MongoDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMongoDatabase(t mockConstructorTestingTNewMongoDatabase) *MongoDatabase {
	mock := &MongoDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
