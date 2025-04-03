package repository

import (
	"context"
	"fmt"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoClientRepository struct {
	clientCollection *mongo.Collection
}

func NewMongoClientRepository(storage *MongoStorage) ClientRepository {
	return &mongoClientRepository{clientCollection: storage.clientCollection}
}

func (r *mongoClientRepository) Create(ctx context.Context, c *model.Client) error {
	return nil
}

type ClientRepository interface {
	// Create(ctx context.Context, c *model.Client) error
	GetOne(ctx context.Context, clientID string) (*model.Client, error)
	GetAll(ctx context.Context, query *model.GetClientsQuery) ([]model.Client, error)
	Count(ctx context.Context, query *model.GetClientsQuery) (int, error)
	// Update(ctx context.Context, c *model.Client) error
	Update(ctx context.Context, clientID string, data bson.D) error
}

func (s *mongoClientRepository) GetAll(ctx context.Context, query *model.GetClientsQuery) ([]model.Client, error) {
	var clients []model.Client

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}
	skip := (query.Page - 1) * query.PageSize

	filter := bson.M{}
	if query.Name != "" {
		filter["data.profile.names"] = bson.M{
			"$regex":   query.Name,
			"$options": "i",
		}
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.PageSize))

	cursor, err := s.clientCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo find error: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return clients, nil
}

func (s *mongoClientRepository) GetOne(ctx context.Context, clientID string) (*model.Client, error) {
	objID, err := bson.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, fmt.Errorf("error parsing object id: %w", err)
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var client model.Client
	res := s.clientCollection.FindOne(ctx, filter)
	err = res.Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no documents found: %w", err)
		}
		return nil, fmt.Errorf("Error occurred while finding client: %w", err)
	}

	return &client, nil
}

// func (s *MongoStorage) Create(ctx context.Context, c *model.Client) error {
// 	res, err := s.clientCollection.InsertOne(ctx, c)
// 	if err != nil {
// 		return fmt.Errorf("Failed to insert client: %w", err)
// 	}
// 	log.Println(res)
// 	return nil
// }

// func (s *MongoStorage) Update(ctx context.Context, c *model.Client) error {
// 	// coll := s.clientCollection
// 	return nil
// }

func (s *mongoClientRepository) Count(ctx context.Context, query *model.GetClientsQuery) (int, error) {
	count, err := s.clientCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("mongo count error: %w", err)
	}
	return int(count), nil
}

func (s *mongoClientRepository) Update(ctx context.Context, clientID string, data bson.D) error {
	objID, err := bson.ObjectIDFromHex(clientID)
	if err != nil {
		return fmt.Errorf("error parsing object id: %w", err)
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			bson.E{Key: "data", Value: data},
		}},
	}

	result, err := s.clientCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("mongo update error: %w", err)
	}

	fmt.Printf("matched: %d, modified: %d\n", result.MatchedCount, result.ModifiedCount)
	// TODO: log the update
	return nil
}
