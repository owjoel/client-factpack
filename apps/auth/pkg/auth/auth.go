// Package auth is used to initialize the Cognito client and define the user groups.
package auth

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

/*
AdminGroup is the group name for admin users
AgentGroup is the group name for agent users
*/
var (
	AdminGroup = "admin"
	AgentGroup = "agent"
)

// Init initializes the Cognito client with a configurable loader function.
func Init(loadConfig func(context.Context) (aws.Config, error)) *cognitoidentityprovider.Client {
	cfg, err := loadConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cognito connected")
	return cognitoidentityprovider.NewFromConfig(cfg)
}
