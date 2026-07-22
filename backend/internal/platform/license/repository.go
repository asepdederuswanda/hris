package license

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database License.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByID mencari lisensi berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*License, error) {
	var l License
	if err := r.db.Where("id = ?", id).First(&l).Error; err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}
	return &l, nil
}

// FindByCompanyID mencari lisensi aktif untuk company.
func (r *Repository) FindByCompanyID(companyID uuid.UUID) (*License, error) {
	var l License
	if err := r.db.Where("company_id = ? AND status = ?", companyID, LicenseActive).First(&l).Error; err != nil {
		return nil, fmt.Errorf("active license not found for company: %w", err)
	}
	return &l, nil
}

// FindByLicenseKey mencari lisensi berdasarkan license key.
func (r *Repository) FindByLicenseKey(key string) (*License, error) {
	var l License
	if err := r.db.Where("license_key = ?", key).First(&l).Error; err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}
	return &l, nil
}

// FindAll mengembalikan semua lisensi dengan pagination.
func (r *Repository) FindAll(page, perPage int) ([]License, int64, error) {
	var licenses []License
	var total int64

	query := r.db.Model(&License{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count licenses: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&licenses).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list licenses: %w", err)
	}

	return licenses, total, nil
}

// Create menyimpan lisensi baru.
func (r *Repository) Create(license *License) error {
	if err := r.db.Create(license).Error; err != nil {
		return fmt.Errorf("failed to create license: %w", err)
	}
	return nil
}

// Update mengupdate lisensi.
func (r *Repository) Update(license *License) error {
	if err := r.db.Save(license).Error; err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}
	return nil
}
