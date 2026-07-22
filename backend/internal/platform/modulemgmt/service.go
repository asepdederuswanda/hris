package modulemgmt

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// Service untuk business logic Module Management.
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

// CreateModule mendaftarkan modul baru di platform.
func (s *Service) CreateModule(req CreateModuleRequest) (*ModuleResponse, error) {
	// Cek duplikasi slug
	if existing, _ := s.repo.FindBySlug(req.Slug); existing != nil {
		return nil, fmt.Errorf("module with slug '%s' already exists", req.Slug)
	}

	module := &PlatformModule{
		Name:        req.Name,
		Slug:        req.Slug,
		Version:     req.Version,
		Description: req.Description,
		IsCore:      req.IsCore,
	}

	if err := s.repo.Create(module); err != nil {
		return nil, err
	}

	s.logger.Info("Module registered",
		zap.String("module_id", module.ID.String()),
		zap.String("name", module.Name),
		zap.String("slug", module.Slug),
	)

	response := module.ToResponse()
	return &response, nil
}

// GetModule mengembalikan modul berdasarkan ID.
func (s *Service) GetModule(id string) (*ModuleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}

	module, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := module.ToResponse()
	return &response, nil
}

// ListModules mengembalikan daftar semua modul dengan pagination.
func (s *Service) ListModules(page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	modules, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []ModuleResponse
	for _, m := range modules {
		responses = append(responses, m.ToResponse())
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

// UpdateModule mengupdate modul.
func (s *Service) UpdateModule(id string, req UpdateModuleRequest) (*ModuleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}

	module, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		module.Name = *req.Name
	}
	if req.Version != nil {
		module.Version = *req.Version
	}
	if req.Description != nil {
		module.Description = *req.Description
	}
	if req.IsCore != nil {
		module.IsCore = *req.IsCore
	}

	if err := s.repo.Update(module); err != nil {
		return nil, err
	}

	response := module.ToResponse()
	return &response, nil
}

// ActivateModule mengaktifkan modul untuk company tertentu.
func (s *Service) ActivateModule(moduleID, companyID string) (*CompanyModuleResponse, error) {
	mid, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	// Cek apakah modul sudah terdaftar
	module, err := s.repo.FindByID(mid)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	// Upsert company-module relation
	cm, err := s.repo.UpsertCompanyModule(cid, mid, true)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Module activated for company",
		zap.String("module_id", moduleID),
		zap.String("module_name", module.Name),
		zap.String("company_id", companyID),
	)

	return &CompanyModuleResponse{
		CompanyID:   cm.CompanyID.String(),
		ModuleID:    cm.ModuleID.String(),
		ModuleName:  module.Name,
		Enabled:     cm.Enabled,
		ActivatedAt: cm.ActivatedAt,
	}, nil
}

// DeactivateModule menonaktifkan modul untuk company tertentu.
func (s *Service) DeactivateModule(moduleID, companyID string) (*CompanyModuleResponse, error) {
	mid, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	module, err := s.repo.FindByID(mid)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	cm, err := s.repo.UpsertCompanyModule(cid, mid, false)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Module deactivated for company",
		zap.String("module_id", moduleID),
		zap.String("module_name", module.Name),
		zap.String("company_id", companyID),
	)

	return &CompanyModuleResponse{
		CompanyID:   cm.CompanyID.String(),
		ModuleID:    cm.ModuleID.String(),
		ModuleName:  module.Name,
		Enabled:     cm.Enabled,
		ActivatedAt: cm.ActivatedAt,
	}, nil
}

// ListCompanyModules mengembalikan daftar modul yang terdaftar untuk company.
func (s *Service) ListCompanyModules(companyID string) ([]CompanyModuleResponse, error) {
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	modules, err := s.repo.FindCompanyModules(cid)
	if err != nil {
		return nil, err
	}

	return modules, nil
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
