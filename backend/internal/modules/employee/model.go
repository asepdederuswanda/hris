package employee

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Employee merepresentasikan data inti karyawan.
// Catatan: Tidak menggunakan DeletedAt (soft delete) karena SQL migration
// tidak memiliki kolom deleted_at di tabel employee.
type Employee struct {
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID      string     `gorm:"type:varchar(50);not null;index" json:"employee_id"`
	NIK             *string    `gorm:"type:varchar(16);index" json:"nik,omitempty"`
	FamilyID        *string    `gorm:"type:varchar(16)" json:"family_id,omitempty"`
	Name            string     `gorm:"type:varchar(255);not null;index" json:"name"`
	MotherName      *string    `gorm:"type:varchar(255)" json:"mother_name,omitempty"`
	Gender          *string    `gorm:"type:varchar(10)" json:"gender,omitempty"`
	NationalityType *string    `gorm:"type:varchar(10)" json:"nationality_type,omitempty"`
	NationalityID   *string    `gorm:"type:char(2);index" json:"nationality_id,omitempty"`
	POB             *string    `gorm:"type:varchar(255)" json:"pob,omitempty"`
	DOB             *string    `gorm:"type:date" json:"dob,omitempty"`
	PhoneNumber     *string    `gorm:"type:varchar(255)" json:"phone_number,omitempty"`
	Email           *string    `gorm:"type:varchar(255);uniqueIndex" json:"email,omitempty"`
	LinkedIn        *string    `gorm:"column:linkedin;type:varchar(255)" json:"linkedin,omitempty"`
	Instagram       *string    `gorm:"column:ig;type:varchar(255)" json:"ig,omitempty"`
	ProfilePicture  *string    `gorm:"type:varchar(255)" json:"profile_picture,omitempty"`
	ReligionID      *uuid.UUID `gorm:"type:char(36);index" json:"religion_id,omitempty"`
	MaritalStatusID *uuid.UUID `gorm:"type:char(36);index" json:"marital_status_id,omitempty"`
	Status          string     `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedBy       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relasi
	Addresses         []EmployeeAddress    `gorm:"foreignKey:EmployeeID" json:"addresses,omitempty"`
	EmergencyContacts []EmergencyContact   `gorm:"foreignKey:EmployeeID" json:"emergency_contacts,omitempty"`
	Families          []EmployeeFamily     `gorm:"foreignKey:EmployeeID" json:"families,omitempty"`
	Educations        []EmployeeEducation  `gorm:"foreignKey:EmployeeID" json:"educations,omitempty"`
	Experiences       []EmployeeExperience `gorm:"foreignKey:EmployeeID" json:"experiences,omitempty"`
	Documents         []EmployeeDocument   `gorm:"foreignKey:EmployeeID" json:"documents,omitempty"`
	Insurances        []EmployeeInsurance  `gorm:"foreignKey:EmployeeID" json:"insurances,omitempty"`
	Employments       []Employment         `gorm:"foreignKey:EmployeeID" json:"employments,omitempty"`
}

func (Employee) TableName() string {
	return "employees"
}

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeAddress menyimpan alamat karyawan.
type EmployeeAddress struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID  *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	Type        *string    `gorm:"type:varchar(50)" json:"type,omitempty"`
	Address     *string    `gorm:"type:varchar(255)" json:"address,omitempty"`
	ProvinceID  *string    `gorm:"type:char(2);index" json:"province_id,omitempty"`
	RegencyID   *string    `gorm:"type:char(4);index" json:"regency_id,omitempty"`
	DistrictID  *string    `gorm:"type:char(6);index" json:"district_id,omitempty"`
	VillageID   *string    `gorm:"type:char(10);index" json:"village_id,omitempty"`
	PostalCode  *string    `gorm:"type:varchar(5)" json:"postal_code,omitempty"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (EmployeeAddress) TableName() string {
	return "employee_addresses"
}

func (a *EmployeeAddress) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EmergencyContact menyimpan kontak darurat karyawan.
type EmergencyContact struct {
	ID                 uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID         *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	Name               string     `gorm:"type:varchar(255);not null" json:"name"`
	RelationshipTypeID *uuid.UUID `gorm:"type:char(36);index" json:"relationship_type_id,omitempty"`
	PhoneNumber        string     `gorm:"type:varchar(50);not null" json:"phone_number"`
	Address            *string    `gorm:"type:varchar(255)" json:"address,omitempty"`
	CreatedBy          *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (EmergencyContact) TableName() string {
	return "emergency_contacts"
}

func (c *EmergencyContact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// EmployeeFamily menyimpan data keluarga karyawan.
type EmployeeFamily struct {
	ID                 uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID         *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	NIK                *string    `gorm:"type:varchar(16)" json:"nik,omitempty"`
	Name               string     `gorm:"type:varchar(255);not null" json:"name"`
	DOB                *string    `gorm:"type:date" json:"dob,omitempty"`
	RelationshipTypeID *uuid.UUID `gorm:"type:char(36);index" json:"relationship_type_id,omitempty"`
	EducationID        *uuid.UUID `gorm:"type:char(36);index" json:"education_id,omitempty"`
	CreatedBy          *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (EmployeeFamily) TableName() string {
	return "employee_families"
}

func (f *EmployeeFamily) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// EmployeeEducation menyimpan riwayat pendidikan karyawan.
type EmployeeEducation struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID  *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	EducationID *uuid.UUID `gorm:"type:char(36);index" json:"education_id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Major       *string    `gorm:"type:varchar(255)" json:"major,omitempty"`
	GradYear    *int       `gorm:"column:graduation_year;type:year" json:"graduation_year,omitempty"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (EmployeeEducation) TableName() string {
	return "employee_educations"
}

func (e *EmployeeEducation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeExperience menyimpan riwayat pekerjaan karyawan.
type EmployeeExperience struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	Company    string     `gorm:"type:varchar(255);not null" json:"company"`
	Position   *string    `gorm:"type:varchar(255)" json:"position,omitempty"`
	StartYear  *int       `gorm:"type:year" json:"start_year,omitempty"`
	EndYear    *int       `gorm:"type:year" json:"end_year,omitempty"`
	CreatedBy  *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (EmployeeExperience) TableName() string {
	return "employee_experiences"
}

func (e *EmployeeExperience) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeDocument menyimpan dokumen karyawan.
type EmployeeDocument struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name"`
	File       string     `gorm:"type:varchar(255);not null" json:"file"`
	Note       *string    `gorm:"type:varchar(255)" json:"note,omitempty"`
	CreatedBy  *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (EmployeeDocument) TableName() string {
	return "employee_documents"
}

func (d *EmployeeDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// EmployeeInsurance menyimpan data asuransi karyawan.
type EmployeeInsurance struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	Category   *string    `gorm:"type:varchar(20)" json:"category,omitempty"`
	Number     string     `gorm:"type:varchar(100);not null" json:"number"`
	Name       string     `gorm:"type:varchar(100);not null" json:"name"`
	Type       *string    `gorm:"type:varchar(100)" json:"type,omitempty"`
	CreatedBy  *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (EmployeeInsurance) TableName() string {
	return "employee_insurances"
}

func (i *EmployeeInsurance) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// Employment menyimpan riwayat jabatan karyawan.
type Employment struct {
	ID                   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID           *uuid.UUID `gorm:"type:char(36);index" json:"employee_id,omitempty"`
	OrganizationID       *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	PositionID           *uuid.UUID `gorm:"type:char(36);index" json:"position_id,omitempty"`
	EmploymentStatusID   *uuid.UUID `gorm:"type:char(36);index" json:"employment_status_id,omitempty"`
	DecisionLetterNumber string     `gorm:"type:varchar(50);not null" json:"decision_letter_number"`
	DecisionLetterDate   string     `gorm:"type:date;not null" json:"decision_letter_date"`
	EffectiveDate        string     `gorm:"type:date;not null" json:"effective_date"`
	EffectiveEndDate     *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	CreatedBy            *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (Employment) TableName() string {
	return "employments"
}

func (e *Employment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
