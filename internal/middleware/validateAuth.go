package middleware

import (
	"net/http"
	"strings"

	"lemfi/simplebank/config"
	"lemfi/simplebank/pkg/token"

	"github.com/gin-gonic/gin"
)

type contextKey string

const userContextKey = contextKey("user")

type UserClaimsData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var AnonymousUser = &UserClaimsData{}

func (u *UserClaimsData) IsAnonymous() bool {
	return u == AnonymousUser
}

func ContextSetUser(c *gin.Context, user *UserClaimsData) {
	c.Set(string(userContextKey), user)
}

func ContextGetUser(c *gin.Context) *UserClaimsData {
	user, exists := c.Get(string(userContextKey))
	if !exists {
		return AnonymousUser
	}

	if userData, ok := user.(*UserClaimsData); ok {
		return userData
	}

	return AnonymousUser
}

// ValidateAuth middleware validates the authentication token
func ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			config.Logger.Error("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			config.Logger.Error("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			config.Logger.Error("Empty token after Bearer prefix")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token cannot be empty",
			})
			c.Abort()
			return
		}

		// Validate the token
		tokenMaker := token.GetTokenMaker()
		payload, err := tokenMaker.VerifyToken(tokenString, token.TokenTypeAccessToken)
		if err != nil {
			config.Logger.Error("Token validation failed", "error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Create user claims data
		userData := &UserClaimsData{
			ID:       payload.ID.String(),
			Username: payload.Username,
			Role:     payload.Role,
		}

		// Set user data in context
		ContextSetUser(c, userData)

		// Log successful authentication
		config.Logger.Info("User authenticated successfully",
			"username", userData.Username,
			"userID", userData.ID,
			"endpoint", c.Request.URL.Path,
		)

		c.Next()
	}
}
