package employee

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Employee CRUD
// =========================================================================

func (s *Service) Create(ctx context.Context, req CreateEmployeeRequest) (*EmployeeResponse, error) {
	emp := &Employee{
		EmployeeID: req.EmployeeID,
		Name:       req.Name,
		Status:     "active",
	}

	if req.NIK != nil {
		emp.NIK = req.NIK
	}
	if req.FamilyID != nil {
		emp.FamilyID = req.FamilyID
	}
	if req.MotherName != nil {
		emp.MotherName = req.MotherName
	}
	if req.Gender != nil {
		emp.Gender = req.Gender
	}
	if req.NationalityType != nil {
		emp.NationalityType = req.NationalityType
	}
	if req.NationalityID != nil {
		emp.NationalityID = req.NationalityID
	}
	if req.POB != nil {
		emp.POB = req.POB
	}
	if req.DOB != nil {
		emp.DOB = req.DOB
	}
	if req.PhoneNumber != nil {
		emp.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		emp.Email = req.Email
	}
	if req.LinkedIn != nil {
		emp.LinkedIn = req.LinkedIn
	}
	if req.Instagram != nil {
		emp.Instagram = req.Instagram
	}

	if req.ReligionID != nil && *req.ReligionID != "" {
		id, err := uuid.Parse(*req.ReligionID)
		if err != nil {
			return nil, fmt.Errorf("invalid religion_id: %w", err)
		}
		emp.ReligionID = &id
	}
	if req.MaritalStatusID != nil && *req.MaritalStatusID != "" {
		id, err := uuid.Parse(*req.MaritalStatusID)
		if err != nil {
			return nil, fmt.Errorf("invalid marital_status_id: %w", err)
		}
		emp.MaritalStatusID = &id
	}

	if err := s.repo.CreateEmployee(ctx, emp); err != nil {
		return nil, err
	}

	s.logger.Info("Employee created",
		zap.String("id", emp.ID.String()),
		zap.String("name", emp.Name),
		zap.String("employee_id", emp.EmployeeID),
	)

	response := emp.ToResponse()
	return &response, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*EmployeeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := emp.ToResponse()
	return &response, nil
}

