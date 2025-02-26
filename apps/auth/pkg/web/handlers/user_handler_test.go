package handlers

import (
	"strings"

	"errors"

	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	// "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
			name: "Success - Valid User",
			requestBody: models.SignUpReq{
				Email: "test@example.com",
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
			},
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"Invalid Email"}`,
		},
		{
			name: "Fail - Cognito Error",
			requestBody: models.SignUpReq{
				Email: "error@example.com",
			},
			mockReturnErr:  errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"Failed to sign up user"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil

			if tc.mockExpected {
				suite.mockService.On("AdminCreateUser", mock.Anything, tc.requestBody).Return(tc.mockReturnErr)
			}

			formData := url.Values{}
			formData.Set("email", tc.requestBody.Email)
			requestBody := formData.Encode()

			req, _ := http.NewRequest(http.MethodPost, "/auth/createUser", strings.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/createUser", suite.handler.CreateUser)

			r.ServeHTTP(w, req)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)

			assert.JSONEq(suite.T(), tc.expectedBody, w.Body.String())

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
			}, mockReturnErr: errors.New("cognito error"),
			mockExpected:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"Internal server error"}`,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockService.ExpectedCalls = nil // reset the expected calls

			if tc.mockExpected {
				suite.mockService.On("ForgetPassword", mock.Anything, mock.Anything).Return(tc.mockReturnErr)
			}

			formData := url.Values{}
			formData.Set("username", tc.requestBody.Username)
			requestBody := formData.Encode()

			req, _ := http.NewRequest(http.MethodPost, "/auth/forgetPassword", strings.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			r := gin.Default()
			r.POST("/auth/forgetPassword", suite.handler.ForgetPassword)

			r.ServeHTTP(w, req)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)

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
		// {
		// 	name: "Success - Valid request body",
		// 	requestBody: models.LoginReq{
		// 		Username: "testUsername",
		// 		Password: "testPassword",
		// 	},
		// 	mockReturn: &models.LoginRes{
		// 		Challenge: "SOFTWARE_TOKEN_MFA",
		// 		Session:   "test-session",
		// 	},
		// 	mockReturnErr:  nil,
		// 	mockExpected:   true,
		// 	expectedStatus: http.StatusOK,
		// 	expectedBody:   `{"challenge":"SOFTWARE_TOKEN_MFA"}`,
		// },
		{
			name:           "Fail - Invalid request body",
			requestBody:    nil,
			mockReturn:     nil,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"Invalid Form Data"}`,
		},
		// {
		// 	name: "Fail - Cognito Error",
		// 	requestBody: models.LoginReq{
		// 		Username: "testUsername",
		// 		Password: "testPassword",
		// 	},
		// 	mockReturn:     nil,
		// 	mockReturnErr:  errors.New("cognito error"),
		// 	mockExpected:   true,
		// 	expectedStatus: http.StatusInternalServerError,
		// 	expectedBody:   `{"status":"Internal server error"}`,
		// },
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
			name: "Fail - Invalid request body",
			requestBody: nil,
			sessionCookie:  true,
			mockReturn:     nil,
			mockReturnErr:  nil,
			mockExpected:   false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"Invalid Form Data"}`,
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
			expectedBody:   `{"status":"Internal server error"}`,
		},
		{
			name: "Fail - Missing session cookie",
			requestBody: models.SetNewPasswordReq{
				Username:    "testUsername",
				NewPassword: "testPassword",
				Session:     "test-session",
			},
			sessionCookie: false,
			mockReturn:    nil,
			mockReturnErr: nil,
			mockExpected:  false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":"Session cookie missing"}`,
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

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
