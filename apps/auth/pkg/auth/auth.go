// Cognito init and config variables needed
package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	AdminGroup = "admin"
	AgentGroup = "agent"
)

func Init() *cognitoidentityprovider.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cognito connected")
	return cognitoidentityprovider.NewFromConfig(cfg)
}