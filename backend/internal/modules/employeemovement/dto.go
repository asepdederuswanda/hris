package employeemovement

import "time"

// =========================================================================
// Employee Movement DTOs
// =========================================================================

type CreateMovementRequest struct {
	EmployeeID             string  `json:"employee_id" binding:"required,uuid"`
	MovementType           string  `json:"movement_type" binding:"required,oneof=promotion demotion mutation contract_extension status_change retirement offboarding other"`
	FromEmploymentID       *string `json:"from_employment_id" binding:"omitempty,uuid"`
	ToEmploymentID         *string `json:"to_employment_id" binding:"omitempty,uuid"`
	FromOrganizationID     *string `json:"from_organization_id" binding:"omitempty,uuid"`
	ToOrganizationID       *string `json:"to_organization_id" binding:"omitempty,uuid"`
	FromPositionID         *string `json:"from_position_id" binding:"omitempty,uuid"`
	ToPositionID           *string `json:"to_position_id" binding:"omitempty,uuid"`
	FromEmploymentStatusID *string `json:"from_employment_status_id" binding:"omitempty,uuid"`
	ToEmploymentStatusID   *string `json:"to_employment_status_id" binding:"omitempty,uuid"`
	Reason                 *string `json:"reason"`
	DecisionLetterNumber   string  `json:"decision_letter_number" binding:"required"`
	DecisionLetterDate     string  `json:"decision_letter_date" binding:"required"`
	EffectiveDate          string  `json:"effective_date" binding:"required"`
	Notes                  *string `json:"notes"`
}

type UpdateMovementRequest struct {
	MovementType           *string `json:"movement_type" binding:"omitempty,oneof=promotion demotion mutation contract_extension status_change retirement offboarding other"`
	ToOrganizationID       *string `json:"to_organization_id" binding:"omitempty,uuid"`
	ToPositionID           *string `json:"to_position_id" binding:"omitempty,uuid"`
	ToEmploymentStatusID   *string `json:"to_employment_status_id" binding:"omitempty,uuid"`
	Reason                 *string `json:"reason"`
	DecisionLetterNumber   *string `json:"decision_letter_number"`
	DecisionLetterDate     *string `json:"decision_letter_date"`
	EffectiveDate          *string `json:"effective_date"`
	Status                 *string `json:"status" binding:"omitempty,oneof=draft approved executed cancelled"`
	Notes                  *string `json:"notes"`
}

type MovementResponse struct {
	ID                   string     `json:"id"`
	EmployeeID           string     `json:"employee_id"`
	MovementType         string     `json:"movement_type"`
	FromEmploymentID     *string    `json:"from_employment_id,omitempty"`
	ToEmploymentID       *string    `json:"to_employment_id,omitempty"`
	FromOrganizationID   *string    `json:"from_organization_id,omitempty"`
	ToOrganizationID     *string    `json:"to_organization_id,omitempty"`
	FromPositionID       *string    `json:"from_position_id,omitempty"`
	ToPositionID         *string    `json:"to_position_id,omitempty"`
	FromEmploymentStatusID *string  `json:"from_employment_status_id,omitempty"`
	ToEmploymentStatusID   *string  `json:"to_employment_status_id,omitempty"`
	Reason               *string    `json:"reason,omitempty"`
	DecisionLetterNumber string     `json:"decision_letter_number"`
	DecisionLetterDate   string     `json:"decision_letter_date"`
	EffectiveDate        string     `json:"effective_date"`
	Status               string     `json:"status"`
	Notes                *string    `json:"notes,omitempty"`
	ApprovedBy           *string    `json:"approved_by,omitempty"`
	ApprovedAt           *time.Time `json:"approved_at,omitempty"`
	ExecutedBy           *string    `json:"executed_by,omitempty"`
	ExecutedAt           *time.Time `json:"executed_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type PaginatedMovementResponse struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"total_pages"`
}

// =========================================================================
// Employee Contract DTOs
// =========================================================================

type CreateContractRequest struct {
	EmployeeID          string  `json:"employee_id" binding:"required,uuid"`
	ContractNumber      string  `json:"contract_number" binding:"required"`
	ContractType        string  `json:"contract_type" binding:"required,oneof=pkwt pkwtt daily other"`
	StartDate           string  `json:"start_date" binding:"required"`
	EndDate             *string `json:"end_date" binding:"omitempty"`
	PreviousContractID  *string `json:"previous_contract_id" binding:"omitempty,uuid"`
	DecisionLetterNumber *string `json:"decision_letter_number"`
	Notes               *string `json:"notes"`
	DocumentURL         *string `json:"document_url"`
}

