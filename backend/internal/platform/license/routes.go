package license

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint License Management ke router group.
// Semua endpoint memerlukan JWT authentication + RBAC permission.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	protected := rg.Group("")
	protected.Use(authMW)
	protected.Use(rbacMW)
	{
		licenses := protected.Group("/licenses")
		{
			licenses.GET("", handler.ListLicenses)
			licenses.POST("", handler.CreateLicense)
			licenses.GET("/:id", handler.GetLicense)
			licenses.PUT("/:id", handler.UpdateLicense)
		}
	}
}
