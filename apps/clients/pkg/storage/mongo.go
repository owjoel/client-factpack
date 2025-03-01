package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	database   = "client-factpack"
	collection = "clients"
)

type MongoStorage struct {
	*mongo.Database
	userCollection *mongo.Collection
}

func InitMongo() *MongoStorage {
	uri := config.MongoURI
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database(database)
	coll := db.Collection(collection)
	return &MongoStorage{db, coll}
}

func (s *MongoStorage) GetAll(ctx context.Context) ([]model.Client, error) {
	var clients []model.Client
	c, err := s.userCollection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []model.Client{}, fmt.Errorf("no documents found: %w", err)
		}
		return nil, fmt.Errorf("Error occurred while finding client: %w", err)
	}
	if err = c.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("Error decoding client documents: %w", err)
	}

	return clients, nil
}

func (s *MongoStorage) Get(ctx context.Context, clientID string) (*model.Client, error) {
	objID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, fmt.Errorf("error parsing object id: %w", err)
	}
	log.Println(objID)
	// filterID := bson.M{"_id": objID}

	filter := bson.M{"_id": objID}
	// log.Println(objID.Hex())
	// log.Println(filterID)

	var client model.Client
	res := s.userCollection.FindOne(ctx, filter)
	log.Println("Hello")
	log.Println()
	err = res.Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Client{}, fmt.Errorf("no documents found: %w", err)
		}
		return nil, fmt.Errorf("Error occurred while finding client: %w", err)
	}

	return &client, nil
}

func (s *MongoStorage) Create(ctx context.Context, c *model.Client) error {
	return nil
}

func (s *MongoStorage) Update(ctx context.Context, c *model.Client) error {
	// coll := s.userCollection
	return nil
}
