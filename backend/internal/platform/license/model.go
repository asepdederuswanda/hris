package license

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LicensePlanType enum untuk tipe plan.
type LicensePlanType string

const (
	PlanFree       LicensePlanType = "free"
	PlanBasic      LicensePlanType = "basic"
	PlanPro        LicensePlanType = "pro"
	PlanEnterprise LicensePlanType = "enterprise"
)

// LicenseStatus enum untuk status lisensi.
type LicenseStatus string

const (
	LicenseActive    LicenseStatus = "active"
	LicenseExpired   LicenseStatus = "expired"
	LicenseSuspended LicenseStatus = "suspended"
	LicenseCancelled LicenseStatus = "cancelled"
)

// License merepresentasikan lisensi untuk sebuah company.
type License struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID    uuid.UUID      `gorm:"type:char(36);not null;index" json:"company_id"`
	LicenseKey   string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"license_key"`
	PlanType     string         `gorm:"type:varchar(50);not null" json:"plan_type"`
	MaxEmployees int            `gorm:"default:0" json:"max_employees"`
	MaxModules   int            `gorm:"default:0" json:"max_modules"`
	StartDate    time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time      `gorm:"type:date;not null" json:"end_date"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (License) TableName() string {
	return "licenses"
}

func (l *License) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
