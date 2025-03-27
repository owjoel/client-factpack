package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"

	// "github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/errors"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services"
	"github.com/owjoel/client-factpack/apps/auth/pkg/utils"
)

// UserHandler represents the handler for user operations.
type UserHandler struct {
	service      services.UserInterface
}

// New creates a new user handler.
func New(service services.UserInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// HealthCheck is a basic health check
//
//	@Summary		ping
//	@Description	Basic health check
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	models.StatusRes	"Connection status"
//	@Router			/health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.StatusRes{Status: "Connection successful"})
}

// CreateUser registers user with Cognito user pool via email and password
//
//	@Summary		Create Users
//	@Description	Admin registers user with Cognito user pool via email. Cognito sends an email with a temporary password to the user.
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			email		formData	string	true	"User's email address"
//	@Success		200			{object}	models.StatusRes
//	@Failure		400			{object}	models.StatusRes
//	@Failure		500			{object}	models.StatusRes
//	@Router			/auth/createUser [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.SignUpReq
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	if err := h.service.AdminCreateUser(c.Request.Context(), req); err != nil {
		log.Printf("Error creating user: %v", err)
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}
	c.JSON(http.StatusCreated, models.StatusRes{Status: "Success"})
}

// ForgetPassword sends a reset password email to the user
//
//	@Summary		Forget Password
//	@Description	Forget password
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.ForgetPasswordReq	true	"Username"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/forgetPassword [post]
func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req models.ForgetPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	if err := h.service.ForgetPassword(c.Request.Context(), req); err != nil {
		utils.ErrorResponse(c, errors.CognitoErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "If you have an account, you will receive an email with instructions on how to reset your password."})
}

// UserLogin handles user login
//
//	@Summary		Login
//	@Description	Cognito SSO login using username and password, returns the next auth challenge, either
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.LoginReq	true	"Username, Password"
//	@Success		200		{object}	models.AuthChallengeRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/login [post]
func (h *UserHandler) UserLogin(c *gin.Context) {
	var req models.LoginReq

	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	res, err := h.service.UserLogin(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, errors.CognitoErrorHandler(err))
		return
	}

	// TODO: return some token probs
	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.AuthChallengeRes{Challenge: res.Challenge})
}

// UserInitialChangePassword handles first-time login password change
//
//	@Summary		Change Password for first-time Login
//	@Description	Users are required to change password on first time login, using their username and password sent via email.
//	@Description	Submit The user's username and new password to respond to this auth challenge.
//	@Description	Request must contain "session" cookie containing the session token to respond to the challenge
//	@Description	On success, responds with next auth challenge, which should be to set up MFA
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.SetNewPasswordReq	true	"Username, New Password"
//	@Success		200		{object}	models.AuthChallengeRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Router			/auth/changePassword [post]
func (h *UserHandler) UserInitialChangePassword(c *gin.Context) {
	var req models.SetNewPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}
	req.Session = session

	res, err := h.service.SetNewPassword(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}

	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.AuthChallengeRes{Challenge: res.Challenge})
}

// UserSetupMFA retrieves an OTP token for MFA setup
//
//	@Summary		Get OTP Token for setting up TOTP authenticator
//	@Description	Submit GET query to cognito to obtain an OTP token.
//	@Description	The user can use this token to set up their authenticator app, either through QR code or by manual keying in of the token.
//	@Description	Request must contain "session" cookie containing the session token to respond to the challenge
//	@Description	On success, the token is returned, and the cookie is updated for the next auth step
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	models.SetupMFARes
//	@Failure		401	{object}	models.StatusRes
//	@Failure		500	{object}	models.StatusRes
//	@Router			/auth/setupMFA [get]
func (h *UserHandler) UserSetupMFA(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	res, err := h.service.SetupMFA(c.Request.Context(), session)
	if err != nil {
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}

	c.SetCookie("session", res.Session, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.SetupMFARes{Token: res.Token})
}

// UserVerifyMFA verifies the MFA code
//
//	@Summary		Verify initial code from authenticator app
//	@Description	User submits the code from their authenticator app to verify the TOTP setup
//	@Description	Request must contain "session" cookie containing the session token to respond to the challenge
//	@Description	On success, the user can proceed to sign in again
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.VerifyMFAReq	true	"TOTP Code"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		500		{object}	models.StatusRes
//	@Router			/auth/verifyMFA [post]
func (h *UserHandler) UserVerifyMFA(c *gin.Context) {
	var req models.VerifyMFAReq

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("@@@@@binding error:", err)
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}
	req.Session = session

	if err := h.service.VerifyMFA(c.Request.Context(), req); err != nil {
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}

	c.JSON(http.StatusOK, models.StatusRes{Status: "Success"})
}

// UserLoginMFA processes MFA login
//
//	@Summary		Submit user TOTP code from authenticator app for all subsequent log ins.
//	@Description	Responds to Cognito auth challenge after successful credential sign in
//	@Description	Request must contain "session" cookie containing the session token to respond to the challenge
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.SignInMFAReq	true	"Username, TOTP Code"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		500		{object}	models.StatusRes
//	@Router			/auth/loginMFA [post]
func (h *UserHandler) UserLoginMFA(c *gin.Context) {
	var req models.SignInMFAReq
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}
	req.Session = session

	auth, err := h.service.SignInMFA(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}

	if auth.Challenge != "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	c.SetCookie("access_token", *auth.Result.AccessToken, 3600, "/", config.Host, false, true)
	c.SetCookie("id_token", *auth.Result.IdToken, 3600, "/", config.Host, false, true)
	c.JSON(http.StatusOK, models.StatusRes{Status: "Login Successful"})
}

// ConfirmForgetPassword verifies the OTP for password reset
//
//	@Summary		Confirm Forget Password
//	@Description	Submit Cognito OTP sent to user's email to proceed with password reset
//	@Tags			auth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			request	formData	models.ConfirmForgetPasswordReq	true	"OTP Code"
//	@Success		200		{object}	models.StatusRes
//	@Failure		400		{object}	models.StatusRes
//	@Failure		401		{object}	models.StatusRes
//	@Failure		403		{object}	models.StatusRes
//	@Failure		404		{object}	models.StatusRes
//	@Router			/auth/confirmForgetPassword [post]
func (h *UserHandler) ConfirmForgetPassword(c *gin.Context) {
	var req models.ConfirmForgetPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidInput)
		return
	}

	if err := h.service.ConfirmForgetPassword(c.Request.Context(), req); err != nil {
		utils.ErrorResponse(c, errors.CognitoErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Successfully reset password"})
}

// UserLogout logs out the user by clearing cookies
// Logout clears the access token and ID token by expiring their cookies
// @Summary Logout User
// @Description Clears the session by expiring the cookies containing the JWT tokens
// @Tags auth
// @Produce json
// @Success 200 {object} models.StatusRes
//
//	@Router			/auth/logout [post]
func (h *UserHandler) UserLogout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", config.Host, false, true)
	c.SetCookie("id_token", "", -1, "/", config.Host, false, true)
	c.SetCookie("session", "", -1, "/", config.Host, false, true)

	c.JSON(http.StatusOK, gin.H{"status": "Logout successful"})
}
