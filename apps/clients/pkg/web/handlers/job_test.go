package handlers_test

import (
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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobHandlerTestSuite struct {
	suite.Suite
	mockSvc *mocks.JobServiceInterface
	handler *handlers.JobHandler
	router  *gin.Engine
}

func (suite *JobHandlerTestSuite) SetupTest() {
	suite.mockSvc = new(mocks.JobServiceInterface)
	suite.handler = handlers.NewJobHandler(suite.mockSvc)

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.GET("/jobs", suite.handler.GetAllJobs)
	suite.router.GET("/jobs/:id", suite.handler.GetJob)
}

func (suite *JobHandlerTestSuite) TestGetJob_MissingID() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	suite.handler.GetJob(c)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *JobHandlerTestSuite) TestGetJob_Success() {
	id := primitive.NewObjectID()
	job := &model.Job{
		PrefectFlowID: "prefect-flow-id-123",
		Type:          model.Scrape,
		Status:        model.JobStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Input:         bson.M{"url": "https://example.com"},
		ScrapeResult:  bson.NewObjectID(),
		MatchResults:  []model.MatchResult{},
		Logs:          []model.JobLog{},
	}
	suite.mockSvc.On("GetJob", mock.Anything, id.Hex()).Return(job, nil)

	req, _ := http.NewRequest("GET", "/jobs/"+id.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockSvc.AssertExpectations(suite.T())
}

func (suite *JobHandlerTestSuite) TestGetJob_ErrorCases() {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"InvalidInput", errorx.ErrInvalidInput, http.StatusBadRequest},
		{"NotFound", errorx.ErrNotFound, http.StatusNotFound},
		{"GenericError", assert.AnError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockSvc.ExpectedCalls = nil
			suite.mockSvc.On("GetJob", mock.Anything, "invalidID").
				Return((*model.Job)(nil), tt.err)

			req, _ := http.NewRequest("GET", "/jobs/invalidID", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expected, w.Code)
			suite.mockSvc.AssertExpectations(suite.T())
		})
	}
}

func (suite *JobHandlerTestSuite) TestGetAllJobs_Success() {
	suite.mockSvc.On("GetAllJobs", mock.Anything, mock.Anything).
		Return(1, []model.Job{{Status: model.JobStatusPending}}, nil)

	req, _ := http.NewRequest("GET", "/jobs?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockSvc.AssertExpectations(suite.T())
}

func (suite *JobHandlerTestSuite) TestGetAllJobs_BindError() {
	req, _ := http.NewRequest("GET", "/jobs?page=abc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *JobHandlerTestSuite) TestGetAllJobs_ServiceErrors() {
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
		suite.mockSvc.On("GetAllJobs", mock.Anything, mock.Anything).Return(0, nil, tt.err)

		req, _ := http.NewRequest("GET", "/jobs?page=1", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), tt.expected, w.Code)
	}
}

func TestJobHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(JobHandlerTestSuite))
}
