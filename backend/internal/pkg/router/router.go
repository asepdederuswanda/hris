// Package router menyediakan setup router Gin dan registrasi
// semua modul (platform & tenant) ke dalam grup route yang sesuai.
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/module"
)

// Config untuk router setup.
type Config struct {
	Mode string // debug, release, test
}

// Setup membuat Gin engine baru, mengatur middleware global,
// dan mendaftarkan semua modul ke route group yang sesuai.
//
// Platform modules: /api/v1/platform/*
// Tenant modules:   /api/v1/tenant/*
func Setup(cfg Config, authMiddleware, tenantMiddleware, rbacMiddleware gin.HandlerFunc,
	corsMiddleware gin.HandlerFunc, loggerMiddleware gin.HandlerFunc,
	recoveryMiddleware gin.HandlerFunc,
	platformModules []module.ModuleRegistration,
	tenantModules []module.ModuleRegistration,
) *gin.Engine {

	// Set Gin mode
	gin.SetMode(cfg.Mode)

	r := gin.New()

	// Global middleware
	r.Use(recoveryMiddleware)
	r.Use(loggerMiddleware)
	r.Use(corsMiddleware)

	// Health check (no auth)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "hris-platform"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Platform routes (no tenant context)
	platformGroup := r.Group("/api/v1/platform")
	{
		for _, reg := range platformModules {
			reg.Module.RegisterRoutes(platformGroup)
		}
	}

	// Tenant routes (with auth + tenant + RBAC middleware)
	tenantGroup := r.Group("/api/v1/tenant")
	tenantGroup.Use(authMiddleware)
	tenantGroup.Use(tenantMiddleware)
	tenantGroup.Use(rbacMiddleware)
	{
		for _, reg := range tenantModules {
			reg.Module.RegisterRoutes(tenantGroup)
		}
	}

	return r
}
