package license

import "time"

// CreateLicenseRequest untuk membuat lisensi baru.
type CreateLicenseRequest struct {
	CompanyID    string `json:"company_id" binding:"required"`
	PlanType     string `json:"plan_type" binding:"required,oneof=free basic pro enterprise"`
	MaxEmployees int    `json:"max_employees" binding:"omitempty,min=0"`
	MaxModules   int    `json:"max_modules" binding:"omitempty,min=0"`
	StartDate    string `json:"start_date" binding:"required"`        // format: YYYY-MM-DD
	EndDate      string `json:"end_date" binding:"required"`          // format: YYYY-MM-DD
}

// UpdateLicenseRequest untuk update lisensi.
type UpdateLicenseRequest struct {
	PlanType     *string `json:"plan_type,omitempty" binding:"omitempty,oneof=free basic pro enterprise"`
	MaxEmployees *int    `json:"max_employees,omitempty" binding:"omitempty,min=0"`
	MaxModules   *int    `json:"max_modules,omitempty" binding:"omitempty,min=0"`
	StartDate    *string `json:"start_date,omitempty"`                // format: YYYY-MM-DD
	EndDate      *string `json:"end_date,omitempty"`                  // format: YYYY-MM-DD
	Status       *string `json:"status,omitempty" binding:"omitempty,oneof=active expired suspended cancelled"`
}

// LicenseResponse untuk response data lisensi.
type LicenseResponse struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	LicenseKey   string    `json:"license_key"`
	PlanType     string    `json:"plan_type"`
	MaxEmployees int       `json:"max_employees"`
	MaxModules   int       `json:"max_modules"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PaginatedResponse untuk response pagination.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// ToResponse mengonversi License ke LicenseResponse.
func (l *License) ToResponse() LicenseResponse {
	return LicenseResponse{
		ID:           l.ID.String(),
		CompanyID:    l.CompanyID.String(),
		LicenseKey:   l.LicenseKey,
		PlanType:     l.PlanType,
		MaxEmployees: l.MaxEmployees,
		MaxModules:   l.MaxModules,
		StartDate:    l.StartDate.Format("2006-01-02"),
		EndDate:      l.EndDate.Format("2006-01-02"),
		Status:       l.Status,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}
