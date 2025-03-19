package service

import (
    "context"
    "errors"
	"fmt"
    "testing"


    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"

    "github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/storage"
    "github.com/owjoel/client-factpack/apps/clients/pkg/storage/mocks"
	

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientServiceTestSuite struct {
    suite.Suite
    mockStorage *mocks.ClientInterface
    service     *ClientService
}


func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockStorage = new(mocks.ClientInterface)
	storage.SetInstanceClient(suite.mockStorage)
	suite.service = NewClientService(nil)


	// fmt.Printf("Service: %+v\n", suite.service)
	// fmt.Printf("Storage Singleton: %+v\n", storage.GetInstance())
	// fmt.Printf("Storage Singleton Client: %+v\n", storage.GetInstance().Client)
}

// CreateClient Test


func (suite *ClientServiceTestSuite) TestCreateClient_Success() {
	suite.mockStorage.
		On("Create", mock.Anything, mock.AnythingOfType("*model.Client")).
		Return(nil)

	newClient := &model.Client{
		Profile: model.Profile{
			Name:        "Donald Trump",
			Age:         30,
			Nationality: "USA",
		},
	}

	err := suite.service.CreateClient(context.Background(), newClient)

	assert.NoError(suite.T(), err)
	suite.mockStorage.AssertExpectations(suite.T())
}



func (suite *ClientServiceTestSuite) TestCreateClient_Error() {
	suite.mockStorage.
		On("Create", mock.Anything, mock.AnythingOfType("*model.Client")).
		Return(errors.New("some DB error"))

	newClient := &model.Client{
		Profile: model.Profile{
			Name:        "Donald Trump",
			Age:         30,
			Nationality: "USA",
		},
	}

	err := suite.service.CreateClient(context.Background(), newClient)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "some DB error")
	suite.mockStorage.AssertExpectations(suite.T())
}


// GetClient Test


func (suite *ClientServiceTestSuite) TestGetClient_Success() {
	mockClient := &model.Client{
		ID: bson.NewObjectID(),
		Profile: model.Profile{
			Name: "John Doe",
		},
	}

	suite.mockStorage.
		On("Get", mock.Anything, "client123").
		Return(mockClient, nil)

	retrieved, err := suite.service.GetClient(context.Background(), "client123")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockClient, retrieved)
	suite.mockStorage.AssertExpectations(suite.T())
}



func (suite *ClientServiceTestSuite) TestGetClient_Error() {
	suite.mockStorage.
		On("Get", mock.Anything, "notfound").
		Return(nil, errors.New("no documents found"))

	retrieved, err := suite.service.GetClient(context.Background(), "notfound")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "no documents found")

	//  Instead of assert.Nil, check if it's empty
	assert.Equal(suite.T(), &model.Client{}, retrieved)
	suite.mockStorage.AssertExpectations(suite.T())
}


// GetAllClients


func (suite *ClientServiceTestSuite) TestGetAllClients_Success() {
	mockClients := []model.Client{
		{
			ID: bson.NewObjectID(),
			Profile: model.Profile{
				Name: "Alice",
			},
		},
		{
			ID: bson.NewObjectID(),
			Profile: model.Profile{
				Name: "Bob",
			},
		},
	}

	suite.mockStorage.
		On("GetAll", mock.Anything).
		Return(mockClients, nil)

	retrieved, err := suite.service.GetAllClients(context.Background())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockClients, retrieved)
	suite.mockStorage.AssertExpectations(suite.T())
}


func (suite *ClientServiceTestSuite) TestGetAllClients_Error() {
	suite.mockStorage.
		On("GetAll", mock.Anything).
		Return(nil, errors.New("DB failure"))

	retrieved, err := suite.service.GetAllClients(context.Background())

	assert.Error(suite.T(), err, "Expected an error when DB fails")
	assert.Nil(suite.T(), retrieved, "Expected no clients returned on DB error")
	assert.Contains(suite.T(), err.Error(), "DB failure")

	suite.mockStorage.AssertExpectations(suite.T())
}

func TestClientServiceSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
