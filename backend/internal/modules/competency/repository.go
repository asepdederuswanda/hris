package competency

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

// =========================================================================
// Competency CRUD
// =========================================================================

func (r *Repository) CreateCompetency(ctx context.Context, c *Competency) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(c).Error
}

func (r *Repository) FindCompetencyByID(ctx context.Context, id uuid.UUID) (*Competency, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c Competency
	if err := db.First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) FindAllCompetencies(ctx context.Context, page, perPage int) ([]Competency, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []Competency
	var total int64

	query := db.Model(&Competency{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetency(ctx context.Context, c *Competency) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(c).Error
}

func (r *Repository) DeleteCompetency(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Competency{}).Error
}

// =========================================================================
// CompetenceValue CRUD (legacy)
// =========================================================================

func (r *Repository) CreateCompetenceValue(ctx context.Context, v *CompetenceValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(v).Error
}

func (r *Repository) FindCompetenceValueByID(ctx context.Context, id uuid.UUID) (*CompetenceValue, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var v CompetenceValue
	if err := db.First(&v, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competence value not found: %w", err)
	}
	return &v, nil
}

func (r *Repository) FindAllCompetenceValues(ctx context.Context, page, perPage int) ([]CompetenceValue, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetenceValue
	var total int64

	query := db.Model(&CompetenceValue{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetenceValue(ctx context.Context, v *CompetenceValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(v).Error
}

func (r *Repository) DeleteCompetenceValue(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetenceValue{}).Error
}

// =========================================================================
// CompetencyValue CRUD (structured)
// =========================================================================

func (r *Repository) CreateCompetencyValue(ctx context.Context, v *CompetencyValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(v).Error
}

func (r *Repository) FindCompetencyValueByID(ctx context.Context, id uuid.UUID) (*CompetencyValue, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var v CompetencyValue
	if err := db.First(&v, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency value not found: %w", err)
	}
	return &v, nil
}

func (r *Repository) FindAllCompetencyValues(ctx context.Context, page, perPage int) ([]CompetencyValue, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyValue
	var total int64

	query := db.Model(&CompetencyValue{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("type ASC, level ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetencyValue(ctx context.Context, v *CompetencyValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(v).Error
}

func (r *Repository) DeleteCompetencyValue(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyValue{}).Error
}

// =========================================================================
// CompetencyEvent CRUD
// =========================================================================

func (r *Repository) CreateCompetencyEvent(ctx context.Context, e *CompetencyEvent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(e).Error
}

func (r *Repository) FindCompetencyEventByID(ctx context.Context, id uuid.UUID) (*CompetencyEvent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var e CompetencyEvent
	if err := db.First(&e, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency event not found: %w", err)
	}
	return &e, nil
}

func (r *Repository) FindAllCompetencyEvents(ctx context.Context, page, perPage int) ([]CompetencyEvent, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyEvent
	var total int64

	query := db.Model(&CompetencyEvent{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period_year DESC, period_number ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetencyEvent(ctx context.Context, e *CompetencyEvent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(e).Error
}

func (r *Repository) DeleteCompetencyEvent(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyEvent{}).Error
}

// =========================================================================
// CompetencyEventTarget CRUD
// =========================================================================

func (r *Repository) CreateCompetencyEventTarget(ctx context.Context, t *CompetencyEventTarget) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (r *Repository) FindCompetencyEventTargetByID(ctx context.Context, id uuid.UUID) (*CompetencyEventTarget, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t CompetencyEventTarget
	if err := db.Preload("CompetencyEvent").First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency event target not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) FindAllCompetencyEventTargets(ctx context.Context, page, perPage int) ([]CompetencyEventTarget, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyEventTarget
	var total int64

	query := db.Model(&CompetencyEventTarget{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetencyEventTarget(ctx context.Context, t *CompetencyEventTarget) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(t).Error
}

func (r *Repository) DeleteCompetencyEventTarget(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyEventTarget{}).Error
}

// =========================================================================
// CompetencyScore CRUD
// =========================================================================

func (r *Repository) CreateCompetencyScore(ctx context.Context, s *CompetencyScore) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) FindCompetencyScoreByID(ctx context.Context, id uuid.UUID) (*CompetencyScore, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s CompetencyScore
	if err := db.Preload("Details").First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency score not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) FindAllCompetencyScores(ctx context.Context, page, perPage int) ([]CompetencyScore, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyScore
	var total int64

	query := db.Model(&CompetencyScore{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetencyScore(ctx context.Context, s *CompetencyScore) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

func (r *Repository) DeleteCompetencyScore(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyScore{}).Error
}

// =========================================================================
// CompetencyScoreDetail CRUD
// =========================================================================

func (r *Repository) CreateCompetencyScoreDetail(ctx context.Context, d *CompetencyScoreDetail) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

func (r *Repository) FindCompetencyScoreDetailByID(ctx context.Context, id uuid.UUID) (*CompetencyScoreDetail, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var d CompetencyScoreDetail
	if err := db.First(&d, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency score detail not found: %w", err)
	}
	return &d, nil
}

func (r *Repository) FindAllCompetencyScoreDetails(ctx context.Context, scoreID uuid.UUID, page, perPage int) ([]CompetencyScoreDetail, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyScoreDetail
	var total int64

	query := db.Model(&CompetencyScoreDetail{}).Where("competency_score_id = ?", scoreID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("type ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCompetencyScoreDetail(ctx context.Context, d *CompetencyScoreDetail) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(d).Error
}

func (r *Repository) DeleteCompetencyScoreDetail(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyScoreDetail{}).Error
}
