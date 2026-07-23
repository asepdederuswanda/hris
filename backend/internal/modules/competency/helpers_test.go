package competency

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all competency models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all 7 models
	if err := db.AutoMigrate(
		&Competency{},
		&CompetenceValue{},
		&CompetencyValue{},
		&CompetencyEvent{},
		&CompetencyEventTarget{},
		&CompetencyScore{},
		&CompetencyScoreDetail{},
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

// createTestCompetency inserts a test competency master and returns it.
func createTestCompetency(ctx context.Context, repo *Repository) *Competency {
	c := &Competency{
		Name: "Leadership",
	}
	if err := repo.CreateCompetency(ctx, c); err != nil {
		panic(fmt.Sprintf("failed to create test competency: %v", err))
	}
	return c
}

// createTestCompetencyEvent inserts a test competency event and returns it.
func createTestCompetencyEvent(ctx context.Context, repo *Repository) *CompetencyEvent {
	e := &CompetencyEvent{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	}
	if err := repo.CreateCompetencyEvent(ctx, e); err != nil {
		panic(fmt.Sprintf("failed to create test competency event: %v", err))
	}
	return e
}

// createTestCompetencyScore inserts a test competency score and returns it.
func createTestCompetencyScore(ctx context.Context, repo *Repository) *CompetencyScore {
	s := &CompetencyScore{
		OrganizationID:          uuid.New(),
		TechnicalGapPercentage:  10.5,
		ManagerialGapPercentage: 15.2,
		TotalGapPercentage:      12.0,
		TotalGradePercentage:    88.0,
	}
	if err := repo.CreateCompetencyScore(ctx, s); err != nil {
		panic(fmt.Sprintf("failed to create test competency score: %v", err))
	}
	return s
}

// createTestOrgID returns a UUID string for organization_id in tests.
func createTestOrgID() string {
	return uuid.New().String()
}

// createTestUUID returns a UUID string for general test use.
func createTestUUID() string {
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

// float64Ptr returns a pointer to the given float64.
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}
