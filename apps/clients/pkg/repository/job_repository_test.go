package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)

type JobRepositorySuite struct {
	suite.Suite
	repo    repository.JobRepository
	storage *repository.MongoStorage
	cleanup func()
	ctx     context.Context
}

func (s *JobRepositorySuite) SetupSuite() {
	s.storage, s.cleanup = repository.NewTestMongoStorage(s.T())
	s.repo = repository.NewMongoJobRepository(s.storage)
	s.ctx = context.TODO()
}

func (s *JobRepositorySuite) TearDownSuite() {
	s.cleanup()
}

func (s *JobRepositorySuite) SetupTest() {
	_, err := s.storage.JobCollection().DeleteMany(s.ctx, map[string]any{})
	s.Require().NoError(err)
}

func (s *JobRepositorySuite) TestCreateAndGetOne() {
	_, err := s.storage.JobCollection().DeleteMany(s.ctx, map[string]any{})
	s.Require().NoError(err)

	job := &model.Job{
		Type:      model.Scrape,
		Status:    model.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.Create(s.ctx, job)
	s.Require().NoError(err)
	s.NotEmpty(id)

	fetched, err := s.repo.GetOne(s.ctx, id)
	s.Require().NoError(err)
	s.Equal(job.Type, fetched.Type)
	s.Equal(job.Status, fetched.Status)
}

func (s *JobRepositorySuite) TestGetAll() {
	_, err := s.storage.JobCollection().DeleteMany(s.ctx, map[string]any{})
	s.Require().NoError(err)

	jobs := []*model.Job{
		{Type: model.Scrape, Status: model.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Type: model.Match, Status: model.JobStatusCompleted, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Type: model.Scrape, Status: model.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, job := range jobs {
		_, err := s.repo.Create(s.ctx, job)
		s.Require().NoError(err)
	}

	query := &model.GetJobsQuery{
		Status:   model.JobStatusPending,
		Page:     1,
		PageSize: 10,
	}

	result, err := s.repo.GetAll(s.ctx, query)
	s.Require().NoError(err)
	s.Len(result, 2)
	for _, j := range result {
		s.Equal(model.JobStatusPending, j.Status)
	}
}

func (s *JobRepositorySuite) TestCount() {
	_, err := s.storage.JobCollection().DeleteMany(s.ctx, map[string]any{})
	s.Require().NoError(err)

	// Insert 2 pending, 1 completed
	jobs := []*model.Job{
		{Status: model.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Status: model.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Status: model.JobStatusCompleted, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, job := range jobs {
		_, err := s.repo.Create(s.ctx, job)
		s.Require().NoError(err)
	}

	query := &model.GetJobsQuery{Status: model.JobStatusPending}
	count, err := s.repo.Count(s.ctx, query)
	s.Require().NoError(err)
	s.Equal(2, count)
}

func TestJobRepositorySuite(t *testing.T) {
	suite.Run(t, new(JobRepositorySuite))
}
