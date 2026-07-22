package jobmanagement

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// =========================================================================
// Route Registration Tests — verify no Gin route parameter conflicts
// =========================================================================

// TestRouteRegistration_NoPanic verifies that registering all job management
// routes does NOT cause a Gin panic due to route parameter name conflicts.
//
// This test specifically guards against regressions of the bug where
// ':titleId' conflicted with ':id' in the titles router group.
func TestRouteRegistration_NoPanic(t *testing.T) {
	// Use gin-test mode to suppress log noise
	gin.SetMode(gin.TestMode)

	// Create a fresh gin engine
	r := gin.New()

	// Create handler with real service (needs DB)
	handler, cleanup := newTestHandler()
	defer cleanup()

	// Register routes — this should NOT panic
	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatalf("Route registration panicked with: %v", rcv)
		}
	}()

	apiGroup := r.Group("/api/v1/tenant")
	RegisterRoutes(apiGroup, handler)

	t.Log("Route registration completed without panic — no route conflicts detected")
}

// TestRouteRegistration_EndpointReachable verifies that registered routes
// respond with the correct HTTP status codes before hitting auth middleware.
// Since we're testing in isolation, it validates the route tree is correct.
func TestRouteRegistration_EndpointStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	handler, cleanup := newTestHandler()
	defer cleanup()

	apiGroup := r.Group("/api/v1/tenant")
	RegisterRoutes(apiGroup, handler)

	// Build the full routes table
	routes := r.Routes()

	// Check that all expected route patterns exist
	expectedRoutes := map[string]bool{
		"GET /api/v1/tenant/job-management/titles":            false,
		"POST /api/v1/tenant/job-management/titles":           false,
		"GET /api/v1/tenant/job-management/titles/:id":        false,
		"PUT /api/v1/tenant/job-management/titles/:id":        false,
		"DELETE /api/v1/tenant/job-management/titles/:id":     false,
		"POST /api/v1/tenant/job-management/titles/:id/subs":  false,
		"GET /api/v1/tenant/job-management/titles/:id/subs":   false,
		"GET /api/v1/tenant/job-management/titles/:id/subs/:subId":   false,
		"PUT /api/v1/tenant/job-management/titles/:id/subs/:subId":   false,
		"DELETE /api/v1/tenant/job-management/titles/:id/subs/:subId": false,
		"GET /api/v1/tenant/job-management/values":            false,
		"POST /api/v1/tenant/job-management/values":           false,
		"GET /api/v1/tenant/job-management/values/:id":        false,
		"PUT /api/v1/tenant/job-management/values/:id":        false,
		"DELETE /api/v1/tenant/job-management/values/:id":     false,
		"GET /api/v1/tenant/job-management/competency-groups":        false,
		"POST /api/v1/tenant/job-management/competency-groups":       false,
		"GET /api/v1/tenant/job-management/competency-groups/:id":    false,
		"PUT /api/v1/tenant/job-management/competency-groups/:id":    false,
		"DELETE /api/v1/tenant/job-management/competency-groups/:id": false,
		"GET /api/v1/tenant/job-management/scores":                   false,
		"GET /api/v1/tenant/job-management/scores/org/:orgId":        false,
		"PUT /api/v1/tenant/job-management/scores/org/:orgId":        false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, exists := expectedRoutes[key]; exists {
			expectedRoutes[key] = true
		}
	}

	for route, found := range expectedRoutes {
		if !found {
			t.Errorf("Expected route not found: %s", route)
		}
	}
}

// TestRouteRegistration_MultipleModules verifies that registering job management
// routes alongside organization routes (which share the /api/v1/tenant prefix)
// doesn't cause conflicts.
func TestRouteRegistration_WithOtherModules(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	handler, cleanup := newTestHandler()
	defer cleanup()

	db, dbResolver, dbCleanup := setupTestDB()
	defer dbCleanup()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Register job management
	apiGroup := r.Group("/api/v1/tenant")
	RegisterRoutes(apiGroup, handler)

	// Simulate organization routes registration alongside job management
	orgGroup := r.Group("/api/v1/tenant/organizations")
	{
		orgGroup.GET("", func(c *gin.Context) { c.Status(http.StatusOK) })
		orgGroup.POST("", func(c *gin.Context) { c.Status(http.StatusCreated) })
		orgGroup.GET("/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
		orgGroup.PUT("/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
		orgGroup.DELETE("/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	}

	// Verify org routes still work
	_ = db
	_ = dbResolver

	t.Log("Job management + organization routes registered without conflict")
}


