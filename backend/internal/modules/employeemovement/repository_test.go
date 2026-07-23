package employeemovement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Employee Movement Repository Tests
// =========================================================================

func TestRepo_CreateMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	m := &EmployeeMovement{
		EmployeeID:           uuid.New(),
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}

	if err := repo.CreateMovement(ctx, m); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	if m.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindMovementByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	found, err := repo.FindMovementByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindMovementByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
}

func TestRepo_FindMovementByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindMovementByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent movement")
	}
}

func TestRepo_ListMovements_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	for i := 0; i < 5; i++ {
		createTestMovement(repo, empID)
	}

	movements, total, err := repo.ListMovements(ctx, 1, 3)
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(movements) != 3 {
		t.Errorf("expected 3 movements (page 1 of 3), got %d", len(movements))
	}
}

func TestRepo_FindMovementsByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestMovement(repo, emp1)
	createTestMovement(repo, emp1)
	createTestMovement(repo, emp2)

	movements, total, err := repo.FindMovementsByEmployeeID(ctx, emp1, 1, 10)
	if err != nil {
		t.Fatalf("FindMovementsByEmployeeID failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2 for emp1, got %d", total)
	}
	if len(movements) != 2 {
		t.Errorf("expected 2 movements, got %d", len(movements))
	}
}

func TestRepo_UpdateMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	created.Reason = strPtr("Updated reason")
	if err := repo.UpdateMovement(ctx, created); err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Reason == nil || *found.Reason != "Updated reason" {
		t.Errorf("expected reason 'Updated reason', got '%v'", found.Reason)
	}
}

func TestRepo_DeleteMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	if err := repo.DeleteMovement(ctx, created.ID); err != nil {
		t.Fatalf("DeleteMovement failed: %v", err)
	}

	_, err := repo.FindMovementByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting movement")
	}
}

func TestRepo_ApproveMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	approverID := uuid.New()

	if err := repo.ApproveMovement(ctx, created.ID, approverID); err != nil {
		t.Fatalf("ApproveMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusApproved {
		t.Errorf("expected status 'approved', got '%s'", found.Status)
	}
}

func TestRepo_ApproveMovement_NonDraft_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx, created)

	err := repo.ApproveMovement(ctx, created.ID, uuid.New())
	if err == nil {
		t.Fatal("expected error when approving non-draft movement")
	}
}

func TestRepo_ExecuteMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx, created)

	executorID := uuid.New()
	if err := repo.ExecuteMovement(ctx, created.ID, executorID); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", found.Status)
	}
}

func TestRepo_ExecuteMovement_NonApproved_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	// Still draft — cannot execute without approval

	err := repo.ExecuteMovement(ctx, created.ID, uuid.New())
	if err == nil {
		t.Fatal("expected error when executing non-approved movement")
	}
}

func TestRepo_CancelMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	if err := repo.CancelMovement(ctx, created.ID); err != nil {
		t.Fatalf("CancelMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", found.Status)
	}
}

func TestRepo_CancelMovement_Executed_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx, created)

	err := repo.CancelMovement(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error when cancelling executed movement")
	}
}

// =========================================================================
// Employee Contract Repository Tests
// =========================================================================

func TestRepo_CreateContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	c := &EmployeeContract{
		EmployeeID:     uuid.New(),
		ContractNumber: "CTR-001",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
	}

	if err := repo.CreateContract(ctx, c); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	if c.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindContractByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	found, err := repo.FindContractByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindContractByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
}

func TestRepo_FindContractByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindContractByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}

func TestRepo_UpdateContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	created.Status = ContractStatusExpired
	if err := repo.UpdateContract(ctx, created); err != nil {
		t.Fatalf("UpdateContract failed: %v", err)
	}

	found, _ := repo.FindContractByID(ctx, created.ID)
	if found.Status != ContractStatusExpired {
		t.Errorf("expected status 'expired', got '%s'", found.Status)
	}
}

func TestRepo_DeleteContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	if err := repo.DeleteContract(ctx, created.ID); err != nil {
		t.Fatalf("DeleteContract failed: %v", err)
	}

	_, err := repo.FindContractByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting contract")
	}
}

func TestRepo_ListContracts_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	for i := 0; i < 5; i++ {
		createTestContract(repo, empID)
	}

	contracts, total, err := repo.ListContracts(ctx, 1, 2)
	if err != nil {
		t.Fatalf("ListContracts failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(contracts) != 2 {
		t.Errorf("expected 2 contracts (page 1 of 3), got %d", len(contracts))
	}
}

func TestRepo_FindContractsByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestContract(repo, emp1)
	createTestContract(repo, emp1)
	createTestContract(repo, emp1)
	createTestContract(repo, emp2)

	contracts, total, err := repo.FindContractsByEmployeeID(ctx, emp1, 1, 10)
	if err != nil {
		t.Fatalf("FindContractsByEmployeeID failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3 for emp1, got %d", total)
	}
	if len(contracts) != 3 {
		t.Errorf("expected 3 contracts, got %d", len(contracts))
	}
}
