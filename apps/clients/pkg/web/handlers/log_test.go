package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LogHandlerTestSuite struct {
	suite.Suite
	mockSvc *mocks.LogServiceInterface
	handler *handlers.LogHandler
	router  *gin.Engine
}


func (suite *LogHandlerTestSuite) SetupTest() {
	suite.mockSvc = new(mocks.LogServiceInterface)
	suite.handler = handlers.NewLogHandler(suite.mockSvc)

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.RedirectTrailingSlash = false
	suite.router.GET("/logs", suite.handler.GetLogs)
	suite.router.GET("/logs/:id", suite.handler.GetLog)
	suite.router.POST("/logs", suite.handler.CreateLog)
}


func (suite *LogHandlerTestSuite) TestGetLogs_Success() {
	expectedLogs := []model.Log{{
		ClientID:  "123",
		Actor:     "test-actor",
		Operation: model.OperationGet,
		Details:   "test",
		Timestamp: time.Now(),
	}}

	suite.mockSvc.On("GetLogs", mock.Anything, mock.Anything).
		Return(1, expectedLogs, nil)

	req, _ := http.NewRequest("GET", "/logs?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockSvc.AssertExpectations(suite.T())
}

func (suite *LogHandlerTestSuite) TestGetLogs_InvalidQuery() {
	req, _ := http.NewRequest("GET", "/logs?page=abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *LogHandlerTestSuite) TestGetLogs_DefaultPagination() {
	expectedLogs := []model.Log{}
	suite.mockSvc.On("GetLogs", mock.Anything, mock.Anything).Return(1, expectedLogs, nil)

	req, _ := http.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *LogHandlerTestSuite) TestGetLogs_ErrorPaths() {
	tests := []struct {
		err      error
		expected int
	}{
		{errorx.ErrInvalidInput, http.StatusBadRequest},
		{errorx.ErrDependencyFailed, http.StatusServiceUnavailable},
		{assert.AnError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		suite.mockSvc.ExpectedCalls = nil // reset expectations
		suite.mockSvc.On("GetLogs", mock.Anything, mock.Anything).Return(0, nil, tt.err)

		req, _ := http.NewRequest("GET", "/logs?page=1", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), tt.expected, w.Code)
	}
}

func (suite *LogHandlerTestSuite) TestGetLog_Success() {
	log := &model.Log{
		ClientID:  "123",
		Actor:     "test-actor",
		Operation: model.OperationGet,
		Details:   "test",
		Timestamp: time.Now(),
	}
	suite.mockSvc.On("GetLog", mock.Anything, "abc").Return(log, nil)

	req, _ := http.NewRequest("GET", "/logs/abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *LogHandlerTestSuite) TestGetLog_ErrorPaths() {
	tests := []struct {
		err      error
		expected int
	}{
		{errorx.ErrInvalidInput, http.StatusBadRequest},
		{errorx.ErrNotFound, http.StatusNotFound},
		{assert.AnError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		suite.mockSvc.ExpectedCalls = nil
		suite.mockSvc.On("GetLog", mock.Anything, "abc").Return(&model.Log{}, tt.err)

		req, _ := http.NewRequest("GET", "/logs/abc", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), tt.expected, w.Code)
	}
}

// Note: no handler is ever called if ":id" is empty (404 from Gin)
func (suite *LogHandlerTestSuite) TestGetLog_MissingParam() {
	req, _ := http.NewRequest("GET", "/logs/", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *LogHandlerTestSuite) TestCreateLog_Success() {
	suite.mockSvc.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("abc123", nil)

	body := `{"clientId": "123", "actor": "test-actor", "operation": "get", "details": "test", "timestamp": "2021-01-01T00:00:00Z"}`
	req, _ := http.NewRequest("POST", "/logs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *LogHandlerTestSuite) TestCreateLog_InvalidJSON() {
	body := `{"details":` // malformed
	req, _ := http.NewRequest("POST", "/logs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *LogHandlerTestSuite) TestCreateLog_ErrorPaths() {
	tests := []struct {
		err      error
		expected int
	}{
		{errorx.ErrInvalidInput, http.StatusBadRequest},
		{errorx.ErrDependencyFailed, http.StatusServiceUnavailable},
		{assert.AnError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		suite.mockSvc.ExpectedCalls = nil
		suite.mockSvc.On("CreateLog", mock.Anything, mock.AnythingOfType("*model.Log")).Return("", tt.err)

		body := `{"details":"x"}`
		req, _ := http.NewRequest("POST", "/logs", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), tt.expected, w.Code)
	}
}


func TestLogHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(LogHandlerTestSuite))
}
