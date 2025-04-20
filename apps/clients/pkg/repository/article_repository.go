package repository

import (
	"context"
	"fmt"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoArticleRepository struct {
	articleCollection *mongo.Collection
}

func NewMongoArticleRepository(storage *MongoStorage) ArticleRepository {
	return &mongoArticleRepository{articleCollection: storage.articleCollection}
}

type ArticleRepository interface {
	GetAll(ctx context.Context, query *model.GetArticlesReq) ([]model.Article, error)
}

func (r *mongoArticleRepository) GetAll(ctx context.Context, query *model.GetArticlesReq) ([]model.Article, error) {
	var objIDs []bson.ObjectID

	for _, id := range query.ID {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid object ID '%s'", errorx.ErrInvalidInput, id)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	cursor, err := r.articleCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to find articles", err)
	}

	defer cursor.Close(ctx)

	var articles []model.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("%w: failed to decode articles", err)
	}

	return articles, nil
}
