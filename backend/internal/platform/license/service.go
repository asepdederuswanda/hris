package license

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// Service untuk business logic License Management.
type Service struct {
	repo      *Repository
	dbManager *database.Manager
	logger    *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, dbManager *database.Manager, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		dbManager: dbManager,
		logger:    logger,
	}
}

// CreateLicense membuat lisensi baru untuk company.
func (s *Service) CreateLicense(req CreateLicenseRequest) (*LicenseResponse, error) {
	cid, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company_id: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format (use YYYY-MM-DD): %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format (use YYYY-MM-DD): %w", err)
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	maxEmployees := req.MaxEmployees
	if maxEmployees <= 0 {
		switch req.PlanType {
		case "free":
			maxEmployees = 10
		case "basic":
			maxEmployees = 50
		case "pro":
			maxEmployees = 200
		case "enterprise":
			maxEmployees = 0 // unlimited
		}
	}

	maxModules := req.MaxModules
	if maxModules <= 0 {
		switch req.PlanType {
		case "free":
			maxModules = 3
		case "basic":
			maxModules = 8
		case "pro":
			maxModules = 15
		case "enterprise":
			maxModules = 0 // unlimited
		}
	}

	license := &License{
		CompanyID:    cid,
		LicenseKey:   xid.New().String(),
		PlanType:     req.PlanType,
		MaxEmployees: maxEmployees,
		MaxModules:   maxModules,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       string(LicenseActive),
	}

	if err := s.repo.Create(license); err != nil {
		return nil, err
	}

	s.logger.Info("License created",
		zap.String("license_id", license.ID.String()),
		zap.String("company_id", req.CompanyID),
		zap.String("plan_type", req.PlanType),
	)

	response := license.ToResponse()
	return &response, nil
}

// GetLicense mengembalikan lisensi berdasarkan ID.
func (s *Service) GetLicense(id string) (*LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid license id: %w", err)
	}

	license, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := license.ToResponse()
	return &response, nil
}

// ListLicenses mengembalikan daftar lisensi dengan pagination.
func (s *Service) ListLicenses(page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	licenses, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []LicenseResponse
	for _, l := range licenses {
		responses = append(responses, l.ToResponse())
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

// UpdateLicense mengupdate lisensi.
func (s *Service) UpdateLicense(id string, req UpdateLicenseRequest) (*LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid license id: %w", err)
	}

	license, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	if req.PlanType != nil {
		license.PlanType = *req.PlanType
	}
	if req.MaxEmployees != nil {
		license.MaxEmployees = *req.MaxEmployees
	}
	if req.MaxModules != nil {
		license.MaxModules = *req.MaxModules
	}
	if req.Status != nil {
		license.Status = *req.Status
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			license.StartDate = startDate
		}
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			license.EndDate = endDate
		}
	}

	if err := s.repo.Update(license); err != nil {
		return nil, err
	}

	response := license.ToResponse()
	return &response, nil
}
