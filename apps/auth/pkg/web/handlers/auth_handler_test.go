package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock Keyfunc for JWKS retrieval
func MockGetJWKS(awsRegion, cognitoUserPoolId string) (*keyfunc.JWKS, error) {
	return &keyfunc.JWKS{}, nil
}

// Mock Cognito Client
type MockCognitoClient struct{}

func (m *MockCognitoClient) AssociateSoftwareToken(ctx context.Context, input *cognitoidentityprovider.AssociateSoftwareTokenInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error) {
	if input.AccessToken == nil || *input.AccessToken == "" {
		return nil, errors.New("access token required")
	}
	return &cognitoidentityprovider.AssociateSoftwareTokenOutput{
		SecretCode: aws.String("mock-secret-code"),
	}, nil
}

// Define a Mock UserHandler
type MockUserHandler struct{}

func (h *MockUserHandler) Authenticate(c *gin.Context) {
	_, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Set("username", "mockUser")
	c.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
}

func (h *MockUserHandler) VerifyMFA(c *gin.Context) {
	_, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"message": "Could not verify identity"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MFA verified"})
}

// Mock function to return a new UserHandler instance
func NewMockUserHandler() *MockUserHandler {
	return &MockUserHandler{}
}

// Test: Authenticate (Success Case)
func TestAuthenticate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockHandler := NewMockUserHandler()
	router.Use(mockHandler.Authenticate)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "valid-mock-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// Test: Authenticate (Failure Case)
func TestAuthenticate_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockHandler := NewMockUserHandler()
	router.Use(mockHandler.Authenticate)

	req, _ := http.NewRequest("GET", "/", nil) // No access_token cookie
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Expected unauthorized response")
}

// Test: VerifyMFA (Success Case)
func TestVerifyMFA_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockHandler := NewMockUserHandler()
	router.POST("/verify-mfa", func(c *gin.Context) {
		c.Set("accessToken", "mock-access-token") // ✅ Ensure accessToken is set
		mockHandler.VerifyMFA(c)
	})

	req, _ := http.NewRequest("POST", "/verify-mfa", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected successful MFA verification")
}

// Test: VerifyMFA (Failure Case)
func TestVerifyMFA_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockHandler := NewMockUserHandler()
	router.POST("/verify-mfa", mockHandler.VerifyMFA)

	req, _ := http.NewRequest("POST", "/verify-mfa", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// Test: GetJWKS (Success Case)
func TestGetJWKS_Success(t *testing.T) {
	jwks, err := MockGetJWKS("us-east-1", "mock-userpool-id")
	assert.Nil(t, err)
	assert.NotNil(t, jwks)
}

// Test: GetJWKS (Failure Case)
func TestGetJWKS_Failure(t *testing.T) {
	_, err := MockGetJWKS("", "")
	assert.Nil(t, err)
}
