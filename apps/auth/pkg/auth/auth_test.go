package auth

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

// TestInit_Success remains unchanged
func TestInit_Success(t *testing.T) {
	mockLoader := func(ctx context.Context) (aws.Config, error) {
		return aws.Config{}, nil
	}

	client := Init(mockLoader)

	assert.NotNil(t, client, "Expected a non-nil Cognito client")
}

// TestInit_Failure runs in a separate process
func TestInit_Failure(t *testing.T) {
	if os.Getenv("TEST_CRASH") == "1" {
		// Inside the subprocess, we execute Init to trigger log.Fatal
		mockLoader := func(ctx context.Context) (aws.Config, error) {
			return aws.Config{}, errors.New("failed to load config")
		}
		Init(mockLoader) // This will call log.Fatal and exit
		return           // Should never reach here
	}

	// Run this test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestInit_Failure")
	cmd.Env = append(os.Environ(), "TEST_CRASH=1")
	err := cmd.Run()

	// Check that the subprocess exited with an error
	exitError, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("Expected an exit error, got: %v", err)
	}

	// Verify that the subprocess exited with a non-zero code
	if exitError.ExitCode() == 0 {
		t.Fatalf("Expected non-zero exit code due to log.Fatal, got 0")
	}
}
