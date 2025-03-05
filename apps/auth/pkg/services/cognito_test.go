package services

// import (
// 	"context"
// 	"errors"
// 	"os"
// 	"os/exec"
// 	"testing"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
// 	"github.com/stretchr/testify/assert"
// )

// // MockCognitoClient simulates the Cognito client behavior
// type MockCognitoClient struct{}

// // Simulated successful response for AssociateSoftwareToken
// func (m *MockCognitoClient) AssociateSoftwareToken(ctx context.Context, input *cognitoidentityprovider.AssociateSoftwareTokenInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error) {
// 	if input.Session == nil || *input.Session == "" {
// 		return nil, errors.New("invalid session")
// 	}
// 	return &cognitoidentityprovider.AssociateSoftwareTokenOutput{
// 		SecretCode: aws.String("mock-secret-code"),
// 	}, nil
// }

// // TestAssociateToken_Success validates successful token association
// func TestAssociateToken_Success(t *testing.T) {
// 	mockService := &UserService{CognitoClient: &MockCognitoClient{}}

// 	secretCode, err := mockService.associateToken(context.Background(), "valid-session")

// 	assert.Nil(t, err)
// 	assert.Equal(t, "mock-secret-code", secretCode)
// }

// // TestAssociateToken_Failure runs in a separate process to handle fatal errors
// func TestAssociateToken_Failure(t *testing.T) {
// 	if os.Getenv("TEST_CRASH") == "1" {
// 		mockService := &UserService{CognitoClient: &MockCognitoClient{}}

// 		// Intentionally passing an empty session to trigger an error
// 		_, err := mockService.associateToken(context.Background(), "")
// 		if err != nil {
// 			os.Exit(1) // Simulate process failure for log.Fatal
// 		}
// 		return
// 	}

// 	// Run this test in a subprocess
// 	cmd := exec.Command(os.Args[0], "-test.run=TestAssociateToken_Failure")
// 	cmd.Env = append(os.Environ(), "TEST_CRASH=1")
// 	err := cmd.Run()

// 	// Check that the subprocess exited with an error
// 	exitError, ok := err.(*exec.ExitError)
// 	if !ok {
// 		t.Fatalf("Expected an exit error, got: %v", err)
// 	}

// 	// Verify that the subprocess exited with a non-zero code
// 	if exitError.ExitCode() == 0 {
// 		t.Fatalf("Expected non-zero exit code due to log.Fatal, got 0")
// 	}
// }
