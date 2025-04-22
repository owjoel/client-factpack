package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)


func TestMongoLogRepository_CreateAndGet(t *testing.T) {
	storage, cleanup := repository.NewTestMongoStorage(t)
	defer cleanup()

	repo := repository.NewMongoLogRepository(storage)

	logEntry := &model.Log{
		Actor:     "tester",
		ClientID:  "test-id",
		Operation: model.OperationCreate,
		Timestamp: time.Now(),
	}

	id, err := repo.Create(context.TODO(), logEntry)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	fetched, err := repo.GetOne(context.TODO(), id)
	assert.NoError(t, err)
	assert.Equal(t, logEntry.Actor, fetched.Actor)
}

func TestMongoLogRepository_GetAll(t *testing.T) {
	storage, cleanup := repository.NewTestMongoStorage(t)
	defer cleanup()

	repo := repository.NewMongoLogRepository(storage)

	// Insert multiple logs
	logs := []*model.Log{
		{Actor: "tester", Operation: model.OperationCreate, Timestamp: time.Now()},
		{Actor: "tester", Operation: model.OperationUpdate, Timestamp: time.Now().Add(-time.Hour)},
		{Actor: "admin", Operation: model.OperationCreate, Timestamp: time.Now().Add(-2 * time.Hour)},
	}

	for _, l := range logs {
		_, err := repo.Create(context.TODO(), l)
		assert.NoError(t, err)
	}

	// Test filtering by Actor
	query := &model.GetLogsQuery{
		Actor:    "tester",
		Page:     1,
		PageSize: 10,
	}
	result, err := repo.GetAll(context.TODO(), query)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Test filtering by Operation
	query.Operation = model.OperationUpdate
	result, err = repo.GetAll(context.TODO(), query)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, model.OperationUpdate, result[0].Operation)
}

func TestMongoLogRepository_Count(t *testing.T) {
	storage, cleanup := repository.NewTestMongoStorage(t)
	defer cleanup()

	repo := repository.NewMongoLogRepository(storage)

	// Insert some logs
	for i := 0; i < 3; i++ {
		log := &model.Log{
			Actor:     "counter",
			Operation: model.OperationCreate,
			Timestamp: time.Now(),
		}
		_, err := repo.Create(context.TODO(), log)
		assert.NoError(t, err)
	}

	count, err := repo.Count(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
