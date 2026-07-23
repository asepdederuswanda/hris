package employeemovement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Employee Movement Service Tests
// =========================================================================

func TestService_CreateMovement_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuidStr()
	req := CreateMovementRequest{
		EmployeeID:           employeeID,
		MovementType:         "promotion",
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		Reason:               strPtr("Kinerja baik"),
	}

	resp, err := svc.CreateMovement(ctx(), req)
	if err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	if resp.MovementType != "promotion" {
		t.Errorf("expected movement_type 'promotion', got '%s'", resp.MovementType)
	}
	if resp.Status != "draft" {
		t.Errorf("expected default status 'draft', got '%s'", resp.Status)
	}
	if resp.EmployeeID != employeeID {
		t.Errorf("expected employee_id '%s', got '%s'", employeeID, resp.EmployeeID)
	}
}

func TestService_CreateMovement_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateMovementRequest{
		EmployeeID:           "not-a-uuid",
		MovementType:         "promotion",
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}

	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid employee UUID")
	}
}

func TestService_GetMovementByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	found, err := svc.GetMovementByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetMovementByID_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetMovementByID(ctx(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetMovementByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetMovementByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent movement")
	}
}

func TestService_ListMovements_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	for i := 0; i < 3; i++ {
		createTestMovement(repo, employeeID)
	}

	resp, err := svc.ListMovements(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_ListMovementsByEmployee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestMovement(repo, emp1)
	createTestMovement(repo, emp1)
	createTestMovement(repo, emp2)

	resp, err := svc.ListMovementsByEmployee(ctx(), emp1.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListMovementsByEmployee failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 for emp1, got %d", resp.Total)
	}
}

func TestService_UpdateMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	newReason := "Updated reason"
	updated, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		Reason: &newReason,
	})
	if err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	if updated.Reason == nil || *updated.Reason != "Updated reason" {
		t.Errorf("expected reason 'Updated reason', got '%v'", updated.Reason)
	}
}

func TestService_UpdateMovement_NonDraft_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	newReason := "Should fail"
	_, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		Reason: &newReason,
	})
	if err == nil {
		t.Fatal("expected error when updating non-draft movement")
	}
}

func TestService_DeleteMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	if err := svc.DeleteMovement(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteMovement failed: %v", err)
	}

	_, err := svc.GetMovementByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting movement")
	}
}

func TestService_DeleteMovement_NonDraft_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	err := svc.DeleteMovement(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error when deleting non-draft movement")
	}
}

func TestService_ApproveMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	approverID := uuidStr()
	if err := svc.ApproveMovement(ctx(), created.ID.String(), approverID); err != nil {
		t.Fatalf("ApproveMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected status 'approved', got '%s'", m.Status)
	}
}

func TestService_ApproveMovement_AlreadyApproved_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	err := svc.ApproveMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected error when approving already approved movement")
	}
}

func TestService_ExecuteMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	executorID := uuidStr()
	if err := svc.ExecuteMovement(ctx(), created.ID.String(), executorID); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", m.Status)
	}
}

func TestService_ExecuteMovement_Draft_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	// Cannot execute a draft movement (must approve first)
	_, repo, cleanup2 := newTestService()
	defer cleanup2()
	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	// Keep as draft

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected error when executing draft movement")
	}
}

func TestService_CancelMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	if err := svc.CancelMovement(ctx(), created.ID.String()); err != nil {
		t.Fatalf("CancelMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", m.Status)
	}
}

func TestService_CancelMovement_Executed_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	err := svc.CancelMovement(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error when cancelling executed movement")
	}
}

// =========================================================================
// Employee Contract Service Tests
// =========================================================================

func TestService_CreateContract_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuidStr()
	endDate := "2026-12-31"
	req := CreateContractRequest{
		EmployeeID:     employeeID,
		ContractNumber: "CTR-001",
		ContractType:   "pkwt",
		StartDate:      "2026-01-01",
		EndDate:        &endDate,
	}

	resp, err := svc.CreateContract(ctx(), req)
	if err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	if resp.ContractNumber != "CTR-001" {
		t.Errorf("expected contract_number 'CTR-001', got '%s'", resp.ContractNumber)
	}
	if resp.Status != "active" {
		t.Errorf("expected default status 'active', got '%s'", resp.Status)
	}
}

func TestService_CreateContract_InvalidEmployeeID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateContractRequest{
		EmployeeID:     "not-a-uuid",
		ContractNumber: "CTR-001",
		ContractType:   "pkwt",
		StartDate:      "2026-01-01",
	}

	_, err := svc.CreateContract(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid employee UUID")
	}
}

func TestService_GetContractByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	found, err := svc.GetContractByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetContractByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetContractByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetContractByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}

func TestService_ListContracts_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	for i := 0; i < 3; i++ {
		createTestContract(repo, employeeID)
	}

	resp, err := svc.ListContracts(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListContracts failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_ListContractsByEmployee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestContract(repo, emp1)
	createTestContract(repo, emp2)
	createTestContract(repo, emp2)

	resp, err := svc.ListContractsByEmployee(ctx(), emp2.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListContractsByEmployee failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 for emp2, got %d", resp.Total)
	}
}

func TestService_UpdateContract_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	newStatus := "expired"
	updated, err := svc.UpdateContract(ctx(), created.ID.String(), UpdateContractRequest{
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("UpdateContract failed: %v", err)
	}

	if updated.Status != "expired" {
		t.Errorf("expected status 'expired', got '%s'", updated.Status)
	}
}

func TestService_DeleteContract_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	if err := svc.DeleteContract(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteContract failed: %v", err)
	}

	_, err := svc.GetContractByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting contract")
	}
}

func TestService_DeleteContract_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	err := svc.DeleteContract(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}
