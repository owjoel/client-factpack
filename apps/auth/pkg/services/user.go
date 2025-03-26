package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	// awsconfig "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/auth"
)

type UserInterface interface {
	AdminCreateUser(ctx context.Context, r models.SignUpReq) error
	ForgetPassword(ctx context.Context, r models.ForgetPasswordReq) error
	UserLogin(ctx context.Context, r models.LoginReq) (*models.LoginRes, error)
	SetNewPassword(ctx context.Context, r models.SetNewPasswordReq) (*models.SetNewPasswordRes, error)
	SetupMFA(ctx context.Context, session string) (*models.AssociateTokenRes, error)
	SignInMFA(ctx context.Context, r models.SignInMFAReq) (models.AuthenticationRes, error)
	VerifyMFA(ctx context.Context, r models.VerifyMFAReq) error
	ConfirmForgetPassword(ctx context.Context, r models.ConfirmForgetPasswordReq) error
	GetUserRoleFromToken(token string) (string, error)
}

type CognitoClientInterface interface {
	AdminCreateUser(ctx context.Context, params *cip.AdminCreateUserInput, optFns ...func(*cip.Options)) (*cip.AdminCreateUserOutput, error)
	AdminAddUserToGroup(ctx context.Context, params *cip.AdminAddUserToGroupInput, optFns ...func(*cip.Options)) (*cip.AdminAddUserToGroupOutput, error)
	AdminRemoveUserFromGroup(ctx context.Context, params *cip.AdminRemoveUserFromGroupInput, optFns ...func(*cip.Options)) (*cip.AdminRemoveUserFromGroupOutput, error)
	AdminListGroupsForUser(ctx context.Context, params *cip.AdminListGroupsForUserInput, optFns ...func(*cip.Options)) (*cip.AdminListGroupsForUserOutput, error)
	AdminInitiateAuth(ctx context.Context, params *cip.AdminInitiateAuthInput, optFns ...func(*cip.Options)) (*cip.AdminInitiateAuthOutput, error)
	AdminRespondToAuthChallenge(ctx context.Context, params *cip.AdminRespondToAuthChallengeInput, optFns ...func(*cip.Options)) (*cip.AdminRespondToAuthChallengeOutput, error)
	VerifySoftwareToken(ctx context.Context, params *cip.VerifySoftwareTokenInput, optFns ...func(*cip.Options)) (*cip.VerifySoftwareTokenOutput, error)
	RespondToAuthChallenge(ctx context.Context, params *cip.RespondToAuthChallengeInput, optFns ...func(*cip.Options)) (*cip.RespondToAuthChallengeOutput, error)
	ConfirmForgotPassword(ctx context.Context, params *cip.ConfirmForgotPasswordInput, optFns ...func(*cip.Options)) (*cip.ConfirmForgotPasswordOutput, error)
	AssociateSoftwareToken(ctx context.Context, params *cip.AssociateSoftwareTokenInput, optFns ...func(*cip.Options)) (*cip.AssociateSoftwareTokenOutput, error)
	InitiateAuth(ctx context.Context, params *cip.InitiateAuthInput, optFns ...func(*cip.Options)) (*cip.InitiateAuthOutput, error)
	ForgotPassword(ctx context.Context, params *cip.ForgotPasswordInput, optFns ...func(*cip.Options)) (*cip.ForgotPasswordOutput, error)
    GetUser(ctx context.Context, input *cip.GetUserInput, opts ...func(*cip.Options)) (*cip.GetUserOutput, error)
}

// UserService represents the service for user operations.
type UserService struct {
	CognitoClient CognitoClientInterface
}

// NewUserService creates a new user service.
func NewUserService(client CognitoClientInterface) *UserService {
	return &UserService{
		CognitoClient: client,
	}
}


