package repository

import (
	"context"
	"fmt"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobRepository interface {
	Create(ctx context.Context, job *model.Job) (string, error)
	GetOne(ctx context.Context, jobID string) (*model.Job, error)
	GetAll(ctx context.Context, query *model.GetJobsQuery) ([]model.Job, error)
	Count(ctx context.Context, query *model.GetJobsQuery) (int, error)
}

type mongoJobRepository struct {
	jobCollection *mongo.Collection
}

func NewMongoJobRepository(storage *MongoStorage) JobRepository {
	return &mongoJobRepository{jobCollection: storage.jobCollection}
}

func (r *mongoJobRepository) Create(ctx context.Context, job *model.Job) (string, error) {
	res, err := r.jobCollection.InsertOne(ctx, job)
	if err != nil {
		return "", fmt.Errorf("error creating job: %w", err)
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *mongoJobRepository) GetOne(ctx context.Context, jobID string) (*model.Job, error) {
	var job model.Job

	objID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID: %w", err)
	}

	err = r.jobCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("error finding job: %w", err)
	}

	return &job, nil
}

func (r *mongoJobRepository) GetAll(ctx context.Context, query *model.GetJobsQuery) ([]model.Job, error) {
	filter := bson.M{}
	if query.Status != "" {
		filter["status"] = query.Status
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}
	skip := (query.Page - 1) * query.PageSize

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := r.jobCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var jobs []model.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *mongoJobRepository) Count(ctx context.Context, query *model.GetJobsQuery) (int, error) {

	filter := bson.M{}
	if query.Status != "" {
		filter["status"] = query.Status
	}

	count, err := r.jobCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("mongo count error: %w", err)
	}
	return int(count), nil
}
