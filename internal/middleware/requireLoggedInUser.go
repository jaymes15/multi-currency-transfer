package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuthenticatedUser middleware ensures the user is authenticated
func RequireAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the ContextGetUser helper to retrieve the user information from the request context
		user := ContextGetUser(c)

		// If the user is anonymous, then return an unauthorized response
		if user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// User is authenticated, continue to the next handler
		c.Next()
	}
}

// RequireAuthenticatedUserWithRole middleware ensures the user is authenticated and has a specific role
func RequireAuthenticatedUserWithRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		user := ContextGetUser(c)

		if user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Then check if user has the required role
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		// User is authenticated and has the required role, continue
		c.Next()
	}
}
