package employee

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Employee CRUD Tests
// =========================================================================

func TestRepository_CreateEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := &Employee{
		EmployeeID:      "EMP001",
		Name:            "John Doe",
		Gender:          strPtr("M"),
		PhoneNumber:     strPtr("08123456789"),
		Email:           strPtr("john@example.com"),
		Status:          "active",
	}

	if err := repo.CreateEmployee(ctx, emp); err != nil {
		t.Fatalf("CreateEmployee failed: %v", err)
	}

	if emp.ID == uuid.Nil {
		t.Error("expected employee ID to be generated")
	}
}

func TestRepository_FindEmployeeByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	found, err := repo.FindEmployeeByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindEmployeeByID failed: %v", err)
	}

	if found.Name != "Test Employee" {
		t.Errorf("expected name 'Test Employee', got '%s'", found.Name)
	}
	if found.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", found.Status)
	}
}

func TestRepository_FindEmployeeByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindEmployeeByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent employee")
	}
}

func TestRepository_FindEmployeeByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	createTestEmployee(ctx, repo)

	found, err := repo.FindEmployeeByEmployeeID(ctx, "EMP-TEST-001")
	if err != nil {
		t.Fatalf("FindEmployeeByEmployeeID failed: %v", err)
	}

	if found.Name != "Test Employee" {
		t.Errorf("expected name 'Test Employee', got '%s'", found.Name)
	}
}

func TestRepository_FindAllEmployees_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create 3 employees
	for i := 0; i < 3; i++ {
		emp := &Employee{
			EmployeeID: fmt.Sprintf("EMP%03d", i+1),
			Name:       fmt.Sprintf("Employee %d", i+1),
			Status:     "active",
		}
		if err := repo.CreateEmployee(ctx, emp); err != nil {
			t.Fatalf("failed to create employee %d: %v", i, err)
		}
	}

	// Test page 1 with per_page 2
	emps, total, err := repo.FindAllEmployees(ctx, 1, 2)
	if err != nil {
		t.Fatalf("FindAllEmployees failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(emps) != 2 {
		t.Errorf("expected 2 employees on page 1, got %d", len(emps))
	}
}

func TestRepository_UpdateEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	created.Name = "Updated Name"
	if err := repo.UpdateEmployee(ctx, created); err != nil {
		t.Fatalf("UpdateEmployee failed: %v", err)
	}

	found, _ := repo.FindEmployeeByID(ctx, created.ID)
	if found.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", found.Name)
	}
}

func TestRepository_DeleteEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	if err := repo.DeleteEmployee(ctx, created.ID); err != nil {
		t.Fatalf("DeleteEmployee failed: %v", err)
	}

	// Verify it's gone
	_, err := repo.FindEmployeeByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting employee")
	}
}

// =========================================================================
// Address Sub-module Tests
// =========================================================================

func TestRepository_CreateAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	if err := repo.CreateAddress(ctx, addr); err != nil {
		t.Fatalf("CreateAddress failed: %v", err)
	}

	if addr.ID == uuid.Nil {
		t.Error("expected address ID to be generated")
	}
}

func TestRepository_FindAddressByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	found, err := repo.FindAddressByID(ctx, addr.ID)
	if err != nil {
		t.Fatalf("FindAddressByID failed: %v", err)
	}
	if *found.Type != "MAIN" {
		t.Errorf("expected type 'MAIN', got '%s'", *found.Type)
	}
}

func TestRepository_UpdateAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	addr.Type = strPtr("DOMICILE")
	repo.UpdateAddress(ctx, addr)

	found, _ := repo.FindAddressByID(ctx, addr.ID)
	if *found.Type != "DOMICILE" {
		t.Errorf("expected type 'DOMICILE', got '%s'", *found.Type)
	}
}

func TestRepository_DeleteAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	repo.DeleteAddress(ctx, addr.ID)

	_, err := repo.FindAddressByID(ctx, addr.ID)
	if err == nil {
		t.Fatal("expected error after deleting address")
	}
}

// =========================================================================
// Emergency Contact Sub-module Tests
// =========================================================================

