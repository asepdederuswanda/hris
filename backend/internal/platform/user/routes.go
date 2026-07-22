package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint User & Auth ke router group.
// login & refresh bersifat public (no auth), sisanya memerlukan JWT + RBAC.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	// Public: Auth endpoints (no auth required)
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.RefreshToken)

	// Protected: User management (JWT + RBAC required)
	protected := rg.Group("")
	protected.Use(authMW)
	protected.Use(rbacMW)
	{
		users := protected.Group("/users")
		{
			users.GET("", handler.ListUsers)
			users.GET("/:id", handler.GetUser)
			users.POST("", handler.CreateUser)
			users.PUT("/:id", handler.UpdateUser)
		}
	}
}
