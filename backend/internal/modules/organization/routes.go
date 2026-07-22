package organization

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	orgs := rg.Group("/organizations")
	{
		orgs.POST("", handler.Create)
		orgs.GET("", handler.List)
		orgs.GET("/:id", handler.GetByID)
		orgs.PUT("/:id", handler.Update)
		orgs.DELETE("/:id", handler.Delete)
	}
}