func TestRepository_EmergencyContactCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	// Create
	contact := &EmergencyContact{
		EmployeeID:  &emp.ID,
		Name:        "Emergency Contact",
		PhoneNumber: "08111111111",
	}
	if err := repo.CreateEmergencyContact(ctx, contact); err != nil {
		t.Fatalf("CreateEmergencyContact failed: %v", err)
	}

	// Find
	found, err := repo.FindEmergencyContactByID(ctx, contact.ID)
	if err != nil {
		t.Fatalf("FindEmergencyContactByID failed: %v", err)
	}
	if found.Name != "Emergency Contact" {
		t.Errorf("expected name 'Emergency Contact', got '%s'", found.Name)
	}

	// Update
	found.Name = "Updated Contact"
	repo.UpdateEmergencyContact(ctx, found)
	updated, _ := repo.FindEmergencyContactByID(ctx, contact.ID)
	if updated.Name != "Updated Contact" {
		t.Errorf("expected name 'Updated Contact', got '%s'", updated.Name)
	}

	// Delete
	repo.DeleteEmergencyContact(ctx, contact.ID)
	_, err = repo.FindEmergencyContactByID(ctx, contact.ID)
	if err == nil {
		t.Fatal("expected error after deleting emergency contact")
	}
}

// =========================================================================
// Family Sub-module Tests
// =========================================================================

func TestRepository_FamilyCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	// Create
	fam := &EmployeeFamily{
		EmployeeID: &emp.ID,
		Name:       "Family Member",
	}
	if err := repo.CreateFamily(ctx, fam); err != nil {
		t.Fatalf("CreateFamily failed: %v", err)
	}

	// Find
	found, err := repo.FindFamilyByID(ctx, fam.ID)
	if err != nil {
		t.Fatalf("FindFamilyByID failed: %v", err)
	}
	if found.Name != "Family Member" {
		t.Errorf("expected name 'Family Member', got '%s'", found.Name)
	}

	// Update
	found.Name = "Updated Family"
	repo.UpdateFamily(ctx, found)
	updated, _ := repo.FindFamilyByID(ctx, fam.ID)
	if updated.Name != "Updated Family" {
		t.Errorf("expected name 'Updated Family', got '%s'", updated.Name)
	}

	// Delete
	repo.DeleteFamily(ctx, fam.ID)
	_, err = repo.FindFamilyByID(ctx, fam.ID)
	if err == nil {
		t.Fatal("expected error after deleting family")
	}
}

// =========================================================================
// Education Sub-module Tests
// =========================================================================

func TestRepository_EducationCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	edu := &EmployeeEducation{
		EmployeeID: &emp.ID,
		Name:       "S1 Computer Science",
		Major:      strPtr("Computer Science"),
		GradYear:   intPtr(2020),
	}
	if err := repo.CreateEducation(ctx, edu); err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}

	// CRUD cycle
	found, _ := repo.FindEducationByID(ctx, edu.ID)
	if found.Name != "S1 Computer Science" {
		t.Errorf("expected name 'S1 Computer Science', got '%s'", found.Name)
	}

	found.Name = "S2 Data Science"
	repo.UpdateEducation(ctx, found)
	updated, _ := repo.FindEducationByID(ctx, edu.ID)
	if updated.Name != "S2 Data Science" {
		t.Errorf("expected name 'S2 Data Science', got '%s'", updated.Name)
	}

	repo.DeleteEducation(ctx, edu.ID)
	_, err := repo.FindEducationByID(ctx, edu.ID)
	if err == nil {
		t.Fatal("expected error after deleting education")
	}
}

// =========================================================================
// Experience Sub-module Tests
// =========================================================================

