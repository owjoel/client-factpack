package handlers

import (
	"bytes"
	"encoding/json"
	"strings"

	"errors"

	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	// "time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	// "github.com/stretchr/testify/require"
)

type UserHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.UserInterface
	handler     *UserHandler
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.mockService = new(mocks.UserInterface)
	suite.handler = New(suite.mockService)
}

// func TestParseJWT(t *testing.T) {
// 	// Generate a token for testing
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": "testuser",
// 		"exp":      time.Now().Add(time.Hour * 1).Unix(),
// 	})

// 	// Sign it with a known secret
// 	secret := []byte("my-secret")
// 	tokenString, err := token.SignedString(secret)
// 	require.NoError(t, err)

// 	// Test: Valid token with correct keyFunc
// 	validKeyFunc := func(token *jwt.Token) (interface{}, error) {
// 		return secret, nil
// 	}

// 	parsedToken, err := ParseJWT(tokenString, validKeyFunc)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, parsedToken)
// 	assert.True(t, parsedToken.Valid)

// 	// Test: Invalid token (malformed string)
// 	badTokenString := "this.is.not.valid"
// 	_, err = ParseJWT(badTokenString, validKeyFunc)
// 	assert.Error(t, err)

// 	// Test: Valid token but wrong keyFunc (signature verification fails)
// 	invalidKeyFunc := func(token *jwt.Token) (interface{}, error) {
// 		return []byte("wrong-secret"), nil
// 	}

// 	_, err = ParseJWT(tokenString, invalidKeyFunc)
// 	assert.Error(t, err)
// }

