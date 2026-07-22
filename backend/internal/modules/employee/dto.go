package employee

import "time"

// =========================================================================
// Request DTOs — Employee
// =========================================================================

type CreateEmployeeRequest struct {
	EmployeeID      string  `json:"employee_id" binding:"required,max=50"`
	NIK             *string `json:"nik" binding:"omitempty,max=16"`
	FamilyID        *string `json:"family_id" binding:"omitempty,max=16"`
	Name            string  `json:"name" binding:"required,max=255"`
	MotherName      *string `json:"mother_name" binding:"omitempty,max=255"`
	Gender          *string `json:"gender" binding:"omitempty,oneof=M F"`
	NationalityType *string `json:"nationality_type" binding:"omitempty,oneof=WNI WNA"`
	NationalityID   *string `json:"nationality_id" binding:"omitempty,len=2"`
	POB             *string `json:"pob" binding:"omitempty,max=255"`
	DOB             *string `json:"dob" binding:"omitempty"`
	PhoneNumber     *string `json:"phone_number" binding:"omitempty,max=255"`
	Email           *string `json:"email" binding:"omitempty,email"`
	LinkedIn        *string `json:"linkedin" binding:"omitempty,max=255"`
	Instagram       *string `json:"ig" binding:"omitempty,max=255"`
	ReligionID      *string `json:"religion_id" binding:"omitempty"`
	MaritalStatusID *string `json:"marital_status_id" binding:"omitempty"`
}

