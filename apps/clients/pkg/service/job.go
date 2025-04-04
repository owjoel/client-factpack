package service

import (
	"context"
	"fmt"

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
		return "", fmt.Errorf("error creating job: %w", err)
	}

	return id, nil
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*model.Job, error) {
	job, err := s.jobRepository.GetOne(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobService) GetAllJobs(ctx context.Context, query *model.GetJobsQuery) (total int, jobs []model.Job, err error) {
	jobs, err = s.jobRepository.GetAll(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	total, err = s.jobRepository.Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	return total, jobs, nil
}

