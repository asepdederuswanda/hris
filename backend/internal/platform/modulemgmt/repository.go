package modulemgmt

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database Module & CompanyModule.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindBySlug mencari modul berdasarkan slug.
func (r *Repository) FindBySlug(slug string) (*PlatformModule, error) {
	var m PlatformModule
	if err := r.db.Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}
	return &m, nil
}

// FindByID mencari modul berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*PlatformModule, error) {
	var m PlatformModule
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}
	return &m, nil
}

// FindAll mengembalikan semua modul dengan pagination.
func (r *Repository) FindAll(page, perPage int) ([]PlatformModule, int64, error) {
	var modules []PlatformModule
	var total int64

	query := r.db.Model(&PlatformModule{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count modules: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&modules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list modules: %w", err)
	}

	return modules, total, nil
}

// Create menyimpan modul baru.
func (r *Repository) Create(module *PlatformModule) error {
	if err := r.db.Create(module).Error; err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}
	return nil
}

// Update mengupdate modul.
func (r *Repository) Update(module *PlatformModule) error {
	if err := r.db.Save(module).Error; err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}
	return nil
}

// UpsertCompanyModule membuat atau mengupdate relasi company-module.
func (r *Repository) UpsertCompanyModule(companyID, moduleID uuid.UUID, enabled bool) (*CompanyModule, error) {
	now := time.Now()
	cm := &CompanyModule{
		CompanyID:   companyID,
		ModuleID:    moduleID,
		Enabled:     enabled,
		ActivatedAt: &now,
	}

	// Try to find existing
	var existing CompanyModule
	err := r.db.Where("company_id = ? AND module_id = ?", companyID, moduleID).First(&existing).Error
	if err == nil {
		// Update existing
		existing.Enabled = enabled
		if enabled {
			existing.ActivatedAt = &now
		}
		if err := r.db.Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("failed to update company module: %w", err)
		}
		return &existing, nil
	}

	// Create new
	if err := r.db.Create(cm).Error; err != nil {
		return nil, fmt.Errorf("failed to create company module: %w", err)
	}
	return cm, nil
}

// FindCompanyModules mengembalikan semua modul untuk company tertentu.
func (r *Repository) FindCompanyModules(companyID uuid.UUID) ([]CompanyModuleResponse, error) {
	var results []CompanyModuleResponse

	rows, err := r.db.Table("company_modules").
		Select(`company_modules.company_id, company_modules.module_id, 
				modules.name as module_name, company_modules.enabled, 
				company_modules.activated_at`).
		Joins("JOIN modules ON modules.id = company_modules.module_id").
		Where("company_modules.company_id = ?", companyID).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list company modules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cm CompanyModuleResponse
		if err := r.db.ScanRows(rows, &cm); err != nil {
			return nil, fmt.Errorf("failed to scan company module: %w", err)
		}
		results = append(results, cm)
	}

	return results, nil
}
