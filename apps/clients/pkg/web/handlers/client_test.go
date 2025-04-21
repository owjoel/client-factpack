package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientHandlerTestSuite struct {
	suite.Suite
	mockSvc *mocks.ClientServiceInterface
	handler *handlers.ClientHandler
	router  *gin.Engine
}

func (suite *ClientHandlerTestSuite) SetupTest() {
	suite.mockSvc = new(mocks.ClientServiceInterface)
	suite.handler = handlers.NewClientHandler(suite.mockSvc)

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.router.GET("/health", suite.handler.HealthCheck)
	suite.router.GET("/:id", suite.handler.GetClient)
	suite.router.GET("/", suite.handler.GetAllClients)
	suite.router.POST("/scrape", suite.handler.CreateClientByName)
	suite.router.PUT("/:id", suite.handler.UpdateClient)
	suite.router.POST("/:id/match", suite.handler.MatchClient)
	suite.router.POST("/:id/scrape", suite.handler.RescrapeClient)
}

func (suite *ClientHandlerTestSuite) TestHealthCheck() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Connection successful")
}

func (suite *ClientHandlerTestSuite) TestGetClient_MissingID() {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.GetClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Missing id")
}

func (suite *ClientHandlerTestSuite) TestGetClient_ServiceError() {
	suite.mockSvc.On("GetClient", mock.Anything, "abc").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/abc", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	c.Request = req

	suite.handler.GetClient(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not retrieve client")
}

func (suite *ClientHandlerTestSuite) TestCreateClientByName_InvalidJSON() {
	req, _ := http.NewRequest("POST", "/scrape", bytes.NewBufferString("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid request")
}

func (suite *ClientHandlerTestSuite) TestCreateClientByName_MissingName() {
	body := `{"name": ""}`
	req, _ := http.NewRequest("POST", "/scrape", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Missing name")
}

func (suite *ClientHandlerTestSuite) TestCreateClientByName_Success() {
	suite.mockSvc.On("CreateClientByName", mock.Anything, mock.AnythingOfType("*model.CreateClientByNameReq")).
		Return("job123", nil)

	body := `{"name": "OpenAI"}`
	req, _ := http.NewRequest("POST", "/scrape", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "job123")
	suite.mockSvc.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestCreateClientByName_ServiceError() {
	suite.mockSvc.On("CreateClientByName", mock.Anything, mock.Anything).
		Return("", assert.AnError)

	body := `{"name": "ErrorCorp"}`
	req, _ := http.NewRequest("POST", "/scrape", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not create client")
}

func (suite *ClientHandlerTestSuite) TestGetAllClients_Success() {
	suite.mockSvc.On("GetAllClients", mock.Anything, mock.Anything).
		Return(1, []model.Client{{Data: bson.D{{Key: "name", Value: "Test Corp"}}}}, nil)

	req, _ := http.NewRequest("GET", "/?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Test Corp")
}

func (suite *ClientHandlerTestSuite) TestGetAllClients_BindError() {
	req, _ := http.NewRequest("GET", "/?pageSize=bad", nil) // pageSize expects int
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid request parameters")
}

func (suite *ClientHandlerTestSuite) TestGetAllClients_ServiceError() {
	suite.mockSvc.On("GetAllClients", mock.Anything, mock.Anything).
		Return(0, nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not retrieve clients")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_MissingID() {
	req, _ := http.NewRequest("POST", "//match", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	suite.handler.MatchClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Missing id")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_TextOnly_Success() {
	suite.mockSvc.On("MatchClient", mock.Anything, mock.AnythingOfType("*model.MatchClientReq"), "abc").
		Return("job123", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("text", "Company description")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/abc/match", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "job123")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_FileOnly_Success() {
	suite.mockSvc.On("MatchClient", mock.Anything, mock.AnythingOfType("*model.MatchClientReq"), "abc").
		Return("job456", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("test content"))
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/abc/match", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "job456")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_BothFileAndText() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("text", "conflicting input")
	part, _ := writer.CreateFormFile("file", "conflict.txt")
	_, _ = part.Write([]byte("file content"))
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/abc/match", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Provide either a file or raw text")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_FileOpenFail() {
	req := httptest.NewRequest("POST", "/abc/match", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	c.Request = req
	suite.handler.MatchClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Provide either a file or raw text")
}

func (suite *ClientHandlerTestSuite) TestMatchClient_ServiceError() {
	suite.mockSvc.On("MatchClient", mock.Anything, mock.AnythingOfType("*model.MatchClientReq"), "abc").
		Return("job456", assert.AnError)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("test content"))
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/abc/match", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not match client")
}

func (suite *ClientHandlerTestSuite) TestUpdateClient_MissingID() {
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{} // no "id"

	suite.handler.UpdateClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Missing id")
}

func (suite *ClientHandlerTestSuite) TestUpdateClient_BindError() {
	req, _ := http.NewRequest("PUT", "/abc", bytes.NewBufferString(`{"invalid":`)) // malformed JSON
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	suite.handler.UpdateClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid request")
}

func (suite *ClientHandlerTestSuite) TestUpdateClient_ServiceError() {
	suite.mockSvc.On("UpdateClient", mock.Anything, "abc", mock.Anything).Return(assert.AnError)

	body := `{"changes":[{"path":"profile.name","old":"old","new":"new"}]}`
	req, _ := http.NewRequest("PUT", "/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not update client")
}

func (suite *ClientHandlerTestSuite) TestUpdateClient_Success() {
	suite.mockSvc.On("UpdateClient", mock.Anything, "abc", mock.Anything).Return(nil)

	body := `{"changes":[{"path":"profile.name","old":"old","new":"new"}]}`
	req, _ := http.NewRequest("PUT", "/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Client updated")
}

func (suite *ClientHandlerTestSuite) TestRescrapeClient_MissingID() {
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{} // no "id"

	suite.handler.RescrapeClient(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Missing id")
}

func (suite *ClientHandlerTestSuite) TestRescrapeClient_ServiceError() {
	suite.mockSvc.On("RescrapeClient", mock.Anything, "abc").Return(assert.AnError)

	req, _ := http.NewRequest("POST", "/abc/scrape", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Could not rescrape client")
}

func (suite *ClientHandlerTestSuite) TestRescrapeClient_Success() {
	suite.mockSvc.On("RescrapeClient", mock.Anything, "abc").Return(nil)

	req, _ := http.NewRequest("POST", "/abc/scrape", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Client rescraped")
}

func TestClientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ClientHandlerTestSuite))
}
