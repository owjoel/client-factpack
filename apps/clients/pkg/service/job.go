package service

import (
	"context"
	"errors"
	"fmt"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)

type JobService struct {
	jobRepository repository.JobRepository
}

type JobServiceInterface interface {
	CreateJob(ctx context.Context, job *model.Job) (string, error)
	GetJob(ctx context.Context, jobID string) (*model.Job, error)
	GetAllJobs(ctx context.Context, query *model.GetJobsQuery) (total int, jobs []model.Job, err error)
}

func NewJobService(jobRepository repository.JobRepository) *JobService {
	return &JobService{jobRepository: jobRepository}
}

func (s *JobService) CreateJob(ctx context.Context, job *model.Job) (string, error) {
	id, err := s.jobRepository.Create(ctx, job)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return "", err
		}
		return "", fmt.Errorf("%w: error creating job", errorx.ErrInternal)
	}

	return id, nil
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*model.Job, error) {
	job, err := s.jobRepository.GetOne(ctx, jobID)
	if err != nil {
		if errors.Is(err, errorx.ErrInvalidInput) || errors.Is(err, errorx.ErrNotFound) || errors.Is(err, errorx.ErrDependencyFailed) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: error getting job", errorx.ErrInternal)
	}
	return job, nil
}

func (s *JobService) GetAllJobs(ctx context.Context, query *model.GetJobsQuery) (total int, jobs []model.Job, err error) {
	jobs, err = s.jobRepository.GetAll(ctx, query)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("%w: error getting jobs", errorx.ErrInternal)
	}

	total, err = s.jobRepository.Count(ctx, query)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("%w: failed to count jobs", errorx.ErrInternal)
	}

	return total, jobs, nil
}

