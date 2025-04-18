package handlers

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
)

type ClientHandler struct {
	service service.ClientServiceInterface
}

func NewClientHandler(service service.ClientServiceInterface) *ClientHandler {
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
//	@Router			/:id [get]
func (h *ClientHandler) GetClient(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Missing id"})
		return
	}

	client, err := h.service.GetClient(c.Request.Context(), id)
	if err != nil {
		log.Printf("Failed to retrieve client (ID: %s): %v", id, err)
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
//	@Router			/ [get]
func (h *ClientHandler) GetAllClients(c *gin.Context) {
	query := &model.GetClientsQuery{}

	if err := c.ShouldBindQuery(query); err != nil {
		log.Printf("Failed to bind query: %v", err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Invalid request parameters"})
		return
	}

	total, clients, err := h.service.GetAllClients(c.Request.Context(), query)
	if err != nil {
		log.Printf("Failed to retrieve clients: %v", err)
		c.JSON(http.StatusBadRequest, model.StatusRes{Status: "Could not retrieve clients"})
		return
	}

	resp(c, http.StatusOK, model.GetClientsResponse{
		Total: total,
		Data:  clients,
	})
}

// CreateClientByName submits a job to prefect to create a client profile
//
//	@Summary		Create Client By Name
//	@Description	Create a client profile by name
//	@Tags			clients
//	@Accept			application/json
//	@Produce		json
//	@Param			name	body		model.CreateClientByNameReq	true	"Client name"
//	@Success		200	{object}	handlers.Response{data=model.CreateClientByNameRes}
//	@Failure		400	{object}	handlers.Response
//	@Router			/scrape [post]
func (h *ClientHandler) CreateClientByName(c *gin.Context) {
	req := &model.CreateClientByNameReq{}

	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request"})
		return
	}

	if req.Name == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Missing name"})
		return
	}

	id, err := h.service.CreateClientByName(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Could not create client"})
		return
	}

	resp(c, http.StatusOK, model.JobIDRes{JobID: id})
}

// UpdateClient updates a client profile
//
//	@Summary		Update Client
//	@Description	Update a client profile
//	@Tags			clients
//	@Accept			application/json
//	@Produce		json
//	@Param			id	query		string	true	"Hex id used to identify client"
//	@Param			client	body		model.Client	true "Client data"
//	@Success		200	{object}	handlers.Response
//	@Failure		400	{object}	handlers.Response
//	@Router			/:id [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	clientID := c.Param("id")
	req := &model.UpdateClientReq{}
	if clientID == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Missing id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request"})
		return
	}

	err := h.service.UpdateClient(c.Request.Context(), clientID, req.Changes)
	if err != nil {
		log.Printf("Failed to update client: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Could not update client"})
		return
	}

	resp(c, http.StatusOK, model.StatusRes{Status: "Client updated"})
}

func (h *ClientHandler) RescrapeClient(c *gin.Context) {
	clientID := c.Param("id")
	if clientID == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Missing id"})
		return
	}

	err := h.service.RescrapeClient(c.Request.Context(), clientID)
	if err != nil {
		log.Printf("Failed to rescrape client: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Could not rescrape client"})
		return
	}

	resp(c, http.StatusOK, model.StatusRes{Status: "Client rescraped"})
}

func (h *ClientHandler) MatchClient(c *gin.Context) {
	clientID := c.Param("id")
	if clientID == "" {
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Missing id"})
	}

	var fileBytes []byte
	var fileName string

	formFile, fileErr := c.FormFile("file")
	text := c.PostForm("text")

	if (fileErr == nil && text != "") || (fileErr != nil && text == "") {
		resp(c, http.StatusBadRequest, model.ErrorResponse{
			Message: "Provide either a file or raw text, not both",
		})
		return
	}

	if fileErr == nil {
		file, err := formFile.Open()
		if err != nil {
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Failed to open uploaded file"})
			return
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Failed to read uploaded file"})
			return
		}
		fileName = formFile.Filename
	} else {
		fileBytes = []byte(text)
		fileName = "input.txt"
	}

	req := &model.MatchClientReq{
		FileName:  fileName,
		FileBytes: base64.StdEncoding.EncodeToString(fileBytes),
	}

	id, err := h.service.MatchClient(c.Request.Context(), req, clientID)
	if err != nil {
		log.Printf("Failed to match client: %v", err)
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Could not match client"})
		return
	}

	resp(c, http.StatusOK, model.JobIDRes{JobID: id})
}