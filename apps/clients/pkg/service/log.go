package service

import (
	"context"

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
		return 0, nil, err
	}

	total, err := s.logRepository.Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	return total, logs, nil
}

func (s *LogService) GetLog(ctx context.Context, logID string) (*model.Log, error) {
	log, err := s.logRepository.GetOne(ctx, logID)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (s *LogService) CreateLog(ctx context.Context, log *model.Log) (string, error) {
	if log == nil {
		return "", errorx.ErrInvalidInput
	}

	id, err := s.logRepository.Create(ctx, log)
	if err != nil {
		return "", err
	}
	return id, nil
}
