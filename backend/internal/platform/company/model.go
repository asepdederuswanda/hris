package company

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyStatus enum untuk status company.
type CompanyStatus string

const (
	CompanyStatusActive     CompanyStatus = "active"
	CompanyStatusSuspended  CompanyStatus = "suspended"
	CompanyStatusTerminated CompanyStatus = "terminated"
)

// Company merepresentasikan perusahaan/tenant di platform database.
type Company struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	NPWP      *string        `gorm:"type:varchar(16)" json:"npwp,omitempty"`
	NIB       *string        `gorm:"type:varchar(25)" json:"nib,omitempty"`
	Address   *string        `gorm:"type:text" json:"address,omitempty"`
	Email     *string        `gorm:"type:varchar(255)" json:"email,omitempty"`
	Phone     *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Status    CompanyStatus  `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedBy *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}

// BeforeCreate hook untuk generate slug.
func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
