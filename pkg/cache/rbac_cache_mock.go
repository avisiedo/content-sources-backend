// Code generated by mockery v2.20.0. DO NOT EDIT.

package cache

import (
	context "context"

	rbac "github.com/RedHatInsights/rbac-client-go"
	mock "github.com/stretchr/testify/mock"
)

// MockRbacCache is an autogenerated mock type for the RbacCache type
type MockRbacCache struct {
	mock.Mock
}

// GetAccessList provides a mock function with given fields: ctx
func (_m *MockRbacCache) GetAccessList(ctx context.Context) (rbac.AccessList, error) {
	ret := _m.Called(ctx)

	var r0 rbac.AccessList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (rbac.AccessList, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) rbac.AccessList); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rbac.AccessList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAccessList provides a mock function with given fields: ctx, accessList
func (_m *MockRbacCache) SetAccessList(ctx context.Context, accessList rbac.AccessList) error {
	ret := _m.Called(ctx, accessList)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, rbac.AccessList) error); ok {
		r0 = rf(ctx, accessList)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockRbacCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRbacCache creates a new instance of MockRbacCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRbacCache(t mockConstructorTestingTNewMockRbacCache) *MockRbacCache {
	mock := &MockRbacCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
