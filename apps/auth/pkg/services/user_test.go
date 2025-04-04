package services

import (
	"context"
	"errors"
	"fmt"

	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v4"

	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/services/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockCognitoClient *mocks.CognitoClientRepository
	mockUserService   *UserService
}

func generateTestJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("test-secret"))
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockCognitoClient = new(mocks.CognitoClientRepository)
	suite.mockUserService = &UserService{
		CognitoClient: suite.mockCognitoClient,
	}

	// os.Setenv("USER_POOL_ID", "test-user-pool-id")
	// os.Setenv("CLIENT_ID", "test-client-id")
	// os.Setenv("CLIENT_SECRET", "test-client-secret")
	// os.Setenv("AUTO_RESET_PASSWORD", "true")
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestNewUserService() {

	userService := NewUserService(suite.mockCognitoClient)

	suite.NotNil(userService, "Expected userService to be non-nil")
	suite.Equal(suite.mockCognitoClient, userService.CognitoClient, "Expected CognitoClient to be assigned correctly")
}

func (suite *UserServiceTestSuite) TestCreateUsername() {
	tests := []struct {
		name          string
		email         string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "Valid email returns username",
			email:         "john.doe@example.com",
			expectedValue: "john.doe",
			expectError:   false,
		},
		{
			name:          "Invalid email returns error - empty username",
			email:         "@example.com",
			expectedValue: "",
			expectError:   true,
		},
		{
			name:          "Invalid email returns error - missing @",
			email:         "invalidemail",
			expectedValue: "",
			expectError:   true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			username, err := createUsername(tc.email)

			if tc.expectError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedValue, username)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestAdminCreateUser() {
	tests := []struct {
		name              string
		request           models.SignUpReq
		mockCreateReturn  *cognitoidentityprovider.AdminCreateUserOutput
		mockCreateErr     error
		mockAddReturn     *cognitoidentityprovider.AdminAddUserToGroupOutput
		mockAddErr        error
		expectedErrorText string
	}{
		{
			name: "Success - Valid request",
			request: models.SignUpReq{
				Email: "test@example.com",
				Role:  "agent",
			},
			mockCreateReturn: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &types.UserType{
					Username:       aws.String("test"),
					UserCreateDate: aws.Time(time.Now()),
				},
			},
			mockCreateErr:     nil,
			mockAddReturn:     &cognitoidentityprovider.AdminAddUserToGroupOutput{},
			mockAddErr:        nil,
			expectedErrorText: "",
		},
		{
			name: "Fail - User creation error",
			request: models.SignUpReq{
				Email: "test@example.com",
				Role:  "agent",
			},
			mockCreateReturn:  nil,
			mockCreateErr:     fmt.Errorf("user creation failed"),
			mockAddReturn:     nil,
			mockAddErr:        nil,
			expectedErrorText: "error during sign up: user creation failed",
		},
		{
			name: "Success - With group add failure (logs error but continues)",
			request: models.SignUpReq{
				Email: "test@example.com",
				Role:  "agent",
			},
			mockCreateReturn: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &types.UserType{
					Username:       aws.String("test"),
					UserCreateDate: aws.Time(time.Now()),
				},
			},
			mockCreateErr:     nil,
			mockAddReturn:     nil,
			mockAddErr:        fmt.Errorf("group add failed"),
			expectedErrorText: "",
		},
		{
			name: "Fail - CreateUsername returns empty username",
			request: models.SignUpReq{
				Email: "@example.com", // will produce empty username
				Role:  "agent",
			},
			expectedErrorText: "error creating username: username cannot be empty",
		},
		{
			name: "Fail - Missing user role",
			request: models.SignUpReq{
				Email: "test@example.com",
				Role:  "", // triggers "user role is required"
			},
			expectedErrorText: "user role is required",
		},
		{
			name: "Fail - CreateUsername returns error",
			request: models.SignUpReq{
				Email: "invalid-email", // Simulate an email that causes createUsername to fail
				Role:  "agent",
			},
			mockCreateReturn:  nil,
			mockCreateErr:     nil,
			mockAddReturn:     nil,
			mockAddErr:        nil,
			expectedErrorText: "error creating username: invalid email format",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			// Set up mocks for all cases
			if tc.mockCreateReturn != nil || tc.mockCreateErr != nil {
				suite.mockCognitoClient.On("AdminCreateUser", mock.Anything, mock.Anything).
					Return(tc.mockCreateReturn, tc.mockCreateErr)
			}

			if tc.mockAddReturn != nil || tc.mockAddErr != nil {
				suite.mockCognitoClient.On("AdminAddUserToGroup", mock.Anything, mock.Anything).
					Return(tc.mockAddReturn, tc.mockAddErr)
			}

			err := suite.mockUserService.AdminCreateUser(context.Background(), tc.request)

			if tc.expectedErrorText != "" {
				suite.Error(err)
				suite.Contains(err.Error(), tc.expectedErrorText)
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestForgetPassword() {
	tests := []struct {
		name             string
		request          models.ForgetPasswordReq
		mockForgetReturn *cognitoidentityprovider.ForgotPasswordOutput
		mockForgetErr    error
		expectedError    error
	}{
		{
			name: "Success - Valid request",
			request: models.ForgetPasswordReq{
				Username: "test",
			},
			mockForgetReturn: &cognitoidentityprovider.ForgotPasswordOutput{
				CodeDeliveryDetails: &types.CodeDeliveryDetailsType{
					Destination: aws.String("test@example.com"),
				},
			},
			mockForgetErr: nil,
			expectedError: nil,
		},
		{
			name: "Fail - Forgot password error",
			request: models.ForgetPasswordReq{
				Username: "test",
			},
			mockForgetReturn: nil,
			mockForgetErr:    errors.New("forgot password failed"),
			expectedError:    errors.New("forgot password failed"),
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("ForgotPassword", mock.Anything, mock.Anything).Return(tc.mockForgetReturn, tc.mockForgetErr)

			err := suite.mockUserService.ForgetPassword(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestUserLogin() {
	tests := []struct {
		name            string
		request         models.LoginReq
		mockLoginReturn *cognitoidentityprovider.InitiateAuthOutput
		mockLoginErr    error
		expectedError   error
	}{
		{
			name: "Success - Valid request",
			request: models.LoginReq{
				Username: "test",
				Password: "test",
			},
			mockLoginReturn: &cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					AccessToken:  aws.String("test-access-token"),
					RefreshToken: aws.String("test-refresh-token"),
				},
				ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
				Session:       aws.String("test-session"),
			},
			mockLoginErr: nil,
		},
		{
			name: "Fail - Login error",
			request: models.LoginReq{
				Username: "test",
				Password: "test",
			},
			mockLoginReturn: nil,
			mockLoginErr:    errors.New("login failed"),
			expectedError:   fmt.Errorf("failed to initiate auth: %w", errors.New("login failed")),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("InitiateAuth", mock.Anything, mock.Anything).Return(tc.mockLoginReturn, tc.mockLoginErr)

			_, err := suite.mockUserService.UserLogin(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestSetNewPassword() {
	tests := []struct {
		name                     string
		request                  models.SetNewPasswordReq
		mockSetNewPasswordReturn *cognitoidentityprovider.RespondToAuthChallengeOutput
		mockSetNewPasswordErr    error
		expectedError            error
	}{
		{
			name: "Success - Valid request",
			request: models.SetNewPasswordReq{
				Username:    "test",
				NewPassword: "test-new-password",
				Session:     "test-session",
			},
			mockSetNewPasswordReturn: &cognitoidentityprovider.RespondToAuthChallengeOutput{
				ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
				Session:       aws.String("test-session"),
			},
			mockSetNewPasswordErr: nil,
			expectedError:         nil,
		},
		{
			name: "Fail - Set new password error",
			request: models.SetNewPasswordReq{
				Username:    "test",
				NewPassword: "test-new-password",
				Session:     "test-session",
			},
			mockSetNewPasswordReturn: nil,
			mockSetNewPasswordErr:    errors.New("set new password failed"),
			expectedError:            fmt.Errorf("failed to respond to auth challenge: %w", errors.New("set new password failed")),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("RespondToAuthChallenge", mock.Anything, mock.Anything).Return(tc.mockSetNewPasswordReturn, tc.mockSetNewPasswordErr)

			_, err := suite.mockUserService.SetNewPassword(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestSetupMFA() {
	tests := []struct {
		name                             string
		session                          string
		mockAssociateSoftwareTokenReturn *cognitoidentityprovider.AssociateSoftwareTokenOutput
		mockAssociateSoftwareTokenErr    error
		mockSetupMFAReturn               *models.AssociateTokenRes
		expectedError                    error
	}{
		{
			name:    "Success - Valid request",
			session: "test-session",
			mockAssociateSoftwareTokenReturn: &cognitoidentityprovider.AssociateSoftwareTokenOutput{
				SecretCode: aws.String("test-token"),
				Session:    aws.String("test-session"),
			},
			mockAssociateSoftwareTokenErr: nil,
			mockSetupMFAReturn: &models.AssociateTokenRes{
				Token:   "test-token",
				Session: "test-session",
			},
			expectedError: nil,
		},
		{
			name:                             "Fail - Setup MFA error",
			session:                          "test-session",
			mockAssociateSoftwareTokenReturn: nil,
			mockAssociateSoftwareTokenErr:    errors.New("setup mfa failed"),
			mockSetupMFAReturn: &models.AssociateTokenRes{
				Token:   "test-token",
				Session: "test-session",
			},
			expectedError: fmt.Errorf("error associating token: %w", errors.New("setup mfa failed")),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("AssociateSoftwareToken", mock.Anything, mock.Anything).Return(tc.mockAssociateSoftwareTokenReturn, tc.mockAssociateSoftwareTokenErr)

			res, err := suite.mockUserService.SetupMFA(context.Background(), tc.session)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.mockSetupMFAReturn.Token, res.Token)
				suite.Equal(tc.mockSetupMFAReturn.Session, res.Session)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestGetChallengeName() {
	tests := []struct {
		name          string
		challenge     types.ChallengeNameType
		expectedValue string
	}{
		{
			name:          "NewPasswordRequired challenge",
			challenge:     types.ChallengeNameTypeNewPasswordRequired,
			expectedValue: "NEW_PASSWORD_REQUIRED",
		},
		{
			name:          "MFASetup challenge",
			challenge:     types.ChallengeNameTypeMfaSetup,
			expectedValue: "MFA_SETUP",
		},
		{
			name:          "SoftwareTokenMFA challenge",
			challenge:     types.ChallengeNameTypeSoftwareTokenMfa,
			expectedValue: "SOFTWARE_TOKEN_MFA",
		},
		{
			name:          "Unknown challenge returns empty string",
			challenge:     "UNKNOWN_CHALLENGE",
			expectedValue: "",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			result := getChallengeName(tc.challenge)
			suite.Equal(tc.expectedValue, result)
		})
	}
}

func (suite *UserServiceTestSuite) TestVerifyMFA() {
	tests := []struct {
		name                          string
		request                       models.VerifyMFAReq
		mockVerifySoftwareTokenReturn *cognitoidentityprovider.VerifySoftwareTokenOutput
		mockVerifySoftwareTokenErr    error
		mockVerifyMFAReturn           *models.VerifyMFAReq
		expectedError                 error
	}{
		{
			name: "Success - Valid request",
			request: models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			mockVerifySoftwareTokenReturn: &cognitoidentityprovider.VerifySoftwareTokenOutput{
				Status: "SUCCESS",
			},
			mockVerifySoftwareTokenErr: nil,
			mockVerifyMFAReturn: &models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			expectedError: nil,
		},
		{
			name: "Fail - Verify MFA error",
			request: models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			mockVerifySoftwareTokenReturn: &cognitoidentityprovider.VerifySoftwareTokenOutput{
				Status: "FAILURE",
			},
			mockVerifySoftwareTokenErr: errors.New("verify mfa failed"),
			mockVerifyMFAReturn: &models.VerifyMFAReq{
				Code:    "test-code",
				Session: "test-session",
			},
			expectedError: fmt.Errorf("failed to verify MFA: %w", errors.New("verify mfa failed")),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("VerifySoftwareToken", mock.Anything, mock.Anything).Return(tc.mockVerifySoftwareTokenReturn, tc.mockVerifySoftwareTokenErr)

			err := suite.mockUserService.VerifyMFA(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestSignInMFA() {
	tests := []struct {
		name                string
		request             models.SignInMFAReq
		mockSignInReturn    *cognitoidentityprovider.RespondToAuthChallengeOutput
		mockSignInErr       error
		mockSignInMFAReturn *models.AuthenticationRes
		expectedError       error
	}{
		{
			name: "Success - Valid request",
			request: models.SignInMFAReq{
				Username: "test",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockSignInReturn: &cognitoidentityprovider.RespondToAuthChallengeOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					AccessToken:  aws.String("test-access-token"),
					RefreshToken: aws.String("test-refresh-token"),
				},
				ChallengeName: "",
				Session:       aws.String("test-session"),
			},
			mockSignInErr: nil,
			mockSignInMFAReturn: &models.AuthenticationRes{
				Result: types.AuthenticationResultType{
					AccessToken:  aws.String("test-access-token"),
					RefreshToken: aws.String("test-refresh-token"),
				},
				Challenge: "",
			},
			expectedError: nil,
		},

		{
			name: "Success - Valid request, new password required",
			request: models.SignInMFAReq{
				Username: "test",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockSignInReturn: &cognitoidentityprovider.RespondToAuthChallengeOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					AccessToken:  aws.String("test-access-token"),
					RefreshToken: aws.String("test-refresh-token"),
				},
				ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
				Session:       aws.String("test-session"),
			},
			mockSignInErr: nil,
			mockSignInMFAReturn: &models.AuthenticationRes{
				Result: types.AuthenticationResultType{
					AccessToken:  aws.String("test-access-token"),
					RefreshToken: aws.String("test-refresh-token"),
				},
				Challenge: "NEW_PASSWORD_REQUIRED",
			},
			expectedError: nil,
		},
		{
			name: "Fail - Sign in MFA error",
			request: models.SignInMFAReq{
				Username: "test",
				Code:     "test-code",
				Session:  "test-session",
			},
			mockSignInReturn: nil,
			mockSignInErr:    errors.New("sign in mfa failed"),
			mockSignInMFAReturn: &models.AuthenticationRes{
				Result:    types.AuthenticationResultType{},
				Challenge: "test-challenge",
			},
			expectedError: fmt.Errorf("failed to sign in with MFA: %w", errors.New("sign in mfa failed")),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("RespondToAuthChallenge", mock.Anything, mock.Anything).Return(tc.mockSignInReturn, tc.mockSignInErr)

			res, err := suite.mockUserService.SignInMFA(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.mockSignInMFAReturn.Result, res.Result)
				suite.Equal(tc.mockSignInMFAReturn.Challenge, res.Challenge)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestConfirmForgetPassword() {
	tests := []struct {
		name              string
		request           models.ConfirmForgetPasswordReq
		mockConfirmReturn *cognitoidentityprovider.ConfirmForgotPasswordOutput
		mockConfirmErr    error
		expectedError     error
	}{
		{
			name: "Success - Valid request",
			request: models.ConfirmForgetPasswordReq{
				Username:    "test",
				Code:        "test-code",
				NewPassword: "test-new-password",
			},
			mockConfirmReturn: nil,
			mockConfirmErr:    nil,
			expectedError:     nil,
		},
		{
			name: "Fail - Confirm forget password error",
			request: models.ConfirmForgetPasswordReq{
				Username:    "test",
				Code:        "test-code",
				NewPassword: "test-new-password",
			},
			mockConfirmReturn: nil,
			mockConfirmErr:    errors.New("confirm forget password failed"),
			expectedError:     errors.New("confirm forget password failed"),
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mockCognitoClient.ExpectedCalls = nil

			suite.mockCognitoClient.On("ConfirmForgotPassword", mock.Anything, mock.Anything).Return(tc.mockConfirmReturn, tc.mockConfirmErr)

			err := suite.mockUserService.ConfirmForgetPassword(context.Background(), tc.request)

			if tc.expectedError != nil {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error())
			} else {
				suite.NoError(err)
			}

			suite.mockCognitoClient.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromToken_HasGroups() {
	claims := jwt.MapClaims{
		"cognito:groups": []string{"agent"},
	}
	tokenWithGroups, err := generateTestJWT(claims)
	suite.NoError(err)

	role, err := suite.mockUserService.GetUserRoleFromToken(tokenWithGroups)

	suite.NoError(err)
	suite.Equal("agent", role)
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromToken_NoGroups_FallbackToCognito() {
	claims := jwt.MapClaims{
		"user": "foo",
	}
	tokenWithoutGroups, err := generateTestJWT(claims)
	suite.NoError(err)

	// Mock fallback to Cognito
	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).Return(&cognitoidentityprovider.GetUserOutput{
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("sub"),
				Value: aws.String("mock-username"),
			},
		},
	}, nil)

	suite.mockCognitoClient.On("AdminListGroupsForUser", mock.Anything, mock.Anything).Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []types.GroupType{
			{GroupName: aws.String("mock-group")},
		},
	}, nil)

	role, err := suite.mockUserService.GetUserRoleFromToken(tokenWithoutGroups)

	suite.NoError(err)
	suite.Equal("mock-group", role)
	suite.mockCognitoClient.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromToken_InvalidJWT() {
	invalidToken := "not.a.valid.jwt"

	role, err := suite.mockUserService.GetUserRoleFromToken(invalidToken)

	suite.Error(err)
	suite.Equal("", role)
	suite.Contains(err.Error(), "error parsing JWT token")
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromCognito_GetUserFails() {
	token := "valid-token"

	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).
		Return(nil, fmt.Errorf("failed to get user")).
		Once()

	role, err := suite.mockUserService.GetUserRoleFromCognito(token)

	suite.Equal("", role)
	suite.Error(err)
	suite.Contains(err.Error(), "error retrieving user from Cognito")

	suite.mockCognitoClient.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromCognito_NoSubAttribute() {
	token := "valid-token"

	// GetUser returns no "sub" attribute
	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.GetUserOutput{
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("test@example.com")},
			},
		}, nil).
		Once()

	role, err := suite.mockUserService.GetUserRoleFromCognito(token)

	suite.Equal("", role)
	suite.Error(err)
	suite.Contains(err.Error(), "username not found in Cognito attributes")

	suite.mockCognitoClient.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromCognito_AdminListGroupsFails() {
	token := "valid-token"

	// Mock GetUser with "sub"
	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.GetUserOutput{
			UserAttributes: []types.AttributeType{
				{Name: aws.String("sub"), Value: aws.String("mock-username")},
			},
		}, nil).
		Once()

	// AdminListGroupsForUser returns error
	suite.mockCognitoClient.On("AdminListGroupsForUser", mock.Anything, mock.Anything).
		Return(nil, fmt.Errorf("failed to list groups")).
		Once()

	role, err := suite.mockUserService.GetUserRoleFromCognito(token)

	suite.Equal("", role)
	suite.Error(err)
	suite.Contains(err.Error(), "error retrieving user groups from Cognito")

	suite.mockCognitoClient.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromCognito_NoGroups() {
	token := "valid-token"

	// Mock GetUser with "sub"
	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.GetUserOutput{
			UserAttributes: []types.AttributeType{
				{Name: aws.String("sub"), Value: aws.String("mock-username")},
			},
		}, nil).
		Once()

	// AdminListGroupsForUser returns empty groups
	suite.mockCognitoClient.On("AdminListGroupsForUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups: []types.GroupType{},
		}, nil).
		Once()

	role, err := suite.mockUserService.GetUserRoleFromCognito(token)

	suite.Equal("", role)
	suite.Error(err)
	suite.Contains(err.Error(), "user is not in any Cognito group")

	suite.mockCognitoClient.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUserRoleFromCognito_Success() {
	token := "valid-token"

	// Mock GetUser with "sub"
	suite.mockCognitoClient.On("GetUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.GetUserOutput{
			UserAttributes: []types.AttributeType{
				{Name: aws.String("sub"), Value: aws.String("mock-username")},
			},
		}, nil).
		Once()

	// AdminListGroupsForUser returns one group
	suite.mockCognitoClient.On("AdminListGroupsForUser", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups: []types.GroupType{
				{GroupName: aws.String("mock-group")},
			},
		}, nil).
		Once()

	role, err := suite.mockUserService.GetUserRoleFromCognito(token)

	suite.NoError(err)
	suite.Equal("mock-group", role)

	suite.mockCognitoClient.AssertExpectations(suite.T())
}
