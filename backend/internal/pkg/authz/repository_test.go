package authz

import (
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Role CRUD
// =========================================================================

func TestRepository_CreateRole_Success(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{
		Name: "Test Role",
		Slug: "test-role",
	}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	if role.ID == uuid.Nil {
		t.Error("expected role ID to be set")
	}

	// Verify in DB
	var saved RbacRole
	if err := db.First(&saved, "id = ?", role.ID).Error; err != nil {
		t.Fatalf("failed to find created role: %v", err)
	}
	if saved.Name != "Test Role" {
		t.Errorf("expected name 'Test Role', got %q", saved.Name)
	}
}

func TestRepository_CreateRole_DuplicateSlug(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role1 := &RbacRole{Name: "Role 1", Slug: "duplicate"}
	if err := repo.CreateRole(role1); err != nil {
		t.Fatalf("first CreateRole failed: %v", err)
	}

	role2 := &RbacRole{Name: "Role 2", Slug: "duplicate"}
	if err := repo.CreateRole(role2); err == nil {
		t.Error("expected error for duplicate slug")
	}
}

func TestRepository_FindRoleByID_Found(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Test Role", Slug: "test-role"}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	found, err := repo.FindRoleByID(role.ID)
	if err != nil {
		t.Fatalf("FindRoleByID failed: %v", err)
	}
	if found.Name != "Test Role" {
		t.Errorf("expected name 'Test Role', got %q", found.Name)
	}
}

func TestRepository_FindRoleByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	_, err := repo.FindRoleByID(uuid.New())
	if err == nil {
		t.Error("expected error for non-existent role")
	}
}

func TestRepository_FindRoleBySlug_Found(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Test Role", Slug: "test-role"}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	found, err := repo.FindRoleBySlug("test-role")
	if err != nil {
		t.Fatalf("FindRoleBySlug failed: %v", err)
	}
	if found.Name != "Test Role" {
		t.Errorf("expected name 'Test Role', got %q", found.Name)
	}
}

func TestRepository_FindRoleBySlug_NotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	_, err := repo.FindRoleBySlug("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent slug")
	}
}

func TestRepository_FindAllRoles_Empty(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	roles, err := repo.FindAllRoles()
	if err != nil {
		t.Fatalf("FindAllRoles failed: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected 0 roles, got %d", len(roles))
	}
}

func TestRepository_FindAllRoles_Multiple(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	repo.CreateRole(&RbacRole{Name: "Beta", Slug: "beta"})
	repo.CreateRole(&RbacRole{Name: "Alpha", Slug: "alpha"})

	roles, err := repo.FindAllRoles()
	if err != nil {
		t.Fatalf("FindAllRoles failed: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
	// Should be ordered by name ASC
	if roles[0].Name != "Alpha" || roles[1].Name != "Beta" {
		t.Errorf("expected [Alpha, Beta], got [%s, %s]", roles[0].Name, roles[1].Name)
	}
}

func TestRepository_UpdateRole(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Old Name", Slug: "old-slug"}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	role.Name = "New Name"
	if err := repo.UpdateRole(role); err != nil {
		t.Fatalf("UpdateRole failed: %v", err)
	}

	var saved RbacRole
	db.First(&saved, "id = ?", role.ID)
	if saved.Name != "New Name" {
		t.Errorf("expected name 'New Name', got %q", saved.Name)
	}
}

func TestRepository_DeleteRole_NonSystem(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Custom", Slug: "custom", IsSystem: false}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	if err := repo.DeleteRole(role.ID); err != nil {
		t.Fatalf("DeleteRole failed: %v", err)
	}

	// Verify deleted
	_, err := repo.FindRoleByID(role.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestRepository_DeleteRole_SystemRole_Fails(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "System", Slug: "system", IsSystem: true}
	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	err := repo.DeleteRole(role.ID)
	if err == nil {
		t.Error("expected error when deleting system role")
	}
}

func TestRepository_DeleteRole_NotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	err := repo.DeleteRole(uuid.New())
	if err == nil {
		t.Error("expected error for non-existent role")
	}
}

// =========================================================================
// Permission CRUD
// =========================================================================

func TestRepository_CreatePermission_Success(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	perm := &RbacPermission{
		Resource: "test-resource",
		Action:   "view",
	}
	if err := repo.CreatePermission(perm); err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}
	if perm.ID == uuid.Nil {
		t.Error("expected permission ID to be set")
	}
}

func TestRepository_CreatePermission_Duplicate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	p1 := &RbacPermission{Resource: "test", Action: "view"}
	if err := repo.CreatePermission(p1); err != nil {
		t.Fatalf("first CreatePermission failed: %v", err)
	}

	p2 := &RbacPermission{Resource: "test", Action: "view"}
	if err := repo.CreatePermission(p2); err == nil {
		t.Error("expected error for duplicate resource+action")
	}
}

