package company

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database Company.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create menyimpan company baru ke database.
func (r *Repository) Create(company *Company) error {
	if err := r.db.Create(company).Error; err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// FindByID mencari company berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*Company, error) {
	var company Company
	if err := r.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &company, nil
}

// FindBySlug mencari company berdasarkan slug.
func (r *Repository) FindBySlug(slug string) (*Company, error) {
	var company Company
	if err := r.db.Where("slug = ?", slug).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &company, nil
}

// FindAll mengembalikan semua company dengan pagination.
func (r *Repository) FindAll(page, perPage int) ([]Company, int64, error) {
	var companies []Company
	var total int64

	query := r.db.Model(&Company{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	// Get paginated data
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}

	return companies, total, nil
}

// Update mengupdate company.
func (r *Repository) Update(company *Company) error {
	if err := r.db.Save(company).Error; err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

// SoftDelete melakukan soft delete company.
func (r *Repository) SoftDelete(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&Company{}).Error; err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
