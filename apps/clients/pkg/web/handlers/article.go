package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type ArticleHandler struct {
	service service.ArticleServiceInterface
}

func NewArticleHandler(service service.ArticleServiceInterface) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	var query model.GetArticlesReq

	if err := c.ShouldBindJSON(&query); err != nil {
		resp(c, http.StatusNotFound, model.ErrorResponse{Message: "Invalid request body"})
		return
	}

	articles, err := h.service.GetAllArticles(c.Request.Context(), &query)
	if err != nil {
		resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve articles"})
		return
	}

	resp(c, http.StatusOK, model.GetArticlesRes{Articles: articles})
}
