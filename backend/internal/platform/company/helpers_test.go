package company

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// FakeTenantManager adalah implementasi mock dari TenantManager interface
// untuk keperluan unit test. Tidak memerlukan koneksi database nyata.
type FakeTenantManager struct {
	DriverName            string
	ProvisionTenantFunc   func(companyID, dbName, dbUser, dbPassword, driverType string) (*database.TenantConnection, error)
	DropTenantDBFunc      func(companyID string) error
	RemoveTenantConnFunc  func(companyID string) error
	DeactivateTenantFunc  func(companyID string) error
	ActivateTenantFunc    func(companyID string) error
	SaveTenantConnFunc    func(conn *database.TenantConnection) error
	TenantDBFunc          func(companyID string) (*gorm.DB, error)
}

func (f *FakeTenantManager) Driver() string {
	if f.DriverName != "" {
		return f.DriverName
	}
	return "postgres"
}

func (f *FakeTenantManager) ProvisionTenant(companyID, dbName, dbUser, dbPassword, driverType string) (*database.TenantConnection, error) {
	if f.ProvisionTenantFunc != nil {
		return f.ProvisionTenantFunc(companyID, dbName, dbUser, dbPassword, driverType)
	}
	return &database.TenantConnection{
		ID:        companyID,
		CompanyID: companyID,
		Driver:    driverType,
		DBName:    dbName,
		IsActive:  true,
	}, nil
}

func (f *FakeTenantManager) SaveTenantConnection(conn *database.TenantConnection) error {
	if f.SaveTenantConnFunc != nil {
		return f.SaveTenantConnFunc(conn)
	}
	return nil
}

func (f *FakeTenantManager) TenantDB(companyID string) (*gorm.DB, error) {
	if f.TenantDBFunc != nil {
		return f.TenantDBFunc(companyID)
	}
	return nil, nil
}

func (f *FakeTenantManager) DeactivateTenantConnection(companyID string) error {
	if f.DeactivateTenantFunc != nil {
		return f.DeactivateTenantFunc(companyID)
	}
	return nil
}

func (f *FakeTenantManager) DropTenantDB(companyID string) error {
	if f.DropTenantDBFunc != nil {
		return f.DropTenantDBFunc(companyID)
	}
	return nil
}

func (f *FakeTenantManager) RemoveTenantConnection(companyID string) error {
	if f.RemoveTenantConnFunc != nil {
		return f.RemoveTenantConnFunc(companyID)
	}
	return nil
}

func (f *FakeTenantManager) ActivateTenantConnection(companyID string) error {
	if f.ActivateTenantFunc != nil {
		return f.ActivateTenantFunc(companyID)
	}
	return nil
}

// setupTestDB creates an in-memory SQLite database and auto-migrates Company model.
// Returns the GORM DB and a cleanup function.
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(&Company{}); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

// createTestCompany inserts a test company and returns it.
// Status defaults to "active".
func createTestCompany(db *gorm.DB, name string) *Company {
	c := &Company{
		Name:   name,
		Slug:   name, // will be overwritten by test logic
		Status: CompanyStatusActive,
	}
	if err := db.Create(c).Error; err != nil {
		panic(fmt.Sprintf("failed to create test company: %v", err))
	}
	return c
}

// newTestService creates a Service with SQLite repository and FakeTenantManager.
// Returns the service and a cleanup function.
func newTestService() (*Service, *FakeTenantManager, func()) {
	db, cleanup := setupTestDB()
	repo := NewRepository(db)
	logger, _ := zap.NewDevelopment()
	fakeTM := &FakeTenantManager{}
	svc := NewService(repo, fakeTM, logger)
	return svc, fakeTM, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}
