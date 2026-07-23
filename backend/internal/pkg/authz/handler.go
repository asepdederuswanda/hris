package authz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler untuk endpoint manajemen RBAC.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler RBAC baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Role Handlers
// =========================================================================

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	resp, err := h.service.CreateRole(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": roles})
}

func (h *Handler) GetRole(c *gin.Context) {
	resp, err := h.service.GetRole(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	resp, err := h.service.UpdateRole(c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	if err := h.service.DeleteRole(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role deleted"})
}

// =========================================================================
// Permission Handlers
// =========================================================================

func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	resp, err := h.service.CreatePermission(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListPermissions(c *gin.Context) {
	perms, err := h.service.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": perms})
}

func (h *Handler) DeletePermission(c *gin.Context) {
	if err := h.service.DeletePermission(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission deleted"})
}

// =========================================================================
// Role-Permission Assignment Handlers
// =========================================================================

func (h *Handler) AssignPermission(c *gin.Context) {
	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if err := h.service.AssignPermission(c.Param("id"), req.PermissionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission assigned to role"})
}

func (h *Handler) RevokePermission(c *gin.Context) {
	if err := h.service.RevokePermission(c.Param("id"), c.Param("permissionId")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission revoked from role"})
}