type UpdateEmployeeRequest struct {
	NIK             *string `json:"nik" binding:"omitempty,max=16"`
	FamilyID        *string `json:"family_id" binding:"omitempty,max=16"`
	Name            *string `json:"name" binding:"omitempty,max=255"`
	MotherName      *string `json:"mother_name" binding:"omitempty,max=255"`
	Gender          *string `json:"gender" binding:"omitempty,oneof=M F"`
	NationalityType *string `json:"nationality_type" binding:"omitempty,oneof=WNI WNA"`
	NationalityID   *string `json:"nationality_id" binding:"omitempty,len=2"`
	POB             *string `json:"pob" binding:"omitempty,max=255"`
	DOB             *string `json:"dob" binding:"omitempty"`
	PhoneNumber     *string `json:"phone_number" binding:"omitempty,max=255"`
	Email           *string `json:"email" binding:"omitempty,email"`
	LinkedIn        *string `json:"linkedin" binding:"omitempty,max=255"`
	Instagram       *string `json:"ig" binding:"omitempty,max=255"`
	ReligionID      *string `json:"religion_id" binding:"omitempty"`
	MaritalStatusID *string `json:"marital_status_id" binding:"omitempty"`
	Status          *string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// =========================================================================
// Request DTOs — Sub-modules
// =========================================================================

type CreateAddressRequest struct {
	Type       string  `json:"type" binding:"required,oneof=MAIN DOMICILE"`
	Address    string  `json:"address" binding:"required"`
	ProvinceID *string `json:"province_id"`
	RegencyID  *string `json:"regency_id"`
	DistrictID *string `json:"district_id"`
	VillageID  *string `json:"village_id"`
	PostalCode *string `json:"postal_code" binding:"omitempty,max=5"`
}

type UpdateAddressRequest struct {
	Type       *string `json:"type" binding:"omitempty,oneof=MAIN DOMICILE"`
	Address    *string `json:"address"`
	ProvinceID *string `json:"province_id"`
	RegencyID  *string `json:"regency_id"`
	DistrictID *string `json:"district_id"`
	VillageID  *string `json:"village_id"`
	PostalCode *string `json:"postal_code" binding:"omitempty,max=5"`
}

type CreateEmergencyContactRequest struct {
	Name               string  `json:"name" binding:"required,max=255"`
	RelationshipTypeID *string `json:"relationship_type_id"`
	PhoneNumber        string  `json:"phone_number" binding:"required,max=50"`
	Address            *string `json:"address"`
}

type UpdateEmergencyContactRequest struct {
	Name               *string `json:"name" binding:"omitempty,max=255"`
	RelationshipTypeID *string `json:"relationship_type_id"`
	PhoneNumber        *string `json:"phone_number" binding:"omitempty,max=50"`
	Address            *string `json:"address"`
}

type CreateFamilyRequest struct {
	NIK                *string `json:"nik" binding:"omitempty,max=16"`
	Name               string  `json:"name" binding:"required,max=255"`
	DOB                *string `json:"dob"`
	RelationshipTypeID *string `json:"relationship_type_id"`
	EducationID        *string `json:"education_id"`
}

type UpdateFamilyRequest struct {
	NIK                *string `json:"nik" binding:"omitempty,max=16"`
	Name               *string `json:"name" binding:"omitempty,max=255"`
	DOB                *string `json:"dob"`
	RelationshipTypeID *string `json:"relationship_type_id"`
	EducationID        *string `json:"education_id"`
}

type CreateEducationRequest struct {
	EducationID *string `json:"education_id"`
	Name        string  `json:"name" binding:"required,max=255"`
	Major       *string `json:"major"`
	GradYear    *int    `json:"graduation_year"`
}

type UpdateEducationRequest struct {
	EducationID *string `json:"education_id"`
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Major       *string `json:"major"`
	GradYear    *int    `json:"graduation_year"`
}

type CreateExperienceRequest struct {
	Company   string  `json:"company" binding:"required,max=255"`
	Position  *string `json:"position"`
	StartYear *int    `json:"start_year"`
	EndYear   *int    `json:"end_year"`
}

type UpdateExperienceRequest struct {
	Company   *string `json:"company" binding:"omitempty,max=255"`
	Position  *string `json:"position"`
	StartYear *int    `json:"start_year"`
	EndYear   *int    `json:"end_year"`
}

type CreateDocumentRequest struct {
	Name string  `json:"name" binding:"required,max=255"`
	File string  `json:"file" binding:"required,max=255"`
	Note *string `json:"note"`
}

type UpdateDocumentRequest struct {
	Name *string `json:"name" binding:"omitempty,max=255"`
	File *string `json:"file" binding:"omitempty,max=255"`
	Note *string `json:"note"`
}

type CreateInsuranceRequest struct {
	Category *string `json:"category" binding:"omitempty,oneof='BPJS' 'Non BPJS'"`
	Number   string  `json:"number" binding:"required,max=100"`
	Name     string  `json:"name" binding:"required,max=100"`
	Type     *string `json:"type"`
}

type UpdateInsuranceRequest struct {
	Category *string `json:"category" binding:"omitempty,oneof='BPJS' 'Non BPJS'"`
	Number   *string `json:"number" binding:"omitempty,max=100"`
	Name     *string `json:"name" binding:"omitempty,max=100"`
	Type     *string `json:"type"`
}

type CreateEmploymentRequest struct {
	OrganizationID      *string `json:"organization_id"`
	PositionID          *string `json:"position_id"`
	EmploymentStatusID  *string `json:"employment_status_id"`
	DecisionLetterNumber string `json:"decision_letter_number" binding:"required,max=50"`
	DecisionLetterDate  string  `json:"decision_letter_date" binding:"required"`
	EffectiveDate       string  `json:"effective_date" binding:"required"`
	EffectiveEndDate    *string `json:"effective_end_date"`
}

type UpdateEmploymentRequest struct {
	OrganizationID      *string `json:"organization_id"`
	PositionID          *string `json:"position_id"`
	EmploymentStatusID  *string `json:"employment_status_id"`
	DecisionLetterNumber *string `json:"decision_letter_number" binding:"omitempty,max=50"`
	DecisionLetterDate  *string `json:"decision_letter_date"`
	EffectiveDate       *string `json:"effective_date"`
	EffectiveEndDate    *string `json:"effective_end_date"`
}

// =========================================================================
// Response DTOs
// =========================================================================

type AddressResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type,omitempty"`
	Address    string `json:"address,omitempty"`
	ProvinceID string `json:"province_id,omitempty"`
	RegencyID  string `json:"regency_id,omitempty"`
	DistrictID string `json:"district_id,omitempty"`
	VillageID  string `json:"village_id,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type EmergencyContactResponse struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	RelationshipTypeID string `json:"relationship_type_id,omitempty"`
	PhoneNumber        string `json:"phone_number"`
	Address            string `json:"address,omitempty"`
}

type FamilyResponse struct {
	ID                 string `json:"id"`
	NIK                string `json:"nik,omitempty"`
	Name               string `json:"name"`
	DOB                string `json:"dob,omitempty"`
	RelationshipTypeID string `json:"relationship_type_id,omitempty"`
	EducationID        string `json:"education_id,omitempty"`
}

type EducationResponse struct {
	ID         string `json:"id"`
	EducationID string `json:"education_id,omitempty"`
	Name       string `json:"name"`
	Major      string `json:"major,omitempty"`
	GradYear   int    `json:"graduation_year,omitempty"`
}

type ExperienceResponse struct {
	ID        string `json:"id"`
	Company   string `json:"company"`
	Position  string `json:"position,omitempty"`
	StartYear int    `json:"start_year,omitempty"`
	EndYear   int    `json:"end_year,omitempty"`
}

type DocumentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	File string `json:"file"`
	Note string `json:"note,omitempty"`
}

type InsuranceResponse struct {
	ID       string `json:"id"`
	Category string `json:"category,omitempty"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Type     string `json:"type,omitempty"`
}

