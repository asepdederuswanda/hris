package organization

import "time"

type CreateOrganizationRequest struct {
	OrganizationSummaryID *string `json:"organization_summary_id"`
	Code                  string  `json:"code" binding:"required,max=10"`
	Nomenclature          string  `json:"nomenclature" binding:"required,max=255"`
	ParentID              *string `json:"parent_id"`
	ZoneID                *string `json:"zone_id"`
	JobFamilyID           *string `json:"job_family_id"`
	GradingID             *string `json:"grading_id"`
}

type UpdateOrganizationRequest struct {
	Code         *string `json:"code" binding:"omitempty,max=10"`
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=255"`
	ParentID     *string `json:"parent_id"`
	ZoneID       *string `json:"zone_id"`
	JobFamilyID  *string `json:"job_family_id"`
	GradingID    *string `json:"grading_id"`
}

type OrganizationResponse struct {
	ID                   string     `json:"id"`
	Code                 string     `json:"code"`
	FullCode             string     `json:"full_code"`
	Nomenclature         string     `json:"nomenclature"`
	ParentID             *string    `json:"parent_id,omitempty"`
	Level                int        `json:"level"`
	SortOrder            int        `json:"sort_order"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Children             []OrganizationResponse `json:"children,omitempty"`
}

func (o *Organization) ToResponse() OrganizationResponse {
	resp := OrganizationResponse{
		ID:          o.ID.String(),
		Code:        o.Code,
		FullCode:    o.FullCode,
		Nomenclature: o.Nomenclature,
		Level:       o.Level,
		SortOrder:   o.SortOrder,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	if o.ParentID != nil {
		pid := o.ParentID.String()
		resp.ParentID = &pid
	}

	// Convert children recursively
	if len(o.Children) > 0 {
		resp.Children = make([]OrganizationResponse, 0, len(o.Children))
		for _, child := range o.Children {
			resp.Children = append(resp.Children, child.ToResponse())
		}
	}

	return resp
}
