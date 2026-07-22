package jobmanagement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// 9.1 Job Management Titles (Jenis Jabatan)
// =========================================================================
type JobTitle struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name         *string    `gorm:"type:varchar(100)" json:"name,omitempty"`
	Descriptions *string    `gorm:"type:text" json:"descriptions,omitempty"`
	Status       *int8      `gorm:"type:tinyint" json:"status,omitempty"`
	CreatedBy    *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy    *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Subs []JobTitleSub `gorm:"foreignKey:JobManagementTitleID" json:"subs,omitempty"`
}

func (JobTitle) TableName() string { return "job_management_titles" }

func (t *JobTitle) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.2 Job Management Title Subs (Sub Jenis Jabatan)
// =========================================================================
type JobTitleSub struct {
	ID                    uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobManagementTitleID  *uuid.UUID `gorm:"type:char(36);index" json:"job_management_title_id,omitempty"`
	JobManagementTitleName *string   `gorm:"type:varchar(100)" json:"job_management_title_name,omitempty"`
	Name                  *string    `gorm:"type:varchar(100)" json:"name,omitempty"`
	Descriptions          *string    `gorm:"type:text" json:"descriptions,omitempty"`
	Status                *int8      `gorm:"type:tinyint" json:"status,omitempty"`
	CreatedBy             *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy             *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (JobTitleSub) TableName() string { return "job_management_title_subs" }

func (s *JobTitleSub) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.3 Job Management Values (Nilai Jabatan)
// =========================================================================
type JobValue struct {
	ID                      uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobManagementTitleSubID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_title_sub_id,omitempty"`
	JobManagementTitleSubName *string  `gorm:"type:varchar(100)" json:"job_management_title_sub_name,omitempty"`
	Type                    string     `gorm:"type:varchar(255);not null" json:"type"`
	Level                   *int       `gorm:"type:int" json:"level,omitempty"`
	Descriptions            *string    `gorm:"type:text" json:"descriptions,omitempty"`
	Note                    *string    `gorm:"type:text" json:"note,omitempty"`
	Sort                    *int       `gorm:"type:int" json:"sort,omitempty"`
	CreatedBy               *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy               *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

func (JobValue) TableName() string { return "job_management_values" }

func (v *JobValue) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.4 Job Management Objectives (Tujuan Jabatan)
// =========================================================================
type JobObjective struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Objective      *string    `gorm:"type:text" json:"objective,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobObjective) TableName() string { return "job_management_objectives" }

func (o *JobObjective) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.5 Job Management Identifications (Identitas Jabatan)
// =========================================================================
type JobIdentification struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	GradingID      uuid.UUID  `gorm:"type:char(36);not null" json:"grading_id"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobIdentification) TableName() string { return "job_management_identifications" }

func (i *JobIdentification) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.6 Job Management Responsibilities (Tanggung Jawab)
// =========================================================================
type JobResponsibility struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID    *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature      string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode          string     `gorm:"type:varchar(20);not null" json:"full_code"`
	MainTask          *string    `gorm:"type:text" json:"main_task,omitempty"`
	Activities        *string    `gorm:"type:text" json:"activities,omitempty"`
	Outputs           *string    `gorm:"type:text" json:"outputs,omitempty"`
	SuccessIndicators *string    `gorm:"type:text" json:"success_indicators,omitempty"`
	CreatedBy         *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (JobResponsibility) TableName() string { return "job_management_responsibilities" }

func (r *JobResponsibility) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.7 Job Management Education Experiences (Pendidikan & Pengalaman)
// =========================================================================
type JobEducationExperience struct {
	ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID                  *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                    string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                        string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueEducationID   *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_education_id,omitempty"`
	JobManagementValueExperienceID  *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_experience_id,omitempty"`
	CreatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
}

func (JobEducationExperience) TableName() string { return "job_management_education_experiences" }

func (e *JobEducationExperience) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.8 Job Management HR Authorities (Kewenangan SDM)
// =========================================================================
type JobHRAuthority struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Description    *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobHRAuthority) TableName() string { return "job_management_hr_authorities" }

func (a *JobHRAuthority) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.9 Job Management Operational Authorities (Kewenangan Operasional)
// =========================================================================
type JobOperationalAuthority struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Description    *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobOperationalAuthority) TableName() string { return "job_management_operational_authorities" }

