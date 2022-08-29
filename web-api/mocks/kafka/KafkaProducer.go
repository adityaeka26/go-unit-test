// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// KafkaProducer is an autogenerated mock type for the KafkaProducer type
type KafkaProducer struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: topic, msg
func (_m *KafkaProducer) SendMessage(topic string, msg string) error {
	ret := _m.Called(topic, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(topic, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewKafkaProducer interface {
	mock.TestingT
	Cleanup(func())
}

// NewKafkaProducer creates a new instance of KafkaProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKafkaProducer(t mockConstructorTestingTNewKafkaProducer) *KafkaProducer {
	mock := &KafkaProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}