package repository

import (
	"log"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	database   = "client-factpack"
	collection = "clients"
	templates  = "templates"
	jobs       = "jobs"
	logs       = "logs"
)

type MongoStorage struct {
	*mongo.Database
	clientCollection *mongo.Collection
	jobCollection    *mongo.Collection
	logCollection    *mongo.Collection
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
	clientColl := db.Collection(collection)
	jobColl := db.Collection(jobs)
	logColl := db.Collection(logs)
	return &MongoStorage{db, clientColl, jobColl, logColl}
}
