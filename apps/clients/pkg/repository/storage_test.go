package repository

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"testing"

// 	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
// 	"github.com/owjoel/client-factpack/apps/clients/pkg/storage/mocks"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// type DatabaseTestSuite struct {
// 	suite.Suite
// 	mockStorage *mocks.ClientRepository
// 	mockClient  *model.Client
// }

// func (suite *DatabaseTestSuite) TestSetInstanceClient() {
// 	mockClient := new(mocks.ClientRepository)
// 	SetInstanceClient(mockClient)
// 	assert.Equal(suite.T(), mockClient, GetInstance().Client, "SetInstanceClient should set the ClientRepository correctly")
// }

// func (suite *DatabaseTestSuite) SetupSuite() {
// 	suite.mockStorage = new(mocks.ClientRepository)
// 	db = &storage{Client: suite.mockStorage}

// 	if err := json.Unmarshal([]byte(mocks.MockClientJSON), &suite.mockClient); err != nil {
// 		log.Fatalf("Failed to unmarshal mock client: %v", err)
// 	}
// }

// func (suite *DatabaseTestSuite) TestGetInstance() {
// 	instance := GetInstance()
// 	suite.Equal(db, instance, "GetInstance should return the same storage instance")
// }

// func (suite *DatabaseTestSuite) TestCreateClient() {
// 	suite.mockStorage.On("Create", mock.Anything, suite.mockClient).Return(nil)

// 	err := db.Client.Create(context.Background(), suite.mockClient)
// 	assert.NoError(suite.T(), err, "Should create a client successfully")
// 	suite.mockStorage.AssertExpectations(suite.T())
// }

// func (suite *DatabaseTestSuite) TestGetClient() {
// 	suite.mockStorage.On("Get", mock.Anything, "12345").Return(suite.mockClient, nil)

// 	res, err := db.Client.Get(context.Background(), "12345")
// 	assert.NoError(suite.T(), err, "Should retrieve the client successfully")
// 	assert.Equal(suite.T(), suite.mockClient, res, "Retrieved client should match expected client")
// 	suite.mockStorage.AssertExpectations(suite.T())
// }

// func (suite *DatabaseTestSuite) TestGetAllClients() {
// 	clients := []model.Client{
// 		*suite.mockClient,
// 		*suite.mockClient,
// 	}
// 	suite.mockStorage.On("GetAll", mock.Anything).Return(clients, nil)

// 	res, err := db.Client.GetAll(context.Background())
// 	assert.NoError(suite.T(), err, "Should retrieve all clients successfully")
// 	assert.Equal(suite.T(), clients, res, "Retrieved clients should match expected clients")
// 	suite.mockStorage.AssertExpectations(suite.T())
// }

// func (suite *DatabaseTestSuite) TestGetClient_NotFound() {
// 	suite.mockStorage.On("Get", mock.Anything, "99999").Return(nil, mongo.ErrNoDocuments)

// 	res, err := db.Client.Get(context.Background(), "99999")
// 	assert.Error(suite.T(), err, "Should return an error if the client is not found")
// 	assert.Nil(suite.T(), res, "Returned client should be nil if not found")
// 	suite.mockStorage.AssertExpectations(suite.T())
// }

// func TestDatabaseSuite(t *testing.T) {
// 	suite.Run(t, new(DatabaseTestSuite))
// }
