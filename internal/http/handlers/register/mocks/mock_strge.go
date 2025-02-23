// Code generated by mockery 2.52.3. DO NOT EDIT.

package register

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Mockstrge is an autogenerated mock type for the strge type
type Mockstrge struct {
	mock.Mock
}

type Mockstrge_Expecter struct {
	mock *mock.Mock
}

func (_m *Mockstrge) EXPECT() *Mockstrge_Expecter {
	return &Mockstrge_Expecter{mock: &_m.Mock}
}

// AddUser provides a mock function with given fields: ctx, email, username, password
func (_m *Mockstrge) AddUser(ctx context.Context, email string, username string, password string) error {
	ret := _m.Called(ctx, email, username, password)

	if len(ret) == 0 {
		panic("no return value specified for AddUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, email, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Mockstrge_AddUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUser'
type Mockstrge_AddUser_Call struct {
	*mock.Call
}

// AddUser is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - username string
//   - password string
func (_e *Mockstrge_Expecter) AddUser(ctx interface{}, email interface{}, username interface{}, password interface{}) *Mockstrge_AddUser_Call {
	return &Mockstrge_AddUser_Call{Call: _e.mock.On("AddUser", ctx, email, username, password)}
}

func (_c *Mockstrge_AddUser_Call) Run(run func(ctx context.Context, email string, username string, password string)) *Mockstrge_AddUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *Mockstrge_AddUser_Call) Return(_a0 error) *Mockstrge_AddUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Mockstrge_AddUser_Call) RunAndReturn(run func(context.Context, string, string, string) error) *Mockstrge_AddUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockstrge creates a new instance of Mockstrge. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockstrge(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mockstrge {
	mock := &Mockstrge{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
