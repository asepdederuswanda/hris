package company

import "time"

// CreateCompanyRequest DTO untuk create company.
type CreateCompanyRequest struct {
	Name    string  `json:"name" binding:"required,min=3,max=255"`
	NPWP    *string `json:"npwp" binding:"omitempty,len=16"`
	NIB     *string `json:"nib" binding:"omitempty,max=25"`
	Address *string `json:"address"`
	Email   *string `json:"email" binding:"omitempty,email"`
	Phone   *string `json:"phone" binding:"omitempty,max=20"`
}

// UpdateCompanyRequest DTO untuk update company.
type UpdateCompanyRequest struct {
	Name    *string `json:"name" binding:"omitempty,min=3,max=255"`
	NPWP    *string `json:"npwp" binding:"omitempty,len=16"`
	NIB     *string `json:"nib" binding:"omitempty,max=25"`
	Address *string `json:"address"`
	Email   *string `json:"email" binding:"omitempty,email"`
	Phone   *string `json:"phone" binding:"omitempty,max=20"`
}

// CompanyResponse DTO untuk response company.
type CompanyResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	NPWP      *string    `json:"npwp,omitempty"`
	NIB       *string    `json:"nib,omitempty"`
	Address   *string    `json:"address,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// PaginatedResponse DTO untuk response pagination.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// ToResponse mengonversi Company model ke CompanyResponse.
func (c *Company) ToResponse() CompanyResponse {
	return CompanyResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		NPWP:      c.NPWP,
		NIB:       c.NIB,
		Address:   c.Address,
		Email:     c.Email,
		Phone:     c.Phone,
		Status:    string(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
