package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler untuk HTTP endpoints Company.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create menangani POST /api/v1/platform/companies
func (h *Handler) Create(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	response, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Company created successfully",
	})
}

// GetByID menangani GET /api/v1/platform/companies/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// List menangani GET /api/v1/platform/companies
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.List(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update menangani PUT /api/v1/platform/companies/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	response, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Company updated successfully",
	})
}

// Delete menangani DELETE /api/v1/platform/companies/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Company deleted successfully",
	})
}

// Suspend menangani POST /api/v1/platform/companies/:id/suspend
func (h *Handler) Suspend(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Suspend(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "SUSPEND_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Company suspended successfully — tenant connection deactivated",
	})
}

// Activate menangani POST /api/v1/platform/companies/:id/activate
func (h *Handler) Activate(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Activate(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ACTIVATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Company activated successfully — tenant connection reactivated",
	})
}

// Terminate menangani POST /api/v1/platform/companies/:id/terminate
func (h *Handler) Terminate(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Terminate(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TERMINATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Company terminated — tenant database dropped and connection removed",
	})
}

// Backup menangani POST /api/v1/platform/companies/:id/backup
func (h *Handler) Backup(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement actual backup logic (Phase 2)
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Backup initiated (Phase 2 - not yet implemented)",
	})
}

// Restore menangani POST /api/v1/platform/companies/:id/restore
func (h *Handler) Restore(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement actual restore logic (Phase 2)
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Restore initiated (Phase 2 - not yet implemented)",
	})
}
