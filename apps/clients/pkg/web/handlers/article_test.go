package handlers_test

import (
	"bytes"
	"encoding/json"
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

type ArticleHandlerTestSuite struct {
	suite.Suite
	mockSvc *mocks.ArticleServiceInterface
	handler *handlers.ArticleHandler
	router  *gin.Engine
}

func (suite *ArticleHandlerTestSuite) SetupTest() {
	suite.mockSvc = new(mocks.ArticleServiceInterface)
	suite.handler = handlers.NewArticleHandler(suite.mockSvc)

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	
	suite.router.POST("/articles", suite.handler.GetAllArticles)
}

func (suite *ArticleHandlerTestSuite) TestGetAllArticles_Success() {

	id1 := bson.NewObjectID()
	id2 := bson.NewObjectID()
	id3 := bson.NewObjectID()

	suite.mockSvc.On("GetAllArticles", mock.Anything, mock.Anything).Return([]model.Article{
		{ID: id1, Title: "Article 1", URL: "https://example.com/1", Summary: "Summary 1"},
		{ID: id2, Title: "Article 2", URL: "https://example.com/2", Summary: "Summary 2"},
		{ID: id3, Title: "Article 3", URL: "https://example.com/3", Summary: "Summary 3"},
	}, nil)

	w := httptest.NewRecorder()
	body := model.GetArticlesReq{
		ID: []string{id1.Hex(), id2.Hex(), id3.Hex()},
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonBody))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ArticleHandlerTestSuite) TestGetAllArticles_BindError() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid request body")
}

func TestArticleHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}
