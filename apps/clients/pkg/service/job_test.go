package service_test

import (
	"context"
	"testing"
	"time"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobServiceTestSuite struct {
	suite.Suite
	mockRepo   *mocks.JobRepository
	jobService *service.JobService
}

func (suite *JobServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.JobRepository)
	suite.jobService = service.NewJobService(suite.mockRepo)
}

func (suite *JobServiceTestSuite) TestCreateJob_Success() {
	job := &model.Job{
		PrefectFlowID: "prefect-flow-id-123",
		Type:          model.Scrape,
		Status:        model.JobStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Input:         bson.M{"url": "https://example.com"},
		ScrapeResult:  bson.NewObjectID(),
		MatchResults:  []model.MatchResult{},
		Logs:          []model.JobLog{},
	}

	suite.mockRepo.On("Create", mock.Anything, job).Return("test-job-id", nil)

	id, err := suite.jobService.CreateJob(context.Background(), job)

	suite.NoError(err)
	suite.Equal("test-job-id", id)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestCreateJob_RepoReturnsDependencyFailed() {
	job := &model.Job{}
	suite.mockRepo.On("Create", mock.Anything, job).Return("", errorx.ErrDependencyFailed)

	id, err := suite.jobService.CreateJob(context.Background(), job)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestCreateJob_RepoReturnsError() {
	job := &model.Job{}
	suite.mockRepo.On("Create", mock.Anything, job).Return("", assert.AnError)

	id, err := suite.jobService.CreateJob(context.Background(), job)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetJob_Success() {
	expectedJob := &model.Job{
		PrefectFlowID: "prefect-flow-id-123",
		Type:          model.Scrape,
		Status:        model.JobStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Input:         bson.M{"url": "https://example.com"},
		ScrapeResult:  bson.NewObjectID(),
		MatchResults:  []model.MatchResult{},
		Logs:          []model.JobLog{},
	}
	suite.mockRepo.On("GetOne", mock.Anything, "job-id").Return(expectedJob, nil)

	job, err := suite.jobService.GetJob(context.Background(), "job-id")

	suite.NoError(err)
	suite.Equal(expectedJob, job)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetJob_RepoReturnsNotFound() {
	suite.mockRepo.On("GetOne", mock.Anything, "bad-id").Return(nil, errorx.ErrNotFound)

	job, err := suite.jobService.GetJob(context.Background(), "bad-id")

	suite.Error(err)
	suite.Nil(job)
	suite.ErrorIs(err, errorx.ErrNotFound)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetJob_RepoReturnsInvalidInput() {
	suite.mockRepo.On("GetOne", mock.Anything, "bad-id").Return(nil, errorx.ErrInvalidInput)

	job, err := suite.jobService.GetJob(context.Background(), "bad-id")

	suite.Error(err)
	suite.Nil(job)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetJob_RepoReturnsDependencyFailed() {
	suite.mockRepo.On("GetOne", mock.Anything, "bad-id").Return(nil, errorx.ErrDependencyFailed)

	job, err := suite.jobService.GetJob(context.Background(), "bad-id")

	suite.Error(err)
	suite.Nil(job)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetJob_RepoReturnsError() {
	suite.mockRepo.On("GetOne", mock.Anything, "bad-id").Return(nil, assert.AnError)

	job, err := suite.jobService.GetJob(context.Background(), "bad-id")

	suite.Error(err)
	suite.Nil(job)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetAllJobs_Success() {
	query := &model.GetJobsQuery{}
	mockJobs := []model.Job{
		{
			PrefectFlowID: "prefect-flow-id-123",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
		{
			PrefectFlowID: "prefect-flow-id-456",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
	}

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(mockJobs, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(2, nil)

	total, jobs, err := suite.jobService.GetAllJobs(context.Background(), query)

	suite.NoError(err)
	suite.Equal(2, total)
	suite.Equal(mockJobs, jobs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetAllJobs_RepoReturnsDependencyFailed() {
	query := &model.GetJobsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, errorx.ErrDependencyFailed)

	total, jobs, err := suite.jobService.GetAllJobs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.Zero(total)
	suite.Nil(jobs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetAllJobs_RepoReturnsError() {
	query := &model.GetJobsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, assert.AnError)

	total, jobs, err := suite.jobService.GetAllJobs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.Zero(total)
	suite.Nil(jobs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetAllJobs_RepoCountReturnsDependencyFailed() {
	mockJobs := []model.Job{
		{
			PrefectFlowID: "prefect-flow-id-123",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
		{
			PrefectFlowID: "prefect-flow-id-456",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
	}

	query := &model.GetJobsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(mockJobs, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(0, errorx.ErrDependencyFailed)

	total, jobs, err := suite.jobService.GetAllJobs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.Zero(total)
	suite.Nil(jobs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *JobServiceTestSuite) TestGetAllJobs_RepoCountReturnsError() {
	mockJobs := []model.Job{
		{
			PrefectFlowID: "prefect-flow-id-123",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
		{
			PrefectFlowID: "prefect-flow-id-456",
			Type:          model.Scrape,
			Status:        model.JobStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Input:         bson.M{"url": "https://example.com"},
			ScrapeResult:  bson.NewObjectID(),
			MatchResults:  []model.MatchResult{},
			Logs:          []model.JobLog{},
		},
	}
	query := &model.GetJobsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(mockJobs, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(0, assert.AnError)

	total, jobs, err := suite.jobService.GetAllJobs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.Zero(total)
	suite.Nil(jobs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestJobServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JobServiceTestSuite))
}
