package repository

import (
	"context"
	"fmt"
	"log"


	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
)

type mongoClientRepository struct {
	clientCollection *mongo.Collection
}

func NewMongoClientRepository(storage *MongoStorage) ClientRepository {
	return &mongoClientRepository{clientCollection: storage.clientCollection}
}

type ClientRepository interface {
	Create(ctx context.Context, c *model.Client) (string, error)
	GetOne(ctx context.Context, clientID string) (*model.Client, error)
	GetAll(ctx context.Context, query *model.GetClientsQuery) ([]model.Client, error)
	Count(ctx context.Context, query *model.GetClientsQuery) (int, error)
	Update(ctx context.Context, clientID string, update bson.D) error
	GetClientNameByID(ctx context.Context, clientID string) (string, error)
}

func (r *mongoClientRepository) Create(ctx context.Context, c *model.Client) (string, error) {
	result, err := r.clientCollection.InsertOne(ctx, c)
	if err != nil {
		return "", fmt.Errorf("%w: mongo insert error", errorx.ErrDependencyFailed)
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("%w: failed to cast inserted ID to ObjectID", errorx.ErrInternal)
	}

	log.Printf("[MongoDB] Inserted client with ID: %s", insertedID.Hex())
	return insertedID.Hex(), nil
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

	if query.Sort {
		opts.SetSort(bson.D{{Key: "metadata.updatedAt", Value: -1}})
	}

	cursor, err := s.clientCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: mongo find error", errorx.ErrDependencyFailed)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("%w: decode error", errorx.ErrInternal)
	}

	return clients, nil
}

func (s *mongoClientRepository) GetOne(ctx context.Context, clientID string) (*model.Client, error) {
	objID, err := bson.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: error parsing object id", errorx.ErrInvalidInput)
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var client model.Client
	res := s.clientCollection.FindOne(ctx, filter)
	err = res.Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: no documents found", errorx.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: error occurred while finding client", errorx.ErrDependencyFailed)
	}

	return &client, nil
}

func (s *mongoClientRepository) Count(ctx context.Context, query *model.GetClientsQuery) (int, error) {
	filter := bson.M{}

	if query.Name != "" {
		filter["data.profile.names"] = bson.M{
			"$regex":   query.Name,
			"$options": "i",
		}
	}

	count, err := s.clientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%w: mongo count error", errorx.ErrDependencyFailed)
	}
	return int(count), nil
}

func (s *mongoClientRepository) Update(ctx context.Context, clientID string, update bson.D) error {
	objID, err := bson.ObjectIDFromHex(clientID)
	if err != nil {
		return fmt.Errorf("%w: error parsing object id", errorx.ErrInvalidInput)
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	updateDoc := bson.D{{Key: "$set", Value: update}}

	result, err := s.clientCollection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("%w: mongo update error", errorx.ErrDependencyFailed)
	}

	if result.MatchedCount == 0 {
		log.Printf("No client found: %s", clientID)
	}
	if result.ModifiedCount == 0 {
		log.Printf("No updates made for client: %s", clientID)
	}

	return nil
}

func (s *mongoClientRepository) GetClientNameByID(ctx context.Context, clientID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(clientID)
	if err != nil {
		return "", fmt.Errorf("%w: invalid client ID", errorx.ErrInvalidInput)
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	projection := bson.D{{Key: "data.profile.names", Value: 1}}

	var result struct {
		Data struct {
			Profile struct {
				Names []string `bson:"names"`
			} `bson:"profile"`
		} `bson:"data"`
	}

	err = s.clientCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("%w: failed to fetch client name", errorx.ErrDependencyFailed)
	}

	if len(result.Data.Profile.Names) == 0 {
		return "", fmt.Errorf("%w: client has no names", errorx.ErrNotFound)
	}

	return result.Data.Profile.Names[0], nil
}
