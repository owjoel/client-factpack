package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ArticleRepositorySuite struct {
	suite.Suite
	storage *repository.MongoStorage
	repo    repository.ArticleRepository
	cleanup func()
	ctx     context.Context
}

func (s *ArticleRepositorySuite) SetupSuite() {
	s.storage, s.cleanup = repository.NewTestMongoStorage(s.T())
	s.repo = repository.NewMongoArticleRepository(s.storage)
	s.ctx = context.TODO()
}

func (s *ArticleRepositorySuite) TearDownSuite() {
	s.cleanup()
}

func (s *ArticleRepositorySuite) SetupTest() {
	_, err := s.storage.ArticleCollection().DeleteMany(s.ctx, bson.M{})
	s.Require().NoError(err)
}

func (s *ArticleRepositorySuite) TestGetAll() {
	// Insert dummy articles
	articles := []model.Article{
		{Title: "A1", URL: "U1", Summary: "S1"},
		{Title: "A2", URL: "U2", Summary: "S2"},
	}

	res, err := s.storage.ArticleCollection().InsertMany(s.ctx, articles)
	s.Require().NoError(err)
	s.Len(res.InsertedIDs, 2)

	var ids []string
	for _, id := range res.InsertedIDs {
		ids = append(ids, id.(bson.ObjectID).Hex())
	}

	req := &model.GetArticlesReq{ID: ids}
	fetched, err := s.repo.GetAll(s.ctx, req)
	s.Require().NoError(err)
	s.Len(fetched, 2)
	s.ElementsMatch([]string{"A1", "A2"}, []string{fetched[0].Title, fetched[1].Title})
}

func TestArticleRepositorySuite(t *testing.T) {
	suite.Run(t, new(ArticleRepositorySuite))
}
