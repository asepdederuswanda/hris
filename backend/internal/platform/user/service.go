package user

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
)

// Service untuk business logic user & authentication.
type Service struct {
	repo        *Repository
	authManager *auth.Manager
	logger      *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, authManager *auth.Manager, logger *zap.Logger) *Service {
	return &Service{
		repo:        repo,
		authManager: authManager,
		logger:      logger,
	}
}

// Login memvalidasi kredensial dan mengembalikan JWT token.
func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT tokens
	companyID := ""
	if user.CompanyID != nil {
		companyID = user.CompanyID.String()
	}

	accessToken, refreshToken, err := s.authManager.GenerateTokenPair(
		user.ID.String(),
		companyID,
		user.Email,
		user.Name,
		user.Role,
		[]string{}, // TODO: load permissions
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login
	_ = s.repo.UpdateLastLogin(user.ID)

	s.logger.Info("Platform user logged in",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("role", user.Role),
	)

	response := user.ToResponse()
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
		User:         response,
	}, nil
}

// RefreshToken memvalidasi refresh token dan mengembalikan access token baru.
func (s *Service) RefreshToken(refreshToken string) (*RefreshTokenResponse, error) {
	claims, err := s.authManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Generate new access token
	newAccessToken, err := s.authManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	user, err := s.repo.FindByEmail(claims.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User:        user.ToResponse(),
	}, nil
}

// CreateUser membuat platform user baru.
func (s *Service) CreateUser(req CreateUserRequest) (*UserResponse, error) {
	// Cek apakah email sudah terdaftar
	if existing, _ := s.repo.FindByEmail(req.Email); existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &PlatformUser{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         req.Role,
		IsActive:     true,
	}

	// Set company ID if provided
	if req.CompanyID != nil && *req.CompanyID != "" {
		cid, err := uuid.Parse(*req.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("invalid company_id: %w", err)
		}
		user.CompanyID = &cid
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	s.logger.Info("Platform user created",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("role", user.Role),
	)

	response := user.ToResponse()
	return &response, nil
}

// GetUser mengembalikan user berdasarkan ID.
func (s *Service) GetUser(id string) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ListUsers mengembalikan daftar platform user dengan pagination.
func (s *Service) ListUsers(page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	users, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []UserResponse
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser mengupdate data user.
func (s *Service) UpdateUser(id string, req UpdateUserRequest) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// PaginatedResponse untuk response pagination (reused across modules).
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// EnsureMigrate menjalankan auto-migrasi untuk model user.
func (s *Service) EnsureMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&PlatformUser{})
}

// EnsureSeed menjalankan seeder untuk super admin default.
func (s *Service) EnsureSeed(db *gorm.DB) error {
	// Cek apakah sudah ada super admin
	var count int64
	db.Model(&PlatformUser{}).Where("role = ?", RoleSuperAdmin).Count(&count)
	if count > 0 {
		return nil // already seeded
	}

	// Buat default super admin (hanya untuk development)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	superAdmin := &PlatformUser{
		Email:        "superadmin@hris-platform.com",
		PasswordHash: string(hashedPassword),
		Name:         "Super Admin",
		Role:         string(RoleSuperAdmin),
		IsActive:     true,
	}
	// Generate UUID langsung tanpa hook
	if superAdmin.ID == uuid.Nil {
		superAdmin.ID = uuid.New()
	}

	// Gunakan db session tanpa default values
	if err := db.Session(&gorm.Session{}).Create(superAdmin).Error; err != nil {
		return fmt.Errorf("failed to seed super admin: %w", err)
	}

	return nil
}
