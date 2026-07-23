package employeemovement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk database operations Employee Movement & Career Management.
type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

// NewRepository membuat Repository baru.
func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Employee Movement
// =========================================================================

func (r *Repository) CreateMovement(ctx context.Context, m *EmployeeMovement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(m).Error; err != nil {
		return fmt.Errorf("failed to create employee movement: %w", err)
	}
	return nil
}

func (r *Repository) FindMovementByID(ctx context.Context, id uuid.UUID) (*EmployeeMovement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var m EmployeeMovement
	if err := db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, fmt.Errorf("employee movement not found: %w", err)
	}
	return &m, nil
}

func (r *Repository) FindMovementsByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]EmployeeMovement, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var movements []EmployeeMovement
	var total int64

	query := db.Model(&EmployeeMovement{}).Where("employee_id = ?", employeeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employee movements: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&movements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employee movements: %w", err)
	}
	return movements, total, nil
}

func (r *Repository) ListMovements(ctx context.Context, page, perPage int) ([]EmployeeMovement, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var movements []EmployeeMovement
	var total int64

	query := db.Model(&EmployeeMovement{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movements: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&movements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list movements: %w", err)
	}
	return movements, total, nil
}

func (r *Repository) UpdateMovement(ctx context.Context, m *EmployeeMovement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(m).Error; err != nil {
		return fmt.Errorf("failed to update employee movement: %w", err)
	}
	return nil
}

func (r *Repository) DeleteMovement(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&EmployeeMovement{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee movement not found")
	}
	return nil
}

// =========================================================================
// Employee Contract
// =========================================================================

func (r *Repository) CreateContract(ctx context.Context, c *EmployeeContract) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(c).Error; err != nil {
		return fmt.Errorf("failed to create employee contract: %w", err)
	}
	return nil
}

func (r *Repository) FindContractByID(ctx context.Context, id uuid.UUID) (*EmployeeContract, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c EmployeeContract
	if err := db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, fmt.Errorf("employee contract not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) FindContractsByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]EmployeeContract, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var contracts []EmployeeContract
	var total int64

	query := db.Model(&EmployeeContract{}).Where("employee_id = ?", employeeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employee contracts: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employee contracts: %w", err)
	}
	return contracts, total, nil
}

func (r *Repository) ListContracts(ctx context.Context, page, perPage int) ([]EmployeeContract, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var contracts []EmployeeContract
	var total int64

	query := db.Model(&EmployeeContract{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list contracts: %w", err)
	}
	return contracts, total, nil
}

func (r *Repository) UpdateContract(ctx context.Context, c *EmployeeContract) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(c).Error; err != nil {
		return fmt.Errorf("failed to update employee contract: %w", err)
	}
	return nil
}

func (r *Repository) DeleteContract(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&EmployeeContract{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee contract: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee contract not found")
	}
	return nil
}

// =========================================================================
// Approval flows
// =========================================================================

func (r *Repository) ApproveMovement(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.Model(&EmployeeMovement{}).
		Where("id = ? AND status = ?", id, MovementStatusDraft).
		Updates(map[string]interface{}{
			"status":      MovementStatusApproved,
			"approved_by": approvedBy.String(),
			"approved_at": now,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to approve movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("movement not found or not in draft status")
	}
	return nil
}

func (r *Repository) ExecuteMovement(ctx context.Context, id uuid.UUID, executedBy uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.Model(&EmployeeMovement{}).
		Where("id = ? AND status = ?", id, MovementStatusApproved).
		Updates(map[string]interface{}{
			"status":      MovementStatusExecuted,
			"executed_by": executedBy.String(),
			"executed_at": now,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to execute movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("movement not found or not in approved status")
	}
	return nil
}

func (r *Repository) CancelMovement(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Model(&EmployeeMovement{}).
		Where("id = ? AND status IN ?", id, []MovementStatus{MovementStatusDraft, MovementStatusApproved}).
		Update("status", MovementStatusCancelled)
	if result.Error != nil {
		return fmt.Errorf("failed to cancel movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("movement not found or already executed/cancelled")
	}
	return nil
}

// ExtendContract membuat kontrak baru sebagai perpanjangan dari kontrak sebelumnya.
func (r *Repository) ExtendContract(ctx context.Context, newContract *EmployeeContract, previousID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Set previous contract as extended
	if err := tx.Model(&EmployeeContract{}).
		Where("id = ?", previousID).
		Update("status", ContractStatusExtended).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update previous contract status: %w", err)
	}

	// Create new contract
	if err := tx.Create(newContract).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new contract: %w", err)
	}

	return tx.Commit().Error
}
