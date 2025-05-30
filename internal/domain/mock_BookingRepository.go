// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package domain

import (
	"context"
	"time"

	mock "github.com/stretchr/testify/mock"
)

// NewMockBookingRepository creates a new instance of MockBookingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBookingRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBookingRepository {
	mock := &MockBookingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockBookingRepository is an autogenerated mock type for the BookingRepository type
type MockBookingRepository struct {
	mock.Mock
}

type MockBookingRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBookingRepository) EXPECT() *MockBookingRepository_Expecter {
	return &MockBookingRepository_Expecter{mock: &_m.Mock}
}

// Find provides a mock function for the type MockBookingRepository
func (_mock *MockBookingRepository) Find(ctx context.Context, bookingID string) (*Booking, error) {
	ret := _mock.Called(ctx, bookingID)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *Booking
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (*Booking, error)); ok {
		return returnFunc(ctx, bookingID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) *Booking); ok {
		r0 = returnFunc(ctx, bookingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Booking)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, bookingID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockBookingRepository_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockBookingRepository_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - bookingID string
func (_e *MockBookingRepository_Expecter) Find(ctx interface{}, bookingID interface{}) *MockBookingRepository_Find_Call {
	return &MockBookingRepository_Find_Call{Call: _e.mock.On("Find", ctx, bookingID)}
}

func (_c *MockBookingRepository_Find_Call) Run(run func(ctx context.Context, bookingID string)) *MockBookingRepository_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockBookingRepository_Find_Call) Return(booking *Booking, err error) *MockBookingRepository_Find_Call {
	_c.Call.Return(booking, err)
	return _c
}

func (_c *MockBookingRepository_Find_Call) RunAndReturn(run func(ctx context.Context, bookingID string) (*Booking, error)) *MockBookingRepository_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindForDateRange provides a mock function for the type MockBookingRepository
func (_mock *MockBookingRepository) FindForDateRange(ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time) ([]*Booking, error) {
	ret := _mock.Called(ctx, campsiteID, startDate, endDate)

	if len(ret) == 0 {
		panic("no return value specified for FindForDateRange")
	}

	var r0 []*Booking
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) ([]*Booking, error)); ok {
		return returnFunc(ctx, campsiteID, startDate, endDate)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) []*Booking); ok {
		r0 = returnFunc(ctx, campsiteID, startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Booking)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string, time.Time, time.Time) error); ok {
		r1 = returnFunc(ctx, campsiteID, startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockBookingRepository_FindForDateRange_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindForDateRange'
type MockBookingRepository_FindForDateRange_Call struct {
	*mock.Call
}

// FindForDateRange is a helper method to define mock.On call
//   - ctx context.Context
//   - campsiteID string
//   - startDate time.Time
//   - endDate time.Time
func (_e *MockBookingRepository_Expecter) FindForDateRange(ctx interface{}, campsiteID interface{}, startDate interface{}, endDate interface{}) *MockBookingRepository_FindForDateRange_Call {
	return &MockBookingRepository_FindForDateRange_Call{Call: _e.mock.On("FindForDateRange", ctx, campsiteID, startDate, endDate)}
}

func (_c *MockBookingRepository_FindForDateRange_Call) Run(run func(ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time)) *MockBookingRepository_FindForDateRange_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		var arg2 time.Time
		if args[2] != nil {
			arg2 = args[2].(time.Time)
		}
		var arg3 time.Time
		if args[3] != nil {
			arg3 = args[3].(time.Time)
		}
		run(
			arg0,
			arg1,
			arg2,
			arg3,
		)
	})
	return _c
}

func (_c *MockBookingRepository_FindForDateRange_Call) Return(bookings []*Booking, err error) *MockBookingRepository_FindForDateRange_Call {
	_c.Call.Return(bookings, err)
	return _c
}

func (_c *MockBookingRepository_FindForDateRange_Call) RunAndReturn(run func(ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time) ([]*Booking, error)) *MockBookingRepository_FindForDateRange_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function for the type MockBookingRepository
func (_mock *MockBookingRepository) Insert(ctx context.Context, booking *Booking) error {
	ret := _mock.Called(ctx, booking)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *Booking) error); ok {
		r0 = returnFunc(ctx, booking)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockBookingRepository_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockBookingRepository_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - ctx context.Context
//   - booking *Booking
func (_e *MockBookingRepository_Expecter) Insert(ctx interface{}, booking interface{}) *MockBookingRepository_Insert_Call {
	return &MockBookingRepository_Insert_Call{Call: _e.mock.On("Insert", ctx, booking)}
}

func (_c *MockBookingRepository_Insert_Call) Run(run func(ctx context.Context, booking *Booking)) *MockBookingRepository_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *Booking
		if args[1] != nil {
			arg1 = args[1].(*Booking)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockBookingRepository_Insert_Call) Return(err error) *MockBookingRepository_Insert_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockBookingRepository_Insert_Call) RunAndReturn(run func(ctx context.Context, booking *Booking) error) *MockBookingRepository_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function for the type MockBookingRepository
func (_mock *MockBookingRepository) Update(ctx context.Context, booking *Booking) error {
	ret := _mock.Called(ctx, booking)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *Booking) error); ok {
		r0 = returnFunc(ctx, booking)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockBookingRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockBookingRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - booking *Booking
func (_e *MockBookingRepository_Expecter) Update(ctx interface{}, booking interface{}) *MockBookingRepository_Update_Call {
	return &MockBookingRepository_Update_Call{Call: _e.mock.On("Update", ctx, booking)}
}

func (_c *MockBookingRepository_Update_Call) Run(run func(ctx context.Context, booking *Booking)) *MockBookingRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *Booking
		if args[1] != nil {
			arg1 = args[1].(*Booking)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockBookingRepository_Update_Call) Return(err error) *MockBookingRepository_Update_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockBookingRepository_Update_Call) RunAndReturn(run func(ctx context.Context, booking *Booking) error) *MockBookingRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}
