package monitoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/cache"
	"github.com/inthros/hris-platform/internal/pkg/database"
)

// Handler untuk HTTP endpoints Platform Monitoring.
type Handler struct {
	dbManager    *database.Manager
	cacheManager *cache.Cache
	logger       *zap.Logger
}

// NewHandler membuat Handler baru.
func NewHandler(dbManager *database.Manager, cacheManager *cache.Cache, logger *zap.Logger) *Handler {
	return &Handler{
		dbManager:    dbManager,
		cacheManager: cacheManager,
		logger:       logger,
	}
}

// =========================================================================
// GET /api/v1/platform/monitoring/health
// =========================================================================
// HealthCheck mengembalikan status kesehatan sistem termasuk:
//   - Status koneksi semua database (platform + tenants)
//   - Ringkasan connection pool (total open, idle, in-use, waiting)
//   - Status Redis cache connection
//
// Response: 200 jika semua sehat, 503 jika ada yang tidak sehat.

func (h *Handler) HealthCheck(c *gin.Context) {
	dbHealth := h.dbManager.HealthCheck()
	poolStats := h.dbManager.PoolStats()

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

	// DB connection status
	dbStatus := make(map[string]string)
	for name, err := range dbHealth {
		if err != nil {
			dbStatus[name] = "unhealthy: " + err.Error()
		} else {
			dbStatus[name] = "connected"
		}
	}

	// Cache health check
	cacheStatus := "connected"
	if h.cacheManager != nil {
		if err := h.cacheManager.Ping(c.Request.Context()); err != nil {
			cacheStatus = "unhealthy: " + err.Error()
			allHealthy = false
			statusCode = http.StatusServiceUnavailable
			status = "degraded"
		}
	} else {
		cacheStatus = "not_initialized"
	}

	// Pool summary: aggregate across all connections
	poolSummary := poolSummaryStats(poolStats)

	c.JSON(statusCode, gin.H{
		"success":     true,
		"status":      status,
		"service":     "hris-platform",
		"database":    dbStatus,
		"cache":       cacheStatus,
		"pool_stats":  poolSummary,
		"pool_detail": poolStats,
	})
}

// =========================================================================
// GET /api/v1/platform/monitoring/pool
// =========================================================================
// PoolStats mengembalikan detail connection pool untuk setiap koneksi
// database (platform + semua tenant aktif).
//
// Berguna untuk:
//   - Debug connection leak (lihat InUse vs Idle)
//   - Monitoring pool pressure (lihat WaitCount, WaitDuration)
//   - Capacity planning (lihat MaxIdleClosed, MaxLifetimeClosed)

func (h *Handler) PoolStats(c *gin.Context) {
	stats := h.dbManager.PoolStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_connections": len(stats),
			"connections":       stats,
			"summary":           poolSummaryStats(stats),
		},
	})
}

// =========================================================================
// GET /api/v1/platform/monitoring/tenants
// =========================================================================
// TenantHealth mengembalikan daftar tenant beserta status koneksi
// dan statistik pool masing-masing.

func (h *Handler) TenantHealth(c *gin.Context) {
	dbHealth := h.dbManager.HealthCheck()
	poolStats := h.dbManager.PoolStats()

	var tenants []map[string]interface{}
	for name, err := range dbHealth {
		if len(name) <= 7 || name[:7] != "tenant:" {
			continue
		}

		companyID := name[7:]
		status := "healthy"
		if err != nil {
			status = "unhealthy"
		}

		tenantData := gin.H{
			"company_id": companyID,
			"status":     status,
			"error":      errString(err),
		}

		// Attach pool stats jika tersedia
		if stat, ok := poolStats[name]; ok {
			tenantData["pool"] = gin.H{
				"open":              stat.Open,
				"idle":              stat.Idle,
				"in_use":            stat.InUse,
				"max_open":          stat.MaxOpen,
				"wait_count":        stat.WaitCount,
				"wait_duration":     stat.WaitDuration,
				"max_idle_closed":   stat.MaxIdleClosed,
				"max_lifetime_closed": stat.MaxLifetimeClosed,
			}
		}

		tenants = append(tenants, tenantData)
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

// =========================================================================
// GET /api/v1/platform/monitoring/tenants/:id
// =========================================================================
// TenantDetail mengembalikan status kesehatan dan pool statistik
// untuk satu tenant tertentu.

func (h *Handler) TenantDetail(c *gin.Context) {
	id := c.Param("id")

	dbHealth := h.dbManager.HealthCheck()
	poolStats := h.dbManager.PoolStats()

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

	data := gin.H{
		"company_id": id,
		"status":     status,
		"error":      errorMsg,
	}

	// Attach detailed pool stats
	if stat, ok := poolStats[key]; ok {
		data["pool"] = stat
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// =========================================================================
// Helpers
// =========================================================================

// poolSummaryStats menghitung aggregate dari semua pool stats.
func poolSummaryStats(stats map[string]*database.PoolStat) gin.H {
	totalOpen := 0
	totalIdle := 0
	totalInUse := 0
	totalWaitCount := int64(0)
	connectedCount := 0

	for _, s := range stats {
		totalOpen += s.Open
		totalIdle += s.Idle
		totalInUse += s.InUse
		totalWaitCount += s.WaitCount
		if s.Open > 0 || s.Idle > 0 {
			connectedCount++
		}
	}

	return gin.H{
		"total_connections": len(stats),
		"active_connections": connectedCount,
		"total_open":         totalOpen,
		"total_idle":         totalIdle,
		"total_in_use":       totalInUse,
		"total_wait_count":   totalWaitCount,
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