func (s *Service) List(ctx context.Context, page, perPage int) (*ListResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	employees, total, err := s.repo.FindAllEmployees(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []EmployeeResponse
	for _, e := range employees {
		responses = append(responses, e.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &ListResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateEmployeeRequest) (*EmployeeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		emp.Name = *req.Name
	}
	if req.NIK != nil {
		emp.NIK = req.NIK
	}
	if req.FamilyID != nil {
		emp.FamilyID = req.FamilyID
	}
	if req.MotherName != nil {
		emp.MotherName = req.MotherName
	}
	if req.Gender != nil {
		emp.Gender = req.Gender
	}
	if req.NationalityType != nil {
		emp.NationalityType = req.NationalityType
	}
	if req.NationalityID != nil {
		emp.NationalityID = req.NationalityID
	}
	if req.POB != nil {
		emp.POB = req.POB
	}
	if req.DOB != nil {
		emp.DOB = req.DOB
	}
	if req.PhoneNumber != nil {
		emp.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		emp.Email = req.Email
	}
	if req.LinkedIn != nil {
		emp.LinkedIn = req.LinkedIn
	}
	if req.Instagram != nil {
		emp.Instagram = req.Instagram
	}
	if req.ReligionID != nil && *req.ReligionID != "" {
		id, err := uuid.Parse(*req.ReligionID)
		if err != nil {
			return nil, fmt.Errorf("invalid religion_id: %w", err)
		}
		emp.ReligionID = &id
	}
	if req.MaritalStatusID != nil && *req.MaritalStatusID != "" {
		id, err := uuid.Parse(*req.MaritalStatusID)
		if err != nil {
			return nil, fmt.Errorf("invalid marital_status_id: %w", err)
		}
		emp.MaritalStatusID = &id
	}
	if req.Status != nil {
		emp.Status = *req.Status
	}

	if err := s.repo.UpdateEmployee(ctx, emp); err != nil {
		return nil, err
	}

	response := emp.ToResponse()
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee id: %w", err)
	}
	return s.repo.DeleteEmployee(ctx, uid)
}

// =========================================================================
// Sub-module CRUD: Addresses
// =========================================================================

func (s *Service) CreateAddress(ctx context.Context, employeeID string, req CreateAddressRequest) (*AddressResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	addr := &EmployeeAddress{
		EmployeeID: &empUID,
		Type:       &req.Type,
		Address:    &req.Address,
	}
	if req.ProvinceID != nil {
		addr.ProvinceID = req.ProvinceID
	}
	if req.RegencyID != nil {
		addr.RegencyID = req.RegencyID
	}
	if req.DistrictID != nil {
		addr.DistrictID = req.DistrictID
	}
	if req.VillageID != nil {
		addr.VillageID = req.VillageID
	}
	if req.PostalCode != nil {
		addr.PostalCode = req.PostalCode
	}

	if err := s.repo.CreateAddress(ctx, addr); err != nil {
		return nil, err
	}

	response := toAddressResponse(addr)
	return &response, nil
}

func (s *Service) UpdateAddress(ctx context.Context, employeeID, addressID string, req UpdateAddressRequest) (*AddressResponse, error) {
	addrUID, err := uuid.Parse(addressID)
	if err != nil {
		return nil, fmt.Errorf("invalid address id: %w", err)
	}

	addr, err := s.repo.FindAddressByID(ctx, addrUID)
	if err != nil {
		return nil, err
	}

	if req.Type != nil {
		addr.Type = req.Type
	}
	if req.Address != nil {
		addr.Address = req.Address
	}
	if req.ProvinceID != nil {
		addr.ProvinceID = req.ProvinceID
	}
	if req.RegencyID != nil {
		addr.RegencyID = req.RegencyID
	}
	if req.DistrictID != nil {
		addr.DistrictID = req.DistrictID
	}
	if req.VillageID != nil {
		addr.VillageID = req.VillageID
	}
	if req.PostalCode != nil {
		addr.PostalCode = req.PostalCode
	}

	if err := s.repo.UpdateAddress(ctx, addr); err != nil {
		return nil, err
	}

	response := toAddressResponse(addr)
	return &response, nil
}

func (s *Service) DeleteAddress(ctx context.Context, employeeID, addressID string) error {
	addrUID, err := uuid.Parse(addressID)
	if err != nil {
		return fmt.Errorf("invalid address id: %w", err)
	}
	return s.repo.DeleteAddress(ctx, addrUID)
}

// =========================================================================
// Sub-module CRUD: Emergency Contacts
// =========================================================================

func (s *Service) CreateEmergencyContact(ctx context.Context, employeeID string, req CreateEmergencyContactRequest) (*EmergencyContactResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	contact := &EmergencyContact{
		EmployeeID:  &empUID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		contact.RelationshipTypeID = &id
	}
	if req.Address != nil {
		contact.Address = req.Address
	}

	if err := s.repo.CreateEmergencyContact(ctx, contact); err != nil {
		return nil, err
	}

	response := toEmergencyContactResponse(contact)
	return &response, nil
}

func (s *Service) UpdateEmergencyContact(ctx context.Context, employeeID, contactID string, req UpdateEmergencyContactRequest) (*EmergencyContactResponse, error) {
	contUID, err := uuid.Parse(contactID)
	if err != nil {
		return nil, fmt.Errorf("invalid contact id: %w", err)
	}

	contact, err := s.repo.FindEmergencyContactByID(ctx, contUID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		contact.PhoneNumber = *req.PhoneNumber
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		contact.RelationshipTypeID = &id
	}
	if req.Address != nil {
		contact.Address = req.Address
	}

	if err := s.repo.UpdateEmergencyContact(ctx, contact); err != nil {
		return nil, err
	}

	response := toEmergencyContactResponse(contact)
	return &response, nil
}

func (s *Service) DeleteEmergencyContact(ctx context.Context, employeeID, contactID string) error {
	contUID, err := uuid.Parse(contactID)
	if err != nil {
		return fmt.Errorf("invalid contact id: %w", err)
	}
	return s.repo.DeleteEmergencyContact(ctx, contUID)
}

// =========================================================================
// Sub-module CRUD: Families
// =========================================================================

func (s *Service) CreateFamily(ctx context.Context, employeeID string, req CreateFamilyRequest) (*FamilyResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	fam := &EmployeeFamily{
		EmployeeID: &empUID,
		Name:       req.Name,
	}
	if req.NIK != nil {
		fam.NIK = req.NIK
	}
	if req.DOB != nil {
		fam.DOB = req.DOB
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		fam.RelationshipTypeID = &id
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		fam.EducationID = &id
	}

	if err := s.repo.CreateFamily(ctx, fam); err != nil {
		return nil, err
	}

	response := toFamilyResponse(fam)
	return &response, nil
}

func (s *Service) UpdateFamily(ctx context.Context, employeeID, familyID string, req UpdateFamilyRequest) (*FamilyResponse, error) {
	famUID, err := uuid.Parse(familyID)
	if err != nil {
		return nil, fmt.Errorf("invalid family id: %w", err)
	}

	fam, err := s.repo.FindFamilyByID(ctx, famUID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		fam.Name = *req.Name
	}
	if req.NIK != nil {
		fam.NIK = req.NIK
	}
	if req.DOB != nil {
		fam.DOB = req.DOB
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		fam.RelationshipTypeID = &id
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		fam.EducationID = &id
	}

	if err := s.repo.UpdateFamily(ctx, fam); err != nil {
		return nil, err
	}

	response := toFamilyResponse(fam)
	return &response, nil
}

func (s *Service) DeleteFamily(ctx context.Context, employeeID, familyID string) error {
	famUID, err := uuid.Parse(familyID)
	if err != nil {
		return fmt.Errorf("invalid family id: %w", err)
	}
	return s.repo.DeleteFamily(ctx, famUID)
}

// =========================================================================
// Sub-module CRUD: Educations
// =========================================================================

func (s *Service) CreateEducation(ctx context.Context, employeeID string, req CreateEducationRequest) (*EducationResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	edu := &EmployeeEducation{
		EmployeeID: &empUID,
		Name:       req.Name,
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		edu.EducationID = &id
	}
	if req.Major != nil {
		edu.Major = req.Major
	}
	if req.GradYear != nil {
		edu.GradYear = req.GradYear
	}

	if err := s.repo.CreateEducation(ctx, edu); err != nil {
		return nil, err
	}

	response := toEducationResponse(edu)
	return &response, nil
}

func (s *Service) UpdateEducation(ctx context.Context, employeeID, educationID string, req UpdateEducationRequest) (*EducationResponse, error) {
	eduUID, err := uuid.Parse(educationID)
	if err != nil {
		return nil, fmt.Errorf("invalid education id: %w", err)
	}

	edu, err := s.repo.FindEducationByID(ctx, eduUID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		edu.Name = *req.Name
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		edu.EducationID = &id
	}
	if req.Major != nil {
		edu.Major = req.Major
	}
	if req.GradYear != nil {
		edu.GradYear = req.GradYear
	}

	if err := s.repo.UpdateEducation(ctx, edu); err != nil {
		return nil, err
	}

	response := toEducationResponse(edu)
	return &response, nil
}

func (s *Service) DeleteEducation(ctx context.Context, employeeID, educationID string) error {
	eduUID, err := uuid.Parse(educationID)
	if err != nil {
		return fmt.Errorf("invalid education id: %w", err)
	}
	return s.repo.DeleteEducation(ctx, eduUID)
}

// =========================================================================
// Sub-module CRUD: Experiences
// =========================================================================

func (s *Service) CreateExperience(ctx context.Context, employeeID string, req CreateExperienceRequest) (*ExperienceResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	exp := &EmployeeExperience{
		EmployeeID: &empUID,
		Company:    req.Company,
	}
	if req.Position != nil {
		exp.Position = req.Position
	}
	if req.StartYear != nil {
		exp.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		exp.EndYear = req.EndYear
	}

	if err := s.repo.CreateExperience(ctx, exp); err != nil {
		return nil, err
	}

	response := toExperienceResponse(exp)
	return &response, nil
}

func (s *Service) UpdateExperience(ctx context.Context, employeeID, experienceID string, req UpdateExperienceRequest) (*ExperienceResponse, error) {
	expUID, err := uuid.Parse(experienceID)
	if err != nil {
		return nil, fmt.Errorf("invalid experience id: %w", err)
	}

	exp, err := s.repo.FindExperienceByID(ctx, expUID)
	if err != nil {
		return nil, err
	}

	if req.Company != nil {
		exp.Company = *req.Company
	}
	if req.Position != nil {
		exp.Position = req.Position
	}
	if req.StartYear != nil {
		exp.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		exp.EndYear = req.EndYear
	}

	if err := s.repo.UpdateExperience(ctx, exp); err != nil {
		return nil, err
	}

	response := toExperienceResponse(exp)
	return &response, nil
}

func (s *Service) DeleteExperience(ctx context.Context, employeeID, experienceID string) error {
	expUID, err := uuid.Parse(experienceID)
	if err != nil {
		return fmt.Errorf("invalid experience id: %w", err)
	}
	return s.repo.DeleteExperience(ctx, expUID)
}

// =========================================================================
// Sub-module CRUD: Documents
// =========================================================================

func (s *Service) CreateDocument(ctx context.Context, employeeID string, req CreateDocumentRequest) (*DocumentResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	doc := &EmployeeDocument{
		EmployeeID: &empUID,
		Name:       req.Name,
		File:       req.File,
	}
	if req.Note != nil {
		doc.Note = req.Note
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	response := toDocumentResponse(doc)
	return &response, nil
}

func (s *Service) UpdateDocument(ctx context.Context, employeeID, documentID string, req UpdateDocumentRequest) (*DocumentResponse, error) {
	docUID, err := uuid.Parse(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	doc, err := s.repo.FindDocumentByID(ctx, docUID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		doc.Name = *req.Name
	}
	if req.File != nil {
		doc.File = *req.File
	}
	if req.Note != nil {
		doc.Note = req.Note
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	response := toDocumentResponse(doc)
	return &response, nil
}

func (s *Service) DeleteDocument(ctx context.Context, employeeID, documentID string) error {
	docUID, err := uuid.Parse(documentID)
	if err != nil {
		return fmt.Errorf("invalid document id: %w", err)
	}
	return s.repo.DeleteDocument(ctx, docUID)
}

// =========================================================================
// Sub-module CRUD: Insurances
// =========================================================================

func (s *Service) CreateInsurance(ctx context.Context, employeeID string, req CreateInsuranceRequest) (*InsuranceResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	ins := &EmployeeInsurance{
		EmployeeID: &empUID,
		Number:     req.Number,
		Name:       req.Name,
	}
	if req.Category != nil {
		ins.Category = req.Category
	}
	if req.Type != nil {
		ins.Type = req.Type
	}

	if err := s.repo.CreateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	response := toInsuranceResponse(ins)
	return &response, nil
}

func (s *Service) UpdateInsurance(ctx context.Context, employeeID, insuranceID string, req UpdateInsuranceRequest) (*InsuranceResponse, error) {
	insUID, err := uuid.Parse(insuranceID)
	if err != nil {
		return nil, fmt.Errorf("invalid insurance id: %w", err)
	}

	ins, err := s.repo.FindInsuranceByID(ctx, insUID)
	if err != nil {
		return nil, err
	}

	if req.Number != nil {
		ins.Number = *req.Number
	}
	if req.Name != nil {
		ins.Name = *req.Name
	}
	if req.Category != nil {
		ins.Category = req.Category
	}
	if req.Type != nil {
		ins.Type = req.Type
	}

	if err := s.repo.UpdateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	response := toInsuranceResponse(ins)
	return &response, nil
}

func (s *Service) DeleteInsurance(ctx context.Context, employeeID, insuranceID string) error {
	insUID, err := uuid.Parse(insuranceID)
	if err != nil {
		return fmt.Errorf("invalid insurance id: %w", err)
	}
	return s.repo.DeleteInsurance(ctx, insUID)
}

// =========================================================================
// Sub-module CRUD: Employments
// =========================================================================

func (s *Service) CreateEmployment(ctx context.Context, employeeID string, req CreateEmploymentRequest) (*EmploymentResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	empl := &Employment{
		EmployeeID:          &empUID,
		DecisionLetterNumber: req.DecisionLetterNumber,
		DecisionLetterDate:   req.DecisionLetterDate,
		EffectiveDate:        req.EffectiveDate,
	}
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		id, _ := uuid.Parse(*req.OrganizationID)
		empl.OrganizationID = &id
	}
	if req.PositionID != nil && *req.PositionID != "" {
		id, _ := uuid.Parse(*req.PositionID)
		empl.PositionID = &id
	}
	if req.EmploymentStatusID != nil && *req.EmploymentStatusID != "" {
		id, _ := uuid.Parse(*req.EmploymentStatusID)
		empl.EmploymentStatusID = &id
	}
	if req.EffectiveEndDate != nil {
		empl.EffectiveEndDate = req.EffectiveEndDate
	}

	if err := s.repo.CreateEmployment(ctx, empl); err != nil {
		return nil, err
	}

	response := toEmploymentResponse(empl)
	return &response, nil
}

func (s *Service) UpdateEmployment(ctx context.Context, employeeID, employmentID string, req UpdateEmploymentRequest) (*EmploymentResponse, error) {
	emplUID, err := uuid.Parse(employmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid employment id: %w", err)
	}

	empl, err := s.repo.FindEmploymentByID(ctx, emplUID)
	if err != nil {
		return nil, err
	}

	if req.DecisionLetterNumber != nil {
		empl.DecisionLetterNumber = *req.DecisionLetterNumber
	}
	if req.DecisionLetterDate != nil {
		empl.DecisionLetterDate = *req.DecisionLetterDate
	}
	if req.EffectiveDate != nil {
		empl.EffectiveDate = *req.EffectiveDate
	}
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		id, _ := uuid.Parse(*req.OrganizationID)
		empl.OrganizationID = &id
	}
	if req.PositionID != nil && *req.PositionID != "" {
		id, _ := uuid.Parse(*req.PositionID)
		empl.PositionID = &id
	}
	if req.EmploymentStatusID != nil && *req.EmploymentStatusID != "" {
		id, _ := uuid.Parse(*req.EmploymentStatusID)
		empl.EmploymentStatusID = &id
	}
	if req.EffectiveEndDate != nil {
		empl.EffectiveEndDate = req.EffectiveEndDate
	}

	if err := s.repo.UpdateEmployment(ctx, empl); err != nil {
		return nil, err
	}

	response := toEmploymentResponse(empl)
	return &response, nil
}

func (s *Service) DeleteEmployment(ctx context.Context, employeeID, employmentID string) error {
	emplUID, err := uuid.Parse(employmentID)
	if err != nil {
		return fmt.Errorf("invalid employment id: %w", err)
	}
	return s.repo.DeleteEmployment(ctx, emplUID)
}
