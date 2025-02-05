package services

// import (
// 	"context"
// 	"fmt"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
// )

// func (s *UserService) associateToken(ctx context.Context, session string) (string, error) {
// 	input := &cip.AssociateSoftwareTokenInput{
// 		Session: aws.String(session),
// 	}
// 	res, err := s.CognitoClient.AssociateSoftwareToken(ctx, input)
// 	if err != nil {
// 		return "", fmt.Errorf("error associating token: %w", err)
// 	}
// 	// tokenString := fmt.Sprintf("%v", token)
// 	// input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
// 	// 	AccessToken: aws.String(tokenString),
// 	// }
// 	return *res.SecretCode, nil
// }
