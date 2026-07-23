package authz

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan endpoint RBAC management ke router group platform.
// Endpoint:
//   GET/POST /roles                    → List / Create role
//   GET/PUT/DELETE /roles/:id          → Get / Update / Delete role
//   GET/POST /permissions               → List / Create permission
//   DELETE /permissions/:id             → Delete permission
//   POST /roles/:id/permissions         → Assign permission to role
//   DELETE /roles/:id/permissions/:permissionId → Revoke permission from role
//   POST /rbac/reload                   → Reload enforcer from database
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	rbac := rg.Group("/rbac")
	{
		// Roles
		rbac.GET("/roles", handler.ListRoles)
		rbac.POST("/roles", handler.CreateRole)
		rbac.GET("/roles/:id", handler.GetRole)
		rbac.PUT("/roles/:id", handler.UpdateRole)
		rbac.DELETE("/roles/:id", handler.DeleteRole)

		// Permissions
		rbac.GET("/permissions", handler.ListPermissions)
		rbac.POST("/permissions", handler.CreatePermission)
		rbac.DELETE("/permissions/:id", handler.DeletePermission)

		// Role-Permission assignments
		rbac.POST("/roles/:id/permissions", handler.AssignPermission)
		rbac.DELETE("/roles/:id/permissions/:permissionId", handler.RevokePermission)
	}
}
