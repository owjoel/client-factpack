package auth_test

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
)

// Mock Service for Account Management
type MockAccountService struct{}

// Simulating successful user signup
func (m *MockAccountService) SignUp(user models.SignUpReq) error {
	if user.Email == "existing@gmail.com" {
		return errors.New("user already exists")
	}
	return nil
}

// Simulating user login
func (m *MockAccountService) Login(req models.LoginReq) (models.LoginRes, error) {
	if req.Username == "testuser" && req.Password == "correctpassword" {
		return models.LoginRes{Challenge: "", Session: "valid-session"}, nil
	}
	return models.LoginRes{}, errors.New("invalid credentials")
}

// Simulating forgotten password flow
func (m *MockAccountService) ForgetPassword(req models.ForgetPasswordReq) error {
	if req.Username == "unknown" {
		return errors.New("user not found")
	}
	return nil
}

// Simulating MFA verification
func (m *MockAccountService) VerifyMFA(req models.VerifyMFAReq) error {
	if req.Code == "123456" {
		return nil
	}
	return errors.New("invalid MFA code")
}

// ✅ **Unit Test: Sign Up**
func TestSignUp(t *testing.T) {
	service := &MockAccountService{}

	t.Run("Successful Sign Up", func(t *testing.T) {
		req := models.SignUpReq{Email: "newuser@gmail.com"}
		err := service.SignUp(req)
		assert.Nil(t, err) // Expecting no error
	})

	t.Run("Sign Up with Existing Email", func(t *testing.T) {
		req := models.SignUpReq{Email: "existing@gmail.com"}
		err := service.SignUp(req)
		assert.NotNil(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})
}

// ✅ **Unit Test: Login**
func TestLogin(t *testing.T) {
	service := &MockAccountService{}

	t.Run("Successful Login", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "correctpassword"}
		res, err := service.Login(req)
		assert.Nil(t, err)
		assert.Equal(t, "valid-session", res.Session)
	})

	t.Run("Invalid Login", func(t *testing.T) {
		req := models.LoginReq{Username: "testuser", Password: "wrongpassword"}
		_, err := service.Login(req)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

// ✅ **Unit Test: Forgot Password**
func TestForgetPassword(t *testing.T) {
	service := &MockAccountService{}

	t.Run("Valid Username", func(t *testing.T) {
		req := models.ForgetPasswordReq{Username: "validuser"}
		err := service.ForgetPassword(req)
		assert.Nil(t, err)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		req := models.ForgetPasswordReq{Username: "unknown"}
		err := service.ForgetPassword(req)
		assert.NotNil(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

// ✅ **Unit Test: Verify MFA**
func TestVerifyMFA(t *testing.T) {
	service := &MockAccountService{}

	t.Run("Correct MFA Code", func(t *testing.T) {
		req := models.VerifyMFAReq{Code: "123456"}
		err := service.VerifyMFA(req)
		assert.Nil(t, err)
	})

	t.Run("Incorrect MFA Code", func(t *testing.T) {
		req := models.VerifyMFAReq{Code: "654321"}
		err := service.VerifyMFA(req)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid MFA code", err.Error())
	})
}
