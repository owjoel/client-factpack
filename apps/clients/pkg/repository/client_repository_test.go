package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)

type ClientRepositorySuite struct {
	suite.Suite
	repo    repository.ClientRepository
	storage *repository.MongoStorage
	cleanup func()
	ctx     context.Context
}

func (s *ClientRepositorySuite) SetupSuite() {
	s.storage, s.cleanup = repository.NewTestMongoStorage(s.T())
	s.repo = repository.NewMongoClientRepository(s.storage)
	s.ctx = context.TODO()
}

func (s *ClientRepositorySuite) TearDownSuite() {
	s.cleanup()
}

func (s *ClientRepositorySuite) SetupTest() {
	_, err := s.storage.ClientCollection().DeleteMany(s.ctx, bson.D{})
	s.Require().NoError(err)
}

func (s *ClientRepositorySuite) TestCreateAndGetOne() {
	client := &model.Client{
		Data: bson.D{
			{Key: "profile", Value: bson.D{{Key: "names", Value: bson.A{"Alice Smith"}}}},
		},
	}

	id, err := s.repo.Create(s.ctx, client)
	s.Require().NoError(err)
	s.NotEmpty(id)

	fetched, err := s.repo.GetOne(s.ctx, id)
	s.Require().NoError(err)
	name, ok := extractName(fetched.Data)
	s.True(ok)
	s.Equal("Alice Smith", name)
}

func (s *ClientRepositorySuite) TestGetAllAndCount() {
	clients := []*model.Client{
		{Data: bson.D{{Key: "profile", Value: bson.D{{Key: "names", Value: bson.A{"Alice Smith"}}}}}},
		{Data: bson.D{{Key: "profile", Value: bson.D{{Key: "names", Value: bson.A{"Bob Lee"}}}}}},
	}

	for _, c := range clients {
		_, err := s.repo.Create(s.ctx, c)
		s.Require().NoError(err)
	}

	query := &model.GetClientsQuery{
		Name:     "Alice",
		Page:     1,
		PageSize: 10,
	}

	fetched, err := s.repo.GetAll(s.ctx, query)
	s.Require().NoError(err)
	s.Len(fetched, 1)
	name, ok := extractName(fetched[0].Data)
	s.True(ok)
	s.Equal("Alice Smith", name)

	count, err := s.repo.Count(s.ctx, query)
	s.Require().NoError(err)
	s.Equal(1, count)
}

func (s *ClientRepositorySuite) TestUpdate() {
	client := &model.Client{
		Data: bson.D{
			{Key: "profile", Value: bson.D{{Key: "names", Value: bson.A{"Initial Name"}}}},
		},
	}

	id, err := s.repo.Create(s.ctx, client)
	s.Require().NoError(err)

	update := bson.D{{Key: "data.profile.names", Value: bson.A{"Updated Name"}}}
	err = s.repo.Update(s.ctx, id, update)
	s.Require().NoError(err)

	fetched, err := s.repo.GetOne(s.ctx, id)
	s.Require().NoError(err)
	name, ok := extractName(fetched.Data)
	s.True(ok)
	s.Equal("Updated Name", name)
}

func (s *ClientRepositorySuite) TestGetClientNameByID() {
	client := &model.Client{
		Data: bson.D{
			{Key: "profile", Value: bson.D{{Key: "names", Value: bson.A{"Jane Doe"}}}},
		},
	}

	id, err := s.repo.Create(s.ctx, client)
	s.Require().NoError(err)

	name, err := s.repo.GetClientNameByID(s.ctx, id)
	s.Require().NoError(err)
	s.Equal("Jane Doe", name)
}

func extractName(data bson.D) (string, bool) {
	for _, elem := range data {
		if elem.Key == "profile" {
			if profileDoc, ok := elem.Value.(bson.D); ok {
				for _, p := range profileDoc {
					if p.Key == "names" {
						if names, ok := p.Value.(bson.A); ok && len(names) > 0 {
							if name, ok := names[0].(string); ok {
								return name, true
							}
						}
					}
				}
			}
		}
	}
	return "", false
}

func TestClientRepositorySuite(t *testing.T) {
	suite.Run(t, new(ClientRepositorySuite))
}