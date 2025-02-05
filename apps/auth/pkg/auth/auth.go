// Package auth is used to initialize the Cognito client and define the user groups.
package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
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

// Init initializes the Cognito client
func Init() *cognitoidentityprovider.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cognito connected")
	return cognitoidentityprovider.NewFromConfig(cfg)
}
