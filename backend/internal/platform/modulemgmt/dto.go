package modulemgmt

import "time"

// CreateModuleRequest untuk mendaftarkan modul baru.
type CreateModuleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Slug        string `json:"slug" binding:"required,min=2,max=100"`
	Version     string `json:"version" binding:"required,max=20"`
	Description string `json:"description,omitempty"`
	IsCore      bool   `json:"is_core,omitempty"`
}

// UpdateModuleRequest untuk update modul.
type UpdateModuleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Version     *string `json:"version,omitempty" binding:"omitempty,max=20"`
	Description *string `json:"description,omitempty"`
	IsCore      *bool   `json:"is_core,omitempty"`
}

// ModuleResponse untuk response data modul.
type ModuleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Version     string    `json:"version"`
	Description string    `json:"description,omitempty"`
	IsCore      bool      `json:"is_core"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToggleModuleRequest untuk activate/deactivate modul untuk company.
type ToggleModuleRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
}

// CompanyModuleResponse untuk response company-module association.
type CompanyModuleResponse struct {
	CompanyID   string     `json:"company_id"`
	ModuleID    string     `json:"module_id"`
	ModuleName  string     `json:"module_name"`
	Enabled     bool       `json:"enabled"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

// ToResponse mengonversi PlatformModule ke ModuleResponse.
func (m *PlatformModule) ToResponse() ModuleResponse {
	return ModuleResponse{
		ID:          m.ID.String(),
		Name:        m.Name,
		Slug:        m.Slug,
		Version:     m.Version,
		Description: m.Description,
		IsCore:      m.IsCore,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
