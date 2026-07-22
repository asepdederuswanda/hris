package monitoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// Handler untuk HTTP endpoints Platform Monitoring.
type Handler struct {
	dbManager *database.Manager
	logger    *zap.Logger
}

// NewHandler membuat Handler baru.
func NewHandler(dbManager *database.Manager, logger *zap.Logger) *Handler {
	return &Handler{
		dbManager: dbManager,
		logger:    logger,
	}
}

// HealthCheck menangani GET /api/v1/platform/monitoring/health
func (h *Handler) HealthCheck(c *gin.Context) {
	dbHealth := h.dbManager.HealthCheck()

	allHealthy := true
	for _, err := range dbHealth {
		if err != nil {
			allHealthy = false
			break
		}
	}

	statusCode := http.StatusOK
	status := "healthy"
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
		status = "degraded"
	}

	// Build db status map
	dbStatus := make(map[string]string)
	for name, err := range dbHealth {
		if err != nil {
			dbStatus[name] = "unhealthy: " + err.Error()
		} else {
			dbStatus[name] = "connected"
		}
	}

	c.JSON(statusCode, gin.H{
		"success":  true,
		"status":   status,
		"service":  "hris-platform",
		"database": dbStatus,
	})
}

// TenantHealth menangani GET /api/v1/platform/monitoring/tenants
func (h *Handler) TenantHealth(c *gin.Context) {
	dbHealth := h.dbManager.HealthCheck()

	var tenants []map[string]interface{}
	for name, err := range dbHealth {
		if len(name) > 7 && name[:7] == "tenant:" {
			status := "healthy"
			if err != nil {
				status = "unhealthy"
			}
			tenants = append(tenants, gin.H{
				"company_id": name[7:],
				"status":     status,
				"error":      errString(err),
			})
		}
	}

	if tenants == nil {
		tenants = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_connections": len(tenants),
			"tenants":           tenants,
		},
	})
}

// TenantDetail menangani GET /api/v1/platform/monitoring/tenants/:id
func (h *Handler) TenantDetail(c *gin.Context) {
	id := c.Param("id")

	dbHealth := h.dbManager.HealthCheck()

	key := "tenant:" + id
	err, exists := dbHealth[key]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Tenant connection not found",
			},
		})
		return
	}

	status := "healthy"
	errorMsg := ""
	if err != nil {
		status = "unhealthy"
		errorMsg = err.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"company_id": id,
			"status":     status,
			"error":      errorMsg,
		},
	})
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
