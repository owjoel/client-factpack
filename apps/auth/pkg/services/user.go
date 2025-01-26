package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
func (s *UserService) CreateUser(ctx context.Context, r models.SignUpRequest) error {

	username, err := createUsername(r.Email)
	if err != nil {
		return fmt.Errorf("error creating username: %w", err)
	}
	
	output, err := s.CognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(config.UserPoolId),
		Username: aws.String(username),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(r.Email)},
		},
	})
	if err != nil {
		return fmt.Errorf("error durring sign up: %w", err)
	}
	fmt.Printf("User %s created at %v\n", username, output.User.UserCreateDate)
	return nil
}

func (s *UserService) VerifyUser() {}

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