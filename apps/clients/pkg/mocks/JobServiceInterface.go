// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	mock "github.com/stretchr/testify/mock"
)

// JobServiceInterface is an autogenerated mock type for the JobServiceInterface type
type JobServiceInterface struct {
	mock.Mock
}

// CreateJob provides a mock function with given fields: ctx, job
func (_m *JobServiceInterface) CreateJob(ctx context.Context, job *model.Job) (string, error) {
	ret := _m.Called(ctx, job)

	if len(ret) == 0 {
		panic("no return value specified for CreateJob")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Job) (string, error)); ok {
		return rf(ctx, job)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Job) string); ok {
		r0 = rf(ctx, job)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Job) error); ok {
		r1 = rf(ctx, job)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllJobs provides a mock function with given fields: ctx, query
func (_m *JobServiceInterface) GetAllJobs(ctx context.Context, query *model.GetJobsQuery) (int, []model.Job, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetAllJobs")
	}

	var r0 int
	var r1 []model.Job
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetJobsQuery) (int, []model.Job, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetJobsQuery) int); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetJobsQuery) []model.Job); ok {
		r1 = rf(ctx, query)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]model.Job)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.GetJobsQuery) error); ok {
		r2 = rf(ctx, query)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetJob provides a mock function with given fields: ctx, jobID
func (_m *JobServiceInterface) GetJob(ctx context.Context, jobID string) (*model.Job, error) {
	ret := _m.Called(ctx, jobID)

	if len(ret) == 0 {
		panic("no return value specified for GetJob")
	}

	var r0 *model.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Job, error)); ok {
		return rf(ctx, jobID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Job); ok {
		r0 = rf(ctx, jobID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, jobID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJobServiceInterface creates a new instance of JobServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJobServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *JobServiceInterface {
	mock := &JobServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
