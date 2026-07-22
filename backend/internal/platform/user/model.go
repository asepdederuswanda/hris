package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlatformUserRole enum untuk role user.
type PlatformUserRole string

const (
	RoleSuperAdmin   PlatformUserRole = "super_admin"
	RoleCompanyAdmin PlatformUserRole = "company_admin"
)

// PlatformUser merepresentasikan user yang mengelola platform.
type PlatformUser struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID    *uuid.UUID     `gorm:"type:char(36);index" json:"company_id,omitempty"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Role         string         `gorm:"type:varchar(50);default:'admin'" json:"role"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PlatformUser) TableName() string {
	return "platform_users"
}

func (u *PlatformUser) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
