package user

import "time"

// LoginRequest untuk login platform admin.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse untuk response login.
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest untuk refresh access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse untuk response refresh token.
type RefreshTokenResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// CreateUserRequest untuk membuat platform user baru.
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	Name      string  `json:"name" binding:"required,min=1"`
	Role      string  `json:"role" binding:"required,oneof=super_admin company_admin"`
	CompanyID *string `json:"company_id,omitempty"`
}

// UpdateUserRequest untuk update platform user.
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=1"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=super_admin company_admin"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UserResponse untuk response data user.
type UserResponse struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CompanyID *string    `json:"company_id,omitempty"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse mengonversi PlatformUser ke UserResponse.
func (u *PlatformUser) ToResponse() UserResponse {
	resp := UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLoginAt,
	}
	if u.CompanyID != nil {
		cid := u.CompanyID.String()
		resp.CompanyID = &cid
	}
	return resp
}
