package repository

import (
	"log"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	database   = "client-factpack"
	article    = "articles"
	collection = "clients"
	templates  = "templates"
	jobs       = "jobs"
	logs       = "logs"
)

type MongoStorage struct {
	*mongo.Database
	articleCollection *mongo.Collection
	clientCollection  *mongo.Collection
	jobCollection     *mongo.Collection
	logCollection     *mongo.Collection
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
	articleColl := db.Collection(article)
	clientColl := db.Collection(collection)
	jobColl := db.Collection(jobs)
	logColl := db.Collection(logs)
	return &MongoStorage{db, articleColl, clientColl, jobColl, logColl}
}

func (s *MongoStorage) JobCollection() *mongo.Collection {
	return s.jobCollection
}

func (s *MongoStorage) ArticleCollection() *mongo.Collection {
	return s.articleCollection
}

func (s *MongoStorage) ClientCollection() *mongo.Collection {
	return s.clientCollection
}
