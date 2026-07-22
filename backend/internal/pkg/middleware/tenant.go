package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TenantRequired memastikan bahwa request memiliki company_id
// (berasal dari JWT claims yang sudah di-set oleh AuthJWT middleware).
//
// Middleware ini:
// 1. Validasi company_id ada di gin context (dari JWT claims)
// 2. Propagate company_id ke request context agar bisa diakses
//    oleh service/repository melalui c.Request.Context()
func TenantRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.GetString("company_id")
		if companyID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TENANT_REQUIRED",
					"message": "Tenant context is required",
				},
			})
			return
		}

		// Propagate company_id ke request context agar bisa diakses
		// oleh service/repository via c.Request.Context()
		ctx := context.WithValue(c.Request.Context(), "company_id", companyID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetCompanyID adalah helper untuk mengambil company_id dari context.
func GetCompanyID(c *gin.Context) string {
	return c.GetString("company_id")
}

// GetUserID adalah helper untuk mengambil user_id dari context.
func GetUserID(c *gin.Context) string {
	return c.GetString("user_id")
}