func (s *UserService) AdminCreateUser(ctx context.Context, r models.SignUpReq) error {
	username, err := createUsername(r.Email)
	if err != nil {
		return fmt.Errorf("error creating username: %w", err)
	}
	if r.Role == "" {
		return fmt.Errorf("user role is required")
	}

	output, err := s.CognitoClient.AdminCreateUser(ctx, &cip.AdminCreateUserInput{
		UserPoolId: aws.String(config.UserPoolID),
		Username:   aws.String(username),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(r.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
	if err != nil {
		return fmt.Errorf("error during sign up: %w", err)
	}
	log.Printf("User %s created at %v\n", username, output.User.UserCreateDate)
	
	// Add User to Group. Allow fail, add user in through AWS console
	_, err = s.CognitoClient.AdminAddUserToGroup(ctx, &cip.AdminAddUserToGroupInput{
		GroupName: aws.String(r.Role),
		UserPoolId: aws.String(config.UserPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		log.Printf("Unable to add user %s into group \"%s\"\n", username, auth.AgentGroup)
	}
	return nil
}

// SignUpUser registers user with Cognito user pool via email and password
// func (s *UserService) SignUpUser(ctx context.Context, r models.SignUpReq) error {
// 	username, err := createUsername((r.Email))
// 	if err != nil {
// 		return fmt.Errorf("error create username: %w", err)
// 	}

// 	input := &cip.SignUpInput{
// 		ClientId:   aws.String(config.ClientID),
// 		Username:   aws.String(username),
// 		Password:   aws.String(r.Password),
// 		SecretHash: aws.String(CalculateSecretHash(username)),
// 		UserAttributes: []types.AttributeType{{
// 			Name: aws.String("email"), Value: aws.String(r.Email),
// 		}},
// 	}

// 	_, err = s.CognitoClient.SignUp(ctx, input)
// 	if err != nil {
// 		return fmt.Errorf("error during sign up: %w", err)
// 	}
// 	log.Printf("User %s created", username)

// 	// Add User to Group. Allow fail, add user in through AWS console
// 	_, err = s.CognitoClient.AdminAddUserToGroup(ctx, &cip.AdminAddUserToGroupInput{
// 		GroupName:  aws.String(auth.AdminGroup),
// 		UserPoolId: aws.String(config.UserPoolID),
// 		Username:   aws.String(username),
// 	})
// 	if err != nil {
// 		log.Printf("Unable to add user %s into group \"%s\"\n", username, auth.AgentGroup)
// 	}
// 	return nil
// }

func createUsername(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format")
	}
	if parts[0] == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	return parts[0], nil
}

// CalculateSecretHash calculates the secret hash for Cognito
func CalculateSecretHash(username string) string {
	message := username + config.ClientID
	h := hmac.New(sha256.New, []byte(config.ClientSecret))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ForgetPassword sends a password reset code to the user's email
func (s *UserService) ForgetPassword(ctx context.Context, r models.ForgetPasswordReq) error {
	username := r.Username

	input := &cip.ForgotPasswordInput{
		ClientId:   aws.String(config.ClientID),
		Username:   aws.String(username),
		SecretHash: aws.String(CalculateSecretHash(username)),
	}

	_, err := s.CognitoClient.ForgotPassword(context.Background(), input)
	if err != nil {
		fmt.Printf("failed to initiate password reset: %s", err)
		return err
	}

	return nil
}

// UserLogin authenticates user with Cognito user pool via email and password
func (s *UserService) UserLogin(ctx context.Context, r models.LoginReq) (*models.LoginRes, error) {

	input := &cip.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    r.Username,
			"PASSWORD":    r.Password,
			"SECRET_HASH": CalculateSecretHash(r.Username),
		},
		ClientId: aws.String(config.ClientID),
	}

	// returns tokens on success, see https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
	response, err := s.CognitoClient.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate auth: %w", err)
	}

	return &models.LoginRes{
		Challenge: getChallengeName(response.ChallengeName),
		Session:   *response.Session,
	}, nil

	// Return based on challenges
	// if response.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
	// 	fmt.Println("New password required") // TODO: return Session token or store it in db

	// 	// TODO: (FOR TESTING ONLY) remove in future
	// 	err := s.handleNewPasswordChallenge(ctx, r.Username, "Password@1", *response.Session)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println("Password updated successfully.")
	// } else {
	// 	fmt.Println("Authentication successful.")
	// }

	// If challenge is to set up MFA on first time login, get OTP and send back as QR code
	// if response.ChallengeName == types.ChallengeNameTypeMfaSetup {
	// 	fmt.Println("Setup MFA required")

	// 	secretKey, err := s.AssociateToken(ctx, *response.Session)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
}

// SetNewPassword is used to set a new password for the user and check for next auth challenge
func (s *UserService) SetNewPassword(ctx context.Context, r models.SetNewPasswordReq) (*models.SetNewPasswordRes, error) {

	challengeInput := &cip.RespondToAuthChallengeInput{
		ClientId:      aws.String(config.ClientID),
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		Session:       aws.String(r.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":     r.Username,
			"NEW_PASSWORD": r.NewPassword,
			"SECRET_HASH":  CalculateSecretHash(r.Username),
		},
	}

	res, err := s.CognitoClient.RespondToAuthChallenge(ctx, challengeInput)
	if err != nil {
		return nil, fmt.Errorf("failed to respond to auth challenge: %s", err)
	}

	return &models.SetNewPasswordRes{
		Challenge: getChallengeName(res.ChallengeName),
		Session:   *res.Session,
	}, nil
}

// SetupMFA is used to associate the user's auth challenge session with the account
// Returns an OTP which can be returned as a QR to the client
func (s *UserService) SetupMFA(ctx context.Context, session string) (*models.AssociateTokenRes, error) {
	input := &cip.AssociateSoftwareTokenInput{
		Session: aws.String(session),
	}
	res, err := s.CognitoClient.AssociateSoftwareToken(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error associating token: %w", err)
	}
	return &models.AssociateTokenRes{
		Token:   *res.SecretCode,
		Session: *res.Session,
	}, nil
}

// VerifyMFA is used to verify the MFA code from the user's authentication app
// Submit code from user's authentication app. Cognito will update MFA settings, but does not complete authentication.
// Nothing is returned to user on success.
func (s *UserService) VerifyMFA(ctx context.Context, r models.VerifyMFAReq) error {
	input := &cip.VerifySoftwareTokenInput{
		Session:  aws.String(r.Session),
		UserCode: aws.String(r.Code),
	}
	res, err := s.CognitoClient.VerifySoftwareToken(ctx, input)
	if err != nil || res.Status != "SUCCESS" {
		return fmt.Errorf("failed to verify MFA: %w", err)
	}
	return nil
}

// SignInMFA is used to sign in with MFA
func (s *UserService) SignInMFA(ctx context.Context, r models.SignInMFAReq) (models.AuthenticationRes, error) {
	input := &cip.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeSoftwareTokenMfa,
		ClientId:      aws.String(config.ClientID),
		Session:       aws.String(r.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":                r.Username,
			"SECRET_HASH":             CalculateSecretHash(r.Username),
			"SOFTWARE_TOKEN_MFA_CODE": r.Code,
		},
	}
	res, err := s.CognitoClient.RespondToAuthChallenge(ctx, input)
	if err != nil {
		return models.AuthenticationRes{}, fmt.Errorf("failed to sign in with MFA: %w", err)
	}
	return models.AuthenticationRes{Result: *res.AuthenticationResult, Challenge: getChallengeName(res.ChallengeName)}, nil
}

// ConfirmForgetPassword is used to confirm the password reset.
func (s *UserService) ConfirmForgetPassword(ctx context.Context, r models.ConfirmForgetPasswordReq) error {
	input := &cip.ConfirmForgotPasswordInput{
		ClientId:         aws.String(config.ClientID),
		Username:         aws.String(r.Username),
		Password:         aws.String(r.NewPassword),
		ConfirmationCode: aws.String(r.Code),
		SecretHash:       aws.String(CalculateSecretHash(r.Username)),
	}

	_, err := s.CognitoClient.ConfirmForgotPassword(ctx, input)
	if err != nil {
		fmt.Printf("failed to confirm password reset: %s\n", err)
		return err
	}
	return nil
}

// func (s *UserService) handleNewPasswordChallenge(ctx context.Context, username, newPassword, session string) error {
// 	// flag for testing
// 	autoResetEnabled, err := strconv.ParseBool(config.AutoResetPassword)
// 	if err != nil {
// 		return fmt.Errorf("invalid AUTO_RESET_PASSWORD value: %w", err)
// 	}

// 	if !autoResetEnabled {
// 		return fmt.Errorf("auto-reset password is disabled")
// 	}

// 	challengeInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
// 		ClientId:      aws.String(config.ClientId),
// 		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
// 		Session:       aws.String(session),
// 		ChallengeResponses: map[string]string{
// 			"USERNAME":     username,
// 			"NEW_PASSWORD": newPassword,
// 			"SECRET_HASH": CalculateSecretHash(username),
// 		},
// 	}

// 	_, err = s.CognitoClient.RespondToAuthChallenge(ctx, challengeInput)
// 	if err != nil {
// 		fmt.Printf("failed to respond to auth challenge: %s", err)
// 		return err
// 	}

// 	return nil
// }

// Returns the string value of the auth challenge that should be returned to the client
func getChallengeName(challenge types.ChallengeNameType) string {
	if challenge == types.ChallengeNameTypeNewPasswordRequired {
		return "NEW_PASSWORD_REQUIRED"
	} else if challenge == types.ChallengeNameTypeMfaSetup {
		return "MFA_SETUP"
	} else if challenge == types.ChallengeNameTypeSoftwareTokenMfa {
		return "SOFTWARE_TOKEN_MFA"
	}
	return ""
}

// Extract user role from JWT token
func (s *UserService) GetUserRoleFromToken(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse JWT token (without verification)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("error parsing JWT token: %v", err)
	}

	// Check if role exists in JWT
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, exists := claims["cognito:groups"]; exists {
			groupList := role.([]interface{}) // Cognito stores groups as an array
			if len(groupList) > 0 {
				return groupList[0].(string), nil // Return first group
			}
		}
	}

	// If no role in JWT, fallback to Cognito API
	return s.GetUserRoleFromCognito(tokenString)
}

// GetUserRoleFromCognito fetches the user's group membership from AWS Cognito
func (s *UserService) GetUserRoleFromCognito(token string) (string, error) {
	// Retrieve username from Cognito token
	resp, err := s.CognitoClient.GetUser(context.TODO(), &cip.GetUserInput{
		AccessToken: aws.String(token),
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving user from Cognito: %v", err)
	}

	var username string
	for _, attr := range resp.UserAttributes {
		if *attr.Name == "sub" { // Cognito uses "sub" as the unique user ID
			username = *attr.Value
			break
		}
	}
	if username == "" {
		return "", fmt.Errorf("username not found in Cognito attributes")
	}

	// Call ListGroupsForUser to get Cognito group memberships
	groupResp, err := s.CognitoClient.AdminListGroupsForUser(context.TODO(), &cip.AdminListGroupsForUserInput{
		UserPoolId: aws.String(config.UserPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving user groups from Cognito: %v", err)
	}

	//  Extract the first group name as the user's role
	log.Println(groupResp.Groups)
	if len(groupResp.Groups) > 0 {
		return *groupResp.Groups[0].GroupName, nil
	}

	return "", fmt.Errorf("user is not in any Cognito group")
}
