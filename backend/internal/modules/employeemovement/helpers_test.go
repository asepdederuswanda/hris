package employeemovement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&EmployeeMovement{},
		&EmployeeContract{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// newTestService creates a Service with in-memory SQLite repository.
func newTestService() (*Service, *Repository, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// createTestMovement inserts a test employee movement for the given employee.
func createTestMovement(repo *Repository, employeeID uuid.UUID) *EmployeeMovement {
	ctx := context.Background()
	m := &EmployeeMovement{
		EmployeeID:           employeeID,
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: fmt.Sprintf("SK-%s", uuid.New().String()[:8]),
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		Status:               MovementStatusDraft,
	}
	if err := repo.CreateMovement(ctx, m); err != nil {
		panic(fmt.Sprintf("failed to create test movement: %v", err))
	}
	return m
}

// createTestContract inserts a test employee contract for the given employee.
func createTestContract(repo *Repository, employeeID uuid.UUID) *EmployeeContract {
	ctx := context.Background()
	c := &EmployeeContract{
		EmployeeID:     employeeID,
		ContractNumber: fmt.Sprintf("CTR-%s", uuid.New().String()[:8]),
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		EndDate:        strPtr("2026-12-31"),
		Status:         ContractStatusActive,
	}
	if err := repo.CreateContract(ctx, c); err != nil {
		panic(fmt.Sprintf("failed to create test contract: %v", err))
	}
	return c
}

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}
