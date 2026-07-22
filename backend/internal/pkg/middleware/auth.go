// Package middleware menyediakan Gin middleware untuk
// authentication, authorization, tenant resolution, logging, CORS, dll.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/auth"
)

// AuthJWT mengembalikan Gin middleware yang memvalidasi JWT token
// dari header Authorization: Bearer <token>.
//
// Setelah validasi, middleware menyimpan claims ke dalam gin.Context:
//   - "user_id"
//   - "company_id"
//   - "email"
//   - "role"
//   - "permissions"
func AuthJWT(authManager *auth.Manager, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Missing authorization header",
				},
			})
			return
		}

		// Expect: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid authorization header format",
				},
			})
			return
		}

		tokenString := parts[1]
		claims, err := authManager.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Invalid JWT token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired token",
				},
			})
			return
		}

		// Verify it's an access token
		if claims.TokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid token type",
				},
			})
			return
		}

		// Set claims into context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("company_id", claims.CompanyID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}
