package service_test

import (
	"context"
	"testing"

	"github.com/owjoel/client-factpack/apps/clients/config"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientServiceTestSuite struct {
	suite.Suite
	clientService *service.ClientService
	mockRepo      *mocks.ClientRepository
	mockLog       *mocks.LogServiceInterface
	mockJob       *mocks.JobServiceInterface
	mockPrefect   *mocks.PrefectFlowRunnerInterface
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.ClientRepository)
	suite.mockLog = new(mocks.LogServiceInterface)
	suite.mockJob = new(mocks.JobServiceInterface)
	suite.mockPrefect = new(mocks.PrefectFlowRunnerInterface)
	suite.clientService = service.NewClientService(suite.mockRepo, suite.mockJob, suite.mockLog, suite.mockPrefect)
}

func (suite *ClientServiceTestSuite) TestGetClient() {
	clientID := "test-client-id"
	expectedClient := &model.Client{}

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(expectedClient, nil)
	suite.mockLog.On("CreateLog", mock.Anything, mock.Anything).Return("", nil)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.NoError(err)
	suite.Equal(expectedClient, client)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error() {
	clientID := "test-client-id"

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, assert.AnError)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.Error(err)
	suite.Nil(client)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error_NotFound() {
	clientID := "test-client-id"

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, errorx.ErrNotFound)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.Error(err)
	suite.Nil(client)
	suite.ErrorIs(err, errorx.ErrNotFound)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error_DependencyFailed() {
	clientID := "test-client-id"

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, errorx.ErrDependencyFailed)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.Error(err)
	suite.Nil(client)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error_InvalidInput() {
	clientID := "test-client-id"

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, errorx.ErrInvalidInput)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.Error(err)
	suite.Nil(client)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error_CreateLogError() {
	clientID := "test-client-id"
	expectedClient := &model.Client{}

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(expectedClient, nil)
	suite.mockLog.On("CreateLog", mock.Anything, mock.Anything).Return("", assert.AnError)

	client, err := suite.clientService.GetClient(context.Background(), clientID)

	suite.NoError(err)
	suite.Equal(expectedClient, client)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetAllClients() {
	query := &model.GetClientsQuery{}
	expectedClients := []model.Client{}
	expectedTotal := 10

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(expectedClients, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(expectedTotal, nil)

	total, clients, err := suite.clientService.GetAllClients(context.Background(), query)

	suite.NoError(err)
	suite.Equal(expectedTotal, total)
	suite.Equal(expectedClients, clients)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetAllClients_Error() {
	query := &model.GetClientsQuery{}

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, assert.AnError)

	total, clients, err := suite.clientService.GetAllClients(context.Background(), query)

	suite.Error(err)
	suite.Equal(0, total)
	suite.Nil(clients)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetAllClients_Error_DependencyFailed() {
	query := &model.GetClientsQuery{}

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, errorx.ErrDependencyFailed)

	total, clients, err := suite.clientService.GetAllClients(context.Background(), query)

	suite.Error(err)
	suite.Equal(0, total)
	suite.Nil(clients)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetAllClients_Error_Count() {
	query := &model.GetClientsQuery{}

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(0, assert.AnError)

	total, clients, err := suite.clientService.GetAllClients(context.Background(), query)

	suite.Error(err)
	suite.Equal(0, total)
	suite.Nil(clients)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetAllClients_Error_CountDependencyFailed() {
	query := &model.GetClientsQuery{}

	suite.mockRepo.On("GetAll", mock.Anything, query).Return(nil, nil)
	suite.mockRepo.On("Count", mock.Anything, query).Return(0, errorx.ErrDependencyFailed)

	total, clients, err := suite.clientService.GetAllClients(context.Background(), query)

	suite.Error(err)
	suite.Equal(0, total)
	suite.Nil(clients)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	expectedClientID := "client-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(expectedClientID, nil)

	suite.mockPrefect.
		On("Trigger", config.PrefectScrapeFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
			return params["job_id"] == expectedJobID &&
				params["client_id"] == expectedClientID &&
				params["target"] == req.Name &&
				params["username"] == username
		})).
		Return(nil)

	suite.mockLog.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("test-log-id", nil)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.NoError(err)
	suite.Equal(expectedJobID, jobID)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName_CreateJobError() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, assert.AnError)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.Error(err)
	suite.Equal("", jobID)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName_CreateJobDependencyFailed() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, errorx.ErrDependencyFailed)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.Error(err)
	suite.Equal("", jobID)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName_CreateClientError() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	username := "test-user"
	expectedClientID := "client-id"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(expectedClientID, assert.AnError)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.Error(err)
	suite.Equal("", jobID)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName_CreateClientDependencyFailed() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	username := "test-user"
	expectedClientID := "client-id"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(expectedClientID, errorx.ErrDependencyFailed)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.Error(err)
	suite.Equal("", jobID)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClientByName_TriggerError() {
	req := &model.CreateClientByNameReq{Name: "Test Client"}
	expectedJobID := "job-id"
	username := "test-user"
	expectedClientID := "client-id"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(expectedClientID, nil)
	suite.mockPrefect.On("Trigger", config.PrefectScrapeFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == expectedJobID &&
			params["client_id"] == expectedClientID &&
			params["target"] == req.Name &&
			params["username"] == username
	})).Return(assert.AnError)

	jobID, err := suite.clientService.CreateClientByName(ctx, req)

	suite.Error(err)
	suite.Equal("", jobID)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("GetClientNameByID", mock.Anything, clientID).Return("Test Client", nil)
	suite.mockPrefect.On("Trigger", config.PrefectScrapeFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == expectedJobID &&
			params["client_id"] == clientID &&
			params["target"] == "Test Client" &&
			params["username"] == username
	})).Return(nil)
	suite.mockLog.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("test-log-id", nil)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_CreateJobError() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, assert.AnError)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_CreateJobDependencyFailed() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, errorx.ErrDependencyFailed)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_GetClientNameByIDError() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("GetClientNameByID", mock.Anything, clientID).Return("", assert.AnError)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_GetClientNameByIDDependencyFailed() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("GetClientNameByID", mock.Anything, clientID).Return("", errorx.ErrDependencyFailed)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_TriggerError() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("GetClientNameByID", mock.Anything, clientID).Return("Test Client", nil)
	suite.mockPrefect.On("Trigger", config.PrefectScrapeFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == expectedJobID &&
			params["client_id"] == clientID &&
			params["target"] == "Test Client" &&
			params["username"] == username
	})).Return(assert.AnError)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestRescrapeClient_CreateLogError() {
	clientID := "test-client-id"
	expectedJobID := "job-id"
	username := "test-user"

	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return(expectedJobID, nil)
	suite.mockRepo.On("GetClientNameByID", mock.Anything, clientID).Return("Test Client", nil)
	suite.mockPrefect.On("Trigger", config.PrefectScrapeFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == expectedJobID &&
			params["client_id"] == clientID &&
			params["target"] == "Test Client" &&
			params["username"] == username
	})).Return(nil)
	suite.mockLog.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("", assert.AnError)

	err := suite.clientService.RescrapeClient(ctx, clientID)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	expectedUpdate := bson.D{{Key: "data.test-path-name", Value: "New Name"}}
	expectedClient := &model.Client{}

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(expectedClient, nil)
	suite.mockRepo.On("Update", mock.Anything, clientID, expectedUpdate).Return(nil)

	suite.mockLog.On("CreateLog", mock.Anything, mock.MatchedBy(func(log *model.Log) bool {
		return log.ClientID == clientID &&
			log.Actor == username &&
			log.Operation == model.OperationUpdate &&
			log.Details != ""
	})).Return("test-log-id", nil)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_GetOneError() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, assert.AnError)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_NotFound() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, nil)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrNotFound)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_GetOneDependencyFailed() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(nil, errorx.ErrDependencyFailed)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_UpdateError() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(&model.Client{}, nil)
	suite.mockRepo.On("Update", mock.Anything, clientID, mock.Anything).Return(assert.AnError)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_UpdateDependencyFailed() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(&model.Client{}, nil)
	suite.mockRepo.On("Update", mock.Anything, clientID, mock.Anything).Return(errorx.ErrDependencyFailed)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_UpdateInvalidInput() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{}
	ctx := context.WithValue(context.Background(), "username", username)
	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(&model.Client{}, nil)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_CreateLogError() {
	clientID := "test-client-id"
	username := "test-user"
	changes := []model.SimpleChanges{
		{
			Path: "test-path-name",
			Old:  "Old Name",
			New:  "New Name",
		},
	}
	ctx := context.WithValue(context.Background(), "username", username)
	suite.mockRepo.On("GetOne", mock.Anything, clientID).Return(&model.Client{}, nil)
	suite.mockRepo.On("Update", mock.Anything, clientID, mock.Anything).Return(nil)
	suite.mockLog.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("", assert.AnError)

	err := suite.clientService.UpdateClient(ctx, clientID, changes)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockLog.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestMatchClient() {
	clientID := "test-client-id"
	username := "test-user"
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return("job-id", nil)

	suite.mockPrefect.On("Trigger", config.PrefectMatchFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == "job-id" &&
			params["file_name"] == "test-file-name" &&
			params["file_bytes"] == "test-file-bytes" &&
			params["target_id"] == clientID &&
			params["username"] == username
	})).Return(nil)

	jobID, err := suite.clientService.MatchClient(ctx, &model.MatchClientReq{
		FileName:  "test-file-name",
		FileBytes: "test-file-bytes",
	}, clientID)

	suite.NoError(err)
	suite.Equal("job-id", jobID)
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestMatchClient_CreateJobError() {
	clientID := "test-client-id"
	username := "test-user"
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return("", assert.AnError)

	jobID, err := suite.clientService.MatchClient(ctx, &model.MatchClientReq{
		FileName:  "test-file-name",
		FileBytes: "test-file-bytes",
	}, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.Empty(jobID)
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestMatchClient_CreateJobDependencyFailed() {
	clientID := "test-client-id"
	username := "test-user"
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return("", errorx.ErrDependencyFailed)
	
	jobID, err := suite.clientService.MatchClient(ctx, &model.MatchClientReq{
		FileName:  "test-file-name",
		FileBytes: "test-file-bytes",
	}, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrDependencyFailed)
	suite.Empty(jobID)
	suite.mockJob.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestMatchClient_TriggerError() {
	clientID := "test-client-id"
	username := "test-user"
	ctx := context.WithValue(context.Background(), "username", username)

	suite.mockJob.On("CreateJob", mock.Anything, mock.Anything).Return("job-id", nil)
	suite.mockPrefect.On("Trigger", config.PrefectMatchFlowID, mock.MatchedBy(func(params map[string]interface{}) bool {
		return params["job_id"] == "job-id" &&
			params["file_name"] == "test-file-name" &&
			params["file_bytes"] == "test-file-bytes" &&
			params["target_id"] == clientID &&
			params["username"] == username
	})).Return(assert.AnError)

	jobID, err := suite.clientService.MatchClient(ctx, &model.MatchClientReq{
		FileName:  "test-file-name",
		FileBytes: "test-file-bytes",
	}, clientID)

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInternal)
	suite.Empty("", jobID)
	suite.mockJob.AssertExpectations(suite.T())
	suite.mockPrefect.AssertExpectations(suite.T())
}

func TestClientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
