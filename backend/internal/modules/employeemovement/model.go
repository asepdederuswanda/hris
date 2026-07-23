package employeemovement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MovementType enum untuk tipe pergerakan karyawan.
type MovementType string

const (
	MovementTypePromotion         MovementType = "promotion"
	MovementTypeDemotion          MovementType = "demotion"
	MovementTypeMutation          MovementType = "mutation"
	MovementTypeContractExtension MovementType = "contract_extension"
	MovementTypeStatusChange      MovementType = "status_change"
	MovementTypeRetirement        MovementType = "retirement"
	MovementTypeOffboarding       MovementType = "offboarding"
	MovementTypeOther             MovementType = "other"
)

// MovementStatus enum untuk status pergerakan.
type MovementStatus string

const (
	MovementStatusDraft     MovementStatus = "draft"
	MovementStatusApproved  MovementStatus = "approved"
	MovementStatusExecuted  MovementStatus = "executed"
	MovementStatusCancelled MovementStatus = "cancelled"
)

// ContractType enum untuk tipe kontrak.
type ContractType string

const (
	ContractTypePKWT  ContractType = "pkwt"
	ContractTypePKWTT ContractType = "pkwtt"
	ContractTypeDaily ContractType = "daily"
	ContractTypeOther ContractType = "other"
)

// ContractStatus enum untuk status kontrak.
type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusExpired   ContractStatus = "expired"
	ContractStatusExtended  ContractStatus = "extended"
	ContractStatusTerminated ContractStatus = "terminated"
)

// =========================================================================
// 12.1 EmployeeMovement (Riwayat Pergerakan Karyawan)
// =========================================================================

type EmployeeMovement struct {
	ID                   uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID           uuid.UUID      `gorm:"type:char(36);not null;index:idx_emp_mvmt_employee" json:"employee_id"`
	MovementType         MovementType   `gorm:"type:varchar(50);not null;index:idx_emp_mvmt_type" json:"movement_type"`
	FromEmploymentID     *uuid.UUID     `gorm:"type:char(36);index" json:"from_employment_id,omitempty"`
	ToEmploymentID       *uuid.UUID     `gorm:"type:char(36);index" json:"to_employment_id,omitempty"`
	FromOrganizationID   *uuid.UUID     `gorm:"type:char(36);index:idx_emp_mvmt_from_org" json:"from_organization_id,omitempty"`
	ToOrganizationID     *uuid.UUID     `gorm:"type:char(36);index:idx_emp_mvmt_to_org" json:"to_organization_id,omitempty"`
	FromPositionID       *uuid.UUID     `gorm:"type:char(36);index:idx_emp_mvmt_from_pos" json:"from_position_id,omitempty"`
	ToPositionID         *uuid.UUID     `gorm:"type:char(36);index:idx_emp_mvmt_to_pos" json:"to_position_id,omitempty"`
	FromEmploymentStatusID *uuid.UUID   `gorm:"type:char(36);index" json:"from_employment_status_id,omitempty"`
	ToEmploymentStatusID   *uuid.UUID   `gorm:"type:char(36);index" json:"to_employment_status_id,omitempty"`
	Reason               *string        `gorm:"type:text" json:"reason,omitempty"`
	DecisionLetterNumber string         `gorm:"type:varchar(50);not null" json:"decision_letter_number"`
	DecisionLetterDate   string         `gorm:"type:date;not null" json:"decision_letter_date"`
	EffectiveDate        string         `gorm:"type:date;not null;index:idx_emp_mvmt_effective" json:"effective_date"`
	Status               MovementStatus `gorm:"type:varchar(20);default:draft;index:idx_emp_mvmt_status" json:"status"`
	Notes                *string        `gorm:"type:text" json:"notes,omitempty"`
	ApprovedBy           *uuid.UUID     `gorm:"type:char(36)" json:"approved_by,omitempty"`
	ApprovedAt           *time.Time     `gorm:"type:timestamp" json:"approved_at,omitempty"`
	ExecutedBy           *uuid.UUID     `gorm:"type:char(36)" json:"executed_by,omitempty"`
	ExecutedAt           *time.Time     `gorm:"type:timestamp" json:"executed_at,omitempty"`
	CreatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (EmployeeMovement) TableName() string {
	return "employee_movements"
}

func (m *EmployeeMovement) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 12.2 EmployeeContract (PKWT & Perjanjian Kerja)
// =========================================================================

type EmployeeContract struct {
	ID                  uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID          uuid.UUID      `gorm:"type:char(36);not null;index:idx_emp_ctrct_employee" json:"employee_id"`
	ContractNumber      string         `gorm:"type:varchar(50);not null" json:"contract_number"`
	ContractType        ContractType   `gorm:"type:varchar(20);not null;index:idx_emp_ctrct_type" json:"contract_type"`
	StartDate           string         `gorm:"type:date;not null" json:"start_date"`
	EndDate             *string        `gorm:"type:date" json:"end_date,omitempty"`
	ExtensionCount      int            `gorm:"type:int;default:0" json:"extension_count"`
	PreviousContractID  *uuid.UUID     `gorm:"type:char(36);index" json:"previous_contract_id,omitempty"`
	DecisionLetterNumber *string       `gorm:"type:varchar(50)" json:"decision_letter_number,omitempty"`
	Notes               *string        `gorm:"type:text" json:"notes,omitempty"`
	DocumentURL         *string        `gorm:"type:varchar(255)" json:"document_url,omitempty"`
	Status              ContractStatus `gorm:"type:varchar(20);default:active;index:idx_emp_ctrct_status" json:"status"`
	CreatedBy           *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (EmployeeContract) TableName() string {
	return "employee_contracts"
}

func (c *EmployeeContract) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
