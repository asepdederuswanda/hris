package competency

import "time"

// =========================================================================
// Request DTOs — Competency
// =========================================================================

type CreateCompetencyRequest struct {
	Name       string  `json:"name" binding:"required,max=255"`
	Field      *string `json:"field" binding:"omitempty,max=255"`
	Cluster    *string `json:"cluster" binding:"omitempty,max=255"`
	Definition *string `json:"definition"`
}

type UpdateCompetencyRequest struct {
	Name       *string `json:"name" binding:"omitempty,max=255"`
	Field      *string `json:"field" binding:"omitempty,max=255"`
	Cluster    *string `json:"cluster" binding:"omitempty,max=255"`
	Definition *string `json:"definition"`
}

// =========================================================================
// Request DTOs — CompetenceValue (legacy)
// =========================================================================

type CreateCompetenceValueRequest struct {
	Type        *string `json:"type" binding:"omitempty,max=255"`
	Level       *int    `json:"level"`
	Name        string  `json:"name" binding:"required,max=255"`
	Point       *int    `json:"point"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type UpdateCompetenceValueRequest struct {
	Type        *string `json:"type" binding:"omitempty,max=255"`
	Level       *int    `json:"level"`
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Point       *int    `json:"point"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

// =========================================================================
// Request DTOs — CompetencyValue (structured)
// =========================================================================

type CreateCompetencyValueRequest struct {
	Type        string  `json:"type" binding:"required,max=255"`
	Name        string  `json:"name" binding:"required,max=255"`
	Slug        string  `json:"slug" binding:"required,max=255"`
	Level       int     `json:"level" binding:"required"`
	Code        *string `json:"code" binding:"omitempty,max=255"`
	Description *string `json:"description"`
}

type UpdateCompetencyValueRequest struct {
	Type        *string `json:"type" binding:"omitempty,max=255"`
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Slug        *string `json:"slug" binding:"omitempty,max=255"`
	Level       *int    `json:"level"`
	Code        *string `json:"code" binding:"omitempty,max=255"`
	Description *string `json:"description"`
}

// =========================================================================
// Request DTOs — CompetencyEvent
// =========================================================================

type CreateCompetencyEventRequest struct {
	Type         string `json:"type" binding:"required,oneof=auto manual"`
	PeriodType   string `json:"period_type" binding:"required,oneof=annual semester quarter"`
	PeriodYear   int    `json:"period_year" binding:"required"`
	PeriodNumber *int   `json:"period_number"`
	Status       string `json:"status" binding:"omitempty,oneof=draft active closed"`
}

type UpdateCompetencyEventRequest struct {
	Type         *string `json:"type" binding:"omitempty,oneof=auto manual"`
	PeriodType   *string `json:"period_type" binding:"omitempty,oneof=annual semester quarter"`
	PeriodYear   *int    `json:"period_year"`
	PeriodNumber *int    `json:"period_number"`
	Status       *string `json:"status" binding:"omitempty,oneof=draft active closed"`
}

// =========================================================================
// Request DTOs — CompetencyEventTarget
// =========================================================================

type CreateCompetencyEventTargetRequest struct {
	CompetencyEventID  string `json:"competency_event_id" binding:"required"`
	OrganizationID     string `json:"organization_id" binding:"required"`
	EmployeeID         *string `json:"employee_id"`
	MissingSelf        *int   `json:"missing_self"`
	MissingSuperior    *int   `json:"missing_superior"`
	MissingPeer        *int   `json:"missing_peer"`
	MissingSubordinate *int   `json:"missing_subordinate"`
}

type UpdateCompetencyEventTargetRequest struct {
	EmployeeID         *string `json:"employee_id"`
	MissingSelf        *int    `json:"missing_self"`
	MissingSuperior    *int    `json:"missing_superior"`
	MissingPeer        *int    `json:"missing_peer"`
	MissingSubordinate *int    `json:"missing_subordinate"`
}

// =========================================================================
// Request DTOs — CompetencyScore
// =========================================================================

type CreateCompetencyScoreRequest struct {
	OrganizationID          string  `json:"organization_id" binding:"required"`
	EmployeeID              *string `json:"employee_id"`
	TechnicalGapPercentage  float64 `json:"technical_gap_percentage"`
	ManagerialGapPercentage float64 `json:"managerial_gap_percentage"`
	TotalGapPercentage      float64 `json:"total_gap_percentage"`
	TotalGradePercentage    float64 `json:"total_grade_percentage"`
	CompetencyEventID       *string `json:"competency_event_id"`
}

type UpdateCompetencyScoreRequest struct {
	TechnicalGapPercentage  *float64 `json:"technical_gap_percentage"`
	ManagerialGapPercentage *float64 `json:"managerial_gap_percentage"`
	TotalGapPercentage      *float64 `json:"total_gap_percentage"`
	TotalGradePercentage    *float64 `json:"total_grade_percentage"`
	CompetencyEventID       *string  `json:"competency_event_id"`
}

// =========================================================================
// Request DTOs — CompetencyScoreDetail
// =========================================================================

type CreateCompetencyScoreDetailRequest struct {
	CompetencyScoreID     string  `json:"competency_score_id" binding:"required"`
	CompetencyID          string  `json:"competency_id" binding:"required"`
	Type                  string  `json:"type" binding:"required,oneof=technical managerial"`
	StandardLevel         *int    `json:"standard_level"`
	StandardWeight        float64 `json:"standard_weight"`
	EmployeeLevel         *int    `json:"employee_level"`
	GapPercentage         float64 `json:"gap_percentage"`
	WeightedGapPercentage float64 `json:"weighted_gap_percentage"`
}

type UpdateCompetencyScoreDetailRequest struct {
	Type                  *string  `json:"type" binding:"omitempty,oneof=technical managerial"`
	StandardLevel         *int     `json:"standard_level"`
	StandardWeight        *float64 `json:"standard_weight"`
	EmployeeLevel         *int     `json:"employee_level"`
	GapPercentage         *float64 `json:"gap_percentage"`
	WeightedGapPercentage *float64 `json:"weighted_gap_percentage"`
}

// =========================================================================
// Response DTOs
// =========================================================================

type CompetencyResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Field      string    `json:"field,omitempty"`
	Cluster    string    `json:"cluster,omitempty"`
	Definition string    `json:"definition,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CompetenceValueResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type,omitempty"`
	Level       int       `json:"level,omitempty"`
	Name        string    `json:"name"`
	Point       int       `json:"point,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CompetencyValueResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Level       int       `json:"level"`
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CompetencyEventResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	PeriodType   string    `json:"period_type"`
	PeriodYear   int       `json:"period_year"`
	PeriodNumber int       `json:"period_number,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CompetencyEventTargetResponse struct {
	ID                 string    `json:"id"`
	CompetencyEventID  string    `json:"competency_event_id"`
	OrganizationID     string    `json:"organization_id"`
	EmployeeID         string    `json:"employee_id,omitempty"`
	MissingSelf        int       `json:"missing_self"`
	MissingSuperior    int       `json:"missing_superior"`
	MissingPeer        int       `json:"missing_peer"`
	MissingSubordinate int       `json:"missing_subordinate"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CompetencyScoreResponse struct {
	ID                       string    `json:"id"`
	OrganizationID           string    `json:"organization_id"`
	EmployeeID               string    `json:"employee_id,omitempty"`
	TechnicalGapPercentage   float64   `json:"technical_gap_percentage"`
	ManagerialGapPercentage  float64   `json:"managerial_gap_percentage"`
	TotalGapPercentage       float64   `json:"total_gap_percentage"`
	TotalGradePercentage     float64   `json:"total_grade_percentage"`
	CompetencyEventID        string    `json:"competency_event_id,omitempty"`
	AssessedAt               *time.Time `json:"assessed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type CompetencyScoreDetailResponse struct {
	ID                    string    `json:"id"`
	CompetencyScoreID     string    `json:"competency_score_id"`
	CompetencyID          string    `json:"competency_id"`
	Type                  string    `json:"type"`
	StandardLevel         int       `json:"standard_level,omitempty"`
	StandardWeight        float64   `json:"standard_weight"`
	EmployeeLevel         int       `json:"employee_level,omitempty"`
	GapPercentage         float64   `json:"gap_percentage"`
	WeightedGapPercentage float64   `json:"weighted_gap_percentage"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
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

func (c *Competency) ToResponse() CompetencyResponse {
	r := CompetencyResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.Field != nil {
		r.Field = *c.Field
	}
	if c.Cluster != nil {
		r.Cluster = *c.Cluster
	}
	if c.Definition != nil {
		r.Definition = *c.Definition
	}
	return r
}

func (v *CompetenceValue) ToResponse() CompetenceValueResponse {
	r := CompetenceValueResponse{
		ID:        v.ID.String(),
		Name:      v.Name,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.Type != nil {
		r.Type = *v.Type
	}
	if v.Level != nil {
		r.Level = *v.Level
	}
	if v.Point != nil {
		r.Point = *v.Point
	}
	if v.Description != nil {
		r.Description = *v.Description
	}
	return r
}

func (v *CompetencyValue) ToResponse() CompetencyValueResponse {
	r := CompetencyValueResponse{
		ID:        v.ID.String(),
		Type:      v.Type,
		Name:      v.Name,
		Slug:      v.Slug,
		Level:     v.Level,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.Code != nil {
		r.Code = *v.Code
	}
	if v.Description != nil {
		r.Description = *v.Description
	}
	return r
}

func (e *CompetencyEvent) ToResponse() CompetencyEventResponse {
	r := CompetencyEventResponse{
		ID:         e.ID.String(),
		Type:       e.Type,
		PeriodType: e.PeriodType,
		PeriodYear: e.PeriodYear,
		Status:     e.Status,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
	if e.PeriodNumber != nil {
		r.PeriodNumber = *e.PeriodNumber
	}
	return r
}

func (t *CompetencyEventTarget) ToResponse() CompetencyEventTargetResponse {
	r := CompetencyEventTargetResponse{
		ID:                 t.ID.String(),
		CompetencyEventID:  t.CompetencyEventID.String(),
		OrganizationID:     t.OrganizationID.String(),
		MissingSelf:        t.MissingSelf,
		MissingSuperior:    t.MissingSuperior,
		MissingPeer:        t.MissingPeer,
		MissingSubordinate: t.MissingSubordinate,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
	if t.EmployeeID != nil {
		r.EmployeeID = t.EmployeeID.String()
	}
	return r
}

func (s *CompetencyScore) ToResponse() CompetencyScoreResponse {
	r := CompetencyScoreResponse{
		ID:                       s.ID.String(),
		OrganizationID:           s.OrganizationID.String(),
		TechnicalGapPercentage:   s.TechnicalGapPercentage,
		ManagerialGapPercentage:  s.ManagerialGapPercentage,
		TotalGapPercentage:       s.TotalGapPercentage,
		TotalGradePercentage:     s.TotalGradePercentage,
		AssessedAt:               s.AssessedAt,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
	}
	if s.EmployeeID != nil {
		r.EmployeeID = s.EmployeeID.String()
	}
	if s.CompetencyEventID != nil {
		r.CompetencyEventID = s.CompetencyEventID.String()
	}
	return r
}

func (d *CompetencyScoreDetail) ToResponse() CompetencyScoreDetailResponse {
	r := CompetencyScoreDetailResponse{
		ID:                    d.ID.String(),
		CompetencyScoreID:     d.CompetencyScoreID.String(),
		CompetencyID:          d.CompetencyID.String(),
		Type:                  d.Type,
		StandardWeight:        d.StandardWeight,
		GapPercentage:         d.GapPercentage,
		WeightedGapPercentage: d.WeightedGapPercentage,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
	if d.StandardLevel != nil {
		r.StandardLevel = *d.StandardLevel
	}
	if d.EmployeeLevel != nil {
		r.EmployeeLevel = *d.EmployeeLevel
	}
	return r
}
