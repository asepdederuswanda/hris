package employee

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all employee models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&Employee{},
		&EmployeeAddress{},
		&EmergencyContact{},
		&EmployeeFamily{},
		&EmployeeEducation{},
		&EmployeeExperience{},
		&EmployeeDocument{},
		&EmployeeInsurance{},
		&Employment{},
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

// createTestEmployee inserts a test employee and returns it.
func createTestEmployee(ctx context.Context, repo *Repository) *Employee {
	emp := &Employee{
		EmployeeID: "EMP-TEST-001",
		Name:       "Test Employee",
		Gender:     strPtr("M"),
		Status:     "active",
	}
	if err := repo.CreateEmployee(ctx, emp); err != nil {
		panic(fmt.Sprintf("failed to create test employee: %v", err))
	}
	return emp
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}


