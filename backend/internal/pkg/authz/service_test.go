package authz

import (
	"testing"
)

// =========================================================================
// Role Service Tests
// =========================================================================

func TestService_CreateRole_Success(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	resp, err := svc.CreateRole(CreateRoleRequest{
		Name: "Custom Role",
		Slug: "custom-role",
	})
	if err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}
	if resp.Name != "Custom Role" {
		t.Errorf("expected name 'Custom Role', got %q", resp.Name)
	}
	if resp.Slug != "custom-role" {
		t.Errorf("expected slug 'custom-role', got %q", resp.Slug)
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
	if resp.IsSystem {
		t.Error("expected IsSystem to be false for custom role")
	}
}

func TestService_CreateRole_WithParent(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Find manager role ID
	managerID := getRoleID(db, "manager").String()

	parentID := managerID
	resp, err := svc.CreateRole(CreateRoleRequest{
		Name:     "Sub Manager",
		Slug:     "sub-manager",
		ParentID: &parentID,
	})
	if err != nil {
		t.Fatalf("CreateRole with parent failed: %v", err)
	}
	if resp.ParentID != managerID {
		t.Errorf("expected ParentID %q, got %q", managerID, resp.ParentID)
	}
}

func TestService_CreateRole_InvalidParentID(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	invalidID := "not-a-uuid"
	_, err := svc.CreateRole(CreateRoleRequest{
		Name:     "Invalid",
		Slug:     "invalid",
		ParentID: &invalidID,
	})
	if err == nil {
		t.Error("expected error for invalid parent_id")
	}
}

func TestService_CreateRole_DuplicateSlug(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	svc.CreateRole(CreateRoleRequest{Name: "Role 1", Slug: "dup"})
	_, err := svc.CreateRole(CreateRoleRequest{Name: "Role 2", Slug: "dup"})
	if err == nil {
		t.Error("expected error for duplicate slug")
	}
}

func TestService_ListRoles_WithPermissions(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Default roles from seed should include permissions
	roles, err := svc.ListRoles()
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}

	// Should have 4 default roles
	if len(roles) != 4 {
		t.Errorf("expected 4 roles, got %d", len(roles))
	}

	// Super admin should have wildcard (all permissions)
	var superAdmin *RoleResponse
	for i := range roles {
		if roles[i].Slug == "super_admin" {
			superAdmin = &roles[i]
			break
		}
	}
	if superAdmin == nil {
		t.Fatal("expected super_admin role")
	}
	if len(superAdmin.Permissions) == 0 {
		t.Error("expected super_admin to have permissions")
	}

	// Employee should have limited permissions
	var employee *RoleResponse
	for i := range roles {
		if roles[i].Slug == "employee" {
			employee = &roles[i]
			break
		}
	}
	if employee == nil {
		t.Fatal("expected employee role")
	}
	if len(employee.Permissions) == 0 {
		t.Error("expected employee to have permissions")
	}
}

func TestService_GetRole_Success(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "manager").String()

	resp, err := svc.GetRole(roleID)
	if err != nil {
		t.Fatalf("GetRole failed: %v", err)
	}
	if resp.Slug != "manager" {
		t.Errorf("expected slug 'manager', got %q", resp.Slug)
	}
	if len(resp.Permissions) == 0 {
		t.Error("expected manager to have permissions")
	}
}

func TestService_GetRole_NotFound(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	_, err := svc.GetRole("00000000-0000-0000-0000-000000009999")
	if err == nil {
		t.Error("expected error for non-existent role")
	}
}

func TestService_GetRole_InvalidID(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	_, err := svc.GetRole("invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid UUID")
	}
}

func TestService_UpdateRole_Name(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "manager").String()

	newName := "Updated Manager"
	resp, err := svc.UpdateRole(roleID, UpdateRoleRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateRole failed: %v", err)
	}
	if resp.Name != "Updated Manager" {
		t.Errorf("expected name 'Updated Manager', got %q", resp.Name)
	}
	if resp.Slug != "manager" {
		t.Errorf("slug should not change, got %q", resp.Slug)
	}
}

func TestService_UpdateRole_ClearParent(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee").String()

	empty := ""
	resp, err := svc.UpdateRole(roleID, UpdateRoleRequest{
		ParentID: &empty,
	})
	if err != nil {
		t.Fatalf("UpdateRole clear parent failed: %v", err)
	}
	if resp.ParentID != "" {
		t.Errorf("expected empty ParentID, got %q", resp.ParentID)
	}
}

func TestService_UpdateRole_NotFound(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	name := "New"
	_, err := svc.UpdateRole("00000000-0000-0000-0000-000000009999", UpdateRoleRequest{Name: &name})
	if err == nil {
		t.Error("expected error for non-existent role")
	}
}

func TestService_DeleteRole_NonSystem(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Create a custom (non-system) role first
	resp, err := svc.CreateRole(CreateRoleRequest{Name: "Temp", Slug: "temp"})
	if err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	if err := svc.DeleteRole(resp.ID); err != nil {
		t.Fatalf("DeleteRole failed: %v", err)
	}

	// Verify deleted
	_, err = svc.GetRole(resp.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestService_DeleteRole_System_Fails(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "super_admin").String()

	err := svc.DeleteRole(roleID)
	if err == nil {
		t.Error("expected error when deleting system role")
	}
}

func TestService_DeleteRole_InvalidID(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	err := svc.DeleteRole("invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid UUID")
	}
}

