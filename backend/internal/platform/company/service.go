package company

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/migrator"
)

// TenantManager mendefinisikan interface untuk operasi lifecycle tenant
// yang digunakan oleh Company Service. Memungkinkan mocking di unit test
// tanpa memerlukan koneksi database nyata.
type TenantManager interface {
	Driver() string
	ProvisionTenant(companyID, dbName, dbUser, dbPassword, driverType string) (*database.TenantConnection, error)
	SaveTenantConnection(conn *database.TenantConnection) error
	TenantDB(companyID string) (*gorm.DB, error)
	DeactivateTenantConnection(companyID string) error
	DropTenantDB(companyID string) error
	RemoveTenantConnection(companyID string) error
	ActivateTenantConnection(companyID string) error
}

// Service untuk business logic Company.
type Service struct {
	repo      *Repository
	dbManager TenantManager
	logger    *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, dbManager TenantManager, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		dbManager: dbManager,
		logger:    logger,
	}
}

// Create membuat company baru dan melakukan provisioning tenant database.
func (s *Service) Create(req CreateCompanyRequest) (*CompanyResponse, error) {
	// Generate slug dari name
	companySlug := slug.Make(req.Name)

	// Cek duplikasi slug
	if existing, _ := s.repo.FindBySlug(companySlug); existing != nil {
		companySlug = fmt.Sprintf("%s-%s", companySlug, uuid.New().String()[:8])
	}

	company := &Company{
		Name:    req.Name,
		Slug:    companySlug,
		NPWP:    req.NPWP,
		NIB:     req.NIB,
		Address: req.Address,
		Email:   req.Email,
		Phone:   req.Phone,
		Status:  CompanyStatusActive,
	}

	if err := s.repo.Create(company); err != nil {
		return nil, err
	}

	s.logger.Info("Company created",
		zap.String("company_id", company.ID.String()),
		zap.String("name", company.Name),
		zap.String("slug", company.Slug),
	)

	// Tenant provisioning: create database + run migrations
	if err := s.provisionTenant(company); err != nil {
		// Company tetap created, tapi provisioning gagal — log error
		s.logger.Error("Tenant provisioning failed",
			zap.String("company_id", company.ID.String()),
			zap.Error(err),
		)
		// Set company status to suspended karena tenant gagal di-provision
		company.Status = CompanyStatusSuspended
		_ = s.repo.Update(company)
	} else {
		s.logger.Info("Tenant provisioning completed",
			zap.String("company_id", company.ID.String()),
		)
	}

	response := company.ToResponse()
	return &response, nil
}

// provisionTenant membuat database tenant dan menjalankan migrations.
func (s *Service) provisionTenant(company *Company) error {
	s.logger.Info("Starting tenant provisioning",
		zap.String("company_id", company.ID.String()),
		zap.String("company_name", company.Name),
	)

	// 1. Generate database name dari company slug
	dbName := fmt.Sprintf("hris_%s", company.Slug)

	// 2. Create database via superuser connection
	// Untuk development, tenant DB menggunakan kredensial superuser yang sama
	// Production: buat dedicated database user per tenant
	conn, err := s.dbManager.ProvisionTenant(
		company.ID.String(),
		dbName,
		"root",   // username — development: root user
		"",       // password — development: empty
		s.dbManager.Driver(),
	)
	if err != nil {
		return fmt.Errorf("failed to create tenant database: %w", err)
	}

	// 3. Simpan TenantConnection ke platform DB
	if err := s.dbManager.SaveTenantConnection(conn); err != nil {
		return fmt.Errorf("failed to save tenant connection: %w", err)
	}

	// 4. Dapatkan koneksi GORM ke tenant database
	tenantDB, err := s.dbManager.TenantDB(company.ID.String())
	if err != nil {
		return fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	// 5. Jalankan tenant SQL migrations (pilih dialect sesuai driver)
	s.logger.Info("Running tenant SQL migrations...")
	tenantRoot := migrator.TenantRootPath(s.dbManager.Driver())
	tenantMigrator := migrator.New(tenantDB, s.logger, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		return fmt.Errorf("tenant migration failed: %w", err)
	}

	s.logger.Info("Tenant provisioning completed successfully",
		zap.String("db_name", dbName),
	)

	return nil
}

// MigrateTenantDB menjalankan migration pada tenant database yang sudah ada.
// Digunakan untuk upgrade schema tenant yang sudah ada.
func (s *Service) MigrateTenantDB(companyID string, db *gorm.DB) error {
	tenantRoot := migrator.TenantRootPath(s.dbManager.Driver())
	tenantMigrator := migrator.New(db, s.logger, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		return fmt.Errorf("tenant migration failed for company %s: %w", companyID, err)
	}
	return nil
}

// GetByID mengembalikan company berdasarkan ID.
func (s *Service) GetByID(id string) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := company.ToResponse()
	return &response, nil
}

