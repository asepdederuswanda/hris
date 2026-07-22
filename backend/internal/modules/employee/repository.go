package employee

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
// Employee
// =========================================================================

func (r *Repository) CreateEmployee(ctx context.Context, emp *Employee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(emp).Error
}

func (r *Repository) FindEmployeeByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employee
	q := db.Preload("Addresses").Preload("EmergencyContacts").
		Preload("Families").Preload("Educations").
		Preload("Experiences").Preload("Documents").
		Preload("Insurances").Preload("Employments")
	if err := q.First(&emp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return &emp, nil
}

func (r *Repository) FindEmployeeByEmployeeID(ctx context.Context, employeeID string) (*Employee, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employee
	if err := db.First(&emp, "employee_id = ?", employeeID).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return &emp, nil
}

func (r *Repository) FindAllEmployees(ctx context.Context, page, perPage int) ([]Employee, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var employees []Employee
	var total int64

	query := db.Model(&Employee{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *Repository) UpdateEmployee(ctx context.Context, emp *Employee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(emp).Error
}

func (r *Repository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Employee{}).Error
}

// =========================================================================
// Addresses (nested under employee)
// =========================================================================

func (r *Repository) CreateAddress(ctx context.Context, addr *EmployeeAddress) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(addr).Error
}

func (r *Repository) FindAddressByID(ctx context.Context, id uuid.UUID) (*EmployeeAddress, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var addr EmployeeAddress
	if err := db.First(&addr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
	}
	return &addr, nil
}

func (r *Repository) UpdateAddress(ctx context.Context, addr *EmployeeAddress) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(addr).Error
}

func (r *Repository) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeAddress{}).Error
}

// =========================================================================
// Emergency Contacts
// =========================================================================

func (r *Repository) CreateEmergencyContact(ctx context.Context, contact *EmergencyContact) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(contact).Error
}

func (r *Repository) FindEmergencyContactByID(ctx context.Context, id uuid.UUID) (*EmergencyContact, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var contact EmergencyContact
	if err := db.First(&contact, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("emergency contact not found: %w", err)
	}
	return &contact, nil
}

func (r *Repository) UpdateEmergencyContact(ctx context.Context, contact *EmergencyContact) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(contact).Error
}

func (r *Repository) DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmergencyContact{}).Error
}

// =========================================================================
// Families
// =========================================================================

func (r *Repository) CreateFamily(ctx context.Context, fam *EmployeeFamily) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(fam).Error
}

func (r *Repository) FindFamilyByID(ctx context.Context, id uuid.UUID) (*EmployeeFamily, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var fam EmployeeFamily
	if err := db.First(&fam, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("family not found: %w", err)
	}
	return &fam, nil
}

func (r *Repository) UpdateFamily(ctx context.Context, fam *EmployeeFamily) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(fam).Error
}

func (r *Repository) DeleteFamily(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeFamily{}).Error
}

// =========================================================================
// Educations
// =========================================================================

func (r *Repository) CreateEducation(ctx context.Context, edu *EmployeeEducation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(edu).Error
}

func (r *Repository) FindEducationByID(ctx context.Context, id uuid.UUID) (*EmployeeEducation, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var edu EmployeeEducation
	if err := db.First(&edu, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("education not found: %w", err)
	}
	return &edu, nil
}

func (r *Repository) UpdateEducation(ctx context.Context, edu *EmployeeEducation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(edu).Error
}

func (r *Repository) DeleteEducation(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeEducation{}).Error
}

// =========================================================================
// Experiences
// =========================================================================

func (r *Repository) CreateExperience(ctx context.Context, exp *EmployeeExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(exp).Error
}

func (r *Repository) FindExperienceByID(ctx context.Context, id uuid.UUID) (*EmployeeExperience, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var exp EmployeeExperience
	if err := db.First(&exp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("experience not found: %w", err)
	}
	return &exp, nil
}

func (r *Repository) UpdateExperience(ctx context.Context, exp *EmployeeExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(exp).Error
}

func (r *Repository) DeleteExperience(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeExperience{}).Error
}

// =========================================================================
// Documents
// =========================================================================

func (r *Repository) CreateDocument(ctx context.Context, doc *EmployeeDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(doc).Error
}

func (r *Repository) FindDocumentByID(ctx context.Context, id uuid.UUID) (*EmployeeDocument, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var doc EmployeeDocument
	if err := db.First(&doc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	return &doc, nil
}

func (r *Repository) UpdateDocument(ctx context.Context, doc *EmployeeDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(doc).Error
}

func (r *Repository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeDocument{}).Error
}

// =========================================================================
// Insurances
// =========================================================================

func (r *Repository) CreateInsurance(ctx context.Context, ins *EmployeeInsurance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(ins).Error
}

func (r *Repository) FindInsuranceByID(ctx context.Context, id uuid.UUID) (*EmployeeInsurance, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ins EmployeeInsurance
	if err := db.First(&ins, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("insurance not found: %w", err)
	}
	return &ins, nil
}

func (r *Repository) UpdateInsurance(ctx context.Context, ins *EmployeeInsurance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(ins).Error
}

func (r *Repository) DeleteInsurance(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeInsurance{}).Error
}

// =========================================================================
// Employments
// =========================================================================

func (r *Repository) CreateEmployment(ctx context.Context, emp *Employment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(emp).Error
}

func (r *Repository) FindEmploymentByID(ctx context.Context, id uuid.UUID) (*Employment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employment
	if err := db.First(&emp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employment not found: %w", err)
	}
	return &emp, nil
}

func (r *Repository) UpdateEmployment(ctx context.Context, emp *Employment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(emp).Error
}

func (r *Repository) DeleteEmployment(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Employment{}).Error
}
