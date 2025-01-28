package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/config"
)

func AuthMiddleware(requiredTokenUse string,
	jwks *keyfunc.JWKS) gin.HandlerFunc {

	awsDefaultRegion := config.AwsRegion
	cognitoUserPoolId := config.UserPoolId
	cognitoAppClientId := config.ClientId

	return func(c *gin.Context) {

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
	}
}
