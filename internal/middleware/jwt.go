package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/octguy/auth-sqlc/internal/service"
	"github.com/octguy/auth-sqlc/pkg/response"
)

// UserIDKey is the context key for the authenticated user's ID.
const UserIDKey = "userID"

// JWTAuth validates the Bearer token and injects the user ID into the context.
func JWTAuth(authSvc service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, http.StatusUnauthorized, "format: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID) // uuid.UUID is JSON-marshallable, so no need to convert to string
		c.Next()
	}
}
