package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// PaginatedResponse DTO untuk response pagination.
type PaginatedResponse struct {
	Success    bool                   `json:"success"`
	Data       []OrganizationResponse `json:"data"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	Total      int64                  `json:"total"`
	TotalPages int                    `json:"total_pages"`
}

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, req CreateOrganizationRequest) (*OrganizationResponse, error) {
	org := &Organization{
		Code:         req.Code,
		Nomenclature: req.Nomenclature,
	}

	// Parse optional foreign keys with error handling
	if req.OrganizationSummaryID != nil && *req.OrganizationSummaryID != "" {
		id, err := uuid.Parse(*req.OrganizationSummaryID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization_summary_id: %w", err)
		}
		org.OrganizationSummaryID = &id
	}
	if req.ParentID != nil && *req.ParentID != "" {
		id, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		org.ParentID = &id
	}

	// Generate full_code and level based on parent
	if org.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, *org.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent not found: %w", err)
		}
		org.FullCode = parent.FullCode + org.Code
		org.Level = parent.Level + 1
	} else {
		org.FullCode = org.Code
		org.Level = 0
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}

	s.logger.Info("Organization created",
		zap.String("id", org.ID.String()),
		zap.String("code", org.FullCode),
		zap.String("nomenclature", org.Nomenclature),
	)

	response := org.ToResponse()
	return &response, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*OrganizationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	org, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := org.ToResponse()
	return &response, nil
}

// List mengembalikan daftar organisasi dengan pagination.
func (s *Service) List(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	orgs, total, err := s.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []OrganizationResponse
	for _, o := range orgs {
		responses = append(responses, o.ToResponse())
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

func (s *Service) GetTree(ctx context.Context) ([]OrganizationResponse, error) {
	tree, err := s.repo.FindTree(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]OrganizationResponse, 0, len(tree))
	for _, org := range tree {
		responses = append(responses, org.ToResponse())
	}
	return responses, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateOrganizationRequest) (*OrganizationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	org, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		org.Code = *req.Code
	}
	if req.Nomenclature != nil {
		org.Nomenclature = *req.Nomenclature
	}

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, err
	}

	response := org.ToResponse()
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.SoftDelete(ctx, uid)
}
