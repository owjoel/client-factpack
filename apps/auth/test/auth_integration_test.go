package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services/mocks"
	"github.com/owjoel/client-factpack/apps/auth/pkg/web/handlers"
	cipTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

)

func setupRouter(mockService *mocks.UserInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := handlers.New(mockService)

	api := router.Group("/api/v1")
	authGroup := api.Group("/auth")

	authGroup.POST("/createUser", handler.CreateUser)
	authGroup.POST("/forgetPassword", handler.ForgetPassword)
	authGroup.POST("/confirmForgetPassword", handler.ConfirmForgetPassword)
	authGroup.POST("/login", handler.UserLogin)
	authGroup.POST("/logout", handler.UserLogout)
	authGroup.POST("/changePassword", handler.UserInitialChangePassword)
	authGroup.GET("/setupMFA", handler.UserSetupMFA)
	authGroup.POST("/verifyMFA", handler.UserVerifyMFA)
	authGroup.POST("/loginMFA", handler.UserLoginMFA)
	authGroup.GET("/username", handler.GetUsername)
	authGroup.GET("/role", handler.GetUserRole)

	return router
}

func TestCreateUserIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	signUpReq := models.SignUpReq{
		Email: "newuser@example.com",
		Role:  "user",
	}
	mockService.On("AdminCreateUser", mock.Anything, signUpReq).Return(nil)

	router := setupRouter(mockService)

	body, _ := json.Marshal(signUpReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/createUser", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestCreateUserIntegration_Failed(t *testing.T) {
	mockService := new(mocks.UserInterface)
	reqData := models.SignUpReq{
		Email: "fail@example.com",
		Role:  "user",
	}
	mockService.On("AdminCreateUser", mock.Anything, reqData).Return(assert.AnError)

	router := setupRouter(mockService)

	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/createUser", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestForgetPasswordIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	reqData := models.ForgetPasswordReq{Username: "forgotuser"}
	mockService.On("ForgetPassword", mock.Anything, reqData).Return(nil)

	router := setupRouter(mockService)

	form := url.Values{}
	form.Set("username", reqData.Username)

	req, _ := http.NewRequest("POST", "/api/v1/auth/forgetPassword", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestConfirmForgetPasswordIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	reqData := models.ConfirmForgetPasswordReq{
		Username:    "forgotuser",
		Code:        "ABC123",
		NewPassword: "NewSecure456!",
	}
	mockService.On("ConfirmForgetPassword", mock.Anything, reqData).Return(nil)

	router := setupRouter(mockService)

	form := url.Values{}
	form.Set("username", reqData.Username)
	form.Set("code", reqData.Code)
	form.Set("newPassword", reqData.NewPassword)

	req, _ := http.NewRequest("POST", "/api/v1/auth/confirmForgetPassword", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUserLoginIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	loginReq := models.LoginReq{
		Username: "test@example.com",
		Password: "SuperSecure123!",
	}
	mockService.
		On("UserLogin", mock.Anything, loginReq).
		Return(&models.LoginRes{Challenge: "MFA_SETUP", Session: "mocked-session-token"}, nil)

	router := setupRouter(mockService)

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUserLoginIntegration_MissingFields(t *testing.T) {
	router := setupRouter(nil)

	// Missing password
	payload := `{"username": "user@example.com"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogoutIntegration(t *testing.T) {
	router := setupRouter(nil)

	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Logout successful")
}

func TestChangePasswordIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	reqData := models.SetNewPasswordReq{
		Username:    "testuser",
		NewPassword: "newpass123!",
	}
	mockService.On("SetNewPassword", mock.Anything, models.SetNewPasswordReq{
		Username:    "testuser",
		NewPassword: "newpass123!",
		Session:     "mocked-session",
	}).Return(&models.SetNewPasswordRes{
		Challenge: "MFA_SETUP",
		Session:   "new-session",
	}, nil)	
	router := setupRouter(mockService)

	form := url.Values{}
	form.Set("username", reqData.Username)
	form.Set("newPassword", reqData.NewPassword)

	req, _ := http.NewRequest("POST", "/api/v1/auth/changePassword", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "mocked-session"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSetupMFAIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	mockService.On("SetupMFA", mock.Anything, "mocked-session").
		Return(&models.AssociateTokenRes{Token: "otp-secret", Session: "next-session"}, nil)

	router := setupRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/auth/setupMFA", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "mocked-session"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSetupMFAIntegration_NoSession(t *testing.T) {
	router := setupRouter(nil)

	req, _ := http.NewRequest("GET", "/api/v1/auth/setupMFA", nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestVerifyMFAIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	reqData := models.VerifyMFAReq{Code: "123456", Session: "mocked-session"}
	mockService.On("VerifyMFA", mock.Anything, reqData).Return(nil)

	router := setupRouter(mockService)

	form := url.Values{}
	form.Set("code", reqData.Code)

	req, _ := http.NewRequest("POST", "/api/v1/auth/verifyMFA", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "mocked-session"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestVerifyMFAIntegration_MissingCode(t *testing.T) {
	router := setupRouter(nil)

	form := url.Values{}
	req, _ := http.NewRequest("POST", "/api/v1/auth/verifyMFA", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "mocked-session"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginMFAIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)

	reqData := models.SignInMFAReq{
		Username: "testuser",
		Code:     "123456",
		Session:  "mocked-session",
	}
	mockService.On("SignInMFA", mock.Anything, reqData).
	Return(models.AuthenticationRes{
		Challenge: "",
		Result: cipTypes.AuthenticationResultType{
			AccessToken: aws.String("mock-access-token"),
			IdToken:     aws.String("mock-id-token"),
		},
	}, nil)

	router := setupRouter(mockService)

	form := url.Values{}
	form.Set("username", reqData.Username)
	form.Set("code", reqData.Code)

	req, _ := http.NewRequest("POST", "/api/v1/auth/loginMFA", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: reqData.Session})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserRoleIntegration(t *testing.T) {
	mockService := new(mocks.UserInterface)
	mockService.On("GetUserRoleFromToken", "mocked-token").
		Return("admin", nil)

	router := setupRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/auth/role", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "mocked-token"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "admin")
	mockService.AssertExpectations(t)
}

func TestGetUsernameIntegration(t *testing.T) {
	router := setupRouter(nil)

	req, _ := http.NewRequest("GET", "/api/v1/auth/username", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "mocked-token"})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}
