package service_test

import (
	"context"
	"testing"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGetClient(t *testing.T) {
	clientRepo := new(mocks.ClientRepository)
	jobService := new(mocks.JobServiceInterface)
	clientService := service.NewClientService(clientRepo, jobService)

	clientID := bson.NewObjectID()
	expectedClient := &model.Client{ID: clientID}

	clientRepo.On("GetOne", mock.Anything, clientID.Hex()).Return(expectedClient, nil)

	client, err := clientService.GetClient(context.Background(), clientID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
	clientRepo.AssertExpectations(t)
}

func TestGetAllClients(t *testing.T) {
	clientRepo := new(mocks.ClientRepository)
	jobService := new(mocks.JobServiceInterface)
	clientService := service.NewClientService(clientRepo, jobService)

	query := &model.GetClientsQuery{Page: 1, PageSize: 10}
	expectedClients := []model.Client{{ID: bson.NewObjectID()}, {ID: bson.NewObjectID()}}
	expectedTotal := 2

	clientRepo.On("GetAll", mock.Anything, query).Return(expectedClients, nil)
	clientRepo.On("Count", mock.Anything).Return(expectedTotal, nil)

	total, clients, err := clientService.GetAllClients(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, total)
	assert.Equal(t, expectedClients, clients)
	clientRepo.AssertExpectations(t)
}

func TestCreateClientByName(t *testing.T) {
	clientRepo := new(mocks.ClientRepository)
	jobService := new(mocks.JobServiceInterface)
	clientService := service.NewClientService(clientRepo, jobService)

	req := &model.CreateClientByNameReq{Name: "Test Client"}

	// Assuming CreateClientByName triggers some job or workflow
	// Mock the necessary interactions here

	err := clientService.CreateClientByName(context.Background(), req)

	assert.NoError(t, err)
	clientRepo.AssertExpectations(t)
}

func TestUpdateClient(t *testing.T) {
	clientRepo := new(mocks.ClientRepository)
	jobService := new(mocks.JobServiceInterface)
	clientService := service.NewClientService(clientRepo, jobService)

	clientID := "123"
	data := bson.D{{Key: "name", Value: "Updated Name"}}

	clientRepo.On("GetOne", mock.Anything, clientID).Return(&model.Client{ID: bson.NewObjectID()}, nil)
	clientRepo.On("Update", mock.Anything, clientID, data).Return(nil)

	err := clientService.UpdateClient(context.Background(), clientID, data)

	assert.NoError(t, err)
	clientRepo.AssertExpectations(t)
}
