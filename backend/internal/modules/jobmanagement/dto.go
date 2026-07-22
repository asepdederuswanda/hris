package jobmanagement

import "time"

// =========================================================================
// Request DTOs — Job Titles
// =========================================================================

type CreateJobTitleRequest struct {
	Name         string `json:"name" binding:"required,max=100"`
	Descriptions string `json:"descriptions"`
	Status       int8   `json:"status"`
}

type UpdateJobTitleRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Descriptions *string `json:"descriptions"`
	Status       *int8   `json:"status"`
}

// =========================================================================
// Request DTOs — Job Title Subs
// =========================================================================

type CreateJobTitleSubRequest struct {
	JobManagementTitleID string `json:"job_management_title_id" binding:"required"`
	Name                 string `json:"name" binding:"required,max=100"`
	Descriptions         string `json:"descriptions"`
	Status               int8   `json:"status"`
}

type UpdateJobTitleSubRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Descriptions *string `json:"descriptions"`
	Status       *int8   `json:"status"`
}

// =========================================================================
// Request DTOs — Job Values
// =========================================================================

type CreateJobValueRequest struct {
	JobManagementTitleSubID *string `json:"job_management_title_sub_id"`
	Type                    string  `json:"type" binding:"required"`
	Level                   *int    `json:"level"`
	Descriptions            string  `json:"descriptions"`
	Note                    string  `json:"note"`
	Sort                    *int    `json:"sort"`
}

type UpdateJobValueRequest struct {
	Type         *string `json:"type" binding:"omitempty"`
	Level        *int    `json:"level"`
	Descriptions *string `json:"descriptions"`
	Note         *string `json:"note"`
	Sort         *int    `json:"sort"`
}

// =========================================================================
// Request DTOs — Job Objectives
// =========================================================================

type CreateJobObjectiveRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Objective      string `json:"objective"`
}

type UpdateJobObjectiveRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Objective    *string `json:"objective"`
}

// =========================================================================
// Request DTOs — Job Identifications
// =========================================================================

type CreateJobIdentificationRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	GradingID      string `json:"grading_id" binding:"required"`
}

type UpdateJobIdentificationRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	GradingID    *string `json:"grading_id"`
}

// =========================================================================
// Request DTOs — Job Responsibilities
// =========================================================================

type CreateJobResponsibilityRequest struct {
	OrganizationID    string `json:"organization_id" binding:"required"`
	Nomenclature      string `json:"nomenclature" binding:"required,max=50"`
	FullCode          string `json:"full_code" binding:"required,max=20"`
	MainTask          string `json:"main_task"`
	Activities        string `json:"activities"`
	Outputs           string `json:"outputs"`
	SuccessIndicators string `json:"success_indicators"`
}

type UpdateJobResponsibilityRequest struct {
	Nomenclature      *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode          *string `json:"full_code" binding:"omitempty,max=20"`
	MainTask          *string `json:"main_task"`
	Activities        *string `json:"activities"`
	Outputs           *string `json:"outputs"`
	SuccessIndicators *string `json:"success_indicators"`
}

// =========================================================================
// Request DTOs — Job Education Experiences
// =========================================================================

type CreateJobEducationExperienceRequest struct {
	OrganizationID                 string  `json:"organization_id" binding:"required"`
	Nomenclature                   string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                       string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueEducationID  *string `json:"job_management_value_education_id"`
	JobManagementValueExperienceID *string `json:"job_management_value_experience_id"`
}

type UpdateJobEducationExperienceRequest struct {
	Nomenclature                   *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                       *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueEducationID  *string `json:"job_management_value_education_id"`
	JobManagementValueExperienceID *string `json:"job_management_value_experience_id"`
}

// =========================================================================
// Request DTOs — Job HR Authorities
// =========================================================================

type CreateJobHRAuthorityRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Description    string `json:"description"`
}

type UpdateJobHRAuthorityRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
}

// =========================================================================
// Request DTOs — Job Operational Authorities
// =========================================================================

type CreateJobOperationalAuthorityRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Description    string `json:"description"`
}

type UpdateJobOperationalAuthorityRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
}

// =========================================================================
// Request DTOs — Job Working Activities
// =========================================================================

