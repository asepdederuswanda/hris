package authz

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RbacRole merepresentasikan role dalam database.
type RbacRole struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Slug        string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"slug"`
	Description *string    `gorm:"type:varchar(255)" json:"description,omitempty"`
	ParentID    *uuid.UUID `gorm:"type:char(36)" json:"parent_id,omitempty"`
	IsSystem    bool       `gorm:"type:smallint;not null;default:0" json:"is_system"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi
	Permissions []RbacRolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
}

func (RbacRole) TableName() string {
	return "rbac_roles"
}

func (r *RbacRole) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RbacPermission merepresentasikan satu permission (resource + action).
type RbacPermission struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Resource    string    `gorm:"type:varchar(100);not null;uniqueIndex:uq_rbac_permission" json:"resource"`
	Action      string    `gorm:"type:varchar(50);not null;uniqueIndex:uq_rbac_permission" json:"action"`
	Description *string   `gorm:"type:varchar(255)" json:"description,omitempty"`
	IsSystem    bool      `gorm:"type:smallint;not null;default:0" json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
}

func (RbacPermission) TableName() string {
	return "rbac_permissions"
}

func (p *RbacPermission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// RbacRolePermission menghubungkan role dengan permission (many-to-many).
type RbacRolePermission struct {
	RoleID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:char(36);primaryKey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RbacRolePermission) TableName() string {
	return "rbac_role_permissions"
}
