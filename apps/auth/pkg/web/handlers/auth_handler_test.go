package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/config"
)

// // Mock Keyfunc for JWKS retrieval
// func MockGetJWKS(awsRegion, cognitoUserPoolId string) (*keyfunc.JWKS, error) {
// 	return &keyfunc.JWKS{}, nil
// }

// Mock JWKS server
type MockJWKSServer struct {
	jwksJSON string
	err      error
}

func (m *MockJWKSServer) GetJWKS(awsRegion, cognitoUserPoolId string) (jwk.Set, error) {
	if m.err != nil {
		return nil, m.err
	}
	return jwk.Parse([]byte(m.jwksJSON))
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

type MockUserService struct{}

func (m *MockUserService) GetUserRoleFromToken(token string) (string, error) {
	if token == "valid-token" {
		return "admin", nil
	}
	return "", errors.New("invalid token")
}

func (m *MockUserService) AdminCreateUser(ctx context.Context, r models.SignUpReq) error                      { return nil }
func (m *MockUserService) ConfirmForgetPassword(ctx context.Context, r models.ConfirmForgetPasswordReq) error { return nil }
func (m *MockUserService) ForgetPassword(ctx context.Context, r models.ForgetPasswordReq) error               { return nil }
func (m *MockUserService) SetNewPassword(ctx context.Context, r models.SetNewPasswordReq) (*models.SetNewPasswordRes, error) {
	return &models.SetNewPasswordRes{}, nil
}
func (m *MockUserService) SetupMFA(ctx context.Context, session string) (*models.AssociateTokenRes, error) {
	return &models.AssociateTokenRes{}, nil
}
func (m *MockUserService) SignInMFA(ctx context.Context, r models.SignInMFAReq) (models.AuthenticationRes, error) {
	return models.AuthenticationRes{}, nil
}
func (m *MockUserService) UserLogin(ctx context.Context, r models.LoginReq) (*models.LoginRes, error) {
	return &models.LoginRes{}, nil
}
func (m *MockUserService) VerifyMFA(ctx context.Context, r models.VerifyMFAReq) error { return nil }

func NewTestUserHandler() *UserHandler {
	return &UserHandler{
		service: &MockUserService{},
	}
}

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

func generateExpiredJWT() (string, error) {
	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(-1 * time.Hour).Unix(), // Expired an hour ago
		"token_use":  "access",
		"sub":        "user123",
		"username":   "user123",
		"client_id":  config.ClientID, // ensure this matches
		"cognito:groups": []string{"group1"},
	})

	// Use a dummy key for signing
	secret := []byte("my-secret")

	// Generate token string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateJWTWithoutExp() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_use":  "access",
		"sub":        "user123",
		"username":   "user123",
		"client_id":  config.ClientID,
		"cognito:groups": []string{"group1"},
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
}

func generateJWTWithoutSub() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":             time.Now().Add(1 * time.Hour).Unix(),
		"token_use":       "access",
		"username":        "user123",
		"client_id":       config.ClientID,
		"cognito:groups":  []string{"group1"},
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
}

func generateInvalidTokenUseJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(1 * time.Hour).Unix(), // Valid expiration
		"token_use":  "id",                                 // Should be "access", but we give "id"
		"sub":        "user123",
		"username":   "user123",
		"client_id":  config.ClientID,                      // Match expected client ID
		"cognito:groups": []string{"group1"},
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
}

func generateInvalidClientIDJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
		"token_use":  "access",
		"sub":        "user123",
		"username":   "user123",
		"client_id":  "wrong-client-id", // This won't match config.ClientID
		"cognito:groups": []string{"group1"},
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
}

func generateJWTWithWrongAudience() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":             time.Now().Add(1 * time.Hour).Unix(),
		"token_use":       "id",                          // triggers audience check
		"sub":             "user123",
		"username":        "user123",
		"aud":             []string{},                    // empty audience list
		"cognito:groups":  []string{"group1"},
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
}

func generateInvalidCognitoGroupsJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
		"token_use":  "access",
		"sub":        "user123",
		"username":   "user123",
		"client_id":  config.ClientID, // correct client ID
		"cognito:groups": "not-an-array", // INVALID! Should be an array
	})

	secret := []byte("my-secret")

	return token.SignedString(secret)
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

func TestAuthenticate_GetJWKSFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Overwrite the config to simulate an invalid region or pool ID.
	config.AwsRegion = ""        // Invalid region will break the JWKS URL
	config.UserPoolID = "bad-id" // Doesn't matter, region is already invalid

	router.Use(handler.Authenticate)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "some-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAuthenticate_InvalidJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Restore valid config values (or mock values that point to a working endpoint)
	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"
	config.UserPoolID = "ap-southeast-1_APVqivCav"

	router.Use(handler.Authenticate)

	req, _ := http.NewRequest("GET", "/", nil)

	// Provide an invalid JWT to trigger jwt.Parse() failure
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "this.is.not.a.valid.jwt"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_EmptyAccessTokenCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"
	config.UserPoolID = "ap-southeast-1_APVqivCav"

	router.Use(handler.Authenticate)

	req, _ := http.NewRequest("GET", "/", nil)

	req.AddCookie(&http.Cookie{Name: "access_token", Value: ""})

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code) // expect 403!
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Restore valid config values
	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"          // replace with your valid region
	config.UserPoolID = "ap-southeast-1_APVqivCav" // replace with your valid pool id

	router.Use(handler.Authenticate)

	// Provide a **realistic but expired** JWT token
	expiredToken, err := generateExpiredJWT()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: expiredToken})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_MissingExpClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"          // replace with your valid region
	config.UserPoolID = "ap-southeast-1_APVqivCav" // replace with your valid pool id

	router.Use(handler.Authenticate)

	// Generate a token with no "exp"
	tokenWithoutExp, err := generateJWTWithoutExp()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: tokenWithoutExp})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should fail at claims.GetExpirationTime() and return invalid token error
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_MissingSubjectClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"  
	config.UserPoolID = "ap-southeast-1_APVqivCav"

	router.Use(handler.Authenticate)

	jwtWithoutSub, err := generateJWTWithoutSub()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: jwtWithoutSub})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code) // ✅ Expecting 401
}