func TestRepository_FindPermissionByID_Found(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	perm := &RbacPermission{Resource: "test", Action: "view"}
	if err := repo.CreatePermission(perm); err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	found, err := repo.FindPermissionByID(perm.ID)
	if err != nil {
		t.Fatalf("FindPermissionByID failed: %v", err)
	}
	if found.Resource != "test" {
		t.Errorf("expected resource 'test', got %q", found.Resource)
	}
}

func TestRepository_FindPermissionByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	_, err := repo.FindPermissionByID(uuid.New())
	if err == nil {
		t.Error("expected error for non-existent permission")
	}
}

func TestRepository_FindAllPermissions(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	repo.CreatePermission(&RbacPermission{Resource: "zzz", Action: "view"})
	repo.CreatePermission(&RbacPermission{Resource: "aaa", Action: "create"})

	perms, err := repo.FindAllPermissions()
	if err != nil {
		t.Fatalf("FindAllPermissions failed: %v", err)
	}
	if len(perms) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(perms))
	}
	// Should be ordered by resource ASC, action ASC
	if perms[0].Resource != "aaa" || perms[1].Resource != "zzz" {
		t.Errorf("expected [aaa, zzz], got [%s, %s]", perms[0].Resource, perms[1].Resource)
	}
}

func TestRepository_DeletePermission_NonSystem(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	perm := &RbacPermission{Resource: "custom", Action: "test", IsSystem: false}
	if err := repo.CreatePermission(perm); err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	if err := repo.DeletePermission(perm.ID); err != nil {
		t.Fatalf("DeletePermission failed: %v", err)
	}

	_, err := repo.FindPermissionByID(perm.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestRepository_DeletePermission_System_Fails(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	perm := &RbacPermission{Resource: "system", Action: "test", IsSystem: true}
	if err := repo.CreatePermission(perm); err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	err := repo.DeletePermission(perm.ID)
	if err == nil {
		t.Error("expected error when deleting system permission")
	}
}

// =========================================================================
// Role-Permission Assignment
// =========================================================================

func TestRepository_AssignPermission_Success(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Custom Role", Slug: "custom"}
	perm := &RbacPermission{Resource: "custom", Action: "view"}

	if err := repo.CreateRole(role); err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}
	if err := repo.CreatePermission(perm); err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	if err := repo.AssignPermission(role.ID, perm.ID); err != nil {
		t.Fatalf("AssignPermission failed: %v", err)
	}
}

func TestRepository_AssignPermission_Duplicate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Custom Role", Slug: "custom"}
	perm := &RbacPermission{Resource: "custom", Action: "view"}

	repo.CreateRole(role)
	repo.CreatePermission(perm)

	// First assignment should succeed
	if err := repo.AssignPermission(role.ID, perm.ID); err != nil {
		t.Fatalf("first AssignPermission failed: %v", err)
	}

	// Second assignment should fail (duplicate PK)
	if err := repo.AssignPermission(role.ID, perm.ID); err == nil {
		t.Error("expected error for duplicate assignment")
	}
}

func TestRepository_RevokePermission_Success(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Custom Role", Slug: "custom"}
	perm := &RbacPermission{Resource: "custom", Action: "view"}

	repo.CreateRole(role)
	repo.CreatePermission(perm)

	// Assign
	if err := repo.AssignPermission(role.ID, perm.ID); err != nil {
		t.Fatalf("AssignPermission failed: %v", err)
	}

	// Revoke
	if err := repo.RevokePermission(role.ID, perm.ID); err != nil {
		t.Fatalf("RevokePermission failed: %v", err)
	}

	// Verify empty
	rps, err := repo.FindRolePermissions(role.ID)
	if err != nil {
		t.Fatalf("FindRolePermissions failed: %v", err)
	}
	if len(rps) != 0 {
		t.Errorf("expected 0 role permissions, got %d", len(rps))
	}
}

func TestRepository_FindRolePermissions(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	role := &RbacRole{Name: "Custom Role", Slug: "custom"}
	perm1 := &RbacPermission{Resource: "r1", Action: "view"}
	perm2 := &RbacPermission{Resource: "r2", Action: "create"}

	repo.CreateRole(role)
	repo.CreatePermission(perm1)
	repo.CreatePermission(perm2)

	repo.AssignPermission(role.ID, perm1.ID)
	repo.AssignPermission(role.ID, perm2.ID)

	rps, err := repo.FindRolePermissions(role.ID)
	if err != nil {
		t.Fatalf("FindRolePermissions failed: %v", err)
	}
	if len(rps) != 2 {
		t.Errorf("expected 2 role permissions, got %d", len(rps))
	}
}

// =========================================================================
// AutoMigrate
// =========================================================================

func TestRepository_AutoMigrate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(db)

	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	// Verify tables exist by performing operations
	role := &RbacRole{Name: "Post-Migrate", Slug: "post-migrate"}
	if err := repo.CreateRole(role); err != nil {
		t.Errorf("CreateRole after AutoMigrate failed: %v", err)
	}
}
