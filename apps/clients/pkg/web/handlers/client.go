package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type ClientHandler struct {
	service *service.ClientService
}

func New() *ClientHandler {
	return &ClientHandler{service: service.NewClientService()}
}

func (h *ClientHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.StatusRes{Status: "Connection successful"})
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	client := &model.Client{}
	if err := c.ShouldBindJSON(&client); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve client"})
		return
	}
	err :=h.service.CreateClient(c.Request.Context(), client)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve client"})
		return
	}
	resp(c, http.StatusCreated, "Success")
}

// GetClient retrieves the profile of the client by id
// In this case, mongo's object id string
func (h *ClientHandler) GetClient(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Missing id"})
		return
	}
	client, err := h.service.GetClient(c.Request.Context(), id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve client"})
		return
	}
	resp(c, http.StatusOK, client)
}

func (h *ClientHandler) GetAllClients(c *gin.Context) {
	clients, err := h.service.GetAllClients(c.Request.Context())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve clients"})
	}
	resp(c, http.StatusOK, clients)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {

}

func (h *ClientHandler) DeleteClient(c *gin.Context) {

}
