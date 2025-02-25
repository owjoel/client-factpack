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
	"github.com/owjoel/client-factpack/errors" // Import errors package
)

// Authenticate is a middleware that checks if the user is authenticated by validating the "accessToken" cookie.
func (h *UserHandler) Authenticate(c *gin.Context) {

	requiredTokenUse := "access" // default check for access token
	awsDefaultRegion := config.AwsRegion
	cognitoUserPoolId := config.UserPoolID
	cognitoAppClientId := config.ClientID

	jwks, err := GetJWKS(awsDefaultRegion, cognitoUserPoolId)
	if err != nil {
		log.Printf("Failed to retrieve Cognito JWKS: %s", err)
		errorResponse(c, errors.ErrServerError)
		return
	}

	tokenString, err := c.Cookie("access_token")
	if err != nil || tokenString == "" {
		errorResponse(c, errors.ErrUnauthorized)
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
		errorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Parse JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		errorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Validate token expiration
	expClaim, err := claims.GetExpirationTime()
	if err != nil {
		errorResponse(c, errors.ErrInvalidToken)
		return
	}
	if expClaim.Unix() < time.Now().Unix() {
		errorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Validate token use
	tokenUseClaim, ok := claims["token_use"].(string)
	if !ok || tokenUseClaim != requiredTokenUse {
		errorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Validate subject claim (user identifier)
	subClaim, err := claims.GetSubject()
	if err != nil {
		errorResponse(c, errors.ErrUnauthorized)
		return
	}
	c.Set("username", subClaim)

	// Validate client ID
	var appClientIdClaim string
	if tokenUseClaim == "id" {
		audienceClaims, err := claims.GetAudience()
		if err != nil || len(audienceClaims) == 0 {
			errorResponse(c, errors.ErrUnauthorized)
			return
		}
		appClientIdClaim = audienceClaims[0]
	} else if tokenUseClaim == "access" {
		clientIdClaim, ok := claims["client_id"].(string)
		if !ok {
			errorResponse(c, errors.ErrUnauthorized)
			return
		}
		appClientIdClaim = clientIdClaim
	} else {
		errorResponse(c, errors.ErrUnauthorized)
		return
	}

	if appClientIdClaim != cognitoAppClientId {
		errorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Retrieve Cognito user groups
	userGroupsAttribute, ok := claims["cognito:groups"]
	userGroupsClaims := make([]string, 0)
	if ok {
		switch x := userGroupsAttribute.(type) {
		case []interface{}:
			for _, e := range x {
				userGroupsClaims = append(userGroupsClaims, e.(string))
			}
		default:
			errorResponse(c, errors.ErrUnauthorized)
			return
		}
	}
	c.Set("groups", userGroupsClaims)

	c.Next()
	c.Set("accessToken", token)
	c.Next()
}

// VerifyMFA verifies the user's multi-factor authentication (MFA) setup using their access token.
func (h *UserHandler) VerifyMFA(c *gin.Context) {
	token, exists := c.Get("accessToken")
	if !exists {
		errorResponse(c, errors.ErrUnauthorized)
		return
	}

	tokenString := fmt.Sprintf("%v", token)
	input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(tokenString),
	}
	log.Println(input)
}

// Helper function for standardized error response
func errorResponse(c *gin.Context, err errors.CustomError) {
	c.JSON(err.Status, gin.H{
		"error_code": err.Code,
		"message":    err.Message,
	})
	c.Abort()
}

// GetJWKS retrieves Cognito JSON Web Key Set (JWKS) for verifying JWTs.
func GetJWKS(awsRegion string, cognitoUserPoolId string) (*keyfunc.JWKS, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, cognitoUserPoolId)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	return jwks, nil
}
