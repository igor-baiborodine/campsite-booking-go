// Code generated by mockery v2.43.2. DO NOT EDIT.

package waiter

import mock "github.com/stretchr/testify/mock"

// MockOption is an autogenerated mock type for the Option type
type MockOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: c
func (_m *MockOption) Execute(c *waiterCfg) {
	_m.Called(c)
}

// NewMockOption creates a new instance of MockOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOption {
	mock := &MockOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
