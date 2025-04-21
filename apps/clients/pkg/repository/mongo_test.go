package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewTestMongoStorage(t *testing.T) (*MongoStorage, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "6.0",
		Env:        []string{},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start Mongo container: %s", err)
	}

	hostPort := resource.GetPort("27017/tcp")
	uri := fmt.Sprintf("mongodb://localhost:%s", hostPort)

	var client *mongo.Client
	if err := pool.Retry(func() error {
		client, err = mongo.Connect(options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		return client.Ping(context.TODO(), nil)
	}); err != nil {
		t.Fatalf("MongoDB did not become ready: %s", err)
	}

	db := client.Database("testdb")

	storage := &MongoStorage{
		Database:         db,
		articleCollection: db.Collection("articles"),
		clientCollection:  db.Collection("clients"),
		jobCollection:     db.Collection("jobs"),
		logCollection:     db.Collection("logs"),
	}

	cleanup := func() {
		_ = client.Disconnect(context.TODO())
		_ = pool.Purge(resource)
	}

	return storage, cleanup
}
