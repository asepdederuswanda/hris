package authz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MiddlewareConfig untuk middleware RBAC.
type MiddlewareConfig struct {
	Enforcer *Enforcer
}

// NewMiddleware membuat middleware Gin untuk permission checking.
//
// Middleware ini HARUS ditempatkan SETELAH AuthJWT middleware
// agar claims (role) sudah tersedia di context.
func NewMiddleware(cfg MiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Role sudah di-set oleh AuthJWT middleware
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Missing role in context. Auth middleware required.",
				},
			})
			return
		}

		// Extract resource dari path
		resource := ResourceFromPath(c.FullPath())
		if resource == "" {
			// Path tidak dikenal, allow by default
			c.Next()
			return
		}

		// Extract action dari HTTP method
		action := ActionFromMethod(c.Request.Method)

		// Check permission
		if cfg.Enforcer.Check(role, resource, action) == DecisionDeny {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You don't have permission to perform this action",
					"details": gin.H{
						"role":     role,
						"resource": resource,
						"action":   action,
					},
				},
			})
			return
		}

		c.Next()
	}
}

// RequirePermission adalah middleware factory untuk validasi
// permission spesifik pada suatu route. Digunakan untuk route
// yang memerlukan permission berbeda dari action default.
//
// Contoh:
//
//	router.POST("/companies/:id/suspend", RequirePermission("company", "suspend"), handler.Suspend)
//	router.POST("/companies/:id/activate", RequirePermission("company", "activate"), handler.Activate)
func RequirePermission(enforcer *Enforcer, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Missing role in context.",
				},
			})
			return
		}

		if enforcer.Check(role, resource, action) == DecisionDeny {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You don't have permission to perform this action",
					"details": gin.H{
						"role":     role,
						"resource": resource,
						"action":   action,
					},
				},
			})
			return
		}

		c.Next()
	}
}