func TestRepository_ExperienceCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	exp := &EmployeeExperience{
		EmployeeID: &emp.ID,
		Company:    "PT Test",
		Position:   strPtr("Manager"),
		StartYear:  intPtr(2020),
		EndYear:    intPtr(2024),
	}
	if err := repo.CreateExperience(ctx, exp); err != nil {
		t.Fatalf("CreateExperience failed: %v", err)
	}

	found, _ := repo.FindExperienceByID(ctx, exp.ID)
	if found.Company != "PT Test" {
		t.Errorf("expected company 'PT Test', got '%s'", found.Company)
	}

	// Update
	found.Company = "PT Updated"
	repo.UpdateExperience(ctx, found)
	updated, _ := repo.FindExperienceByID(ctx, exp.ID)
	if updated.Company != "PT Updated" {
		t.Errorf("expected company 'PT Updated', got '%s'", updated.Company)
	}

	// Delete
	repo.DeleteExperience(ctx, exp.ID)
	_, err := repo.FindExperienceByID(ctx, exp.ID)
	if err == nil {
		t.Fatal("expected error after deleting experience")
	}
}

// =========================================================================
// Document Sub-module Tests
// =========================================================================

func TestRepository_DocumentCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	doc := &EmployeeDocument{
		EmployeeID: &emp.ID,
		Name:       "CV",
		File:       "cv.pdf",
		Note:       strPtr("Original CV"),
	}
	if err := repo.CreateDocument(ctx, doc); err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	found, _ := repo.FindDocumentByID(ctx, doc.ID)
	if found.Name != "CV" {
		t.Errorf("expected name 'CV', got '%s'", found.Name)
	}

	found.Name = "Updated CV"
	repo.UpdateDocument(ctx, found)
	updated, _ := repo.FindDocumentByID(ctx, doc.ID)
	if updated.Name != "Updated CV" {
		t.Errorf("expected name 'Updated CV', got '%s'", updated.Name)
	}

	repo.DeleteDocument(ctx, doc.ID)
	_, err := repo.FindDocumentByID(ctx, doc.ID)
	if err == nil {
		t.Fatal("expected error after deleting document")
	}
}

// =========================================================================
// Insurance Sub-module Tests
// =========================================================================

func TestRepository_InsuranceCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	ins := &EmployeeInsurance{
		EmployeeID: &emp.ID,
		Number:     "BPJS001",
		Name:       "BPJS Kesehatan",
		Category:   strPtr("BPJS"),
	}
	if err := repo.CreateInsurance(ctx, ins); err != nil {
		t.Fatalf("CreateInsurance failed: %v", err)
	}

	found, _ := repo.FindInsuranceByID(ctx, ins.ID)
	if found.Name != "BPJS Kesehatan" {
		t.Errorf("expected name 'BPJS Kesehatan', got '%s'", found.Name)
	}

	found.Name = "BPJS Ketenagakerjaan"
	repo.UpdateInsurance(ctx, found)
	updated, _ := repo.FindInsuranceByID(ctx, ins.ID)
	if updated.Name != "BPJS Ketenagakerjaan" {
		t.Errorf("expected name 'BPJS Ketenagakerjaan', got '%s'", updated.Name)
	}

	repo.DeleteInsurance(ctx, ins.ID)
	_, err := repo.FindInsuranceByID(ctx, ins.ID)
	if err == nil {
		t.Fatal("expected error after deleting insurance")
	}
}

// =========================================================================
// Employment Sub-module Tests
// =========================================================================

func TestRepository_EmploymentCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	empl := &Employment{
		EmployeeID:           &emp.ID,
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2024-01-01",
		EffectiveDate:        "2024-01-15",
	}
	if err := repo.CreateEmployment(ctx, empl); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	found, _ := repo.FindEmploymentByID(ctx, empl.ID)
	if found.DecisionLetterNumber != "SK-001" {
		t.Errorf("expected SK 'SK-001', got '%s'", found.DecisionLetterNumber)
	}

	found.DecisionLetterNumber = "SK-002"
	repo.UpdateEmployment(ctx, found)
	updated, _ := repo.FindEmploymentByID(ctx, empl.ID)
	if updated.DecisionLetterNumber != "SK-002" {
		t.Errorf("expected SK 'SK-002', got '%s'", updated.DecisionLetterNumber)
	}

	repo.DeleteEmployment(ctx, empl.ID)
	_, err := repo.FindEmploymentByID(ctx, empl.ID)
	if err == nil {
		t.Fatal("expected error after deleting employment")
	}
}