type UpdateContractRequest struct {
	ContractNumber      *string `json:"contract_number"`
	ContractType        *string `json:"contract_type" binding:"omitempty,oneof=pkwt pkwtt daily other"`
	EndDate             *string `json:"end_date"`
	DecisionLetterNumber *string `json:"decision_letter_number"`
	Notes               *string `json:"notes"`
	DocumentURL         *string `json:"document_url"`
	Status              *string `json:"status" binding:"omitempty,oneof=active expired extended terminated"`
}

type ContractResponse struct {
	ID                  string    `json:"id"`
	EmployeeID          string    `json:"employee_id"`
	ContractNumber      string    `json:"contract_number"`
	ContractType        string    `json:"contract_type"`
	StartDate           string    `json:"start_date"`
	EndDate             *string   `json:"end_date,omitempty"`
	ExtensionCount      int       `json:"extension_count"`
	PreviousContractID  *string   `json:"previous_contract_id,omitempty"`
	DecisionLetterNumber *string  `json:"decision_letter_number,omitempty"`
	Notes               *string   `json:"notes,omitempty"`
	DocumentURL         *string   `json:"document_url,omitempty"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PaginatedContractResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Converter helpers
// =========================================================================

func (m *EmployeeMovement) ToResponse() MovementResponse {
	r := MovementResponse{
		ID:                   m.ID.String(),
		EmployeeID:           m.EmployeeID.String(),
		MovementType:         string(m.MovementType),
		DecisionLetterNumber: m.DecisionLetterNumber,
		DecisionLetterDate:   m.DecisionLetterDate,
		EffectiveDate:        m.EffectiveDate,
		Status:               string(m.Status),
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
	if m.FromEmploymentID != nil {
		s := m.FromEmploymentID.String()
		r.FromEmploymentID = &s
	}
	if m.ToEmploymentID != nil {
		s := m.ToEmploymentID.String()
		r.ToEmploymentID = &s
	}
	if m.FromOrganizationID != nil {
		s := m.FromOrganizationID.String()
		r.FromOrganizationID = &s
	}
	if m.ToOrganizationID != nil {
		s := m.ToOrganizationID.String()
		r.ToOrganizationID = &s
	}
	if m.FromPositionID != nil {
		s := m.FromPositionID.String()
		r.FromPositionID = &s
	}
	if m.ToPositionID != nil {
		s := m.ToPositionID.String()
		r.ToPositionID = &s
	}
	if m.FromEmploymentStatusID != nil {
		s := m.FromEmploymentStatusID.String()
		r.FromEmploymentStatusID = &s
	}
	if m.ToEmploymentStatusID != nil {
		s := m.ToEmploymentStatusID.String()
		r.ToEmploymentStatusID = &s
	}
	if m.Reason != nil {
		r.Reason = m.Reason
	}
	if m.Notes != nil {
		r.Notes = m.Notes
	}
	if m.ApprovedBy != nil {
		s := m.ApprovedBy.String()
		r.ApprovedBy = &s
	}
	if m.ApprovedAt != nil {
		r.ApprovedAt = m.ApprovedAt
	}
	if m.ExecutedBy != nil {
		s := m.ExecutedBy.String()
		r.ExecutedBy = &s
	}
	if m.ExecutedAt != nil {
		r.ExecutedAt = m.ExecutedAt
	}
	return r
}

func (c *EmployeeContract) ToResponse() ContractResponse {
	r := ContractResponse{
		ID:             c.ID.String(),
		EmployeeID:     c.EmployeeID.String(),
		ContractNumber: c.ContractNumber,
		ContractType:   string(c.ContractType),
		StartDate:      c.StartDate,
		ExtensionCount: c.ExtensionCount,
		Status:         string(c.Status),
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
	if c.EndDate != nil {
		r.EndDate = c.EndDate
	}
	if c.PreviousContractID != nil {
		s := c.PreviousContractID.String()
		r.PreviousContractID = &s
	}
	if c.DecisionLetterNumber != nil {
		r.DecisionLetterNumber = c.DecisionLetterNumber
	}
	if c.Notes != nil {
		r.Notes = c.Notes
	}
	if c.DocumentURL != nil {
		r.DocumentURL = c.DocumentURL
	}
	return r
}
