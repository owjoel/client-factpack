package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type ClientHandler struct {
	service *service.ClientService
}

func New(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.StatusRes{Status: "Connection successful"})
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	
}

func (h *ClientHandler) GetClient(c *gin.Context) {

}

func (h *ClientHandler) UpdateClient(c *gin.Context) {

}

func (h *ClientHandler) DeleteClient(c *gin.Context) {

}
