package handlers

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services"
)

type UserHandler struct {
	service *services.UserService
}

// Create new handler object
func New() *UserHandler {
	return &UserHandler{service: services.NewUserService()}
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Connection Successful"})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.SignUpRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error"})
	}

	// Validate email
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid email"})
	}
	
	if err := h.service.CreateUser(c.Request.Context(), req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error"})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "Success"})
}

func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req models.ForgetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error"})
	}
	
	if err := h.service.ForgetPassword(c.Request.Context(), req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "If you have an account, you will receive an email with instructions on how to reset your password."})
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error"})
	}

	if err := h.service.UserLogin(c.Request.Context(), req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error"})
	}

	// TODO: return some token probs
	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}