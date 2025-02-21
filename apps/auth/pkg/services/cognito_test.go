package services_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
)

// Mock Cognito Service
 type MockCognitoService struct{}

// Simulating Cognito User Creation
func (m *MockCognitoService) CreateUser(req models.SignUpReq) error {
	if req.Email == "existing@gmail.com" {
		return errors.New("user already exists")
	}
	return nil
}

// Simulating Cognito Login
func (m *MockCognitoService) ValidateLogin(req models.LoginReq) (models.LoginRes, error) {
	if req.Username == "testuser" && req.Password == "correctpassword" {
		return models.LoginRes{Session: "valid-session"}, nil
	}
	return models.LoginRes{}, errors.New("invalid login")
}

// Simulating Cognito Token Generation
func (m *MockCognitoService) GenerateToken(username string) (string, error) {
	if username == "testuser" {
		return "valid-token", nil
	}
	return "", errors.New("token generation failed")
}

// Unit Test: User Creation
func TestCreateUser(t *testing.T) {
	service := &MockCognitoService{}

	t.Run("Successful User Creation", func(t *testing.T) {
		req := models.SignUpReq{Email: "newuser@gmail.com"}
		err := service.CreateUser(req)
		assert.Nil(t, err)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		req := models.SignUpReq{Email: "existing@gmail.com"}
		err := service.CreateUser(req)
		assert.NotNil(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})

	// t.Run("Invalid Email Format", func(t *testing.T) {
	// 	req := models.SignUpReq{Email: "invalid-email"}
	// 	err := service.CreateUser(req)
	// 	assert.NotNil(t, err)
	// })
	
}

// Unit Test: Login Validation
func TestValidateLogin(t *testing.T) {
	service := &MockCognitoService{}

	t.Run("Valid Login", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "correctpassword"}
		res, err := service.ValidateLogin(req)
		assert.Nil(t, err)
		assert.Equal(t, "valid-session", res.Session)
	})

	t.Run("Invalid Login", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "wrongpassword"}
		_, err := service.ValidateLogin(req)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid login", err.Error())
	})

	t.Run("Empty Username and Password", func(t *testing.T) {
		req := models.LoginReq{Username: "", Password: ""}
		_, err := service.ValidateLogin(req)
		assert.NotNil(t, err)
	})
	
}

// Unit Test: Token Generation
func TestGenerateToken(t *testing.T) {
	service := &MockCognitoService{}

	t.Run("Valid Token Generation", func(t *testing.T) {
		_, err := service.GenerateToken("testuser")
		assert.Nil(t, err)
	})

	t.Run("Invalid Token Generation", func(t *testing.T) {
		_, err := service.GenerateToken("unknownuser")
		assert.NotNil(t, err)
		assert.Equal(t, "token generation failed", err.Error())
	})

	t.Run("Empty Username for Token", func(t *testing.T) {
		_, err := service.GenerateToken("")
		assert.NotNil(t, err)
		assert.Equal(t, "token generation failed", err.Error())
	})
	
}

