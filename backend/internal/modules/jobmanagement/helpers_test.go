package jobmanagement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all job management models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&JobTitle{},
		&JobTitleSub{},
		&JobValue{},
		&JobObjective{},
		&JobIdentification{},
		&JobResponsibility{},
		&JobEducationExperience{},
		&JobHRAuthority{},
		&JobOperationalAuthority{},
		&JobWorkingActivity{},
		&JobWorkingRisk{},
		&JobRelationship{},
		&JobSubordinateControl{},
		&JobAsset{},
		&JobFinancial{},
		&JobPotencyCompetency{},
		&JobScore{},
		&JobCompetencyGroup{},
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

// createTestJobTitle inserts a test job title and returns it.
func createTestJobTitle(ctx context.Context, repo *Repository) *JobTitle {
	name := "Test Title"
	t := &JobTitle{
		Name:   &name,
		Status: int8Ptr(1),
	}
	if err := repo.CreateJobTitle(ctx, t); err != nil {
		panic(fmt.Sprintf("failed to create test job title: %v", err))
	}
	return t
}

// createTestJobValue inserts a test job value with a specific type.
func createTestJobValue(ctx context.Context, repo *Repository, valueType string) *JobValue {
	v := &JobValue{
		Type: valueType,
		Sort: intPtr(1),
	}
	if err := repo.CreateJobValue(ctx, v); err != nil {
		panic(fmt.Sprintf("failed to create test job value: %v", err))
	}
	return v
}

// createTestOrgID returns a UUID string for use as organization_id in tests.
func createTestOrgID() string {
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

// int8Ptr returns a pointer to the given int8.
func int8Ptr(i int8) *int8 {
	return &i
}

// float64Ptr returns a pointer to the given float64.
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}

// uint64Ptr returns a pointer to the given uint64.
func uint64Ptr(u uint64) *uint64 {
	return &u
}
