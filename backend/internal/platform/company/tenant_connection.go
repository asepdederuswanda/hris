package company

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantConnection menyimpan detail koneksi database untuk setiap tenant.
// Setiap tenant bisa menggunakan driver database yang berbeda
// (postgres atau mysql) sesuai kebutuhan.
type TenantConnection struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID uuid.UUID      `gorm:"type:char(36);uniqueIndex;not null;constraint:OnDelete:CASCADE" json:"company_id"`
	Driver    string         `gorm:"type:varchar(20);default:'postgres'" json:"driver"` // postgres | mysql
	Host      string         `gorm:"type:varchar(255);not null;default:'localhost'" json:"host"`
	Port      int            `gorm:"type:integer;not null;default:5432" json:"port"`
	DBName    string         `gorm:"type:varchar(100);not null" json:"db_name"`
	Username  string         `gorm:"type:varchar(100);not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // TODO: encrypt at rest for production
	SSLMode   string         `gorm:"type:varchar(20);default:'require'" json:"sslmode"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TenantConnection) TableName() string {
	return "tenant_connections"
}

func (tc *TenantConnection) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == uuid.Nil {
		tc.ID = uuid.New()
	}
	return nil
}
