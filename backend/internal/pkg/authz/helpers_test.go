package authz

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all RBAC models.
// Returns the GORM DB and a cleanup function.
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&RbacRole{},
		&RbacPermission{},
		&RbacRolePermission{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

// setupTestEnv creates a complete test environment with DB, repository, enforcer, service, handler and logger.
// The enforcer is already loaded with seeded default permissions.
func setupTestEnv() (*gorm.DB, *Repository, *Enforcer, *Service, *Handler, *zap.Logger, func()) {
	db, cleanup := setupTestDB()
	logger := zap.NewNop()
	repo := NewRepository(db)

	// Create enforcer from DB (seeds defaults + loads from DB)
	enforcer, err := NewEnforcerFromDB(db)
	if err != nil {
		cleanup()
		panic(fmt.Sprintf("failed to create enforcer from DB: %v", err))
	}

	service := NewService(repo, enforcer, logger)
	handler := NewHandler(service)

	return db, repo, enforcer, service, handler, logger, cleanup
}

// createTestRole inserts a custom role for testing and returns it.
func createTestRole(db *gorm.DB, name, slug string, parentID *uuid.UUID) *RbacRole {
	role := &RbacRole{
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
	}
	if err := db.Create(role).Error; err != nil {
		panic(fmt.Sprintf("failed to create test role: %v", err))
	}
	return role
}

// createTestPermission inserts a custom permission for testing and returns it.
func createTestPermission(db *gorm.DB, resource, action string) *RbacPermission {
	perm := &RbacPermission{
		Resource: resource,
		Action:   action,
	}
	if err := db.Create(perm).Error; err != nil {
		panic(fmt.Sprintf("failed to create test permission: %v", err))
	}
	return perm
}

// getRoleID retrieves the UUID of a role by slug from the database.
func getRoleID(db *gorm.DB, slug string) uuid.UUID {
	var role RbacRole
	if err := db.Where("slug = ?", slug).First(&role).Error; err != nil {
		panic(fmt.Sprintf("failed to find role %s: %v", slug, err))
	}
	return role.ID
}


