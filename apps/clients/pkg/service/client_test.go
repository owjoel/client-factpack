package service

import (
	"testing"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ClientServiceTestSuite struct {
	suite.Suite
	mockStorage *mocks.ClientInterface
	service     *ClientService
}

func (suite *ClientServiceTestSuite) SetupSuite() {
	suite.mockStorage = new(mocks.ClientInterface)
	suite.service = NewClientService(suite.mockStorage)
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockStorage = new(mocks.ClientInterface)
	suite.service = NewClientService(suite.mockStorage)
}

func (suite *ClientServiceTestSuite) TestCreateClient_Success() {
	suite.mockStorage.On("Create", mock.AnythingOfType("*model.Client")).Return(nil)

	req := &model.CreateClientReq{
		Name:        "Donald Trump",
		Age:         30,
		Nationality: "USA",
	}

	res, err := suite.service.CreateClient(req)

	assert.NoError(suite.T(), err, "CreateClient should not return an error")
	assert.NotNil(suite.T(), res, "Response should not be nil")
	assert.Equal(suite.T(), "Success", res.Status, "Status should be 'Success'")

	suite.mockStorage.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestCreateClient_Error() {
	suite.mockStorage.On("Create", mock.AnythingOfType("*model.Client")).Return(assert.AnError)

	req := &model.CreateClientReq{
		Name:        "Donald Trump",
		Age:         30,
		Nationality: "USA",
	}

	res, err := suite.service.CreateClient(req)

	assert.Error(suite.T(), err, "CreateClient should return an error")
	assert.Equal(suite.T(), "Error creating client profile: assert.AnError general error for testing", err.Error())
	assert.NotNil(suite.T(), res, "Response should not be nil")

	suite.mockStorage.AssertExpectations(suite.T())
}

func TestClientServiceSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
