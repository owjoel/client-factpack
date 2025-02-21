package auth_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
)

// Mock AuthService
 type MockAuthService struct{}

// Simulating Login
func (m *MockAuthService) Login(req models.LoginReq) (models.LoginRes, error) {
	if req.Username == "testuser" && req.Password == "correctpassword" {
		return models.LoginRes{Challenge: "", Session: "valid-session"}, nil
	}
	if req.Username == "unknownuser" {
		return models.LoginRes{}, errors.New("user not found")
	}
	return models.LoginRes{}, errors.New("invalid credentials")
}

// Simulating Password Hashing (Mock Function)
func (m *MockAuthService) HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password too weak")
	}
	return "hashedpassword123", nil
}

// Simulating Token Validation
func (m *MockAuthService) ValidateToken(token string) error {
	if token == "valid-token" {
		return nil
	}
	return errors.New("invalid token")
}

// Unit Test: Login
func TestLogin(t *testing.T) {
	service := &MockAuthService{}

	t.Run("Successful Login", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "correctpassword"}
		res, err := service.Login(req)
		assert.Nil(t, err)
		assert.Equal(t, "valid-session", res.Session)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "wrongpassword"}
		_, err := service.Login(req)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Empty Username or Password", func(t *testing.T) {
		req := models.LoginReq{Username: "", Password: ""}
		_, err := service.Login(req)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid credentials", err.Error()) // Or modify based on expected behavior
	})	
}

// Unit Test: Password Hashing
func TestHashPassword(t *testing.T) {
	service := &MockAuthService{}

	t.Run("Strong Password", func(t *testing.T) {
		_, err := service.HashPassword("StrongPass123")
		assert.Nil(t, err)
	})

	t.Run("Weak Password", func(t *testing.T) {
		_, err := service.HashPassword("123")
		assert.NotNil(t, err)
		assert.Equal(t, "password too weak", err.Error())
	})

	t.Run("Consistent Hashing Output", func(t *testing.T) {
		hash1, err1 := service.HashPassword("StrongPass123")
		hash2, err2 := service.HashPassword("StrongPass123")
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, hash1, hash2) // Ensure hashing is deterministic for same input
	})	
}

// Unit Test: Token Validation
func TestValidateToken(t *testing.T) {
	service := &MockAuthService{}

	t.Run("Valid Token", func(t *testing.T) {
		err := service.ValidateToken("valid-token")
		assert.Nil(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		err := service.ValidateToken("invalid-token")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("Empty Token", func(t *testing.T) {
		err := service.ValidateToken("")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})	
}
