package repository

import (
	"context"
	"errors"
	"fmt"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoLogRepository struct {
	logCollection *mongo.Collection
}

func NewMongoLogRepository(storage *MongoStorage) LogRepository {
	return &mongoLogRepository{logCollection: storage.logCollection}
}

type LogRepository interface {
	Create(ctx context.Context, log *model.Log) (string, error)
	GetAll(ctx context.Context, query *model.GetLogsQuery) ([]model.Log, error)
	GetOne(ctx context.Context, logID string) (*model.Log, error)
	Count(ctx context.Context) (int, error)
}

func (r *mongoLogRepository) Create(ctx context.Context, log *model.Log) (string, error) {
	if log == nil {
		return "", fmt.Errorf("%w: cannot insert nil log", errorx.ErrInvalidInput)
	}

	result, err := r.logCollection.InsertOne(ctx, log)
	if err != nil {
		return "", fmt.Errorf("%w: insert failed", errorx.ErrDependencyFailed)
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("%w: failed to convert inserted ID", errorx.ErrInternal)
	}

	return insertedID.Hex(), nil
}


func (r *mongoLogRepository) GetAll(ctx context.Context, query *model.GetLogsQuery) ([]model.Log, error) {
	filter := bson.M{}

	if query.ClientID != "" {
		objID, err := bson.ObjectIDFromHex(query.ClientID)
		if err != nil {
			return nil, fmt.Errorf("%w: clientID '%s' is not a valid ObjectID", errorx.ErrInvalidInput, query.ClientID)
		}
		filter["clientID"] = objID
	}
	if query.Operation != "" {
		filter["operation"] = query.Operation
	}
	if query.Actor != "" {
		filter["actor"] = query.Actor
	}

	timeFilter := bson.M{}
	if !query.From.IsZero() {
		timeFilter["$gte"] = query.From
	}
	if !query.To.IsZero() {
		timeFilter["$lte"] = query.To
	}
	if len(timeFilter) > 0 {
		filter["timestamp"] = timeFilter
	}

	skip := (query.Page - 1) * query.PageSize
	limit := query.PageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.logCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: mongo find error", errorx.ErrDependencyFailed)
	}

	var logs []model.Log
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("%w: mongo decode error", errorx.ErrInternal)
	}

	return logs, nil
}


func (r *mongoLogRepository) GetOne(ctx context.Context, logID string) (*model.Log, error) {
	objID, err := bson.ObjectIDFromHex(logID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid log ID", errorx.ErrInvalidInput)
	}

	var log model.Log
	err = r.logCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&log)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: log with ID %s not found", errorx.ErrNotFound, logID)
		}
		return nil, fmt.Errorf("%w: failed to query MongoDB", errorx.ErrInternal)
	}

	return &log, nil
}


func (r *mongoLogRepository) Count(ctx context.Context) (int, error) {
	count, err := r.logCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

