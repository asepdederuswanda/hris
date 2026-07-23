package employeemovement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service untuk business logic Employee Movement & Career Management.
type Service struct {
	repo   *Repository
	logger *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// =========================================================================
// Employee Movement
// =========================================================================

// CreateMovement membuat pergerakan karyawan baru.
func (s *Service) CreateMovement(ctx context.Context, req CreateMovementRequest) (*MovementResponse, error) {
	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	movement := &EmployeeMovement{
		EmployeeID:           employeeUUID,
		MovementType:         MovementType(req.MovementType),
		DecisionLetterNumber: req.DecisionLetterNumber,
		DecisionLetterDate:   req.DecisionLetterDate,
		EffectiveDate:        req.EffectiveDate,
		Status:               MovementStatusDraft,
	}

	if req.FromEmploymentID != nil {
		if uid, err := uuid.Parse(*req.FromEmploymentID); err == nil {
			movement.FromEmploymentID = &uid
		}
	}
	if req.ToEmploymentID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentID); err == nil {
			movement.ToEmploymentID = &uid
		}
	}
	if req.FromOrganizationID != nil {
		if uid, err := uuid.Parse(*req.FromOrganizationID); err == nil {
			movement.FromOrganizationID = &uid
		}
	}
	if req.ToOrganizationID != nil {
		if uid, err := uuid.Parse(*req.ToOrganizationID); err == nil {
			movement.ToOrganizationID = &uid
		}
	}
	if req.FromPositionID != nil {
		if uid, err := uuid.Parse(*req.FromPositionID); err == nil {
			movement.FromPositionID = &uid
		}
	}
	if req.ToPositionID != nil {
		if uid, err := uuid.Parse(*req.ToPositionID); err == nil {
			movement.ToPositionID = &uid
		}
	}
	if req.FromEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.FromEmploymentStatusID); err == nil {
			movement.FromEmploymentStatusID = &uid
		}
	}
	if req.ToEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentStatusID); err == nil {
			movement.ToEmploymentStatusID = &uid
		}
	}
	if req.Reason != nil {
		movement.Reason = req.Reason
	}
	if req.Notes != nil {
		movement.Notes = req.Notes
	}

	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement created",
		zap.String("employee_id", req.EmployeeID),
		zap.String("movement_type", req.MovementType),
		zap.String("movement_id", movement.ID.String()),
	)

	response := movement.ToResponse()
	return &response, nil
}

// GetMovementByID mengembalikan pergerakan berdasarkan ID.
func (s *Service) GetMovementByID(ctx context.Context, id string) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := movement.ToResponse()
	return &response, nil
}

// ListMovementsByEmployee mengembalikan daftar pergerakan untuk seorang karyawan.
func (s *Service) ListMovementsByEmployee(ctx context.Context, employeeID string, page, perPage int) (*PaginatedMovementResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	movements, total, err := s.repo.FindMovementsByEmployeeID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []MovementResponse
	for _, m := range movements {
		responses = append(responses, m.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListMovements mengembalikan daftar semua pergerakan dengan pagination.
func (s *Service) ListMovements(ctx context.Context, page, perPage int) (*PaginatedMovementResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	movements, total, err := s.repo.ListMovements(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []MovementResponse
	for _, m := range movements {
		responses = append(responses, m.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateMovement mengupdate pergerakan.
func (s *Service) UpdateMovement(ctx context.Context, id string, req UpdateMovementRequest) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if movement.Status != MovementStatusDraft {
		return nil, fmt.Errorf("cannot update movement with status '%s', only draft movements can be updated", movement.Status)
	}

	if req.MovementType != nil {
		movement.MovementType = MovementType(*req.MovementType)
	}
	if req.ToOrganizationID != nil {
		if uid, err := uuid.Parse(*req.ToOrganizationID); err == nil {
			movement.ToOrganizationID = &uid
		}
	}
	if req.ToPositionID != nil {
		if uid, err := uuid.Parse(*req.ToPositionID); err == nil {
			movement.ToPositionID = &uid
		}
	}
	if req.ToEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentStatusID); err == nil {
			movement.ToEmploymentStatusID = &uid
		}
	}
	if req.Reason != nil {
		movement.Reason = req.Reason
	}
	if req.DecisionLetterNumber != nil {
		movement.DecisionLetterNumber = *req.DecisionLetterNumber
	}
	if req.DecisionLetterDate != nil {
		movement.DecisionLetterDate = *req.DecisionLetterDate
	}
	if req.EffectiveDate != nil {
		movement.EffectiveDate = *req.EffectiveDate
	}
	if req.Status != nil {
		movement.Status = MovementStatus(*req.Status)
	}
	if req.Notes != nil {
		movement.Notes = req.Notes
	}

	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return nil, err
	}

	response := movement.ToResponse()
	return &response, nil
}

// DeleteMovement menghapus pergerakan (hanya draft).
func (s *Service) DeleteMovement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return err
	}

	if movement.Status != MovementStatusDraft {
		return fmt.Errorf("cannot delete movement with status '%s', only draft movements can be deleted", movement.Status)
	}

	return s.repo.DeleteMovement(ctx, uid)
}

// ApproveMovement menyetujui pergerakan.
func (s *Service) ApproveMovement(ctx context.Context, id string, approvedBy string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	approverUUID, err := uuid.Parse(approvedBy)
	if err != nil {
		return fmt.Errorf("invalid approver id: %w", err)
	}

	return s.repo.ApproveMovement(ctx, uid, approverUUID)
}

// ExecuteMovement mengeksekusi pergerakan.
func (s *Service) ExecuteMovement(ctx context.Context, id string, executedBy string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	executorUUID, err := uuid.Parse(executedBy)
	if err != nil {
		return fmt.Errorf("invalid executor id: %w", err)
	}

	return s.repo.ExecuteMovement(ctx, uid, executorUUID)
}

// CancelMovement membatalkan pergerakan.
func (s *Service) CancelMovement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	return s.repo.CancelMovement(ctx, uid)
}

