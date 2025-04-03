// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	bson "go.mongodb.org/mongo-driver/v2/bson"

	mock "github.com/stretchr/testify/mock"

	model "github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
)

// ClientRepository is an autogenerated mock type for the ClientRepository type
type ClientRepository struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, query
func (_m *ClientRepository) Count(ctx context.Context, query *model.GetClientsQuery) (int, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Count")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetClientsQuery) (int, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetClientsQuery) int); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetClientsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, query
func (_m *ClientRepository) GetAll(ctx context.Context, query *model.GetClientsQuery) ([]model.Client, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []model.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetClientsQuery) ([]model.Client, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetClientsQuery) []model.Client); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Client)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetClientsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: ctx, clientID
func (_m *ClientRepository) GetOne(ctx context.Context, clientID string) (*model.Client, error) {
	ret := _m.Called(ctx, clientID)

	if len(ret) == 0 {
		panic("no return value specified for GetOne")
	}

	var r0 *model.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Client, error)); ok {
		return rf(ctx, clientID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Client); ok {
		r0 = rf(ctx, clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Client)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, clientID, data
func (_m *ClientRepository) Update(ctx context.Context, clientID string, data bson.D) error {
	ret := _m.Called(ctx, clientID, data)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bson.D) error); ok {
		r0 = rf(ctx, clientID, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewClientRepository creates a new instance of ClientRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientRepository {
	mock := &ClientRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
