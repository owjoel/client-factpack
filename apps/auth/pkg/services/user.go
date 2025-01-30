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
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/api/models"
	"github.com/owjoel/client-factpack/apps/auth/pkg/auth"
)

// Wrapper object for auth service functions
type UserService struct {
	CognitoClient *cognitoidentityprovider.Client
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

	output, err := s.CognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
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
	_, err = s.CognitoClient.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
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

	input := &cognitoidentityprovider.ForgotPasswordInput{
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

func (s *UserService) UserLogin(ctx context.Context, r models.LoginReq) error {

	input := &cognitoidentityprovider.InitiateAuthInput{
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
		fmt.Printf("failed to initiate auth: %s\n", err)
		return err
	}
	
	if response.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		fmt.Println("New password required") // TODO: return Session token or store it in db


		// TODO: (FOR TESTING ONLY) remove in future
		err := s.handleNewPasswordChallenge(ctx, r.Username, "Password@1", *response.Session)
		if err != nil {
			return err
		}
		fmt.Println("Password updated successfully.")
	} else {
		fmt.Println("Authentication successful.")
	}

	return nil
}

func (s *UserService) ConfirmForgetPassword(ctx context.Context, r models.ConfirmForgetPasswordReq) error {
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
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


func (s *UserService) handleNewPasswordChallenge(ctx context.Context, username, newPassword, session string) error {
	// flag for testing
	autoResetEnabled, err := strconv.ParseBool(config.AutoResetPassword)
	if err != nil {
		return fmt.Errorf("invalid AUTO_RESET_PASSWORD value: %w", err)
	}

	if !autoResetEnabled {
		return fmt.Errorf("auto-reset password is disabled")
	}

	challengeInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:      aws.String(config.ClientId),
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		Session:       aws.String(session),
		ChallengeResponses: map[string]string{
			"USERNAME":     username,
			"NEW_PASSWORD": newPassword,
			"SECRET_HASH": CalculateSecretHash(username),
		},
	}

	_, err = s.CognitoClient.RespondToAuthChallenge(ctx, challengeInput)
	if err != nil {
		fmt.Printf("failed to respond to auth challenge: %s", err)
		return err
	}

	return nil
}