func (suite *UserHandlerTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/health", suite.handler.HealthCheck)

	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	assert.JSONEq(suite.T(), `{"status":"Connection successful"}`, w.Body.String())
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	tests := []struct {
		name           string
		requestBody    models.SignUpReq
		mockReturnErr  error
		mockExpected   bool // if we expect mock service method to be called
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Fail - Invalid JSON Binding",
			requestBody:    models.SignUpReq{}, // Doesn't matter, raw body will be broken
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
		{
			name: "Success - Valid User",
			requestBody: models.SignUpReq{
				Email: "test@example.com",
				Role:  "agent",
			},
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"Success"}`,
		},
		{
			name: "Fail - Invalid Email",
			requestBody: models.SignUpReq{
				Email: "invalid-email",
				Role:  "agent",
			},
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.SignUpReq{
				Email: "error@example.com",
				Role:  "agent",
			},
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			if tc.mockExpected {
				suite.mockService.On("AdminCreateUser", mock.Anything, tc.requestBody).Return(tc.mockReturnErr)
			}

			// Marshal request body into JSON
			requestBody, _ := json.Marshal(tc.requestBody)

			req, _ := http.NewRequest(http.MethodPost, "/auth/createUser", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/createUser", suite.handler.CreateUser)

			r.ServeHTTP(w, req)

			// Check status code
			assert.Equal(suite.T(), tc.expectedStatus, w.Code)

			// Check response body
			assert.JSONEq(suite.T(), tc.expectedBody, w.Body.String())

			// Check if mock expectations were met
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestForgetPassword() {
	tests := []struct {
		name           string
		requestBody    models.ForgetPasswordReq
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.ForgetPasswordReq{
				Username: "testUsername",
			},
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"If you have an account, you will receive an email with instructions on how to reset your password."}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.ForgetPasswordReq{
				Username: "testUsername",
			},
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name:           "Fail - Invalid Request Binding", // NEW CASE
			requestBody:    models.ForgetPasswordReq{},       // doesn't matter here
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil // reset the expected calls

			if tc.mockExpected {
				suite.mockService.On("ForgetPassword", mock.Anything, mock.Anything).Return(tc.mockReturnErr)
			}

			var req *http.Request
			if tc.name == "Fail - Invalid Request Binding" {
				// Send broken JSON to trigger ShouldBind error
				req, _ = http.NewRequest(http.MethodPost, "/auth/forgetPassword", strings.NewReader(`bad-json`))
				req.Header.Set("Content-Type", "application/json")
			} else {
				formData := url.Values{}
				formData.Set("username", tc.requestBody.Username)
				requestBody := formData.Encode()
			
				req, _ = http.NewRequest(http.MethodPost, "/auth/forgetPassword", strings.NewReader(requestBody))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}			

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/forgetPassword", suite.handler.ForgetPassword)

			r.ServeHTTP(w, req)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)

			// Optional: Debug output
			// fmt.Println("Response Body:", w.Body.String())

			assert.JSONEq(suite.T(), tc.expectedBody, w.Body.String())

			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserLogin() {
	tests := []struct {
		name           string
		requestBody    any
		mockReturn     *models.LoginRes
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.LoginReq{
				Username: "testUsername",
				Password: "testPassword",
			},
			mockReturn: &models.LoginRes{
				Challenge: "SOFTWARE_TOKEN_MFA",
				Session:   "test-session",
			},
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"challenge":"SOFTWARE_TOKEN_MFA"}`,
		},
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			mockReturn:     nil,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.LoginReq{
				Username: "testUsername",
				Password: "testPassword",
			},
			mockReturn:     nil,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil
			var requestBody string

			if loginReq, ok := tc.requestBody.(models.LoginReq); ok {
				if tc.mockExpected {
					suite.mockService.On("UserLogin", mock.Anything, loginReq).
						Return(tc.mockReturn, tc.mockReturnErr)
				}
				form := url.Values{}
				form.Add("username", loginReq.Username)
				form.Add("password", loginReq.Password)
				requestBody = form.Encode()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(requestBody))

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/login", suite.handler.UserLogin)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserInitialChangePassword() {
	tests := []struct {
		name           string
		requestBody    any
		sessionCookie  bool
		mockReturn     *models.SetNewPasswordRes
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.SetNewPasswordReq{
				Username:    "testUsername",
				NewPassword: "testPassword",
				Session:     "test-session",
			},
			sessionCookie: true,
			mockReturnErr: nil,
			mockReturn: &models.SetNewPasswordRes{
				Challenge: "SOFTWARE_TOKEN_MFA",
				Session:   "test-session",
			},
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"challenge":"SOFTWARE_TOKEN_MFA"}`,
		},
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			sessionCookie:  true,
			mockReturn:     nil,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.SetNewPasswordReq{
				Username:    "testUsername",
				NewPassword: "testPassword",
				Session:     "test-session",
			},
			sessionCookie:  true,
			mockReturn:     nil,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name: "Fail - Missing session cookie",
			requestBody: models.SetNewPasswordReq{
				Username:    "testUsername",
				NewPassword: "testPassword",
				Session:     "test-session",
			},
			sessionCookie:  false,
			mockReturn:     nil,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error_code":"AUTH_UNAUTHORIZED","message":"Unauthorized access"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			var requestBody string
			if setNewPasswordReq, ok := tc.requestBody.(models.SetNewPasswordReq); ok {
				if tc.mockExpected {
					suite.mockService.On("SetNewPassword", mock.Anything, setNewPasswordReq).
						Return(tc.mockReturn, tc.mockReturnErr)
				}

				form := url.Values{}
				form.Add("username", setNewPasswordReq.Username)
				form.Add("newPassword", setNewPasswordReq.NewPassword)
				requestBody = form.Encode()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/changePassword", strings.NewReader(requestBody))

			if tc.sessionCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: "test-session",
				})
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/changePassword", suite.handler.UserInitialChangePassword)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserSetupMFA() {
	tests := []struct {
		name           string
		mockReturn     *models.AssociateTokenRes
		sessionCookie  bool
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			mockReturn: &models.AssociateTokenRes{
				Token:   "test-token",
				Session: "test-session",
			},
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"test-token"}`,
		},
		{
			name:           "Fail - Cognito Error",
			mockReturn:     nil,
			sessionCookie:  true,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name:           "Fail - Missing session cookie",
			mockReturn:     nil,
			sessionCookie:  false,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error_code":"AUTH_UNAUTHORIZED","message":"Unauthorized access"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			req := httptest.NewRequest(http.MethodGet, "/auth/setupMFA", nil)

			if tc.sessionCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: "test-session",
				})
			}

			if tc.mockExpected {
				suite.mockService.On("SetupMFA", mock.Anything, mock.Anything).Return(tc.mockReturn, tc.mockReturnErr)
			}

			w := httptest.NewRecorder()
			r := gin.Default()
			r.GET("/auth/setupMFA", suite.handler.UserSetupMFA)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserVerifyMFA() {
	tests := []struct {
		name           string
		requestBody    any
		sessionCookie  bool
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Success"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			sessionCookie:  true,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name: "Fail - Missing session cookie",
			requestBody: models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			sessionCookie:  false,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error_code":"AUTH_UNAUTHORIZED","message":"Unauthorized access"}`,
		},
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			var requestBody []byte
			if verifyMFAReq, ok := tc.requestBody.(models.VerifyMFAReq); ok {
				if tc.mockExpected {
					suite.mockService.On("VerifyMFA", mock.Anything, verifyMFAReq).Return(tc.mockReturnErr)
				}
				requestBody, _ = json.Marshal(verifyMFAReq)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/verifyMFA", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			if tc.sessionCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: "test-session",
				})
			}

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/verifyMFA", suite.handler.UserVerifyMFA)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserLoginMFA() {
	mockReturn := &models.AuthenticationRes{
		Result: types.AuthenticationResultType{
			AccessToken:  aws.String("test-access-token"),
			IdToken:      aws.String("test-id-token"),
			RefreshToken: aws.String("test-refresh-token"),
			TokenType:    aws.String("Bearer"),
			ExpiresIn:    3600,
			NewDeviceMetadata: &types.NewDeviceMetadataType{
				DeviceGroupKey: aws.String("test-device-group-key"),
				DeviceKey:      aws.String("test-device-key"),
			},
		},
		Challenge: "",
	}
	tests := []struct {
		name           string
		requestBody    any
		mockReturn     *models.AuthenticationRes
		sessionCookie  bool
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.SignInMFAReq{
				Username: "testUsername",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockReturn:     mockReturn,
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Login Successful"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.SignInMFAReq{
				Username: "testUsername",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockReturn:     mockReturn,
			sessionCookie:  true,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name: "Fail - Missing session cookie",
			requestBody: models.SignInMFAReq{
				Username: "testUsername",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockReturn:     mockReturn,
			sessionCookie:  false,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error_code":"AUTH_UNAUTHORIZED","message":"Unauthorized access"}`,
		},
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			sessionCookie:  true,
			mockReturn:     mockReturn,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
		{
			name: "Fail - Challenge not empty",
			requestBody: models.SignInMFAReq{
				Username: "testUsername",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockReturn: &models.AuthenticationRes{
				Challenge: "SMS_MFA", // Challenge triggers the red block
			},
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusForbidden, // Assuming errors.ErrUnauthorized returns 403
			expectedBody:   `{"error_code":"AUTH_UNAUTHORIZED","message":"Unauthorized access"}`,
		},		
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			var requestBody string
			if signInMFAReq, ok := tc.requestBody.(models.SignInMFAReq); ok {
				if tc.mockExpected {
					suite.mockService.On("SignInMFA", mock.Anything, signInMFAReq).Return(*tc.mockReturn, tc.mockReturnErr)
				}
				form := url.Values{}
				form.Add("username", signInMFAReq.Username)
				form.Add("code", signInMFAReq.Code)
				requestBody = form.Encode()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/loginMFA", strings.NewReader(requestBody))

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tc.sessionCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: "test-session",
				})
			}

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/loginMFA", suite.handler.UserLoginMFA)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestConfirmForgetPassword() {
	tests := []struct {
		name           string
		requestBody    any
		sessionCookie  bool
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Valid request body",
			requestBody: models.ConfirmForgetPasswordReq{
				Username:    "testUsername",
				Code:        "test-code",
				NewPassword: "test-new-password",
			},
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Successfully reset password"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.ConfirmForgetPasswordReq{
				Username:    "testUsername",
				Code:        "test-code",
				NewPassword: "test-new-password",
			},
			sessionCookie:  true,
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_code":"SERVER_ERROR","message":"Internal server error"}`,
		},
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			sessionCookie:  true,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"INPUT_INVALID","message":"Invalid input provided"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			var requestBody string
			if confirmForgetPasswordReq, ok := tc.requestBody.(models.ConfirmForgetPasswordReq); ok {
				if tc.mockExpected {
					suite.mockService.On("ConfirmForgetPassword", mock.Anything, confirmForgetPasswordReq).Return(tc.mockReturnErr)
				}
				form := url.Values{}
				form.Add("username", confirmForgetPasswordReq.Username)
				form.Add("code", confirmForgetPasswordReq.Code)
				form.Add("newPassword", confirmForgetPasswordReq.NewPassword)
				requestBody = form.Encode()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/confirmForgetPassword", strings.NewReader(requestBody))

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tc.sessionCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: "test-session",
				})
			}

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/confirmForgetPassword", suite.handler.ConfirmForgetPassword)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserHandlerTestSuite) TestUserLogout() {
	tests := []struct {
		name           string
		mockReturnErr  error
		mockExpected   bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success - Valid request body",
			mockReturnErr:  nil,
			mockExpected:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Logout successful"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/logout", suite.handler.UserLogout)

			r.ServeHTTP(w, req)

			suite.Equal(tc.expectedStatus, w.Code)
			suite.JSONEq(tc.expectedBody, w.Body.String())
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
