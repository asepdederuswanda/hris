package competency

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// 8.1 Competency (Master kompetensi)
// =========================================================================

type Competency struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Field      *string   `gorm:"type:varchar(255)" json:"field,omitempty"`
	Cluster    *string   `gorm:"type:varchar(255)" json:"cluster,omitempty"`
	Definition *string   `gorm:"type:text" json:"definition,omitempty"`
	CreatedBy  *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Competency) TableName() string {
	return "competencies"
}

func (c *Competency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.2 CompetenceValue (Nilai kompetensi — legacy style)
// =========================================================================

type CompetenceValue struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Type        *string   `gorm:"type:varchar(255)" json:"type,omitempty"`
	Level       *int      `gorm:"type:int" json:"level,omitempty"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Point       *int      `gorm:"type:int" json:"point,omitempty"`
	Description *string   `gorm:"type:varchar(255)" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CompetenceValue) TableName() string {
	return "competence_values"
}

func (v *CompetenceValue) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.3 CompetencyValue (Nilai kompetensi — structured)
// =========================================================================

type CompetencyValue struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Type        string    `gorm:"type:varchar(255);not null;index" json:"type"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Level       int       `gorm:"type:smallint;not null" json:"level"`
	Code        *string   `gorm:"type:varchar(255)" json:"code,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CompetencyValue) TableName() string {
	return "competency_values"
}

func (v *CompetencyValue) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.4 CompetencyEvent (Periode penilaian kompetensi)
// =========================================================================

type CompetencyEvent struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Type         string    `gorm:"type:varchar(20);not null" json:"type"`
	PeriodType   string    `gorm:"type:varchar(20);not null" json:"period_type"`
	PeriodYear   int       `gorm:"type:smallint;not null" json:"period_year"`
	PeriodNumber *int      `gorm:"type:smallint" json:"period_number,omitempty"`
	Status       string    `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedBy    *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy    *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (CompetencyEvent) TableName() string {
	return "competency_events"
}

func (e *CompetencyEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.5 CompetencyEventTarget (Target event kompetensi)
// =========================================================================

type CompetencyEventTarget struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CompetencyEventID  uuid.UUID `gorm:"type:char(36);not null;index" json:"competency_event_id"`
	OrganizationID     uuid.UUID `gorm:"type:char(36);not null;index" json:"organization_id"`
	EmployeeID         *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	MissingSelf        int        `gorm:"type:smallint;default:0" json:"missing_self"`
	MissingSuperior    int        `gorm:"type:smallint;default:0" json:"missing_superior"`
	MissingPeer        int        `gorm:"type:smallint;default:0" json:"missing_peer"`
	MissingSubordinate int        `gorm:"type:smallint;default:0" json:"missing_subordinate"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relasi
	CompetencyEvent *CompetencyEvent `gorm:"foreignKey:CompetencyEventID" json:"competency_event,omitempty"`
}

func (CompetencyEventTarget) TableName() string {
	return "competency_event_targets"
}

func (t *CompetencyEventTarget) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.6 CompetencyScore (Skor penilaian per organisasi)
// =========================================================================

type CompetencyScore struct {
	ID                       uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID           uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_comp_score_org" json:"organization_id"`
	EmployeeID               *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	TechnicalGapPercentage   float64    `gorm:"type:decimal(6,2);default:0" json:"technical_gap_percentage"`
	ManagerialGapPercentage  float64    `gorm:"type:decimal(6,2);default:0" json:"managerial_gap_percentage"`
	TotalGapPercentage       float64    `gorm:"type:decimal(6,2);default:0" json:"total_gap_percentage"`
	TotalGradePercentage     float64    `gorm:"type:decimal(6,2);default:0" json:"total_grade_percentage"`
	CompetencyEventID        *uuid.UUID `gorm:"type:char(36);index" json:"competency_event_id,omitempty"`
	AssessedAt               *time.Time `gorm:"type:timestamp" json:"assessed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`

	// Relasi
	Details []CompetencyScoreDetail `gorm:"foreignKey:CompetencyScoreID" json:"details,omitempty"`
}

func (CompetencyScore) TableName() string {
	return "competency_scores"
}

func (s *CompetencyScore) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 8.7 CompetencyScoreDetail (Detail skor kompetensi)
// =========================================================================

type CompetencyScoreDetail struct {
	ID                    uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CompetencyScoreID     uuid.UUID `gorm:"type:char(36);not null;index" json:"competency_score_id"`
	CompetencyID          uuid.UUID `gorm:"type:char(36);not null;index" json:"competency_id"`
	Type                  string    `gorm:"type:varchar(255);not null" json:"type"`
	StandardLevel         *int      `gorm:"type:smallint" json:"standard_level,omitempty"`
	StandardWeight        float64   `gorm:"type:decimal(6,2);default:0" json:"standard_weight"`
	EmployeeLevel         *int      `gorm:"type:smallint" json:"employee_level,omitempty"`
	GapPercentage         float64   `gorm:"type:decimal(6,2);default:0" json:"gap_percentage"`
	WeightedGapPercentage float64   `gorm:"type:decimal(6,2);default:0" json:"weighted_gap_percentage"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (CompetencyScoreDetail) TableName() string {
	return "competency_score_details"
}

func (d *CompetencyScoreDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