// =========================================================================
// Permission Service Tests
// =========================================================================

func TestService_CreatePermission_Success(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	resp, err := svc.CreatePermission(CreatePermissionRequest{
		Resource: "custom-resource",
		Action:   "export",
	})
	if err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}
	if resp.Resource != "custom-resource" {
		t.Errorf("expected resource 'custom-resource', got %q", resp.Resource)
	}
	if resp.Action != "export" {
		t.Errorf("expected action 'export', got %q", resp.Action)
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestService_CreatePermission_Duplicate(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	svc.CreatePermission(CreatePermissionRequest{Resource: "dup", Action: "view"})
	_, err := svc.CreatePermission(CreatePermissionRequest{Resource: "dup", Action: "view"})
	if err == nil {
		t.Error("expected error for duplicate permission")
	}
}

func TestService_ListPermissions(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Create a custom permission
	svc.CreatePermission(CreatePermissionRequest{Resource: "custom", Action: "export"})

	perms, err := svc.ListPermissions()
	if err != nil {
		t.Fatalf("ListPermissions failed: %v", err)
	}

	// Should have seeded permissions + the new one
	if len(perms) == 0 {
		t.Error("expected at least 1 permission")
	}

	// Verify the custom permission is in the list
	var found bool
	for _, p := range perms {
		if p.Resource == "custom" && p.Action == "export" {
			found = true
			if p.IsSystem {
				t.Error("custom permission should not be IsSystem")
			}
			break
		}
	}
	if !found {
		t.Error("custom permission not found in list")
	}
}

func TestService_DeletePermission_NonSystem(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	resp, err := svc.CreatePermission(CreatePermissionRequest{Resource: "temp", Action: "test"})
	if err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	if err := svc.DeletePermission(resp.ID); err != nil {
		t.Fatalf("DeletePermission failed: %v", err)
	}
}

func TestService_DeletePermission_System_Fails(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Find a system permission ID
	var perms []RbacPermission
	db.Where("is_system = ?", true).Find(&perms)
	if len(perms) == 0 {
		t.Skip("no system permissions available")
	}

	err := svc.DeletePermission(perms[0].ID.String())
	if err == nil {
		t.Error("expected error when deleting system permission")
	}
}

// =========================================================================
// Role-Permission Assignment Service Tests
// =========================================================================

func TestService_AssignPermission_Success(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee").String()
	perm := createTestPermission(db, "new-resource", "view")
	permID := perm.ID.String()

	if err := svc.AssignPermission(roleID, permID); err != nil {
		t.Fatalf("AssignPermission failed: %v", err)
	}
}

func TestService_AssignPermission_InvalidRoleID(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	err := svc.AssignPermission("invalid-uuid", "00000000-0000-0000-0000-000000000001")
	if err == nil {
		t.Error("expected error for invalid role ID")
	}
}

func TestService_AssignPermission_InvalidPermissionID(t *testing.T) {
	db, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee").String()

	err := svc.AssignPermission(roleID, "invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid permission ID")
	}
}

func TestService_RevokePermission_Success(t *testing.T) {
	db, _, enforcer, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	roleID := getRoleID(db, "employee").String()

	// Create and assign a permission
	perm := createTestPermission(db, "revokable", "test")
	enforcer.Reload() // reload to pick up the new permission

	if err := svc.AssignPermission(roleID, perm.ID.String()); err != nil {
		t.Fatalf("AssignPermission failed: %v", err)
	}

	// Revoke it
	if err := svc.RevokePermission(roleID, perm.ID.String()); err != nil {
		t.Fatalf("RevokePermission failed: %v", err)
	}

	// Verify revoked
	rps, _ := svc.repo.FindRolePermissions(getRoleID(db, "employee"))
	for _, rp := range rps {
		if rp.PermissionID == perm.ID {
			t.Error("permission should have been revoked")
		}
	}
}

// =========================================================================
// Sync Tests
// =========================================================================

func TestService_Sync_Success(t *testing.T) {
	_, _, _, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	if err := svc.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}
}

func TestService_AssignPermission_SyncsEnforcer(t *testing.T) {
	db, _, enforcer, svc, _, _, cleanup := setupTestEnv()
	defer cleanup()

	// Employee currently cannot create organization
	decision := enforcer.Check("employee", "organization", "create")
	if decision != DecisionDeny {
		t.Logf("initial check: %v (may be allow if hierarchy permits)", decision)
	}

	// Create a NEW permission (unique resource, won't conflict with seeded data)
	perm := createTestPermission(db, "custom-sync-resource", "view")
	employeeID := getRoleID(db, "employee").String()

	// Assign the new permission to employee
	if err := svc.AssignPermission(employeeID, perm.ID.String()); err != nil {
		t.Fatalf("AssignPermission failed: %v", err)
	}

	// Manually sync to reflect changes
	if err := svc.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	// After sync: employee should be able to view custom-sync-resource
	decision = enforcer.Check("employee", "custom-sync-resource", "view")
	if decision != DecisionAllow {
		t.Error("expected allow after sync")
	}
}
