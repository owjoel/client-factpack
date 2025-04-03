package repository

import (
	"context"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobRepository interface {
	Create(ctx context.Context, job *model.Job) error
	GetOne(ctx context.Context, jobID string) (*model.Job, error)
	GetAll(ctx context.Context, query *model.GetJobsQuery) ([]model.Job, error)
}

type mongoJobRepository struct {
	jobCollection *mongo.Collection
}

func NewMongoJobRepository(storage *MongoStorage) JobRepository {
	return &mongoJobRepository{jobCollection: storage.jobCollection}
}

func (r *mongoJobRepository) Create(ctx context.Context, job *model.Job) error {
	_, err := r.jobCollection.InsertOne(ctx, job)
	return err
}

func (r *mongoJobRepository) GetOne(ctx context.Context, jobID string) (*model.Job, error) {
	var job model.Job
	err := r.jobCollection.FindOne(ctx, bson.M{"_id": jobID}).Decode(&job)
	return &job, err
}

func (r *mongoJobRepository) GetAll(ctx context.Context, query *model.GetJobsQuery) ([]model.Job, error) {
	filter := bson.M{}
	if query.Status != "" {
		filter["status"] = query.Status
	}

	opts := options.Find().SetSkip(int64(query.Page * query.PageSize)).SetLimit(int64(query.PageSize))
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


