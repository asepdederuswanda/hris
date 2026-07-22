package company

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint Company ke router group.
// Semua endpoint memerlukan JWT authentication + RBAC permission.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	protected := rg.Group("")
	protected.Use(authMW)
	protected.Use(rbacMW)
	{
		companies := protected.Group("/companies")
		{
			companies.POST("", handler.Create)
			companies.GET("", handler.List)
			companies.GET("/:id", handler.GetByID)
			companies.PUT("/:id", handler.Update)
			companies.DELETE("/:id", handler.Delete)
			companies.POST("/:id/suspend", handler.Suspend)
			companies.POST("/:id/activate", handler.Activate)
			companies.POST("/:id/terminate", handler.Terminate)
			companies.POST("/:id/backup", handler.Backup)
			companies.POST("/:id/restore", handler.Restore)
		}
	}
}
