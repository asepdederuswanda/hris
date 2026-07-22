package organization

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization merepresentasikan unit organisasi dalam struktur tree.
// Setiap organization bisa memiliki parent (parent_id) untuk membentuk
// hierarki organisasi.
type Organization struct {
	ID                   uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationSummaryID *uuid.UUID     `gorm:"type:char(36)" json:"organization_summary_id,omitempty"`
	Code                 string          `gorm:"type:varchar(10);not null" json:"code"`
	FullCode             string          `gorm:"type:varchar(50);not null;index" json:"full_code"`
	Nomenclature         string          `gorm:"type:varchar(255);not null" json:"nomenclature"`
	ParentID             *uuid.UUID      `gorm:"type:char(36)" json:"parent_id,omitempty"`
	ZoneID               *uuid.UUID      `gorm:"type:char(36)" json:"zone_id,omitempty"`
	JobFamilyID          *uuid.UUID      `gorm:"type:char(36)" json:"job_family_id,omitempty"`
	GradingID            *uuid.UUID      `gorm:"type:char(36)" json:"grading_id,omitempty"`
	Level                int             `gorm:"default:0" json:"level"`
	SortOrder            int             `gorm:"default:0" json:"sort_order"`
	CreatedBy            *uuid.UUID      `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID      `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	DeletedAt            gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Parent     *Organization  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []Organization `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
