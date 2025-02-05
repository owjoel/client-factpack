package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/config"
)

// Authenticate is a middleware that checks if the user is authenticated by validating the "accessToken" cookie.
// If the cookie is missing or invalid, it returns a 403 Forbidden response.
// Otherwise, it sets the token in the context for downstream handlers to use.
func (h *UserHandler) Authenticate(c *gin.Context) {

	requiredTokenUse := "access" // default check for access token
	awsDefaultRegion := config.AwsRegion
	cognitoUserPoolId := config.UserPoolID
	cognitoAppClientId := config.ClientID

	jwks, err := GetJWKS(awsDefaultRegion, cognitoUserPoolId)
	if err != nil {
		log.Fatalf("Failed to retrieve Cognito JWKS\nError: %s", err)
	}

	tokenString, err := c.Cookie("access_token")
	if err != nil || tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// * Verify the signature of the JWT
	// * Verify that the algorithm used is RS256
	// * Verify that the 'exp' claim exists in the token
	// * Verification of audience 'aud' is taken care later when we examine if the
	//   token is 'id' or 'access'
	// * The issuer (iss) claim should match your user pool. For example, a user
	//   pool created in the us-east-1 region
	//   will have the following iss value: https://cognito-idp.us-east-1.amazonaws.com/<userpoolID>.
	token, err := jwt.Parse(tokenString,
		jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", awsDefaultRegion, cognitoUserPoolId)))
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Attempt to parse the JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Compare the "exp" claim to the current time
	expClaim, err := claims.GetExpirationTime()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	if expClaim.Unix() < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Check the token_use claim.
	// If you are only accepting the access token in your web API operations, its value must be access.
	// If you are only using the ID token, its value must be id.
	// If you are using both ID and access tokens, the token_use claim must be either id or access.
	tokenUseClaim, ok := claims["token_use"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	if tokenUseClaim != requiredTokenUse {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// "sub" claim exists in both ID and Access tokens
	subClaim, err := claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Set("username", subClaim)

	// The "aud" claim in an ID token and the "client_id" claim in an access token should match the app
	// client ID that was created in the Amazon Cognito user pool.
	var appClientIdClaim string
	if tokenUseClaim == "id" {
		audienceClaims, err := claims.GetAudience()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		appClientIdClaim = audienceClaims[0]

	} else if tokenUseClaim == "access" {
		clientIdClaim, ok := claims["client_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		appClientIdClaim = clientIdClaim
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	if appClientIdClaim != cognitoAppClientId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Retrieve any Cognito user groups that the user belongs to
	userGroupsAttribute, ok := claims["cognito:groups"]
	userGroupsClaims := make([]string, 0)
	if ok {
		switch x := userGroupsAttribute.(type) {
		case []interface{}:
			for _, e := range x {
				userGroupsClaims = append(userGroupsClaims, e.(string))
			}
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}

	c.Set("groups", userGroupsClaims)

	c.Next()

	c.Set("accessToken", token)
	c.Next()
}

// VerifyMFA verifies the user's multi-factor authentication (MFA) setup using their access token.
// It retrieves the "accessToken" from the context, and if it doesn't exist, returns a 403 Forbidden response.
// If the token exists, it prepares the input for associating an MFA software token with AWS Cognito.
func (h *UserHandler) VerifyMFA(c *gin.Context) {
	token, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"message": "Could not verify identity"})
	}
	tokenString := fmt.Sprintf("%v", token)
	input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(tokenString),
	}
	log.Println(input)
}

// Helper function for Authenticate middleware
func GetJWKS(awsRegion string, cognitoUserPoolId string) (*keyfunc.JWKS, error) {

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, cognitoUserPoolId)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	return jwks, nil
}
