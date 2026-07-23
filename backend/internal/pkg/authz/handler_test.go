package authz

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupTestRouter creates a Gin router with RBAC routes registered
// for handler test purposes.
func setupTestRouter() (*gin.Engine, func()) {
	_, _, _, _, handler, _, cleanup := setupTestEnv()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1/platform")
	RegisterRoutes(rg, handler)

	return r, cleanup
}

// =========================================================================
// Role Handler Tests
// =========================================================================

func TestHandler_CreateRole_Success(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"name": "Custom Role", "slug": "custom-role"}`
	req := httptest.NewRequest("POST", "/api/v1/platform/rbac/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool          `json:"success"`
		Data    RoleResponse  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success true")
	}
	if resp.Data.Name != "Custom Role" {
		t.Errorf("expected name 'Custom Role', got %q", resp.Data.Name)
	}
}

func TestHandler_CreateRole_ValidationError(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"slug": "test"}`
	req := httptest.NewRequest("POST", "/api/v1/platform/rbac/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateRole_InvalidJSON(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	body := `{invalid json`
	req := httptest.NewRequest("POST", "/api/v1/platform/rbac/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_ListRoles_Success(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/platform/rbac/roles", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool           `json:"success"`
		Data    []RoleResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success true")
	}
	if len(resp.Data) != 4 {
		t.Errorf("expected 4 roles, got %d", len(resp.Data))
	}
}

func TestHandler_GetRole_Success(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "manager").String()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: roleID}}
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler.GetRole(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    RoleResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Slug != "manager" {
		t.Errorf("expected slug 'manager', got %q", resp.Data.Slug)
	}
	if len(resp.Data.Permissions) == 0 {
		t.Error("expected manager to have permissions")
	}
}

func TestHandler_GetRole_NotFound(t *testing.T) {
	_, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "00000000-0000-0000-0000-000000009999"}}
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler.GetRole(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateRole_Success(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "manager").String()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: roleID}}
	body := `{"name": "Updated Manager"}`
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateRole(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    RoleResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Name != "Updated Manager" {
		t.Errorf("expected name 'Updated Manager', got %q", resp.Data.Name)
	}
}

func TestHandler_DeleteRole_Custom(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	role := createTestRole(db, "Custom", "custom", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: role.ID.String()}}
	c.Request = httptest.NewRequest("DELETE", "/", nil)

	handler.DeleteRole(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteRole_System_Fails(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "super_admin").String()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: roleID}}
	c.Request = httptest.NewRequest("DELETE", "/", nil)

	handler.DeleteRole(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Permission Handler Tests
// =========================================================================

func TestHandler_CreatePermission_Success(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"resource": "custom", "action": "export"}`
	req := httptest.NewRequest("POST", "/api/v1/platform/rbac/permissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool               `json:"success"`
		Data    PermissionResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Resource != "custom" {
		t.Errorf("expected resource 'custom', got %q", resp.Data.Resource)
	}
}

func TestHandler_CreatePermission_ValidationError(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"action": "test"}`
	req := httptest.NewRequest("POST", "/api/v1/platform/rbac/permissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListPermissions_Success(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/platform/rbac/permissions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                 `json:"success"`
		Data    []PermissionResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("expected success true")
	}
	if len(resp.Data) == 0 {
		t.Error("expected at least 1 permission")
	}
}

func TestHandler_DeletePermission_Custom(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	perm := createTestPermission(db, "custom", "test")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: perm.ID.String()}}
	c.Request = httptest.NewRequest("DELETE", "/", nil)

	handler.DeletePermission(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Role-Permission Assignment Handler Tests
// =========================================================================

func TestHandler_AssignPermission_Success(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee")
	perm := createTestPermission(db, "assign-test", "view")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: roleID.String()}}
	body := `{"permission_id": "` + perm.ID.String() + `"}`
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AssignPermission(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_AssignPermission_InvalidBody(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: roleID.String()}}
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AssignPermission(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_RevokePermission_Success(t *testing.T) {
	db, _, _, _, handler, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee")
	perm := createTestPermission(db, "revoke-test", "view")
	db.Create(&RbacRolePermission{RoleID: roleID, PermissionID: perm.ID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: roleID.String()},
		{Key: "permissionId", Value: perm.ID.String()},
	}
	c.Request = httptest.NewRequest("DELETE", "/", nil)

	handler.RevokePermission(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}