// =========================================================================
// Employee Contract
// =========================================================================

// CreateContract membuat kontrak karyawan baru.
func (s *Service) CreateContract(ctx context.Context, req CreateContractRequest) (*ContractResponse, error) {
	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	contract := &EmployeeContract{
		EmployeeID:     employeeUUID,
		ContractNumber: req.ContractNumber,
		ContractType:   ContractType(req.ContractType),
		StartDate:      req.StartDate,
		Status:         ContractStatusActive,
	}

	if req.EndDate != nil {
		contract.EndDate = req.EndDate
	}
	if req.PreviousContractID != nil {
		if uid, err := uuid.Parse(*req.PreviousContractID); err == nil {
			contract.PreviousContractID = &uid
		}
	}
	if req.DecisionLetterNumber != nil {
		contract.DecisionLetterNumber = req.DecisionLetterNumber
	}
	if req.Notes != nil {
		contract.Notes = req.Notes
	}
	if req.DocumentURL != nil {
		contract.DocumentURL = req.DocumentURL
	}

	// Jika ada previous_contract_id, gunakan ExtendContract flow
	if contract.PreviousContractID != nil {
		if err := s.repo.ExtendContract(ctx, contract, *contract.PreviousContractID); err != nil {
			return nil, err
		}
		contract.ExtensionCount = 1 // akan dihitung manual oleh caller untuk extension > 1
	} else {
		if err := s.repo.CreateContract(ctx, contract); err != nil {
			return nil, err
		}
	}

	s.logger.Info("Employee contract created",
		zap.String("employee_id", req.EmployeeID),
		zap.String("contract_number", req.ContractNumber),
		zap.String("contract_type", req.ContractType),
	)

	response := contract.ToResponse()
	return &response, nil
}

// GetContractByID mengembalikan kontrak berdasarkan ID.
func (s *Service) GetContractByID(ctx context.Context, id string) (*ContractResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}

	contract, err := s.repo.FindContractByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := contract.ToResponse()
	return &response, nil
}

// ListContractsByEmployee mengembalikan daftar kontrak untuk seorang karyawan.
func (s *Service) ListContractsByEmployee(ctx context.Context, employeeID string, page, perPage int) (*PaginatedContractResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	contracts, total, err := s.repo.FindContractsByEmployeeID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []ContractResponse
	for _, c := range contracts {
		responses = append(responses, c.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedContractResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListContracts mengembalikan daftar semua kontrak dengan pagination.
func (s *Service) ListContracts(ctx context.Context, page, perPage int) (*PaginatedContractResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	contracts, total, err := s.repo.ListContracts(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []ContractResponse
	for _, c := range contracts {
		responses = append(responses, c.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedContractResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateContract mengupdate kontrak.
func (s *Service) UpdateContract(ctx context.Context, id string, req UpdateContractRequest) (*ContractResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}

	contract, err := s.repo.FindContractByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if req.ContractNumber != nil {
		contract.ContractNumber = *req.ContractNumber
	}
	if req.ContractType != nil {
		contract.ContractType = ContractType(*req.ContractType)
	}
	if req.EndDate != nil {
		contract.EndDate = req.EndDate
	}
	if req.DecisionLetterNumber != nil {
		contract.DecisionLetterNumber = req.DecisionLetterNumber
	}
	if req.Notes != nil {
		contract.Notes = req.Notes
	}
	if req.DocumentURL != nil {
		contract.DocumentURL = req.DocumentURL
	}
	if req.Status != nil {
		contract.Status = ContractStatus(*req.Status)
	}

	if err := s.repo.UpdateContract(ctx, contract); err != nil {
		return nil, err
	}

	response := contract.ToResponse()
	return &response, nil
}

// DeleteContract menghapus kontrak.
func (s *Service) DeleteContract(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid contract id: %w", err)
	}

	return s.repo.DeleteContract(ctx, uid)
}
