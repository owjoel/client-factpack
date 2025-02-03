package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/auth"
)

// Wrapper object for auth service functions
type UserService struct {
	CognitoClient *cip.Client
}

func NewUserService() *UserService {
	return &UserService{CognitoClient: auth.Init()}
}

// Admin API for creating accounts on behalf of users
func (s *UserService) AdminCreateUser(ctx context.Context, r models.SignUpReq) error {
	username, err := createUsername(r.Email)
	if err != nil {
		return fmt.Errorf("error creating username: %w", err)
	}

	output, err := s.CognitoClient.AdminCreateUser(ctx, &cip.AdminCreateUserInput{
		UserPoolId: aws.String(config.UserPoolId),
		Username:   aws.String(username),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(r.Email)},
		},
	})
	if err != nil {
		return fmt.Errorf("error during sign up: %w", err)
	}
	log.Printf("User %s created at %v\n", username, output.User.UserCreateDate)

	// Add User to Group. Allow fail, add user in through AWS console
	_, err = s.CognitoClient.AdminAddUserToGroup(ctx, &cip.AdminAddUserToGroupInput{
		GroupName: aws.String(auth.AdminGroup),
		UserPoolId: aws.String(config.UserPoolId),
		Username: aws.String(username),
	})
	if err != nil {
		log.Printf("Unable to add user %s into group \"%s\"\n", username, auth.AgentGroup)
	}
	return nil
}

func createUsername(email string) (string, error) {
	username := strings.Split(email, "@")[0]
	if len(username) == 0 {
		return "", nil
	}
	return username, nil
}

// Calculate secret hash for additional security
func CalculateSecretHash(username string) string {
	message := username + config.ClientId
	h := hmac.New(sha256.New, []byte(config.ClientSecret))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *UserService) ForgetPassword(ctx context.Context, r models.ForgetPasswordReq) error {

	username := r.Username

	input := &cip.ForgotPasswordInput{
		ClientId: aws.String(config.ClientId),
		Username: aws.String(username),
		SecretHash: aws.String(CalculateSecretHash(username)),
	}

	_, err := s.CognitoClient.ForgotPassword(context.Background(), input)
	if err != nil {
		fmt.Printf("failed to initiate password reset: %s", err)
		return err
	}

	fmt.Println("Password reset code sent successfully")
	return nil
}

// User logs in through username, password
func (s *UserService) UserLogin(ctx context.Context, r models.LoginReq) (string, error) {

	input := &cip.InitiateAuthInput{                                                                                                                                    
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": r.Username,
			"PASSWORD": r.Password,
			"SECRET_HASH": CalculateSecretHash(r.Username),
		},
		ClientId:   aws.String(config.ClientId),
	}

	// returns tokens on success, see https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
	response, err := s.CognitoClient.InitiateAuth(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to initiate auth: %w\n", err)
	}

	return getChallengeName(response.ChallengeName), nil
	
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

// Sets new password by user, check for next auth challenge
func (s *UserService) SetNewPassword(ctx context.Context, r models.SetNewPasswordReq) (string, error) {
	autoResetEnabled, err := strconv.ParseBool(config.AutoResetPassword)
	if err != nil {
		return "", fmt.Errorf("invalid AUTO_RESET_PASSWORD value: %w", err)
	}

	if !autoResetEnabled {
		return "", fmt.Errorf("auto-reset password is disabled")
	}

	challengeInput := &cip.RespondToAuthChallengeInput{
		ClientId:      aws.String(config.ClientId),
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		Session:       aws.String(r.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":     r.Username,
			"NEW_PASSWORD": r.NewPassword,
			"SECRET_HASH": CalculateSecretHash(r.Username),
		},
	}

	res, err := s.CognitoClient.RespondToAuthChallenge(ctx, challengeInput)
	if err != nil {
		return "", fmt.Errorf("failed to respond to auth challenge: %s", err)
	}

	return getChallengeName(res.ChallengeName), nil
}

// Associates the user's auth challenge session with account, returns an OTP which can be returned as a QR to the client
func (s *UserService) SetupMFA(ctx context.Context, session string) (string, error) {
	secretKey, err := s.associateToken(ctx, session)
	if err != nil {
		return "", fmt.Errorf("Error associating token: %w", err)
	}
	return secretKey, nil
}

// Submit code from user's authentication app. Cognito will update MFA settings, but does not complete authentication. Nothing is returned to user on success.
func (s *UserService) VerifyMFA(ctx context.Context, r models.VerifyMFAReq) error {
	input := &cip.VerifySoftwareTokenInput{
		Session: aws.String(r.Session),
		UserCode: aws.String(r.Code),
	}
	res, err := s.CognitoClient.VerifySoftwareToken(ctx, input)
	if err != nil || res.Status != "SUCCESS" {
		return fmt.Errorf("failed to verify MFA: %w", err)
	}
	return nil
}

func (s *UserService) SignInMFA(ctx context.Context, r models.SignInMFAReq) (models.AuthenticationRes, error) {
	input := &cip.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeSoftwareTokenMfa,
		ClientId: aws.String(config.ClientId),
		Session: aws.String(r.Session),
		ChallengeResponses: map[string]string{
			"USERNAME": r.Username,
			"SOFTWARE_TOKEN_MFA_CODE": r.Code,
		},
	}
	res, err := s.CognitoClient.RespondToAuthChallenge(ctx, input)
	if err != nil {
		return models.AuthenticationRes{}, fmt.Errorf("failed to sign in with MFA: %w", err)
	}
	return models.AuthenticationRes{Result: *res.AuthenticationResult}, nil
}

func (s *UserService) ConfirmForgetPassword(ctx context.Context, r models.ConfirmForgetPasswordReq) error {
	input := &cip.ConfirmForgotPasswordInput{
		ClientId: aws.String(config.ClientId),
		Username: aws.String(r.Username),
		Password: aws.String(r.NewPassword),
		ConfirmationCode: aws.String(r.Code),
		SecretHash: aws.String(CalculateSecretHash(r.Username)),
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
