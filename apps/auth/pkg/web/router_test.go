package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

)

type MockHandler struct{}

func (m *MockHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func isValidEmail(email string) bool {
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(emailRegex)
    return re.MatchString(email)
}

func (m *MockHandler) CreateUser(c *gin.Context) {
    var req struct {
        Email string `json:"email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || !isValidEmail(req.Email) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (m *MockHandler) ForgetPassword(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}	
    c.JSON(http.StatusOK, gin.H{"message": "Password reset initiated"})
}

func (m *MockHandler) ConfirmForgetPassword(c *gin.Context) {
    var req struct {
        Username    string `json:"username"`
        Code        string `json:"code"`
        NewPassword string `json:"newPassword"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Code == "" || req.NewPassword == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Password reset confirmed"})
}

func (m *MockHandler) UserLogin(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Password == "wrongpassword" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (m *MockHandler) UserInitialChangePassword(c *gin.Context) {
    var req struct {
        Username    string `json:"username"`
        OldPassword string `json:"oldPassword"`
        NewPassword string `json:"newPassword"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "wrongpassword" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (m *MockHandler) UserSetupMFA(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "MFA setup initiated"})
}

func (m *MockHandler) UserVerifyMFA(c *gin.Context) {
    var req struct {
        Code string `json:"code"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Code == "wrongcode" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid MFA code"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "MFA verified"})
}

func (m *MockHandler) UserLoginMFA(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Code     string `json:"code"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.Code == "wrongcode" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid MFA code"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "MFA login successful"})
}

func (m *MockHandler) Authenticate(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
}

func (m *MockHandler) UserLogout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func NewMockRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    handler := &MockHandler{}

    v1API := router.Group("/api/v1")
    v1API.GET("/health", handler.HealthCheck)
    {
        auth := v1API.Group("/auth")
        auth.POST("/createUser", handler.CreateUser)
        auth.POST("/forgetPassword", handler.ForgetPassword)
        auth.POST("/confirmForgetPassword", handler.ConfirmForgetPassword)
        auth.POST("/login", handler.UserLogin)
        auth.POST("/changePassword", handler.UserInitialChangePassword)
        auth.GET("/setupMFA", handler.UserSetupMFA)
        auth.POST("/verifyMFA", handler.UserVerifyMFA)
        auth.POST("/loginMFA", handler.UserLoginMFA)
        auth.GET("/checkUser", handler.Authenticate)
        auth.POST("/logout", handler.UserLogout)
    }

    return router
}

// Test Expected Status Codes for All Routes (Success & Failure)
func TestRoutes(t *testing.T) {
	router := NewMockRouter()

	tests := []struct {
		method string
		url    string
		body   string
		status int
		name   string
	}{
		// Success cases
		{"GET", "/api/v1/health", "", http.StatusOK, "HealthCheck_Success"},
		{"POST", "/api/v1/auth/createUser", `{"email": "user@example.com", "password": "password123"}`, http.StatusCreated, "CreateUser_Success"},
		{"POST", "/api/v1/auth/forgetPassword", `{"username": "testuser"}`, http.StatusOK, "ForgetPassword_Success"},
		{"POST", "/api/v1/auth/confirmForgetPassword", `{"username": "testuser", "code": "123456", "newPassword": "newpassword123"}`, http.StatusOK, "ConfirmForgetPassword_Success"},
		{"POST", "/api/v1/auth/login", `{"username": "testuser", "password": "correctpassword"}`, http.StatusOK, "UserLogin_Success"},
		{"POST", "/api/v1/auth/changePassword", `{"username": "testuser", "oldPassword": "oldpassword", "newPassword": "newpassword123"}`, http.StatusOK, "ChangePassword_Success"},
		{"GET", "/api/v1/auth/setupMFA", "", http.StatusOK, "SetupMFA_Success"},
		{"POST", "/api/v1/auth/verifyMFA", `{"code": "123456"}`, http.StatusOK, "VerifyMFA_Success"},
		{"POST", "/api/v1/auth/loginMFA", `{"username": "testuser", "code": "123456"}`, http.StatusOK, "LoginMFA_Success"},
		{"GET", "/api/v1/auth/checkUser", "", http.StatusOK, "CheckUser_Success"},
		{"POST", "/api/v1/auth/logout", "", http.StatusOK, "UserLogout_Success"},
		
		// Failure cases
		{"GET", "/api/v1/invalid_health", "", http.StatusNotFound, "HealthCheck_Failure"},
		{"POST", "/api/v1/auth/createUser", `{"email": "invalid", "password": "password123"}`, http.StatusBadRequest, "CreateUser_Failure"},
		{"POST", "/api/v1/auth/forgetPassword", `{}`, http.StatusBadRequest, "ForgetPassword_Failure"},
		{"POST", "/api/v1/auth/confirmForgetPassword", `{"username": "", "code": "", "newPassword": ""}`, http.StatusBadRequest, "ConfirmForgetPassword_Failure"},
		{"POST", "/api/v1/auth/login", `{"username": "testuser", "password": "wrongpassword"}`, http.StatusUnauthorized, "UserLogin_Failure"},
		{"POST", "/api/v1/auth/changePassword", `{"username": "testuser", "oldPassword": "wrongpassword", "newPassword": "newpassword123"}`, http.StatusUnauthorized, "ChangePassword_Failure"},
		{"GET", "/api/v1/auth/setupMFA", "", http.StatusUnauthorized, "SetupMFA_Failure"},
		{"POST", "/api/v1/auth/verifyMFA", `{"code": "wrongcode"}`, http.StatusUnauthorized, "VerifyMFA_Failure"},
		{"POST", "/api/v1/auth/loginMFA", `{"username": "testuser", "code": "wrongcode"}`, http.StatusUnauthorized, "LoginMFA_Failure"},
		{"GET", "/api/v1/auth/checkUser", "", http.StatusUnauthorized, "CheckUser_Failure"},
		{"POST", "/api/v1/auth/logout", "", http.StatusUnauthorized, "UserLogout_Failure"},
	}	

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
            req.Header.Set("Content-Type", "application/json")

			if tt.name == "CheckUser_Success" || tt.name == "SetupMFA_Success" || tt.name == "UserLogout_Success" {
				req.Header.Set("Authorization", "Bearer mock-token")
			} else if tt.name == "SetupMFA_Failure" || tt.name == "UserLogout_Failure" {
				req.Header.Del("Authorization") // Ensure header is removed
			}
			
            rec := httptest.NewRecorder()
            router.ServeHTTP(rec, req)

            assert.Equal(t, tt.status, rec.Code)
        })
    }
}