func TestAuthenticate_InvalidAudienceClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"  
	config.UserPoolID = "ap-southeast-1_APVqivCav"

	router.Use(handler.Authenticate)

	jwtWrongAud, err := generateJWTWithWrongAudience()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: jwtWrongAud})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code) // ✅ Expecting 401 here
}

func TestAuthenticate_InvalidTokenUse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Restore valid config values
	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
	}()

	config.AwsRegion = "ap-southeast-1"  
	config.UserPoolID = "ap-southeast-1_APVqivCav"

	router.Use(handler.Authenticate)

	// Generate token with wrong token_use
	invalidTokenUseJWT, err := generateInvalidTokenUseJWT()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: invalidTokenUseJWT})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 401 Unauthorized because token_use is "id" instead of "access"
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_InvalidClientID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Restore valid config values
	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	origClientID := config.ClientID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
		config.ClientID = origClientID
	}()

	config.AwsRegion = "ap-southeast-1"          // valid region
	config.UserPoolID = "ap-southeast-1_APVqivCav" // valid pool ID
	config.ClientID = "f96a955c-3051-7083-80d4-e5030afca209"  // the one your app expects!

	router.Use(handler.Authenticate)

	invalidClientIDToken, err := generateInvalidClientIDJWT()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: invalidClientIDToken})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_InvalidCognitoGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewTestUserHandler()

	// Restore valid config values
	origRegion := config.AwsRegion
	origPoolID := config.UserPoolID
	origClientID := config.ClientID
	defer func() {
		config.AwsRegion = origRegion
		config.UserPoolID = origPoolID
		config.ClientID = origClientID
	}()

	config.AwsRegion = "ap-southeast-1"          // valid region
	config.UserPoolID = "ap-southeast-1_APVqivCav" // valid pool ID
	config.ClientID = "f96a955c-3051-7083-80d4-e5030afca209"  // the one your app expects!

	router.Use(handler.Authenticate)

	invalidGroupsToken, err := generateInvalidCognitoGroupsJWT()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: invalidGroupsToken})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_VerifyMFA_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()

	router.POST("/verify-mfa", func(c *gin.Context) {
		c.Set("accessToken", "mock-access-token")
		h.VerifyMFA(c)
	})

	req, _ := http.NewRequest("POST", "/verify-mfa", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_VerifyMFA_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()

	router.POST("/verify-mfa", func(c *gin.Context) {
		h.VerifyMFA(c)
	})

	req, _ := http.NewRequest("POST", "/verify-mfa", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error_code":"AUTH_UNAUTHORIZED"`)
	assert.Contains(t, rec.Body.String(), `"message":"Unauthorized access"`)
}

// ---------- TEST: GetJWKS function ----------

func TestGetJWKS_RealCall(t *testing.T) {

	awsRegion := "ap-southeast-1"
	cognitoUserPoolId := "ap-southeast-1_APVqivCav"

	jwks, err := GetJWKS(awsRegion, cognitoUserPoolId)

	assert.NoError(t, err)
	assert.NotNil(t, jwks)
}

func TestGetJWKS_InvalidRegion(t *testing.T) {
	// Passing an invalid region should construct an invalid URL
	_, err := GetJWKS("", "dummy-pool-id")

	assert.Error(t, err)
}

func TestGetUsername_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()
	router.GET("/username", func(c *gin.Context) {
		// Simulate middleware setting the username
		c.Set("username", "mockUser")
		h.GetUsername(c)
	})

	req, _ := http.NewRequest("GET", "/username", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"username":"mockUser"`)
}

func TestGetUsername_Unauthorized(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    h := NewTestUserHandler()
    router.GET("/username", func(c *gin.Context) {
        h.GetUsername(c) // No username set!
    })

    req, _ := http.NewRequest("GET", "/username", nil)
    rec := httptest.NewRecorder()

    router.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusForbidden, rec.Code)
    assert.Contains(t, rec.Body.String(), `"error_code":"AUTH_UNAUTHORIZED"`)
    assert.Contains(t, rec.Body.String(), `"message":"Unauthorized access"`)
}

func TestGetUserRole_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()
	router.GET("/role", func(c *gin.Context) {
		c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "valid-token"})
		h.GetUserRole(c)
	})

	req, _ := http.NewRequest("GET", "/role", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "valid-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"user_role":"admin"`)
}

func TestGetUserRole_Unauthorized_NoCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()
	router.GET("/role", func(c *gin.Context) {
		h.GetUserRole(c)
	})

	req, _ := http.NewRequest("GET", "/role", nil) // No cookie
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
    assert.Contains(t, rec.Body.String(), `"error_code":"AUTH_UNAUTHORIZED"`)
    assert.Contains(t, rec.Body.String(), `"message":"Unauthorized access"`)
}

func TestGetUserRole_Unauthorized_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := NewTestUserHandler()
	router.GET("/role", func(c *gin.Context) {
		reqCookie := &http.Cookie{Name: "access_token", Value: "invalid-token"}
		c.Request.AddCookie(reqCookie)
		h.GetUserRole(c)
	})

	req, _ := http.NewRequest("GET", "/role", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "invalid-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
    assert.Contains(t, rec.Body.String(), `"error_code":"AUTH_UNAUTHORIZED"`)
    assert.Contains(t, rec.Body.String(), `"message":"Unauthorized access"`)
}
