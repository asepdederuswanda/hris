package authz

import (
	"testing"
)

// =========================================================================
// NewEnforcer (non-DB) — Default Policy Tests
// =========================================================================

func TestNewEnforcer_NonDB_HasDefaultPolicies(t *testing.T) {
	e := NewEnforcer()

	// Super admin: wildcard access to everything
	assertAllow(t, e, "super_admin", "company", "delete")
	assertAllow(t, e, "super_admin", "employee", "view")
	assertAllow(t, e, "super_admin", "nonexistent_resource", "create")

	// Company admin: platform view-only + tenant full
	assertAllow(t, e, "company_admin", "company", "view")
	assertDeny(t, e, "company_admin", "company", "delete")       // platform: view-only
	assertAllow(t, e, "company_admin", "employee", "delete")      // tenant: full (employee:*)
	assertAllow(t, e, "company_admin", "organization", "create")  // tenant: full (organization:*)

	// Manager: view/create/update (no delete)
	assertAllow(t, e, "manager", "organization", "view")
	assertAllow(t, e, "manager", "employee", "update")
	assertDeny(t, e, "manager", "employee", "delete")   // no delete (employee: view,create,update)

	// Manager does NOT have 'company' in policies → traverses to company_admin who has company:view → ALLOW
	assertAllow(t, e, "manager", "company", "view")

	// Employee: view-only
	assertAllow(t, e, "employee", "organization", "view")
	assertDeny(t, e, "employee", "employee", "create")  // no create (employee: view)

	// Employee does NOT have 'company' in policies → traverses to manager (no company) → company_admin (company:view) → ALLOW
	assertAllow(t, e, "employee", "company", "view")

	// Unknown role: no policies, no hierarchy → DENY
	assertDeny(t, e, "unknown_role", "organization", "view")
}

func TestNewEnforcer_NonDB_RoleHierarchy(t *testing.T) {
	e := NewEnforcer()

	// Manager has no explicit policy for 'payroll' → traverses to company_admin (payroll:*) → ALLOW
	assertAllow(t, e, "manager", "payroll", "view")

	// Employee → Manager (no module) → CompanyAdmin (no module) → SuperAdmin (*:*) → ALLOW
	assertAllow(t, e, "employee", "module", "view")
}

func TestNewEnforcer_NonDB_ActionWildcard(t *testing.T) {
	e := NewEnforcer()

	// Company admin has organization:* → any action allowed
	assertAllow(t, e, "company_admin", "organization", "view")
	assertAllow(t, e, "company_admin", "organization", "create")
	assertAllow(t, e, "company_admin", "organization", "update")
	assertAllow(t, e, "company_admin", "organization", "delete")
}

func TestNewEnforcer_NonDB_ManagerHierarchyInheritsPlatform(t *testing.T) {
	e := NewEnforcer()

	// Manager doesn't have explicit platform policies, inherits from company_admin
	assertAllow(t, e, "manager", "company", "view")
	assertAllow(t, e, "manager", "user", "view")
	assertAllow(t, e, "manager", "license", "view")

	// Manager still has explicit policies for tenant resources
	assertAllow(t, e, "manager", "employee", "view")
	assertAllow(t, e, "manager", "competency", "view")
	assertAllow(t, e, "manager", "jobmanagement", "view")
}

func TestNewEnforcer_NonDB_EmployeeReadOnly(t *testing.T) {
	e := NewEnforcer()

	// Employee can only view
	assertAllow(t, e, "employee", "organization", "view")
	assertAllow(t, e, "employee", "employee", "view")
	assertDeny(t, e, "employee", "employee", "create")
	assertDeny(t, e, "employee", "employee", "update")
	assertDeny(t, e, "employee", "employee", "delete")
}

func TestNewEnforcer_NonDB_ExplicitDenyBlocksHierarchy(t *testing.T) {
	e := NewEnforcer()

	// Manager has employee: view,create,update (explicit). 'delete' is denied at manager level.
	// Even though company_admin has employee:*, the explicit policy at manager blocks hierarchy traversal.
	assertDeny(t, e, "manager", "employee", "delete")

	// Employee has employee:view (explicit). 'create' is denied at employee level.
	// Even though manager/company_admin have create, the explicit policy blocks hierarchy.
	assertDeny(t, e, "employee", "employee", "create")
}

