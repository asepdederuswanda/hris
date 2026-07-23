package authz

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi CRUD RBAC di database platform.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository RBAC baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// =========================================================================
// Role CRUD
// =========================================================================

// CreateRole membuat role baru.
func (r *Repository) CreateRole(role *RbacRole) error {
	return r.db.Create(role).Error
}

// FindRoleByID mencari role berdasarkan UUID.
func (r *Repository) FindRoleByID(id uuid.UUID) (*RbacRole, error) {
	var role RbacRole
	if err := r.db.First(&role, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

// FindRoleBySlug mencari role berdasarkan slug.
func (r *Repository) FindRoleBySlug(slug string) (*RbacRole, error) {
	var role RbacRole
	if err := r.db.First(&role, "slug = ?", slug).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

// FindAllRoles mengembalikan semua roles dengan permission counts.
func (r *Repository) FindAllRoles() ([]RbacRole, error) {
	var roles []RbacRole
	if err := r.db.Preload("Permissions").Order("name ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// UpdateRole mengupdate data role.
func (r *Repository) UpdateRole(role *RbacRole) error {
	return r.db.Save(role).Error
}

// DeleteRole menghapus role (hanya jika bukan system role).
func (r *Repository) DeleteRole(id uuid.UUID) error {
	result := r.db.Where("id = ? AND is_system = ?", id, false).Delete(&RbacRole{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found or is a system role")
	}
	return nil
}

// =========================================================================
// Permission CRUD
// =========================================================================

// CreatePermission membuat permission baru.
func (r *Repository) CreatePermission(perm *RbacPermission) error {
	return r.db.Create(perm).Error
}

// FindPermissionByID mencari permission berdasarkan UUID.
func (r *Repository) FindPermissionByID(id uuid.UUID) (*RbacPermission, error) {
	var perm RbacPermission
	if err := r.db.First(&perm, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}
	return &perm, nil
}

// FindAllPermissions mengembalikan semua permissions.
func (r *Repository) FindAllPermissions() ([]RbacPermission, error) {
	var perms []RbacPermission
	if err := r.db.Order("resource ASC, action ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// DeletePermission menghapus permission (hanya jika bukan system).
func (r *Repository) DeletePermission(id uuid.UUID) error {
	result := r.db.Where("id = ? AND is_system = ?", id, false).Delete(&RbacPermission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("permission not found or is a system permission")
	}
	return nil
}

// =========================================================================
// Role-Permission Assignment
// =========================================================================

// AssignPermission menambahkan permission ke role.
func (r *Repository) AssignPermission(roleID, permissionID uuid.UUID) error {
	rp := RbacRolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(&rp).Error
}

// RevokePermission menghapus permission dari role.
func (r *Repository) RevokePermission(roleID, permissionID uuid.UUID) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&RbacRolePermission{}).Error
}

// FindRolePermissions mengembalikan semua permission IDs untuk suatu role.
func (r *Repository) FindRolePermissions(roleID uuid.UUID) ([]RbacRolePermission, error) {
	var rps []RbacRolePermission
	if err := r.db.Where("role_id = ?", roleID).Find(&rps).Error; err != nil {
		return nil, err
	}
	return rps, nil
}

// =========================================================================
// Migration (AutoMigrate GORM models)
// =========================================================================

// AutoMigrate menjalankan AutoMigrate untuk semua model RBAC.
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&RbacRole{},
		&RbacPermission{},
		&RbacRolePermission{},
	)
}
