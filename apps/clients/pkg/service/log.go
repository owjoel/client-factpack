package service

import (
	"context"
	"errors"
	"fmt"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)

type LogService struct {
	logRepository repository.LogRepository
}

type LogServiceInterface interface {
	GetLogs(ctx context.Context, query *model.GetLogsQuery) (total int, logs []model.Log, err error)
	GetLog(ctx context.Context, logID string) (*model.Log, error)
	CreateLog(ctx context.Context, log *model.Log) (string, error)
}

func NewLogService(logRepository repository.LogRepository) *LogService {
	return &LogService{logRepository: logRepository}
}

func (s *LogService) GetLogs(ctx context.Context, query *model.GetLogsQuery) (int, []model.Log, error) {
	logs, err := s.logRepository.GetAll(ctx, query)
	if err != nil {
		if errors.Is(err, errorx.ErrInvalidInput) {
			return 0, nil, fmt.Errorf("%w: invalid input for log query", errorx.ErrInvalidInput)
		}
		return 0, nil, fmt.Errorf("%w: failed to retrieve logs", errorx.ErrDependencyFailed)
	}

	total, err := s.logRepository.Count(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: failed to count logs", errorx.ErrInternal)
	}

	return total, logs, nil
}

func (s *LogService) GetLog(ctx context.Context, logID string) (*model.Log, error) {
	log, err := s.logRepository.GetOne(ctx, logID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidInput):
			return nil, fmt.Errorf("%w: log ID is not valid", errorx.ErrInvalidInput)
		case errors.Is(err, errorx.ErrNotFound):
			return nil, err
		default:
			return nil, fmt.Errorf("%w: service failed to get log", errorx.ErrInternal)
		}
	}
	return log, nil
}

func (s *LogService) CreateLog(ctx context.Context, log *model.Log) (string, error) {
	if log == nil {
		return "", fmt.Errorf("%w: log is nil", errorx.ErrInvalidInput)
	}

	id, err := s.logRepository.Create(ctx, log)
	if err != nil {
		if errors.Is(err, errorx.ErrInvalidInput) || errors.Is(err, errorx.ErrDependencyFailed) {
			return "", err
		}
		return "", fmt.Errorf("%w: failed to create log", errorx.ErrInternal)
	}
	return id, nil
}