// =========================================================================
// MustCheck Tests
// =========================================================================

func TestMustCheck_Allowed(t *testing.T) {
	e := NewEnforcer()
	err := e.MustCheck("super_admin", "anything", "delete")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestMustCheck_Denied(t *testing.T) {
	e := NewEnforcer()
	// employee has employee:view (explicit), 'delete' not in list → DENY
	err := e.MustCheck("employee", "employee", "delete")
	if err == nil {
		t.Error("expected error for denied permission, got nil")
	}
}

// =========================================================================
// AddPolicy Tests
// =========================================================================

func TestAddPolicy_DynamicAddition(t *testing.T) {
	e := NewEnforcer()

	// Before adding: employee has no policy for 'module', traverses to super_admin → ALLOW
	assertAllow(t, e, "employee", "module", "view")

	// Add a RESTRICTIVE policy to override inheritance
	e.AddPolicy(RoleEmployee, "module", "") // empty = no actions allowed for this resource

	// Now: employee has explicit policy for 'module' with empty actions → DENY
	assertDeny(t, e, "employee", "module", "view")
}

func TestAddPolicy_OverwritesExisting(t *testing.T) {
	e := NewEnforcer()

	// Employee has organization:view (only view via explicit or hierarchy)
	// Let's make employee have explicit policy to test overwrite
	e.AddPolicy(RoleEmployee, "organization", "view")
	assertAllow(t, e, "employee", "organization", "view")
	assertDeny(t, e, "employee", "organization", "create")

	// Overwrite with wildcard
	e.AddPolicy(RoleEmployee, "organization", "*")

	// Now create should be allowed
	assertAllow(t, e, "employee", "organization", "create")
}

// =========================================================================
// NewEnforcerFromDB (DB-backed)
// =========================================================================

func TestNewEnforcerFromDB_SeedsDefaults(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enforcer")
	}

	// Check that roles were seeded
	var roleCount int64
	db.Model(&RbacRole{}).Count(&roleCount)
	if roleCount != 4 {
		t.Errorf("expected 4 roles, got %d", roleCount)
	}

	// Check that permissions were seeded
	var permCount int64
	db.Model(&RbacPermission{}).Count(&permCount)
	if permCount == 0 {
		t.Error("expected permissions to be seeded")
	}

	// Check that role_permissions were seeded
	var rpCount int64
	db.Model(&RbacRolePermission{}).Count(&rpCount)
	if rpCount == 0 {
		t.Error("expected role_permissions to be seeded")
	}
}

func TestNewEnforcerFromDB_IdempotentSeed(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	// First call: seeds defaults
	_, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("first NewEnforcerFromDB failed: %v", err)
	}

	// Second call: should not error, skip seeding
	_, err = NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("second NewEnforcerFromDB should not error: %v", err)
	}

	// Should still have exactly 4 roles
	var roleCount int64
	db.Model(&RbacRole{}).Count(&roleCount)
	if roleCount != 4 {
		t.Errorf("expected 4 roles after idempotent seed, got %d", roleCount)
	}
}

func TestNewEnforcerFromDB_SuperAdminAllResources(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// DB-backed super_admin has explicit permissions for all defined resources
	// (no "*" wildcard, but has every resource with every action)
	assertAllow(t, e, "super_admin", "company", "delete")
	assertAllow(t, e, "super_admin", "user", "create")
	assertAllow(t, e, "super_admin", "employee", "view")
	assertAllow(t, e, "super_admin", "module", "view")
	assertAllow(t, e, "super_admin", "organization", "view")
	assertAllow(t, e, "super_admin", "competency", "delete")
	assertAllow(t, e, "super_admin", "company", "activate")
	assertAllow(t, e, "super_admin", "module", "deactivate")
	assertAllow(t, e, "super_admin", "approval", "view")
	assertAllow(t, e, "super_admin", "monitoring", "view")
	assertAllow(t, e, "super_admin", "leave", "view")
	assertAllow(t, e, "super_admin", "attendance", "delete")
	assertAllow(t, e, "super_admin", "jobmanagement", "update")

	// Non-defined resources are NOT accessible via DB-backed enforcer
	// because super_admin has no parent to traverse to
	assertDeny(t, e, "super_admin", "nonexistent_resource", "view")
}