// List mengembalikan daftar company dengan pagination.
func (s *Service) List(page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	companies, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var responses []CompanyResponse
	for _, c := range companies {
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

// Update mengupdate company.
func (s *Service) Update(id string, req UpdateCompanyRequest) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.NPWP != nil {
		company.NPWP = req.NPWP
	}
	if req.NIB != nil {
		company.NIB = req.NIB
	}
	if req.Address != nil {
		company.Address = req.Address
	}
	if req.Email != nil {
		company.Email = req.Email
	}
	if req.Phone != nil {
		company.Phone = req.Phone
	}

	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	response := company.ToResponse()
	return &response, nil
}

// Delete melakukan soft delete company dan menonaktifkan tenant connection.
func (s *Service) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return err
	}

	// Deactivate tenant connection (jika ada)
	if err := s.dbManager.DeactivateTenantConnection(company.ID.String()); err != nil {
		// Log warning jika tidak ada connection (misal provisioning gagal)
		s.logger.Warn("Failed to deactivate tenant connection during soft delete",
			zap.String("company_id", company.ID.String()),
			zap.Error(err),
		)
	}

	// Soft delete company record
	return s.repo.SoftDelete(uid)
}

// UpdateStatus mengupdate status company (active, suspended, terminated).
func (s *Service) UpdateStatus(id string, status string) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	company.Status = CompanyStatus(status)
	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	s.logger.Info("Company status updated",
		zap.String("company_id", company.ID.String()),
		zap.String("status", status),
	)

	response := company.ToResponse()
	return &response, nil
}

// Suspend menonaktifkan tenant: set status suspended + deactivate tenant connection.
func (s *Service) Suspend(id string) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Cek status saat ini — hanya active yang bisa di-suspend
	if company.Status != CompanyStatusActive {
		return nil, fmt.Errorf("company is not active, current status: %s", company.Status)
	}

	// 1. Deactivate tenant connection (is_active = false)
	if err := s.dbManager.DeactivateTenantConnection(company.ID.String()); err != nil {
		// Log error tapi lanjutkan — connection mungkin sudah tidak ada
		s.logger.Warn("Failed to deactivate tenant connection",
			zap.String("company_id", company.ID.String()),
			zap.Error(err),
		)
	}

	// 2. Update company status
	company.Status = CompanyStatusSuspended
	if err := s.repo.Update(company); err != nil {
		return nil, fmt.Errorf("failed to update company status: %w", err)
	}

	s.logger.Info("Company suspended",
		zap.String("company_id", company.ID.String()),
	)

	response := company.ToResponse()
	return &response, nil
}

// Activate mengaktifkan kembali tenant: set status active + reactivate tenant connection.
func (s *Service) Activate(id string) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Cek status — hanya suspended yang bisa di-activate
	if company.Status != CompanyStatusSuspended {
		return nil, fmt.Errorf("company is not suspended, current status: %s", company.Status)
	}

	// 1. Reactivate tenant connection (is_active = true)
	if err := s.dbManager.ActivateTenantConnection(company.ID.String()); err != nil {
		return nil, fmt.Errorf("failed to reactivate tenant connection: %w", err)
	}

	// 2. Update company status
	company.Status = CompanyStatusActive
	if err := s.repo.Update(company); err != nil {
		return nil, fmt.Errorf("failed to update company status: %w", err)
	}

	s.logger.Info("Company activated",
		zap.String("company_id", company.ID.String()),
	)

	response := company.ToResponse()
	return &response, nil
}

// Terminate menghapus tenant secara permanen: drop database + remove connection + set status terminated.
func (s *Service) Terminate(id string) (*CompanyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	company, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Cek status — active/suspended bisa di-terminate, terminated tidak perlu
	if company.Status == CompanyStatusTerminated {
		return nil, fmt.Errorf("company is already terminated")
	}

	// 1. Drop tenant database
	if err := s.dbManager.DropTenantDB(company.ID.String()); err != nil {
		// Log error tapi lanjutkan — DB mungkin sudah dihapus
		s.logger.Warn("Failed to drop tenant database",
			zap.String("company_id", company.ID.String()),
			zap.Error(err),
		)
	}

	// 2. Remove tenant connection record
	if err := s.dbManager.RemoveTenantConnection(company.ID.String()); err != nil {
		// Log error tapi lanjutkan — connection mungkin sudah tidak ada
		s.logger.Warn("Failed to remove tenant connection",
			zap.String("company_id", company.ID.String()),
			zap.Error(err),
		)
	}

	// 3. Update company status
	company.Status = CompanyStatusTerminated
	if err := s.repo.Update(company); err != nil {
		return nil, fmt.Errorf("failed to update company status: %w", err)
	}

	s.logger.Info("Company terminated",
		zap.String("company_id", company.ID.String()),
	)

	response := company.ToResponse()
	return &response, nil
}
