package company

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// testEnv menyimpan lingkungan test untuk handler tests.
type testEnv struct {
	router  *gin.Engine
	db      *gorm.DB
	fakeTM  *FakeTenantManager
	cleanup func()
}

// setupTestEnv creates a complete test environment for handler tests.
func setupTestEnv() *testEnv {
	gin.SetMode(gin.TestMode)

	db, cleanup := setupTestDB()
	repo := NewRepository(db)
	logger, _ := zap.NewDevelopment()
	fakeTM := &FakeTenantManager{}
	svc := NewService(repo, fakeTM, logger)
	handler := NewHandler(svc)

	r := gin.New()
	rg := r.Group("/api/v1/platform")
	{
		companies := rg.Group("/companies")
		{
			companies.POST("/:id/terminate", handler.Terminate)
			companies.POST("/", handler.Create)
			companies.GET("/:id", handler.GetByID)
		}
	}

	return &testEnv{
		router:  r,
		db:      db,
		fakeTM:  fakeTM,
		cleanup: func() { cleanup(); logger.Sync() },
	}
}

// =========================================================================
// Handler Tests — Terminate Endpoint
// =========================================================================

func TestHandler_Terminate_Success(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	// Create a company directly in DB
	company := createTestCompany(env.db, "To Be Terminated")

	env.fakeTM.DropTenantDBFunc = func(companyID string) error {
		return nil
	}
	env.fakeTM.RemoveTenantConnFunc = func(companyID string) error {
		return nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+company.ID.String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool            `json:"success"`
		Data    CompanyResponse `json:"data"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Data.Status)
	}
	if resp.Message == "" {
		t.Error("expected a message")
	}
}

func TestHandler_Terminate_NotFound(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+uuid.New().String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Terminate_AlreadyTerminated(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	company := createTestCompany(env.db, "Already Terminated")
	company.Status = CompanyStatusTerminated
	if err := env.db.Save(company).Error; err != nil {
		t.Fatalf("failed to update company status: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+company.ID.String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Terminate_InvalidUUID(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/not-a-uuid/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}
