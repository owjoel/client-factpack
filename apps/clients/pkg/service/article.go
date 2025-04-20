package service

import (
	"context"
	"fmt"

	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
)

type ArticleService struct {
	articleRepository repository.ArticleRepository
}

type ArticleServiceInterface interface {
	GetAllArticles(ctx context.Context, query *model.GetArticlesReq) (articles []model.Article, err error)
}

func NewArticleService(articleRepository repository.ArticleRepository) *ArticleService {
	return &ArticleService{articleRepository: articleRepository}
}

func (s *ArticleService) GetAllArticles(ctx context.Context, query *model.GetArticlesReq) (articles []model.Article, err error) {
	articles, err = s.articleRepository.GetAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting articles", errorx.ErrInternal)
	}

	return articles, nil
}
