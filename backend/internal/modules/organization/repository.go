package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

func (r *Repository) Create(ctx context.Context, org *Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(org).Error
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := db.Preload("Parent").First(&org, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return &org, nil
}

func (r *Repository) FindTree(ctx context.Context) ([]Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var roots []Organization
	if err := db.Where("parent_id IS NULL").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order("sort_order ASC").
		Find(&roots).Error; err != nil {
		return nil, fmt.Errorf("failed to load organization tree: %w", err)
	}
	return roots, nil
}

func (r *Repository) FindAll(ctx context.Context, page, perPage int) ([]Organization, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var orgs []Organization
	var total int64

	query := db.Model(&Organization{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (r *Repository) Update(ctx context.Context, org *Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(org).Error
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Organization{}).Error
}
