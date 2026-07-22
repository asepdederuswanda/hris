package modulemgmt

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint Module Management ke router group.
// Semua endpoint memerlukan JWT authentication + RBAC permission.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	protected := rg.Group("")
	protected.Use(authMW)
	protected.Use(rbacMW)
	{
		modules := protected.Group("/modules")
		{
			modules.GET("", handler.ListModules)
			modules.POST("", handler.CreateModule)
			modules.GET("/:id", handler.GetModule)
			modules.PUT("/:id", handler.UpdateModule)
			modules.GET("/:id/companies", handler.ListCompanyModules)
			modules.POST("/:id/activate", handler.ActivateModule)
			modules.POST("/:id/deactivate", handler.DeactivateModule)
		}
	}
}
