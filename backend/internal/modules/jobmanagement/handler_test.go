package jobmanagement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =========================================================================
// Handler Tests — verify parameter extraction and request flow
// =========================================================================

func newTestHandler() (*Handler, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)
	return handler, func() { cleanup(); logger.Sync() }
}

// =========================================================================
// Job Title Sub Handler Tests (9.2) — Verify c.Param("id") extraction
// =========================================================================

func TestHandler_CreateJobTitleSub_Success(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	// First create a job title
	titleID := uuid.New().String()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set path parameter
	c.Params = []gin.Param{{Key: "id", Value: titleID}}

	// Set request body
	body := `{"job_management_title_id":"` + titleID + `","name":"Test Sub","status":1}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateJobTitleSub(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool               `json:"success"`
		Data    JobTitleSubResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if resp.Data.Name != "Test Sub" {
		t.Errorf("expected name 'Test Sub', got '%s'", resp.Data.Name)
	}
}

func TestHandler_CreateJobTitleSub_InvalidJSON(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}

	// Invalid JSON
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateJobTitleSub(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestHandler_CreateJobTitleSub_MissingParam(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Intentionally NOT setting any path params to test empty param extraction

	body := `{"job_management_title_id":"` + uuid.New().String() + `","name":"Test Sub"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateJobTitleSub(c)

	// Empty path param causes UUID parse failure — verifies handler correctly
	// passes c.Param("id") to the service which validates it
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for empty path param, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListJobTitleSubs_Success(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	titleID := uuid.New().String()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: titleID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	handler.ListJobTitleSubs(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                   `json:"success"`
		Data    []JobTitleSubResponse  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	// Should be empty list since no subs exist for this title
	if resp.Data == nil {
		t.Fatal("expected non-nil data array")
	}
}

func TestHandler_ListJobTitleSubs_WithSubs(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	// Create a title and sub through the handler itself (shares the same DB)
	wCreate := httptest.NewRecorder()
	cCreate, _ := gin.CreateTestContext(wCreate)
	body := `{"name":"Handler Title"}`
	cCreate.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	cCreate.Request.Header.Set("Content-Type", "application/json")
	handler.CreateJobTitle(cCreate)

	var createResp struct {
		Success bool             `json:"success"`
		Data    JobTitleResponse `json:"data"`
	}
	if err := json.Unmarshal(wCreate.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}
	titleID := createResp.Data.ID

	// Create a sub through the handler
	wSub := httptest.NewRecorder()
	cSub, _ := gin.CreateTestContext(wSub)
	cSub.Params = []gin.Param{{Key: "id", Value: titleID}}
	subBody := `{"job_management_title_id":"` + titleID + `","name":"Handler Sub"}`
	cSub.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(subBody))
	cSub.Request.Header.Set("Content-Type", "application/json")
	handler.CreateJobTitleSub(cSub)

	if wSub.Code != http.StatusCreated {
		t.Fatalf("failed to create sub: %d - %s", wSub.Code, wSub.Body.String())
	}

	// Now test the handler's ListJobTitleSubs
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: titleID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	handler.ListJobTitleSubs(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		Success bool                   `json:"success"`
		Data    []JobTitleSubResponse  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 sub, got %d", len(resp.Data))
	}
	if resp.Data[0].Name != "Handler Sub" {
		t.Errorf("expected 'Handler Sub', got '%s'", resp.Data[0].Name)
	}
}

// =========================================================================
// Job Title Handler Tests (9.1) — Verify c.Param("id") extraction
// =========================================================================

func TestHandler_GetJobTitleByID(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	// Create a title first
	titleID := uuid.New().String()
	wCreate := httptest.NewRecorder()
	cCreate, _ := gin.CreateTestContext(wCreate)
	cCreate.Params = []gin.Param{{Key: "id", Value: titleID}}
	body := `{"name":"GET Test Title"}`
	cCreate.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	cCreate.Request.Header.Set("Content-Type", "application/json")
	handler.CreateJobTitle(cCreate)

	// Parse created ID from the response
	var createResp struct {
		Success bool             `json:"success"`
		Data    JobTitleResponse `json:"data"`
	}
	if err := json.Unmarshal(wCreate.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}

	// Now test GET by ID
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: createResp.Data.ID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	handler.GetJobTitleByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Data    JobTitleResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Data.Name != "GET Test Title" {
		t.Errorf("expected 'GET Test Title', got '%s'", resp.Data.Name)
	}
}

// =========================================================================
// Job Score Handler Tests (9.17) — Verify c.Param("orgId") extraction
// =========================================================================

func TestHandler_GetJobScoreByOrganization(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	orgID := uuid.New().String()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "orgId", Value: orgID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	handler.GetJobScoreByOrganization(c)

	// Should return 404 since no score exists yet
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent score, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpsertJobScore(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	orgID := uuid.New().String()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "orgId", Value: orgID}}
	body := `{"job_value_with_financial":2500,"job_value_without_financial":1800,"has_financial_authority":true}`
	c.Request, _ = http.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpsertJobScore(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for upsert, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Data    JobScoreResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if resp.Data.JobValueWithFinancial != 2500 {
		t.Errorf("expected 2500, got %d", resp.Data.JobValueWithFinancial)
	}

	// Verify can fetch the newly created score
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Params = []gin.Param{{Key: "orgId", Value: orgID}}
	c2.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	handler.GetJobScoreByOrganization(c2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status 200 after upsert, got %d", w2.Code)
	}
}

// =========================================================================
// Handler Tests — Bad Request / Validation
// =========================================================================

func TestHandler_CreateJobTitle_MissingName(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"descriptions":"missing name field"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateJobTitle(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing name, got %d", w.Code)
	}
}

func TestHandler_DeleteJobTitle_NotFound(t *testing.T) {
	handler, cleanup := newTestHandler()
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/", nil)

	handler.DeleteJobTitle(c)

	// Service treats delete of non-existent title as a no-op (return success)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for delete (no-op on non-existent), got %d", w.Code)
	}
}
