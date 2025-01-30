package handlers

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services"
	"github.com/owjoel/client-factpack/apps/auth/pkg/errors"
)

type UserHandler struct {
	service *services.UserService
}

// Create new handler object
func New() *UserHandler {
	return &UserHandler{service: services.NewUserService()}
}

//	@Summary		ping
//	@Description	Basic health check 
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	models.StatusRes	"Connection status"
//	@Router			/health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.StatusRes{Status: "Connection successful"})
}

//	@Summary		Create Users
//	@Description	Admin registers user with Cognito user pool via email
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			email	formData		models.SignUpReq	true	"User email"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		500		{object}	models.StatusRes
//	@Router			/auth/createUser [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.SignUpReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Error"})
	}

	// Validate email
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Invalid Email"})
	}
	
	if err := h.service.AdminCreateUser(c.Request.Context(), req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, models.StatusRes{Status: "Error"})
	}
	c.JSON(http.StatusCreated, models.StatusRes{Status: "Success"})
}

//	@Summary		Forget Password
//	@Description	Forget password
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData		models.ForgetPasswordReq	true	"Username"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/forgetPassword [post]
func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req models.ForgetPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Error"})
		return
	}
	
	if err := h.service.ForgetPassword(c.Request.Context(), req); err != nil {
		status, message := errors.CognitoErrorHandler(err)
		fmt.Println(status, message)
		c.JSON(status, models.StatusRes{Status: message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "If you have an account, you will receive an email with instructions on how to reset your password."})
}

//	@Summary		Login
//	@Description	Cognito SSO login using username and password
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData		models.LoginReq	true	"Username, Password"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/login [post]
func (h *UserHandler) UserLogin(c *gin.Context) {
	var req models.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Error"})
		return
	}

	if err := h.service.UserLogin(c.Request.Context(), req); err != nil {
		status, message := errors.CognitoErrorHandler(err)
		fmt.Println(status, message)
		c.JSON(status, models.StatusRes{Status: message})
		return
	}

	// TODO: return some token probs
	c.JSON(http.StatusOK, models.StatusRes{Status: "Success"})
}

//	@Summary		Confirm Forget Password
//	@Description	Submit Cognito OTP sent to user's email to proceed with password reset
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData		models.ConfirmForgetPasswordReq	true	"OTP Code"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/confirmForgetPassword [post]
func (h *UserHandler) ConfirmForgetPassword(c *gin.Context) {
	var req models.ConfirmForgetPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Error"})
		return
	}

	if err := h.service.ConfirmForgetPassword(c.Request.Context(), req); err != nil {
		status, message := errors.CognitoErrorHandler(err)
		fmt.Println(status, message)
		c.JSON(status, models.StatusRes{Status: message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Successfully reset password"})
}