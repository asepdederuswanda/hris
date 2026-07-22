package employee

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =========================================================================
// Service Tests (using real SQLite repository for true integration)
// =========================================================================

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() { cleanup(); logger.Sync() }
}

func TestService_CreateEmployee_DefaultStatus(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateEmployeeRequest{
		EmployeeID: "SVC001",
		Name:       "Service Test",
	}

	resp, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", resp.Status)
	}
	if resp.Name != "Service Test" {
		t.Errorf("expected name 'Service Test', got '%s'", resp.Name)
	}
}

func TestService_CreateEmployee_WithOptionalFields(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	gender := "M"
	req := CreateEmployeeRequest{
		EmployeeID:  "SVC002",
		Name:        "Full Employee",
		NIK:         strPtr("1234567890123456"),
		Gender:      &gender,
		PhoneNumber: strPtr("08123456789"),
		Email:       strPtr("full@example.com"),
	}

	resp, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("Create with optional fields failed: %v", err)
	}

	if resp.NIK != "1234567890123456" {
		t.Errorf("expected NIK '1234567890123456', got '%s'", resp.NIK)
	}
	if resp.Gender != "M" {
		t.Errorf("expected Gender 'M', got '%s'", resp.Gender)
	}
	if resp.Email != "full@example.com" {
		t.Errorf("expected email 'full@example.com', got '%s'", resp.Email)
	}
}

func TestService_GetByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateEmployeeRequest{
		EmployeeID: "SVC003",
		Name:       "Get By ID Test",
	}
	created, _ := svc.Create(ctx, req)

	found, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.Name != "Get By ID Test" {
		t.Errorf("expected name 'Get By ID Test', got '%s'", found.Name)
	}
}

func TestService_GetByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent employee")
	}
}

func TestService_List_DefaultPagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	// Create employees
	for i := 0; i < 5; i++ {
		svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("LST%03d", i+1),
			Name:       fmt.Sprintf("List Employee %d", i+1),
		})
	}

	// Test with invalid params (should use defaults)
	resp, err := svc.List(ctx, 0, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
}

func TestService_UpdateEmployee(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC004",
		Name:       "Before Update",
	})

	newName := "After Update"
	updated, err := svc.Update(ctx, created.ID, UpdateEmployeeRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != "After Update" {
		t.Errorf("expected name 'After Update', got '%s'", updated.Name)
	}
}

func TestService_DeleteEmployee(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC005",
		Name:       "To Delete",
	})

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err := svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting employee")
	}
}

// =========================================================================
// Sub-module Service Tests
// =========================================================================

func TestService_CreateAddress(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC010",
		Name:       "Address Owner",
	})

	addr, err := svc.CreateAddress(ctx, emp.ID, CreateAddressRequest{
		Type:    "MAIN",
		Address: "Jl. Service Test No. 1",
	})
	if err != nil {
		t.Fatalf("CreateAddress failed: %v", err)
	}

	if addr.Type != "MAIN" {
		t.Errorf("expected type 'MAIN', got '%s'", addr.Type)
	}
}

func TestService_CreateEmergencyContact(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC011",
		Name:       "Contact Owner",
	})

	contact, err := svc.CreateEmergencyContact(ctx, emp.ID, CreateEmergencyContactRequest{
		Name:        "Emergency Contact",
		PhoneNumber: "08111111111",
	})
	if err != nil {
		t.Fatalf("CreateEmergencyContact failed: %v", err)
	}

	if contact.Name != "Emergency Contact" {
		t.Errorf("expected name 'Emergency Contact', got '%s'", contact.Name)
	}
}

func TestService_CreateFamily(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC012",
		Name:       "Family Owner",
	})

	fam, err := svc.CreateFamily(ctx, emp.ID, CreateFamilyRequest{
		Name: "Family Member",
	})
	if err != nil {
		t.Fatalf("CreateFamily failed: %v", err)
	}

	if fam.Name != "Family Member" {
		t.Errorf("expected name 'Family Member', got '%s'", fam.Name)
	}
}

func TestService_CreateEducation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC013",
		Name:       "Education Owner",
	})

	edu, err := svc.CreateEducation(ctx, emp.ID, CreateEducationRequest{
		Name: "S1 Informatics",
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}

	if edu.Name != "S1 Informatics" {
		t.Errorf("expected name 'S1 Informatics', got '%s'", edu.Name)
	}
}

func TestService_CreateExperience(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC014",
		Name:       "Experience Owner",
	})

	exp, err := svc.CreateExperience(ctx, emp.ID, CreateExperienceRequest{
		Company: "PT Test Corp",
	})
	if err != nil {
		t.Fatalf("CreateExperience failed: %v", err)
	}

	if exp.Company != "PT Test Corp" {
		t.Errorf("expected company 'PT Test Corp', got '%s'", exp.Company)
	}
}

func TestService_CreateDocument(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC015",
		Name:       "Document Owner",
	})

	doc, err := svc.CreateDocument(ctx, emp.ID, CreateDocumentRequest{
		Name: "CV",
		File: "cv.pdf",
	})
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	if doc.Name != "CV" {
		t.Errorf("expected name 'CV', got '%s'", doc.Name)
	}
}

func TestService_CreateInsurance(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC016",
		Name:       "Insurance Owner",
	})

	ins, err := svc.CreateInsurance(ctx, emp.ID, CreateInsuranceRequest{
		Number: "BPJS001",
		Name:   "BPJS Kesehatan",
	})
	if err != nil {
		t.Fatalf("CreateInsurance failed: %v", err)
	}

	if ins.Name != "BPJS Kesehatan" {
		t.Errorf("expected name 'BPJS Kesehatan', got '%s'", ins.Name)
	}
}

func TestService_CreateEmployment(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC017",
		Name:       "Employment Owner",
	})

	empl, err := svc.CreateEmployment(ctx, emp.ID, CreateEmploymentRequest{
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2024-01-01",
		EffectiveDate:        "2024-01-15",
	})
	if err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	if empl.DecisionLetterNumber != "SK-001" {
		t.Errorf("expected SK 'SK-001', got '%s'", empl.DecisionLetterNumber)
	}
}

func TestService_GetEmployeeWithSubModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC020",
		Name:       "Full Data Employee",
	})

	// Add address
	svc.CreateAddress(ctx, emp.ID, CreateAddressRequest{
		Type:    "MAIN",
		Address: "Jl. Test",
	})

	// Add family
	svc.CreateFamily(ctx, emp.ID, CreateFamilyRequest{
		Name: "Spouse",
	})

	// Add document
	svc.CreateDocument(ctx, emp.ID, CreateDocumentRequest{
		Name: "ID Card",
		File: "ktp.pdf",
	})

	// Fetch employee with all sub-modules
	fullEmp, err := svc.GetByID(ctx, emp.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if len(fullEmp.Addresses) != 1 {
		t.Errorf("expected 1 address, got %d", len(fullEmp.Addresses))
	}
	if len(fullEmp.Families) != 1 {
		t.Errorf("expected 1 family, got %d", len(fullEmp.Families))
	}
	if len(fullEmp.Documents) != 1 {
		t.Errorf("expected 1 document, got %d", len(fullEmp.Documents))
	}
}