type EmploymentResponse struct {
	ID                   string `json:"id"`
	OrganizationID       string `json:"organization_id,omitempty"`
	PositionID           string `json:"position_id,omitempty"`
	EmploymentStatusID   string `json:"employment_status_id,omitempty"`
	DecisionLetterNumber string `json:"decision_letter_number"`
	DecisionLetterDate   string `json:"decision_letter_date"`
	EffectiveDate        string `json:"effective_date"`
	EffectiveEndDate     string `json:"effective_end_date,omitempty"`
}

type EmployeeResponse struct {
	ID              string     `json:"id"`
	EmployeeID      string     `json:"employee_id"`
	NIK             string     `json:"nik,omitempty"`
	Name            string     `json:"name"`
	Gender          string     `json:"gender,omitempty"`
	POB             string     `json:"pob,omitempty"`
	DOB             string     `json:"dob,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty"`
	Email           string     `json:"email,omitempty"`
	ReligionID      string     `json:"religion_id,omitempty"`
	MaritalStatusID string     `json:"marital_status_id,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Sub-modules
	Addresses        []AddressResponse        `json:"addresses,omitempty"`
	EmergencyContacts []EmergencyContactResponse `json:"emergency_contacts,omitempty"`
	Families         []FamilyResponse         `json:"families,omitempty"`
	Educations       []EducationResponse      `json:"educations,omitempty"`
	Experiences      []ExperienceResponse     `json:"experiences,omitempty"`
	Documents        []DocumentResponse       `json:"documents,omitempty"`
	Insurances       []InsuranceResponse      `json:"insurances,omitempty"`
	Employments      []EmploymentResponse     `json:"employments,omitempty"`
}

// ListResponse untuk response pagination.
type ListResponse struct {
	Success    bool               `json:"success"`
	Data       []EmployeeResponse `json:"data"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"total_pages"`
}

// =========================================================================
// Converter helpers
// =========================================================================

func toAddressResponse(a *EmployeeAddress) AddressResponse {
	r := AddressResponse{ID: a.ID.String()}
	if a.Type != nil {
		r.Type = *a.Type
	}
	if a.Address != nil {
		r.Address = *a.Address
	}
	if a.ProvinceID != nil {
		r.ProvinceID = *a.ProvinceID
	}
	if a.RegencyID != nil {
		r.RegencyID = *a.RegencyID
	}
	if a.DistrictID != nil {
		r.DistrictID = *a.DistrictID
	}
	if a.VillageID != nil {
		r.VillageID = *a.VillageID
	}
	if a.PostalCode != nil {
		r.PostalCode = *a.PostalCode
	}
	return r
}

func toEmergencyContactResponse(c *EmergencyContact) EmergencyContactResponse {
	r := EmergencyContactResponse{
		ID:          c.ID.String(),
		Name:        c.Name,
		PhoneNumber: c.PhoneNumber,
	}
	if c.RelationshipTypeID != nil {
		r.RelationshipTypeID = c.RelationshipTypeID.String()
	}
	if c.Address != nil {
		r.Address = *c.Address
	}
	return r
}

func toFamilyResponse(f *EmployeeFamily) FamilyResponse {
	r := FamilyResponse{
		ID:   f.ID.String(),
		Name: f.Name,
	}
	if f.NIK != nil {
		r.NIK = *f.NIK
	}
	if f.DOB != nil {
		r.DOB = *f.DOB
	}
	if f.RelationshipTypeID != nil {
		r.RelationshipTypeID = f.RelationshipTypeID.String()
	}
	if f.EducationID != nil {
		r.EducationID = f.EducationID.String()
	}
	return r
}

func toEducationResponse(e *EmployeeEducation) EducationResponse {
	r := EducationResponse{
		ID:   e.ID.String(),
		Name: e.Name,
	}
	if e.EducationID != nil {
		r.EducationID = e.EducationID.String()
	}
	if e.Major != nil {
		r.Major = *e.Major
	}
	if e.GradYear != nil {
		r.GradYear = *e.GradYear
	}
	return r
}

func toExperienceResponse(e *EmployeeExperience) ExperienceResponse {
	r := ExperienceResponse{
		ID:      e.ID.String(),
		Company: e.Company,
	}
	if e.Position != nil {
		r.Position = *e.Position
	}
	if e.StartYear != nil {
		r.StartYear = *e.StartYear
	}
	if e.EndYear != nil {
		r.EndYear = *e.EndYear
	}
	return r
}

func toDocumentResponse(d *EmployeeDocument) DocumentResponse {
	r := DocumentResponse{
		ID:   d.ID.String(),
		Name: d.Name,
		File: d.File,
	}
	if d.Note != nil {
		r.Note = *d.Note
	}
	return r
}

func toInsuranceResponse(i *EmployeeInsurance) InsuranceResponse {
	r := InsuranceResponse{
		ID:     i.ID.String(),
		Number: i.Number,
		Name:   i.Name,
	}
	if i.Category != nil {
		r.Category = *i.Category
	}
	if i.Type != nil {
		r.Type = *i.Type
	}
	return r
}

func toEmploymentResponse(e *Employment) EmploymentResponse {
	r := EmploymentResponse{
		ID:                   e.ID.String(),
		DecisionLetterNumber: e.DecisionLetterNumber,
		DecisionLetterDate:   e.DecisionLetterDate,
		EffectiveDate:        e.EffectiveDate,
	}
	if e.OrganizationID != nil {
		r.OrganizationID = e.OrganizationID.String()
	}
	if e.PositionID != nil {
		r.PositionID = e.PositionID.String()
	}
	if e.EmploymentStatusID != nil {
		r.EmploymentStatusID = e.EmploymentStatusID.String()
	}
	if e.EffectiveEndDate != nil {
		r.EffectiveEndDate = *e.EffectiveEndDate
	}
	return r
}

func (e *Employee) ToResponse() EmployeeResponse {
	r := EmployeeResponse{
		ID:         e.ID.String(),
		EmployeeID: e.EmployeeID,
		Name:       e.Name,
		Status:     string(e.Status),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
	if e.NIK != nil {
		r.NIK = *e.NIK
	}
	if e.Gender != nil {
		r.Gender = *e.Gender
	}
	if e.POB != nil {
		r.POB = *e.POB
	}
	if e.DOB != nil {
		r.DOB = *e.DOB
	}
	if e.PhoneNumber != nil {
		r.PhoneNumber = *e.PhoneNumber
	}
	if e.Email != nil {
		r.Email = *e.Email
	}
	if e.ReligionID != nil {
		r.ReligionID = e.ReligionID.String()
	}
	if e.MaritalStatusID != nil {
		r.MaritalStatusID = e.MaritalStatusID.String()
	}

	// Sub-modules
	if len(e.Addresses) > 0 {
		r.Addresses = make([]AddressResponse, 0, len(e.Addresses))
		for _, a := range e.Addresses {
			r.Addresses = append(r.Addresses, toAddressResponse(&a))
		}
	}
	if len(e.EmergencyContacts) > 0 {
		r.EmergencyContacts = make([]EmergencyContactResponse, 0, len(e.EmergencyContacts))
		for _, c := range e.EmergencyContacts {
			r.EmergencyContacts = append(r.EmergencyContacts, toEmergencyContactResponse(&c))
		}
	}
	if len(e.Families) > 0 {
		r.Families = make([]FamilyResponse, 0, len(e.Families))
		for _, f := range e.Families {
			r.Families = append(r.Families, toFamilyResponse(&f))
		}
	}
	if len(e.Educations) > 0 {
		r.Educations = make([]EducationResponse, 0, len(e.Educations))
		for _, ed := range e.Educations {
			r.Educations = append(r.Educations, toEducationResponse(&ed))
		}
	}
	if len(e.Experiences) > 0 {
		r.Experiences = make([]ExperienceResponse, 0, len(e.Experiences))
		for _, ex := range e.Experiences {
			r.Experiences = append(r.Experiences, toExperienceResponse(&ex))
		}
	}
	if len(e.Documents) > 0 {
		r.Documents = make([]DocumentResponse, 0, len(e.Documents))
		for _, d := range e.Documents {
			r.Documents = append(r.Documents, toDocumentResponse(&d))
		}
	}
	if len(e.Insurances) > 0 {
		r.Insurances = make([]InsuranceResponse, 0, len(e.Insurances))
		for _, ins := range e.Insurances {
			r.Insurances = append(r.Insurances, toInsuranceResponse(&ins))
		}
	}
	if len(e.Employments) > 0 {
		r.Employments = make([]EmploymentResponse, 0, len(e.Employments))
		for _, emp := range e.Employments {
			r.Employments = append(r.Employments, toEmploymentResponse(&emp))
		}
	}

	return r
}