func TestNewEnforcerFromDB_CompanyAdminTenantFull(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// Company admin: platform view-only
	assertAllow(t, e, "company_admin", "company", "view")
	assertDeny(t, e, "company_admin", "company", "delete")
	assertAllow(t, e, "company_admin", "user", "view")
	assertDeny(t, e, "company_admin", "user", "create")

	// Company admin: tenant full access (all actions for tenant resources)
	assertAllow(t, e, "company_admin", "organization", "view")
	assertAllow(t, e, "company_admin", "organization", "delete")
	assertAllow(t, e, "company_admin", "employee", "view")
	assertAllow(t, e, "company_admin", "employee", "delete")
	assertAllow(t, e, "company_admin", "competency", "create")
	assertAllow(t, e, "company_admin", "jobmanagement", "update")
}

func TestNewEnforcerFromDB_ManagerEditNoDelete(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// Manager: view/create/update allowed
	assertAllow(t, e, "manager", "organization", "view")
	assertAllow(t, e, "manager", "employee", "create")
	assertAllow(t, e, "manager", "competency", "update")

	// Manager: delete denied (manager has explicit policies → intentional deny)
	assertDeny(t, e, "manager", "organization", "delete")
	assertDeny(t, e, "manager", "employee", "delete")
	assertDeny(t, e, "manager", "competency", "delete")

	// Manager inherits company.view from company_admin through hierarchy
	assertAllow(t, e, "manager", "company", "view")
	assertAllow(t, e, "manager", "user", "view")

	// Manager: delete for platform resources is still denied via company_admin's intentional deny
	assertDeny(t, e, "manager", "company", "delete")  // company_admin has company:view not delete
}

func TestNewEnforcerFromDB_EmployeeViewOnly(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// Employee: view allowed
	assertAllow(t, e, "employee", "organization", "view")
	assertAllow(t, e, "employee", "employee", "view")
	assertAllow(t, e, "employee", "attendance", "view")
	assertAllow(t, e, "employee", "leave", "view")
	assertAllow(t, e, "employee", "payroll", "view")
	assertAllow(t, e, "employee", "competency", "view")

	// Employee: mutation denied (employee has explicit policies → intentional deny)
	assertDeny(t, e, "employee", "employee", "create")
	assertDeny(t, e, "employee", "organization", "update")
	assertDeny(t, e, "employee", "organization", "delete")
	assertDeny(t, e, "employee", "competency", "create")
}

// =========================================================================
// Reload Tests
// =========================================================================

func TestReload_AfterAddingNewResource(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// Before: employee cannot access a custom resource
	assertDeny(t, e, "employee", "custom-resource", "view")

	// Add a NEW permission (unique resource name, won't conflict with seeded data)
	perm := createTestPermission(db, "custom-resource", "view")
	employeeID := getRoleID(db, "employee")
	if err := db.Create(&RbacRolePermission{RoleID: employeeID, PermissionID: perm.ID}).Error; err != nil {
		t.Fatalf("failed to assign permission: %v", err)
	}

	// Before reload: enforcer still has old policies
	assertDeny(t, e, "employee", "custom-resource", "view")

	// Reload
	if err := e.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// After reload: new policy active
	assertAllow(t, e, "employee", "custom-resource", "view")
}

func TestReload_AfterPermissionRevoked(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e, err := NewEnforcerFromDB(db)
	if err != nil {
		t.Fatalf("NewEnforcerFromDB failed: %v", err)
	}

	// Verify employee has 'organization.view' by default (through explicit policy)
	assertAllow(t, e, "employee", "organization", "view")

	// Revoke ALL permissions from employee
	employeeID := getRoleID(db, "employee")
	db.Where("role_id = ?", employeeID).Delete(&RbacRolePermission{})

	// Reload
	if err := e.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// After reload: employee has no explicit policies
	// Through hierarchy: employee → manager (has organization:view,create,update) → ALLOW
	assertAllow(t, e, "employee", "organization", "view")

	// But mutation still blocked at manager level (intentional deny)
	assertDeny(t, e, "employee", "organization", "delete")
}

func TestReload_WithoutDB_ReturnsError(t *testing.T) {
	e := NewEnforcer()

	err := e.Reload()
	if err == nil {
		t.Error("expected error when Reload called on non-DB enforcer")
	}
}

