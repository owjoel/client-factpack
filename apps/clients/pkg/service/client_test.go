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
func (suite *ClientServiceTestSuite) TestGetClient_Success() {
	expectedClient := &model.Client{
		Name:        "Elon Musk",
		Age:         40,
		Nationality: "Canada",
	}
	suite.mockStorage.On("Get", uint(1)).Return(expectedClient, nil)

	res, err := suite.service.GetClient(1)
	assert.NoError(suite.T(), err, "GetClient should not return an error")
	assert.NotNil(suite.T(), res, "Response should not be nil")
	assert.Equal(suite.T(), expectedClient.Name, res.Name, "Name should match")
	assert.Equal(suite.T(), expectedClient.Age, res.Age, "Age should match")
	assert.Equal(suite.T(), expectedClient.Nationality, res.Nationality, "Nationality should match")

	suite.mockStorage.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestGetClient_Error() {
	suite.mockStorage.On("Get", uint(2)).Return(&model.Client{}, assert.AnError)

	res, err := suite.service.GetClient(2)
	assert.Error(suite.T(), err, "GetClient should return an error")
	assert.Equal(suite.T(), "Error retrieving client profile: assert.AnError general error for testing", err.Error())

	assert.NotNil(suite.T(), res, "Response should not be nil")

	suite.mockStorage.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_Success() {
	suite.mockStorage.On("Update", mock.AnythingOfType("*model.Client")).Return(nil)

	req := model.UpdateClientReq{
		Name:        "Jane Doe",
		Age:         28,
		Nationality: "UK",
	}

	res, err := suite.service.UpdateClient(req)
	assert.NoError(suite.T(), err, "UpdateClient should not return an error")
	assert.NotNil(suite.T(), res, "Response should not be nil")
	assert.Equal(suite.T(), "Success", res.Status, "Status should be 'Success'")

	suite.mockStorage.AssertExpectations(suite.T())
}

func (suite *ClientServiceTestSuite) TestUpdateClient_Error() {
	suite.mockStorage.On("Update", mock.AnythingOfType("*model.Client")).Return(assert.AnError)

	req := model.UpdateClientReq{
		Name:        "Kylie Jenner",
		Age:         28,
		Nationality: "USA",
	}

	res, err := suite.service.UpdateClient(req)
	assert.Error(suite.T(), err, "UpdateClient should return an error")
	assert.Nil(suite.T(), res, "Response should be nil on error")
	assert.Equal(suite.T(), "Error updating client profile: assert.AnError general error for testing", err.Error())

	suite.mockStorage.AssertExpectations(suite.T())
}

func TestClientServiceSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
