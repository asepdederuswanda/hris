package competency

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

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Competency CRUD
// =========================================================================

func (s *Service) CreateCompetency(ctx context.Context, req CreateCompetencyRequest) (*CompetencyResponse, error) {
	ent := &Competency{
		Name: req.Name,
	}
	if req.Field != nil {
		ent.Field = req.Field
	}
	if req.Cluster != nil {
		ent.Cluster = req.Cluster
	}
	if req.Definition != nil {
		ent.Definition = req.Definition
	}

	if err := s.repo.CreateCompetency(ctx, ent); err != nil {
		return nil, err
	}

	s.logger.Info("Competency created", zap.String("id", ent.ID.String()), zap.String("name", ent.Name))
	response := ent.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyByID(ctx context.Context, id string) (*CompetencyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency id: %w", err)
	}
	ent, err := s.repo.FindCompetencyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := ent.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencies(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	list, total, err := s.repo.FindAllCompetencies(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []CompetencyResponse
	for _, c := range list {
		responses = append(responses, c.ToResponse())
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

func (s *Service) UpdateCompetency(ctx context.Context, id string, req UpdateCompetencyRequest) (*CompetencyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency id: %w", err)
	}

	ent, err := s.repo.FindCompetencyByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		ent.Name = *req.Name
	}
	if req.Field != nil {
		ent.Field = req.Field
	}
	if req.Cluster != nil {
		ent.Cluster = req.Cluster
	}
	if req.Definition != nil {
		ent.Definition = req.Definition
	}

	if err := s.repo.UpdateCompetency(ctx, ent); err != nil {
		return nil, err
	}

	response := ent.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetency(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency id: %w", err)
	}
	return s.repo.DeleteCompetency(ctx, uid)
}

// =========================================================================
// CompetenceValue CRUD (legacy)
// =========================================================================

func (s *Service) CreateCompetenceValue(ctx context.Context, req CreateCompetenceValueRequest) (*CompetenceValueResponse, error) {
	v := &CompetenceValue{Name: req.Name}
	if req.Type != nil {
		v.Type = req.Type
	}
	if req.Level != nil {
		v.Level = req.Level
	}
	if req.Point != nil {
		v.Point = req.Point
	}
	if req.Description != nil {
		v.Description = req.Description
	}
	if err := s.repo.CreateCompetenceValue(ctx, v); err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetenceValueByID(ctx context.Context, id string) (*CompetenceValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competence value id: %w", err)
	}
	v, err := s.repo.FindCompetenceValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetenceValues(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	list, total, err := s.repo.FindAllCompetenceValues(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []CompetenceValueResponse
	for _, v := range list {
		responses = append(responses, v.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetenceValue(ctx context.Context, id string, req UpdateCompetenceValueRequest) (*CompetenceValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competence value id: %w", err)
	}
	v, err := s.repo.FindCompetenceValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Type != nil {
		v.Type = req.Type
	}
	if req.Level != nil {
		v.Level = req.Level
	}
	if req.Point != nil {
		v.Point = req.Point
	}
	if req.Description != nil {
		v.Description = req.Description
	}
	if err := s.repo.UpdateCompetenceValue(ctx, v); err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetenceValue(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competence value id: %w", err)
	}
	return s.repo.DeleteCompetenceValue(ctx, uid)
}

// =========================================================================
// CompetencyValue CRUD (structured)
// =========================================================================

func (s *Service) CreateCompetencyValue(ctx context.Context, req CreateCompetencyValueRequest) (*CompetencyValueResponse, error) {
	v := &CompetencyValue{
		Type:  req.Type,
		Name:  req.Name,
		Slug:  req.Slug,
		Level: req.Level,
	}
	if req.Code != nil {
		v.Code = req.Code
	}
	if req.Description != nil {
		v.Description = req.Description
	}
	if err := s.repo.CreateCompetencyValue(ctx, v); err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyValueByID(ctx context.Context, id string) (*CompetencyValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency value id: %w", err)
	}
	v, err := s.repo.FindCompetencyValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencyValues(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllCompetencyValues(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []CompetencyValueResponse
	for _, v := range list {
		responses = append(responses, v.ToResponse())
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetencyValue(ctx context.Context, id string, req UpdateCompetencyValueRequest) (*CompetencyValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency value id: %w", err)
	}
	v, err := s.repo.FindCompetencyValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Type != nil {
		v.Type = *req.Type
	}
	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Slug != nil {
		v.Slug = *req.Slug
	}
	if req.Level != nil {
		v.Level = *req.Level
	}
	if req.Code != nil {
		v.Code = req.Code
	}
	if req.Description != nil {
		v.Description = req.Description
	}
	if err := s.repo.UpdateCompetencyValue(ctx, v); err != nil {
		return nil, err
	}
	response := v.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetencyValue(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency value id: %w", err)
	}
	return s.repo.DeleteCompetencyValue(ctx, uid)
}

// =========================================================================
// CompetencyEvent CRUD
// =========================================================================

func (s *Service) CreateCompetencyEvent(ctx context.Context, req CreateCompetencyEventRequest) (*CompetencyEventResponse, error) {
	e := &CompetencyEvent{
		Type:       req.Type,
		PeriodType: req.PeriodType,
		PeriodYear: req.PeriodYear,
	}
	if req.PeriodNumber != nil {
		e.PeriodNumber = req.PeriodNumber
	}
	if req.Status != "" {
		e.Status = req.Status
	}
	if err := s.repo.CreateCompetencyEvent(ctx, e); err != nil {
		return nil, err
	}
	response := e.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyEventByID(ctx context.Context, id string) (*CompetencyEventResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency event id: %w", err)
	}
	e, err := s.repo.FindCompetencyEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := e.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencyEvents(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllCompetencyEvents(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []CompetencyEventResponse
	for _, e := range list {
		responses = append(responses, e.ToResponse())
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetencyEvent(ctx context.Context, id string, req UpdateCompetencyEventRequest) (*CompetencyEventResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency event id: %w", err)
	}
	e, err := s.repo.FindCompetencyEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Type != nil {
		e.Type = *req.Type
	}
	if req.PeriodType != nil {
		e.PeriodType = *req.PeriodType
	}
	if req.PeriodYear != nil {
		e.PeriodYear = *req.PeriodYear
	}
	if req.PeriodNumber != nil {
		e.PeriodNumber = req.PeriodNumber
	}
	if req.Status != nil {
		e.Status = *req.Status
	}
	if err := s.repo.UpdateCompetencyEvent(ctx, e); err != nil {
		return nil, err
	}
	response := e.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetencyEvent(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency event id: %w", err)
	}
	return s.repo.DeleteCompetencyEvent(ctx, uid)
}

// =========================================================================
// CompetencyEventTarget CRUD
// =========================================================================

func (s *Service) CreateCompetencyEventTarget(ctx context.Context, req CreateCompetencyEventTargetRequest) (*CompetencyEventTargetResponse, error) {
	eventUID, err := uuid.Parse(req.CompetencyEventID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_event_id: %w", err)
	}
	orgUID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}

	t := &CompetencyEventTarget{
		CompetencyEventID: eventUID,
		OrganizationID:    orgUID,
	}
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		uid, _ := uuid.Parse(*req.EmployeeID)
		t.EmployeeID = &uid
	}
	if req.MissingSelf != nil {
		t.MissingSelf = *req.MissingSelf
	}
	if req.MissingSuperior != nil {
		t.MissingSuperior = *req.MissingSuperior
	}
	if req.MissingPeer != nil {
		t.MissingPeer = *req.MissingPeer
	}
	if req.MissingSubordinate != nil {
		t.MissingSubordinate = *req.MissingSubordinate
	}
	if err := s.repo.CreateCompetencyEventTarget(ctx, t); err != nil {
		return nil, err
	}
	response := t.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyEventTargetByID(ctx context.Context, id string) (*CompetencyEventTargetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency event target id: %w", err)
	}
	t, err := s.repo.FindCompetencyEventTargetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := t.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencyEventTargets(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllCompetencyEventTargets(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []CompetencyEventTargetResponse
	for _, t := range list {
		responses = append(responses, t.ToResponse())
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetencyEventTarget(ctx context.Context, id string, req UpdateCompetencyEventTargetRequest) (*CompetencyEventTargetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency event target id: %w", err)
	}
	t, err := s.repo.FindCompetencyEventTargetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		uid, _ := uuid.Parse(*req.EmployeeID)
		t.EmployeeID = &uid
	}
	if req.MissingSelf != nil {
		t.MissingSelf = *req.MissingSelf
	}
	if req.MissingSuperior != nil {
		t.MissingSuperior = *req.MissingSuperior
	}
	if req.MissingPeer != nil {
		t.MissingPeer = *req.MissingPeer
	}
	if req.MissingSubordinate != nil {
		t.MissingSubordinate = *req.MissingSubordinate
	}
	if err := s.repo.UpdateCompetencyEventTarget(ctx, t); err != nil {
		return nil, err
	}
	response := t.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetencyEventTarget(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency event target id: %w", err)
	}
	return s.repo.DeleteCompetencyEventTarget(ctx, uid)
}

// =========================================================================
// CompetencyScore CRUD
// =========================================================================

func (s *Service) CreateCompetencyScore(ctx context.Context, req CreateCompetencyScoreRequest) (*CompetencyScoreResponse, error) {
	orgUID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	sc := &CompetencyScore{
		OrganizationID:          orgUID,
		TechnicalGapPercentage:  req.TechnicalGapPercentage,
		ManagerialGapPercentage: req.ManagerialGapPercentage,
		TotalGapPercentage:      req.TotalGapPercentage,
		TotalGradePercentage:    req.TotalGradePercentage,
	}
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		uid, _ := uuid.Parse(*req.EmployeeID)
		sc.EmployeeID = &uid
	}
	if req.CompetencyEventID != nil && *req.CompetencyEventID != "" {
		uid, _ := uuid.Parse(*req.CompetencyEventID)
		sc.CompetencyEventID = &uid
	}
	if err := s.repo.CreateCompetencyScore(ctx, sc); err != nil {
		return nil, err
	}
	response := sc.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyScoreByID(ctx context.Context, id string) (*CompetencyScoreResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency score id: %w", err)
	}
	sc, err := s.repo.FindCompetencyScoreByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := sc.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencyScores(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllCompetencyScores(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []CompetencyScoreResponse
	for _, sc := range list {
		responses = append(responses, sc.ToResponse())
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetencyScore(ctx context.Context, id string, req UpdateCompetencyScoreRequest) (*CompetencyScoreResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency score id: %w", err)
	}
	sc, err := s.repo.FindCompetencyScoreByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.TechnicalGapPercentage != nil {
		sc.TechnicalGapPercentage = *req.TechnicalGapPercentage
	}
	if req.ManagerialGapPercentage != nil {
		sc.ManagerialGapPercentage = *req.ManagerialGapPercentage
	}
	if req.TotalGapPercentage != nil {
		sc.TotalGapPercentage = *req.TotalGapPercentage
	}
	if req.TotalGradePercentage != nil {
		sc.TotalGradePercentage = *req.TotalGradePercentage
	}
	if req.CompetencyEventID != nil && *req.CompetencyEventID != "" {
		uid, _ := uuid.Parse(*req.CompetencyEventID)
		sc.CompetencyEventID = &uid
	}
	if err := s.repo.UpdateCompetencyScore(ctx, sc); err != nil {
		return nil, err
	}
	response := sc.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetencyScore(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency score id: %w", err)
	}
	return s.repo.DeleteCompetencyScore(ctx, uid)
}

// =========================================================================
// CompetencyScoreDetail CRUD
// =========================================================================

func (s *Service) CreateCompetencyScoreDetail(ctx context.Context, req CreateCompetencyScoreDetailRequest) (*CompetencyScoreDetailResponse, error) {
	scoreUID, err := uuid.Parse(req.CompetencyScoreID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_score_id: %w", err)
	}
	compUID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	d := &CompetencyScoreDetail{
		CompetencyScoreID:     scoreUID,
		CompetencyID:          compUID,
		Type:                  req.Type,
		StandardWeight:        req.StandardWeight,
		GapPercentage:         req.GapPercentage,
		WeightedGapPercentage: req.WeightedGapPercentage,
	}
	if req.StandardLevel != nil {
		d.StandardLevel = req.StandardLevel
	}
	if req.EmployeeLevel != nil {
		d.EmployeeLevel = req.EmployeeLevel
	}
	if err := s.repo.CreateCompetencyScoreDetail(ctx, d); err != nil {
		return nil, err
	}
	response := d.ToResponse()
	return &response, nil
}

func (s *Service) GetCompetencyScoreDetailByID(ctx context.Context, id string) (*CompetencyScoreDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency score detail id: %w", err)
	}
	d, err := s.repo.FindCompetencyScoreDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := d.ToResponse()
	return &response, nil
}

func (s *Service) ListCompetencyScoreDetails(ctx context.Context, scoreID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	uid, err := uuid.Parse(scoreID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_score_id: %w", err)
	}
	list, total, err := s.repo.FindAllCompetencyScoreDetails(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []CompetencyScoreDetailResponse
	for _, d := range list {
		responses = append(responses, d.ToResponse())
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateCompetencyScoreDetail(ctx context.Context, id string, req UpdateCompetencyScoreDetailRequest) (*CompetencyScoreDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid competency score detail id: %w", err)
	}
	d, err := s.repo.FindCompetencyScoreDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Type != nil {
		d.Type = *req.Type
	}
	if req.StandardLevel != nil {
		d.StandardLevel = req.StandardLevel
	}
	if req.StandardWeight != nil {
		d.StandardWeight = *req.StandardWeight
	}
	if req.EmployeeLevel != nil {
		d.EmployeeLevel = req.EmployeeLevel
	}
	if req.GapPercentage != nil {
		d.GapPercentage = *req.GapPercentage
	}
	if req.WeightedGapPercentage != nil {
		d.WeightedGapPercentage = *req.WeightedGapPercentage
	}
	if err := s.repo.UpdateCompetencyScoreDetail(ctx, d); err != nil {
		return nil, err
	}
	response := d.ToResponse()
	return &response, nil
}

func (s *Service) DeleteCompetencyScoreDetail(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid competency score detail id: %w", err)
	}
	return s.repo.DeleteCompetencyScoreDetail(ctx, uid)
}
