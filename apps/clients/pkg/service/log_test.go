package service_test

import (
	"context"
	"testing"
	"time"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LogServiceTestSuite struct {
	suite.Suite
	mockRepo   *mocks.LogRepository
	logService *service.LogService
}

func (suite *LogServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.LogRepository)
	suite.logService = service.NewLogService(suite.mockRepo)
}

func (suite *LogServiceTestSuite) TestCreateLog() {
	log := &model.Log{
		ClientID:  "test-client-id",
		Actor:     "test-actor",
		Operation: model.OperationCreate,
		Details:   "test-details",
		Timestamp: time.Now(),
	}
	suite.mockRepo.On("Create", mock.Anything, log).Return("test-log-id", nil)

	id, err := suite.logService.CreateLog(context.Background(), log)

	suite.NoError(err)
	suite.Equal("test-log-id", id)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *LogServiceTestSuite) TestCreateLog_NilInput() {
	id, err := suite.logService.CreateLog(context.Background(), nil)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
}

func (suite *LogServiceTestSuite) TestCreateLog_RepoReturnsInvalidInput() {
	log := &model.Log{}
	suite.mockRepo.On("Create", mock.Anything, log).Return("", errorx.ErrInvalidInput)

	id, err := suite.logService.CreateLog(context.Background(), log)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
}

func (suite *LogServiceTestSuite) TestCreateLog_RepoReturnsDependencyFailed() {
	log := &model.Log{}
	suite.mockRepo.On("Create", mock.Anything, log).Return("", errorx.ErrDependencyFailed)

	id, err := suite.logService.CreateLog(context.Background(), log)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
}

func (suite *LogServiceTestSuite) TestCreateLog_RepoReturnsError() {
	log := &model.Log{}
	suite.mockRepo.On("Create", mock.Anything, log).Return("", errorx.ErrInternal)

	id, err := suite.logService.CreateLog(context.Background(), log)

	suite.Error(err)
	suite.Empty(id)
	suite.ErrorIs(err, errorx.ErrInternal)
}

func (suite *LogServiceTestSuite) TestGetLog() {
	log := &model.Log{
		ClientID:  "test-client-id",
		Actor:     "test-actor",
		Operation: model.OperationCreate,
		Details:   "test-details",
		Timestamp: time.Now(),
	}
	suite.mockRepo.On("GetOne", mock.Anything, "test-log-id").Return(log, nil)

	log, err := suite.logService.GetLog(context.Background(), "test-log-id")

	suite.NoError(err)
	suite.Equal(log, log)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *LogServiceTestSuite) TestGetLog_InvalidInput() {
	suite.mockRepo.On("GetOne", mock.Anything, "bad-id").Return(nil, errorx.ErrInvalidInput)

	log, err := suite.logService.GetLog(context.Background(), "bad-id")

	suite.Error(err)
	suite.Nil(log)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
}

func (suite *LogServiceTestSuite) TestGetLog_NotFound() {
	suite.mockRepo.On("GetOne", mock.Anything, "not-found").Return(nil, errorx.ErrNotFound)

	log, err := suite.logService.GetLog(context.Background(), "not-found")

	suite.Error(err)
	suite.Nil(log)
	suite.ErrorIs(err, errorx.ErrNotFound)
}

func (suite *LogServiceTestSuite) TestGetLog_UnexpectedError() {
	suite.mockRepo.On("GetOne", mock.Anything, "id").Return(nil, errorx.ErrInternal)

	log, err := suite.logService.GetLog(context.Background(), "id")

	suite.Error(err)
	suite.Nil(log)
	suite.ErrorIs(err, errorx.ErrInternal)
}

func (suite *LogServiceTestSuite) TestGetLogs_UnexpectedError() {
	suite.mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(nil, errorx.ErrInternal)

	total, logs, err := suite.logService.GetLogs(context.Background(), &model.GetLogsQuery{})

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.Zero(total)
	suite.Nil(logs)
}

func (suite *LogServiceTestSuite) TestGetLogs_GetAllFails_InvalidInput() {
	query := &model.GetLogsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, errorx.ErrInvalidInput)

	total, logs, err := suite.logService.GetLogs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
	suite.Zero(total)
	suite.Nil(logs)
}

func (suite *LogServiceTestSuite) TestGetLogs_CountFails() {
	query := &model.GetLogsQuery{}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return([]model.Log{}, nil)
	suite.mockRepo.On("Count", mock.Anything).Return(0, errorx.ErrDependencyFailed)

	total, logs, err := suite.logService.GetLogs(context.Background(), query)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.Zero(total)
	suite.Nil(logs)
}

func (suite *LogServiceTestSuite) TestGetLogs() {
	query := &model.GetLogsQuery{
		Page:     1,
		PageSize: 10,
	}
	logs := []model.Log{
		{
			ClientID:  "test-client-id",
			Actor:     "test-actor",
			Operation: model.OperationCreate,
			Details:   "test-details",
			Timestamp: time.Now(),
		},
	}
	suite.mockRepo.On("GetAll", mock.Anything, query).Return(logs, nil)
	suite.mockRepo.On("Count", mock.Anything).Return(1, nil)

	total, logs, err := suite.logService.GetLogs(context.Background(), query)

	suite.NoError(err)
	suite.Equal(1, total)
	suite.Equal(logs, logs)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestLogServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LogServiceTestSuite))
}