func TestNewEnforcerFromDB_InvalidDB_ReturnsError(t *testing.T) {
	// Use a closed DB connection
	db, cleanup := setupTestDB()
	cleanup() // close immediately

	_, err := NewEnforcerFromDB(db)
	if err == nil {
		t.Error("expected error when using closed DB")
	}
}

// =========================================================================
// ResourceFromPath Tests
// =========================================================================

func TestResourceFromPath_PlatformPaths(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/platform/companies", "company"},
		{"/api/v1/platform/users", "user"},
		{"/api/v1/platform/licenses", "license"},
		{"/api/v1/platform/modules", "module"},
		{"/api/v1/platform/monitoring", "monitoring"},
		{"/api/v1/platform/companies/123", "company"},
		{"/api/v1/platform/users/123/permissions", "user"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ResourceFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("ResourceFromPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestResourceFromPath_TenantPaths(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/tenant/organizations", "organization"},
		{"/api/v1/tenant/employees", "employee"},
		{"/api/v1/tenant/attendances", "attendance"},
		{"/api/v1/tenant/competencies", "competency"},
		{"/api/v1/tenant/job-management/titles", "jobmanagement"},
		{"/api/v1/tenant/employees/123", "employee"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ResourceFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("ResourceFromPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestResourceFromPath_UnknownPath(t *testing.T) {
	tests := []string{
		"",
		"/healthz",
		"/api/v1/unknown/companies",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			result := ResourceFromPath(path)
			if result != "" {
				t.Errorf("ResourceFromPath(%q) = %q, want empty string", path, result)
			}
		})
	}
}

// =========================================================================
// ActionFromMethod Tests
// =========================================================================

func TestActionFromMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "view"},
		{"POST", "create"},
		{"PUT", "update"},
		{"DELETE", "delete"},
		{"PATCH", "patch"},
		{"OPTIONS", "options"},
		{"HEAD", "head"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := ActionFromMethod(tt.method)
			if result != tt.expected {
				t.Errorf("ActionFromMethod(%q) = %q, want %q", tt.method, result, tt.expected)
			}
		})
	}
}

// =========================================================================
// singularize Tests
// =========================================================================

func TestSingularize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Irregular plurals
		{"companies", "company"},
		{"licenses", "license"},
		{"modules", "module"},
		{"users", "user"},
		{"monitoring", "monitoring"},
		{"tenants", "tenant"},
		{"organizations", "organization"},
		{"employees", "employee"},
		{"attendances", "attendance"},
		{"competencies", "competency"},
		{"job-management", "jobmanagement"},
		// -ies → -y
		{"policies", "policy"},
		{"categories", "category"},
		// -s → -
		{"departments", "department"},
		{"positions", "position"},
		{"companies", "company"}, // irregular takes precedence
		// Already singular
		{"company", "company"},
		{"employee", "employee"},
		{"data", "data"},
		// Empty
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := singularize(tt.input)
			if result != tt.expected {
				t.Errorf("singularize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// =========================================================================
// seedDefaults — Edge Cases
// =========================================================================

func TestSeedDefaults_SkipIfDataExists(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	e := &Enforcer{}

	// First seed
	if err := e.seedDefaults(db); err != nil {
		t.Fatalf("first seedDefaults failed: %v", err)
	}

	// Second seed should skip (data already exists)
	if err := e.seedDefaults(db); err != nil {
		t.Fatalf("second seedDefaults should not error: %v", err)
	}

	var count int64
	db.Model(&RbacRole{}).Count(&count)
	if count != 4 {
		t.Errorf("expected 4 roles, got %d", count)
	}
}

// =========================================================================
// Helpers
// =========================================================================

func assertAllow(t *testing.T, e *Enforcer, role, resource, action string) {
	t.Helper()
	result := e.Check(role, resource, action)
	if result != DecisionAllow {
		t.Errorf("Check(%q, %q, %q) = %v, want %v (allow)", role, resource, action, result, DecisionAllow)
	}
}

func assertDeny(t *testing.T, e *Enforcer, role, resource, action string) {
	t.Helper()
	result := e.Check(role, resource, action)
	if result != DecisionDeny {
		t.Errorf("Check(%q, %q, %q) = %v, want %v (deny)", role, resource, action, result, DecisionDeny)
	}
}
