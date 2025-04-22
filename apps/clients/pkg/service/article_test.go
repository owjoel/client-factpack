package service_test

import (
	"context"
	"testing"
	// "time"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/mocks"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ArticleServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.ArticleRepository
	articleService *service.ArticleService
}


func (suite *ArticleServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.ArticleRepository)
	suite.articleService = service.NewArticleService(suite.mockRepo)
}

func (suite *ArticleServiceTestSuite) TestGetAllArticles() {
	articles := []model.Article{
		{ID: bson.NewObjectID(), Title: "Article 1", URL: "https://example.com/1", Summary: "Summary 1"},
		{ID: bson.NewObjectID(), Title: "Article 2", URL: "https://example.com/2", Summary: "Summary 2"},
	}
	suite.mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(articles, nil)
}

func (suite *ArticleServiceTestSuite) TestGetAllArticles_Error() {
	suite.mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	articles, err := suite.articleService.GetAllArticles(context.Background(), &model.GetArticlesReq{
		ID: []string{bson.NewObjectID().Hex(), bson.NewObjectID().Hex()},
	})

	suite.Error(err)
	suite.Nil(articles)
}

func (suite *ArticleServiceTestSuite) TestGetAllArticles_InvalidInput() {
	suite.mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(nil, errorx.ErrInvalidInput)
	articles, err := suite.articleService.GetAllArticles(context.Background(), &model.GetArticlesReq{
		ID: []string{"invalid-id"},
	})

	suite.Error(err)
	suite.ErrorIs(err, errorx.ErrInvalidInput)
	suite.Nil(articles)
}

func TestArticleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleServiceTestSuite))
}

