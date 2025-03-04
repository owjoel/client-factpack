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

func New(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// HealthCheck is a basic health check
//
//	@Summary		ping
//	@Description	Basic health check
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	handlers.Response	"Connection status"
//	@Router			/health [get]
func (h *ClientHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.StatusRes{Status: "Connection successful"})
}

// CreateClient creates a new client profile, given the populated json
//
//	@Summary		Create Clients
//	@Description	Create new client profile, given the populated json
//	@Tags			clients
//	@Accept			application/json
//	@Produce		json
//	@Param			client	body		model.Client	true "Client data"
//	@Success		201		{object}	handlers.Response
//	@Failure		400		{object}	handlers.Response
//	@Failure		500		{object}	handlers.Response
//	@Router			/createProfile [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	client := &model.Client{}
	if err := c.ShouldBindJSON(&client); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve client"})
		return
	}
	err := h.service.CreateClient(c.Request.Context(), client)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve client"})
		return
	}
	resp(c, http.StatusCreated, "Success")
}

// GetClient retrieves the profile of the client by id
//
// In this case, mongo's object id string
//
//	@Summary		Get Client By ID
//	@Description	Retrieve client data by profile id
//	@Tags			clients
//	@Produce		json
//	@Param			id	query		string	true	"Hex id used to identify client"
//	@Success		200	{object}	handlers.Response{data=model.Client}
//	@Failure		400	{object}	handlers.Response
//	@Failure		500	{object}	handlers.Response
//	@Router			/retrieveProfile [get]
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

// GetAllClients retrieves all existing client profiles
//
//	@Summary		Get All Clients
//	@Description	Retrieve all client data
//	@Tags			clients
//	@Produce		json
//	@Success		200	{object}	handlers.Response{data=[]model.Client}
//	@Failure		400	{object}	handlers.Response
//	@Failure		500	{object}	handlers.Response
//	@Router			/retrieveAllProfiles [get]
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