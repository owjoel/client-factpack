package handlers_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/clients/config"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
	"github.com/stretchr/testify/assert"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

var (
	ConfigOverride = func(region, pool, client string) {
		config.AwsRegion = region
		config.UserPoolID = pool
		config.ClientID = client
	}
)

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey = &privateKey.PublicKey
}

func startMockJWKS(t *testing.T) (*httptest.Server, *keyfunc.JWKS) {
	// Generate a dummy JWKS server
	keyID := "test-kid"
	_, err := x509.MarshalPKIXPublicKey(publicKey)
	assert.NoError(t, err)

	jwk := map[string]interface{}{
		"kty": "RSA",
		"alg": "RS256",
		"use": "sig",
		"kid": keyID,
		"n":   base64url(publicKey.N.Bytes()),
		"e":   base64url(big.NewInt(int64(publicKey.E)).Bytes()),
	}

	jwksPayload := map[string]interface{}{
		"keys": []interface{}{jwk},
	}
	jwksBytes, _ := json.Marshal(jwksPayload)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(jwksBytes)
	}))

	jwks, err := keyfunc.Get(server.URL, keyfunc.Options{})
	assert.NoError(t, err)

	return server, jwks
}

func createTestToken(kid string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, _ := token.SignedString(privateKey)
	return signed
}

func setupTestRouter(jwks *keyfunc.JWKS) *gin.Engine {
	r := gin.New()

	// Pass the mock GetJWKS function directly
	mockGetJWKS := func(_, _ string) (*keyfunc.JWKS, error) {
		return jwks, nil
	}

	r.Use(handlers.Authenticate(mockGetJWKS))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return r
}

func TestAuthenticate_Success(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "access",
		"client_id": "mock-client",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	// Correctly override config values used in Authenticate()
	ConfigOverride("mock-region", "mock-pool", "mock-client")

	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthenticate_MissingToken(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil) // no cookie
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidClientID(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "access",
		"client_id": "wrong-client-id",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(-1 * time.Hour).Unix(), // expired
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "access",
		"client_id": "mock-client",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_MissingTokenUse(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"client_id": "mock-client",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_MissingClientIDInAccessToken(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "access",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_IDTokenWithMissingAudience(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "id",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidGroupClaimFormat(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	claims := jwt.MapClaims{
		"exp":            time.Now().Add(1 * time.Hour).Unix(),
		"iss":            fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use":      "access",
		"client_id":      "mock-client",
		"username":       "testuser",
		"sub":            "testuser",
		"cognito:groups": "admin", // not array
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Define a custom claims struct that implements jwt.Claims
type CustomClaims struct {
	NotMap string
}

// Implement the jwt.Claims interface
func (c CustomClaims) Valid() error {
	return nil
}

// Implement the missing GetAudience method
func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// Implement the missing GetExpirationTime method
func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil
}

// Implement the missing GetIssuedAt method
func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

// Implement the missing GetIssuer method
func (c CustomClaims) GetIssuer() (string, error) {
	return "", nil
}

// Implement the missing GetNotBefore method
func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// Implement the missing GetSubject method
func (c CustomClaims) GetSubject() (string, error) {
	return "", nil
}

func TestAuthenticate_InvalidClaimsCast(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	// Create a token with custom Claims type to simulate cast failure
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, CustomClaims{
		NotMap: "oops",
	})
	token.Header["kid"] = "test-kid"
	signed, _ := token.SignedString(privateKey)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: signed})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_MissingExpClaim(t *testing.T) {
	jwksServer, jwks := startMockJWKS(t)
	defer jwksServer.Close()

	// Leave out "exp" claim
	claims := jwt.MapClaims{
		"iss":       fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", "mock-region", "mock-pool"),
		"token_use": "access",
		"client_id": "mock-client",
		"username":  "testuser",
		"sub":       "testuser",
	}
	token := createTestToken("test-kid", claims)

	ConfigOverride("mock-region", "mock-pool", "mock-client")
	r := setupTestRouter(jwks)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Helper: convert bytes to base64 URL (no padding)
func base64url(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
