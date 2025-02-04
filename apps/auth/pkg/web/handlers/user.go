package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/errors"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services"
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
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Invalid Form Data"})
		return
	}


	res, err := h.service.UserLogin(c.Request.Context(), req)
	if err != nil {
		status, message := errors.CognitoErrorHandler(err)
		fmt.Println(status, message)
		c.JSON(status, models.StatusRes{Status: message})
		return
	}

	// TODO: return some token probs
	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.AuthChallengeRes{Challenge: res.Challenge})
}


func (h *UserHandler) UserInitialChangePassword(c *gin.Context) {
	var req models.SetNewPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Invalid Form Data"})
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		log.Println("Missing session cookie for auth challenge")
		c.JSON(http.StatusUnauthorized, models.StatusRes{Status: "Session cookie missing"})
		return
	}
	req.Session = session

	res, err := h.service.SetNewPassword(c.Request.Context(), req)
	if err != nil {
		log.Printf("%v", fmt.Errorf("error while changing user password: %w", err))
	}
	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.AuthChallengeRes{Challenge: res.Challenge})
}


func (h *UserHandler) UserSetupMFA(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil {
		log.Println("Missing session cookie for auth challenge")
		c.JSON(http.StatusUnauthorized, models.StatusRes{Status: "Session cookie missing"})
		return
	}

	res, err := h.service.SetupMFA(c.Request.Context(), session)
	if err != nil {
		log.Printf("%v", fmt.Errorf("error setting up mfa: %w", err))
		c.JSON(http.StatusInternalServerError, models.StatusRes{Status: "Unable to return token. Please check server logs"})
		return
	}
	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.SetupMFARes{Token: res.Token})
}

func (h *UserHandler) UserVerifyMFA(c *gin.Context) {
	var req models.VerifyMFAReq
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Invalid Form Data"})
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		log.Println("Missing session cookie for auth challenge")
		c.JSON(http.StatusUnauthorized, models.StatusRes{Status: "Session cookie missing"})
		return
	}

	req.Session = session

	if err := h.service.VerifyMFA(c.Request.Context(), req); err != nil {
		log.Printf("%v", fmt.Errorf("error verifying MFA: %w", err))
		c.JSON(http.StatusInternalServerError, models.StatusRes{Status:"Could not verify"})
		return
	}

	c.JSON(http.StatusOK, models.StatusRes{Status: "Success"})
}

func (h *UserHandler) UserLoginMFA(c *gin.Context) {
	var req models.SignInMFAReq
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("%v", fmt.Errorf("error binding request: %w", err))
		c.JSON(http.StatusBadRequest, models.StatusRes{Status: "Invalid Form Data"})
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		log.Println("Missing session cookie for auth challenge")
		c.JSON(http.StatusUnauthorized, models.StatusRes{Status: "Session cookie missing"})
		return
	}
	req.Session = session

	auth, err := h.service.SignInMFA(c.Request.Context(), req)
	if err != nil {
		log.Printf("%v", fmt.Errorf("error verifying totp code: %w", err))
		c.JSON(http.StatusInternalServerError, models.StatusRes{Status: "Unable to return token. Please check server logs"})
		return
	}

	if auth.Challenge != "" {
		c.JSON(http.StatusUnauthorized, models.StatusRes{Status: "Unexpected challenge"})
	}
	c.SetCookie("access_token", *auth.Result.AccessToken, 3600, "/", config.Host, false, true)
	c.SetCookie("id_token", *auth.Result.IdToken, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.StatusRes{Status: "Login Successful"})
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