type CreateJobWorkingActivityRequest struct {
	OrganizationID      string  `json:"organization_id" binding:"required"`
	Nomenclature        string  `json:"nomenclature" binding:"required,max=50"`
	FullCode            string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

type UpdateJobWorkingActivityRequest struct {
	Nomenclature        *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode            *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

// =========================================================================
// Request DTOs — Job Working Risks
// =========================================================================

type CreateJobWorkingRiskRequest struct {
	OrganizationID                   string  `json:"organization_id" binding:"required"`
	Nomenclature                     string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                         string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueEnvironmentID  *string `json:"job_management_value_environment_id"`
	JobManagementValueHazardID       *string `json:"job_management_value_hazard_id"`
}

type UpdateJobWorkingRiskRequest struct {
	Nomenclature                     *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                         *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueEnvironmentID  *string `json:"job_management_value_environment_id"`
	JobManagementValueHazardID       *string `json:"job_management_value_hazard_id"`
}

// =========================================================================
// Request DTOs — Job Relationships
// =========================================================================

type CreateJobRelationshipRequest struct {
	OrganizationID                  string  `json:"organization_id" binding:"required"`
	Nomenclature                    string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                        string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueRelationshipID *string `json:"job_management_value_relationship_id"`
	JobManagementValueFrequencyID   *string `json:"job_management_value_frequency_id"`
}

type UpdateJobRelationshipRequest struct {
	Nomenclature                    *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                        *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueRelationshipID *string `json:"job_management_value_relationship_id"`
	JobManagementValueFrequencyID   *string `json:"job_management_value_frequency_id"`
}

// =========================================================================
// Request DTOs — Job Subordinate Controls
// =========================================================================

type CreateJobSubordinateControlRequest struct {
	OrganizationID      string  `json:"organization_id" binding:"required"`
	Nomenclature        string  `json:"nomenclature" binding:"required,max=50"`
	FullCode            string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

type UpdateJobSubordinateControlRequest struct {
	Nomenclature        *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode            *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

// =========================================================================
// Request DTOs — Job Assets
// =========================================================================

type CreateJobAssetRequest struct {
	OrganizationID              string  `json:"organization_id" binding:"required"`
	Nomenclature                string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                    string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueAssetID   *string `json:"job_management_value_asset_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
}

type UpdateJobAssetRequest struct {
	Nomenclature                *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                    *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueAssetID   *string `json:"job_management_value_asset_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
}

// =========================================================================
// Request DTOs — Job Financials
// =========================================================================

type CreateJobFinancialRequest struct {
	OrganizationID              string  `json:"organization_id" binding:"required"`
	Nomenclature                string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                    string  `json:"full_code" binding:"required,max=20"`
	IsAuthorized                bool    `json:"is_authorized"`
	JobManagementValueCashID    *string `json:"job_management_value_cash_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
	JobManagementValueImpactID  *string `json:"job_management_value_impact_id"`
}

type UpdateJobFinancialRequest struct {
	Nomenclature                *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                    *string `json:"full_code" binding:"omitempty,max=20"`
	IsAuthorized                *bool   `json:"is_authorized"`
	JobManagementValueCashID    *string `json:"job_management_value_cash_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
	JobManagementValueImpactID  *string `json:"job_management_value_impact_id"`
}

// =========================================================================
// Request DTOs — Job Potency Competencies
// =========================================================================

type CreateJobPotencyCompetencyRequest struct {
	OrganizationID      string   `json:"organization_id" binding:"required"`
	JobManagementValueID *string  `json:"job_management_value_id"`
	CompetencyID        *string  `json:"competency_id"`
	Weight              *float64 `json:"weight"`
}

type UpdateJobPotencyCompetencyRequest struct {
	JobManagementValueID *string  `json:"job_management_value_id"`
	CompetencyID        *string  `json:"competency_id"`
	Weight              *float64 `json:"weight"`
}

// =========================================================================
// Request DTOs — Job Scores
// =========================================================================

type UpdateJobScoreRequest struct {
	JobValueWithFinancial    *uint64 `json:"job_value_with_financial"`
	JobValueWithoutFinancial *uint64 `json:"job_value_without_financial"`
	HasFinancialAuthority    *bool   `json:"has_financial_authority"`
	Components               *string `json:"components"`
	SubComponentPoints       *string `json:"sub_component_points"`
}

// =========================================================================
// Request DTOs — Job Competency Groups
// =========================================================================

type CreateJobCompetencyGroupRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	Category       string  `json:"category" binding:"required,oneof=technical managerial"`
	Weight         float64 `json:"weight" binding:"required"`
}

type UpdateJobCompetencyGroupRequest struct {
	Category *string  `json:"category" binding:"omitempty,oneof=technical managerial"`
	Weight   *float64 `json:"weight"`
}

// =========================================================================
// Response DTOs
// =========================================================================

type JobTitleResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Descriptions string    `json:"descriptions,omitempty"`
	Status       int8      `json:"status,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Subs         []JobTitleSubResponse `json:"subs,omitempty"`
}

type JobTitleSubResponse struct {
	ID                      string    `json:"id"`
	JobManagementTitleID    string    `json:"job_management_title_id,omitempty"`
	JobManagementTitleName  string    `json:"job_management_title_name,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Descriptions            string    `json:"descriptions,omitempty"`
	Status                  int8      `json:"status,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type JobValueResponse struct {
	ID                         string    `json:"id"`
	JobManagementTitleSubID    string    `json:"job_management_title_sub_id,omitempty"`
	JobManagementTitleSubName  string    `json:"job_management_title_sub_name,omitempty"`
	Type                       string    `json:"type"`
	Level                      int       `json:"level,omitempty"`
	Descriptions               string    `json:"descriptions,omitempty"`
	Note                       string    `json:"note,omitempty"`
	Sort                       int       `json:"sort,omitempty"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type JobObjectiveResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Objective      string    `json:"objective,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobIdentificationResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	GradingID      string    `json:"grading_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobResponsibilityResponse struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id,omitempty"`
	Nomenclature      string    `json:"nomenclature"`
	FullCode          string    `json:"full_code"`
	MainTask          string    `json:"main_task,omitempty"`
	Activities        string    `json:"activities,omitempty"`
	Outputs           string    `json:"outputs,omitempty"`
	SuccessIndicators string    `json:"success_indicators,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type JobEducationExperienceResponse struct {
	ID                              string    `json:"id"`
	OrganizationID                  string    `json:"organization_id,omitempty"`
	Nomenclature                    string    `json:"nomenclature"`
	FullCode                        string    `json:"full_code"`
	JobManagementValueEducationID   string    `json:"job_management_value_education_id,omitempty"`
	JobManagementValueExperienceID  string    `json:"job_management_value_experience_id,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type JobHRAuthorityResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobOperationalAuthorityResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobWorkingActivityResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	Nomenclature        string    `json:"nomenclature"`
	FullCode            string    `json:"full_code"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobWorkingRiskResponse struct {
	ID                              string    `json:"id"`
	OrganizationID                  string    `json:"organization_id,omitempty"`
	Nomenclature                    string    `json:"nomenclature"`
	FullCode                        string    `json:"full_code"`
	JobManagementValueEnvironmentID string    `json:"job_management_value_environment_id,omitempty"`
	JobManagementValueHazardID      string    `json:"job_management_value_hazard_id,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type JobRelationshipResponse struct {
	ID                              string    `json:"id"`
	OrganizationID                  string    `json:"organization_id,omitempty"`
	Nomenclature                    string    `json:"nomenclature"`
	FullCode                        string    `json:"full_code"`
	JobManagementValueRelationshipID string   `json:"job_management_value_relationship_id,omitempty"`
	JobManagementValueFrequencyID   string    `json:"job_management_value_frequency_id,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type JobSubordinateControlResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	Nomenclature        string    `json:"nomenclature"`
	FullCode            string    `json:"full_code"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobAssetResponse struct {
	ID                          string    `json:"id"`
	OrganizationID              string    `json:"organization_id,omitempty"`
	Nomenclature                string    `json:"nomenclature"`
	FullCode                    string    `json:"full_code"`
	JobManagementValueAssetID   string    `json:"job_management_value_asset_id,omitempty"`
	JobManagementValueAuthorityID string  `json:"job_management_value_authority_id,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type JobFinancialResponse struct {
	ID                          string    `json:"id"`
	OrganizationID              string    `json:"organization_id,omitempty"`
	Nomenclature                string    `json:"nomenclature"`
	FullCode                    string    `json:"full_code"`
	IsAuthorized                bool      `json:"is_authorized"`
	JobManagementValueCashID    string    `json:"job_management_value_cash_id,omitempty"`
	JobManagementValueAuthorityID string  `json:"job_management_value_authority_id,omitempty"`
	JobManagementValueImpactID  string    `json:"job_management_value_impact_id,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type JobPotencyCompetencyResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CompetencyID        string    `json:"competency_id,omitempty"`
	Weight              float64   `json:"weight,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobScoreResponse struct {
	ID                       string     `json:"id"`
	OrganizationID           string     `json:"organization_id"`
	JobValueWithFinancial    uint64     `json:"job_value_with_financial"`
	JobValueWithoutFinancial uint64     `json:"job_value_without_financial"`
	HasFinancialAuthority    bool       `json:"has_financial_authority"`
	Components               string     `json:"components,omitempty"`
	SubComponentPoints       string     `json:"sub_component_points,omitempty"`
	CalculatedAt             *time.Time `json:"calculated_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type JobCompetencyGroupResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Category       string    `json:"category"`
	Weight         float64   `json:"weight"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Converter Helpers
// =========================================================================

func toJobTitleResponse(t *JobTitle) JobTitleResponse {
	r := JobTitleResponse{
		ID:        t.ID.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Name != nil {
		r.Name = *t.Name
	}
	if t.Descriptions != nil {
		r.Descriptions = *t.Descriptions
	}
	if t.Status != nil {
		r.Status = *t.Status
	}
	if len(t.Subs) > 0 {
		r.Subs = make([]JobTitleSubResponse, 0, len(t.Subs))
		for _, s := range t.Subs {
			r.Subs = append(r.Subs, toJobTitleSubResponse(&s))
		}
	}
	return r
}

func toJobTitleSubResponse(s *JobTitleSub) JobTitleSubResponse {
	r := JobTitleSubResponse{
		ID:        s.ID.String(),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.JobManagementTitleID != nil {
		r.JobManagementTitleID = s.JobManagementTitleID.String()
	}
	if s.JobManagementTitleName != nil {
		r.JobManagementTitleName = *s.JobManagementTitleName
	}
	if s.Name != nil {
		r.Name = *s.Name
	}
	if s.Descriptions != nil {
		r.Descriptions = *s.Descriptions
	}
	if s.Status != nil {
		r.Status = *s.Status
	}
	return r
}

func toJobValueResponse(v *JobValue) JobValueResponse {
	r := JobValueResponse{
		ID:        v.ID.String(),
		Type:      v.Type,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.JobManagementTitleSubID != nil {
		r.JobManagementTitleSubID = v.JobManagementTitleSubID.String()
	}
	if v.JobManagementTitleSubName != nil {
		r.JobManagementTitleSubName = *v.JobManagementTitleSubName
	}
	if v.Level != nil {
		r.Level = *v.Level
	}
	if v.Descriptions != nil {
		r.Descriptions = *v.Descriptions
	}
	if v.Note != nil {
		r.Note = *v.Note
	}
	if v.Sort != nil {
		r.Sort = *v.Sort
	}
	return r
}

func toJobObjectiveResponse(o *JobObjective) JobObjectiveResponse {
	r := JobObjectiveResponse{
		ID:           o.ID.String(),
		Nomenclature: o.Nomenclature,
		FullCode:     o.FullCode,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
	if o.OrganizationID != nil {
		r.OrganizationID = o.OrganizationID.String()
	}
	if o.Objective != nil {
		r.Objective = *o.Objective
	}
	return r
}

func toJobIdentificationResponse(i *JobIdentification) JobIdentificationResponse {
	r := JobIdentificationResponse{
		ID:           i.ID.String(),
		Nomenclature: i.Nomenclature,
		FullCode:     i.FullCode,
		GradingID:    i.GradingID.String(),
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
	if i.OrganizationID != nil {
		r.OrganizationID = i.OrganizationID.String()
	}
	return r
}

func toJobResponsibilityResponse(r *JobResponsibility) JobResponsibilityResponse {
	resp := JobResponsibilityResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.MainTask != nil {
		resp.MainTask = *r.MainTask
	}
	if r.Activities != nil {
		resp.Activities = *r.Activities
	}
	if r.Outputs != nil {
		resp.Outputs = *r.Outputs
	}
	if r.SuccessIndicators != nil {
		resp.SuccessIndicators = *r.SuccessIndicators
	}
	return resp
}

func toJobEducationExperienceResponse(e *JobEducationExperience) JobEducationExperienceResponse {
	r := JobEducationExperienceResponse{
		ID:           e.ID.String(),
		Nomenclature: e.Nomenclature,
		FullCode:     e.FullCode,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
	if e.OrganizationID != nil {
		r.OrganizationID = e.OrganizationID.String()
	}
	if e.JobManagementValueEducationID != nil {
		r.JobManagementValueEducationID = e.JobManagementValueEducationID.String()
	}
	if e.JobManagementValueExperienceID != nil {
		r.JobManagementValueExperienceID = e.JobManagementValueExperienceID.String()
	}
	return r
}

func toJobHRAuthorityResponse(a *JobHRAuthority) JobHRAuthorityResponse {
	r := JobHRAuthorityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

func toJobOperationalAuthorityResponse(a *JobOperationalAuthority) JobOperationalAuthorityResponse {
	r := JobOperationalAuthorityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

func toJobWorkingActivityResponse(a *JobWorkingActivity) JobWorkingActivityResponse {
	r := JobWorkingActivityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.JobManagementValueID != nil {
		r.JobManagementValueID = a.JobManagementValueID.String()
	}
	return r
}

func toJobWorkingRiskResponse(r *JobWorkingRisk) JobWorkingRiskResponse {
	resp := JobWorkingRiskResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.JobManagementValueEnvironmentID != nil {
		resp.JobManagementValueEnvironmentID = r.JobManagementValueEnvironmentID.String()
	}
	if r.JobManagementValueHazardID != nil {
		resp.JobManagementValueHazardID = r.JobManagementValueHazardID.String()
	}
	return resp
}

func toJobRelationshipResponse(r *JobRelationship) JobRelationshipResponse {
	resp := JobRelationshipResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.JobManagementValueRelationshipID != nil {
		resp.JobManagementValueRelationshipID = r.JobManagementValueRelationshipID.String()
	}
	if r.JobManagementValueFrequencyID != nil {
		resp.JobManagementValueFrequencyID = r.JobManagementValueFrequencyID.String()
	}
	return resp
}

func toJobSubordinateControlResponse(c *JobSubordinateControl) JobSubordinateControlResponse {
	r := JobSubordinateControlResponse{
		ID:           c.ID.String(),
		Nomenclature: c.Nomenclature,
		FullCode:     c.FullCode,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
	if c.OrganizationID != nil {
		r.OrganizationID = c.OrganizationID.String()
	}
	if c.JobManagementValueID != nil {
		r.JobManagementValueID = c.JobManagementValueID.String()
	}
	return r
}

func toJobAssetResponse(a *JobAsset) JobAssetResponse {
	r := JobAssetResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.JobManagementValueAssetID != nil {
		r.JobManagementValueAssetID = a.JobManagementValueAssetID.String()
	}
	if a.JobManagementValueAuthorityID != nil {
		r.JobManagementValueAuthorityID = a.JobManagementValueAuthorityID.String()
	}
	return r
}

func toJobFinancialResponse(f *JobFinancial) JobFinancialResponse {
	r := JobFinancialResponse{
		ID:           f.ID.String(),
		Nomenclature: f.Nomenclature,
		FullCode:     f.FullCode,
		IsAuthorized: f.IsAuthorized,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
	if f.OrganizationID != nil {
		r.OrganizationID = f.OrganizationID.String()
	}
	if f.JobManagementValueCashID != nil {
		r.JobManagementValueCashID = f.JobManagementValueCashID.String()
	}
	if f.JobManagementValueAuthorityID != nil {
		r.JobManagementValueAuthorityID = f.JobManagementValueAuthorityID.String()
	}
	if f.JobManagementValueImpactID != nil {
		r.JobManagementValueImpactID = f.JobManagementValueImpactID.String()
	}
	return r
}

func toJobPotencyCompetencyResponse(c *JobPotencyCompetency) JobPotencyCompetencyResponse {
	r := JobPotencyCompetencyResponse{
		ID:        c.ID.String(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.OrganizationID != nil {
		r.OrganizationID = c.OrganizationID.String()
	}
	if c.JobManagementValueID != nil {
		r.JobManagementValueID = c.JobManagementValueID.String()
	}
	if c.CompetencyID != nil {
		r.CompetencyID = c.CompetencyID.String()
	}
	if c.Weight != nil {
		r.Weight = *c.Weight
	}
	return r
}

func toJobScoreResponse(s *JobScore) JobScoreResponse {
	r := JobScoreResponse{
		ID:                       s.ID.String(),
		JobValueWithFinancial:    s.JobValueWithFinancial,
		JobValueWithoutFinancial: s.JobValueWithoutFinancial,
		HasFinancialAuthority:    s.HasFinancialAuthority,
		CalculatedAt:             s.CalculatedAt,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
	}
	if s.OrganizationID != nil {
		r.OrganizationID = s.OrganizationID.String()
	}
	if s.Components != nil {
		r.Components = *s.Components
	}
	if s.SubComponentPoints != nil {
		r.SubComponentPoints = *s.SubComponentPoints
	}
	return r
}

func toJobCompetencyGroupResponse(g *JobCompetencyGroup) JobCompetencyGroupResponse {
	r := JobCompetencyGroupResponse{
		ID:        g.ID.String(),
		Category:  g.Category,
		Weight:    g.Weight,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.OrganizationID != nil {
		r.OrganizationID = g.OrganizationID.String()
	}
	return r
}
