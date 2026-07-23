package monitoring

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Monitoring ke router group.
// Semua endpoint memerlukan JWT authentication + RBAC permission.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	protected := rg.Group("")
	protected.Use(authMW)
	protected.Use(rbacMW)
	{
		monitoring := protected.Group("/monitoring")
		{
			monitoring.GET("/health", handler.HealthCheck)
			monitoring.GET("/pool", handler.PoolStats)
			monitoring.GET("/tenants", handler.TenantHealth)
			monitoring.GET("/tenants/:id", handler.TenantDetail)
		}
	}
}