func (a *JobOperationalAuthority) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.10 Job Management Working Activities (Aktivitas Kerja)
// =========================================================================
type JobWorkingActivity struct {
	ID                  uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature        string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode            string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CreatedBy           *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (JobWorkingActivity) TableName() string { return "job_management_working_activities" }

func (a *JobWorkingActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.11 Job Management Working Risks (Risiko Kerja)
// =========================================================================
type JobWorkingRisk struct {
	ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID                  *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                    string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                        string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueEnvironmentID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_environment_id,omitempty"`
	JobManagementValueHazardID      *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_hazard_id,omitempty"`
	CreatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
}

func (JobWorkingRisk) TableName() string { return "job_management_working_risks" }

func (r *JobWorkingRisk) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.12 Job Management Relationships (Hubungan Kerja)
// =========================================================================
type JobRelationship struct {
	ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID                  *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                    string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                        string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueRelationshipID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_relationship_id,omitempty"`
	JobManagementValueFrequencyID   *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_frequency_id,omitempty"`
	CreatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
}

func (JobRelationship) TableName() string { return "job_management_relationships" }

func (r *JobRelationship) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.13 Job Management Subordinate Controls (Bawahan yang Dikendalikan)
// =========================================================================
type JobSubordinateControl struct {
	ID                  uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature        string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode            string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CreatedBy           *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (JobSubordinateControl) TableName() string { return "job_management_subordinate_controls" }

func (c *JobSubordinateControl) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.14 Job Management Assets (Aset Jabatan)
// =========================================================================
type JobAsset struct {
	ID                          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID              *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                    string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueAssetID   *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_asset_id,omitempty"`
	JobManagementValueAuthorityID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_authority_id,omitempty"`
	CreatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

func (JobAsset) TableName() string { return "job_management_assets" }

func (a *JobAsset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.15 Job Management Financials (Keuangan Jabatan)
// =========================================================================
type JobFinancial struct {
	ID                          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID              *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                    string     `gorm:"type:varchar(20);not null" json:"full_code"`
	IsAuthorized                bool       `gorm:"type:tinyint(1);not null;default:0" json:"is_authorized"`
	JobManagementValueCashID    *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_cash_id,omitempty"`
	JobManagementValueAuthorityID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_authority_id,omitempty"`
	JobManagementValueImpactID  *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_impact_id,omitempty"`
	CreatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

func (JobFinancial) TableName() string { return "job_management_financials" }

func (f *JobFinancial) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.16 Job Management Potency Competencies (Kompetensi Potensi)
// =========================================================================
type JobPotencyCompetency struct {
	ID                  uuid.UUID   `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID  `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CompetencyID        *uuid.UUID  `gorm:"type:char(36);index" json:"competency_id,omitempty"`
	Weight              *float64    `gorm:"type:decimal(8,2)" json:"weight,omitempty"`
	CreatedBy           *uuid.UUID  `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID  `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

func (JobPotencyCompetency) TableName() string { return "job_management_potency_competencies" }

func (c *JobPotencyCompetency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.17 Job Management Scores (Skor Jabatan)
// =========================================================================
type JobScore struct {
	ID                    uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID        *uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex" json:"organization_id"`
	JobValueWithFinancial   uint64        `gorm:"type:bigint unsigned;not null;default:0" json:"job_value_with_financial"`
	JobValueWithoutFinancial uint64       `gorm:"type:bigint unsigned;not null;default:0" json:"job_value_without_financial"`
	HasFinancialAuthority  bool            `gorm:"type:tinyint(1);not null;default:0" json:"has_financial_authority"`
	Components             *string         `gorm:"type:json" json:"components,omitempty"`
	SubComponentPoints     *string         `gorm:"type:json" json:"sub_component_points,omitempty"`
	CalculatedAt           *time.Time      `json:"calculated_at,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

func (JobScore) TableName() string { return "job_management_scores" }

func (s *JobScore) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.18 Job Management Competency Groups (Bobot Kompetensi per Organisasi)
// =========================================================================
type JobCompetencyGroup struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);not null" json:"organization_id"`
	Category       string     `gorm:"type:varchar(20);not null" json:"category"`
	Weight         float64    `gorm:"type:decimal(8,2);not null;default:0" json:"weight"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobCompetencyGroup) TableName() string { return "job_management_competency_groups" }

func (g *JobCompetencyGroup) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
