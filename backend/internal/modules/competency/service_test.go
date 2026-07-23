package competency

import (
	"context"
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

// =========================================================================
// Competency Service Tests (8.1)
// =========================================================================

func TestService_CreateCompetency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetencyRequest{
		Name:       "Leadership",
		Field:      strPtr("Managerial"),
		Cluster:    strPtr("Core"),
		Definition: strPtr("Ability to lead teams effectively"),
	}

	resp, err := svc.CreateCompetency(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetency failed: %v", err)
	}

	if resp.Name != "Leadership" {
		t.Errorf("expected name 'Leadership', got '%s'", resp.Name)
	}
	if resp.Field != "Managerial" {
		t.Errorf("expected field 'Managerial', got '%s'", resp.Field)
	}
	if resp.Cluster != "Core" {
		t.Errorf("expected cluster 'Core', got '%s'", resp.Cluster)
	}
}

func TestService_CreateCompetency_Minimal(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetencyRequest{Name: "Minimal"}

	resp, err := svc.CreateCompetency(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetency failed: %v", err)
	}
	if resp.Name != "Minimal" {
		t.Errorf("expected name 'Minimal', got '%s'", resp.Name)
	}
}

func TestService_GetCompetencyByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Find Me"})

	found, err := svc.GetCompetencyByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetCompetencyByID failed: %v", err)
	}
	if found.Name != "Find Me" {
		t.Errorf("expected name 'Find Me', got '%s'", found.Name)
	}
}

func TestService_GetCompetencyByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetCompetencyByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent competency")
	}
}

func TestService_ListCompetencies_DefaultPagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.CreateCompetency(ctx, CreateCompetencyRequest{
			Name: "Competency",
		})
	}

	resp, err := svc.ListCompetencies(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListCompetencies failed: %v", err)
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

func TestService_UpdateCompetency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Before"})

	newName := "After"
	updated, err := svc.UpdateCompetency(ctx, created.ID, UpdateCompetencyRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateCompetency failed: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("expected name 'After', got '%s'", updated.Name)
	}
}

func TestService_DeleteCompetency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "To Delete"})

	if err := svc.DeleteCompetency(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetency failed: %v", err)
	}

	_, err := svc.GetCompetencyByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency")
	}
}

// =========================================================================
// CompetenceValue Service Tests (8.2 — legacy)
// =========================================================================

func TestService_CreateCompetenceValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetenceValueRequest{
		Type:        strPtr("score"),
		Level:       intPtr(3),
		Name:        "Excellent",
		Point:       intPtr(5),
		Description: strPtr("Top performance"),
	}

	resp, err := svc.CreateCompetenceValue(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetenceValue failed: %v", err)
	}
	if resp.Name != "Excellent" {
		t.Errorf("expected name 'Excellent', got '%s'", resp.Name)
	}
	if resp.Point != 5 {
		t.Errorf("expected point 5, got %d", resp.Point)
	}
}

func TestService_GetCompetenceValueByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetenceValueByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateCompetenceValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetenceValue(ctx, CreateCompetenceValueRequest{
		Name: "Before",
	})

	newName := "After"
	updated, err := svc.UpdateCompetenceValue(ctx, created.ID, UpdateCompetenceValueRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateCompetenceValue failed: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("expected name 'After', got '%s'", updated.Name)
	}
}

func TestService_DeleteCompetenceValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetenceValue(ctx, CreateCompetenceValueRequest{Name: "To Delete"})

	if err := svc.DeleteCompetenceValue(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetenceValue failed: %v", err)
	}

	_, err := svc.GetCompetenceValueByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competence value")
	}
}

// =========================================================================
// CompetencyValue Service Tests (8.3 — structured)
// =========================================================================

func TestService_CreateCompetencyValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetencyValueRequest{
		Type:        "technical",
		Name:        "Problem Solving",
		Slug:        "problem-solving",
		Level:       4,
		Code:        strPtr("T-004"),
		Description: strPtr("Solve complex problems"),
	}

	resp, err := svc.CreateCompetencyValue(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetencyValue failed: %v", err)
	}
	if resp.Name != "Problem Solving" {
		t.Errorf("expected name 'Problem Solving', got '%s'", resp.Name)
	}
	if resp.Slug != "problem-solving" {
		t.Errorf("expected slug 'problem-solving', got '%s'", resp.Slug)
	}
	if resp.Level != 4 {
		t.Errorf("expected level 4, got %d", resp.Level)
	}
}

func TestService_GetCompetencyValueByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyValueByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateCompetencyValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyValue(ctx, CreateCompetencyValueRequest{
		Type:  "technical",
		Name:  "Before",
		Slug:  "before",
		Level: 3,
	})

	newName := "After"
	updated, err := svc.UpdateCompetencyValue(ctx, created.ID, UpdateCompetencyValueRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateCompetencyValue failed: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("expected name 'After', got '%s'", updated.Name)
	}
}

func TestService_DeleteCompetencyValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyValue(ctx, CreateCompetencyValueRequest{
		Type:  "technical",
		Name:  "To Delete",
		Slug:  "to-delete",
		Level: 2,
	})

	if err := svc.DeleteCompetencyValue(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetencyValue failed: %v", err)
	}

	_, err := svc.GetCompetencyValueByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency value")
	}
}

// =========================================================================
// CompetencyEvent Service Tests (8.4)
// =========================================================================

func TestService_CreateCompetencyEvent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetencyEventRequest{
		Type:         "manual",
		PeriodType:   "annual",
		PeriodYear:   2026,
		PeriodNumber: intPtr(1),
	}

	resp, err := svc.CreateCompetencyEvent(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetencyEvent failed: %v", err)
	}
	if resp.Type != "manual" {
		t.Errorf("expected type 'manual', got '%s'", resp.Type)
	}
	if resp.PeriodType != "annual" {
		t.Errorf("expected period_type 'annual', got '%s'", resp.PeriodType)
	}
	if resp.PeriodYear != 2026 {
		t.Errorf("expected period_year 2026, got %d", resp.PeriodYear)
	}
	if resp.Status != "active" {
		t.Errorf("expected default status 'active', got '%s'", resp.Status)
	}
}

func TestService_CreateCompetencyEvent_DefaultStatus(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "auto",
		PeriodType: "semester",
		PeriodYear: 2025,
	})
	if err != nil {
		t.Fatalf("CreateCompetencyEvent failed: %v", err)
	}
	if resp.Status != "active" {
		t.Errorf("expected default status 'active', got '%s'", resp.Status)
	}
}

func TestService_GetCompetencyEventByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyEventByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateCompetencyEvent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	})

	newStatus := "closed"
	updated, err := svc.UpdateCompetencyEvent(ctx, created.ID, UpdateCompetencyEventRequest{
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("UpdateCompetencyEvent failed: %v", err)
	}
	if updated.Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", updated.Status)
	}
}

func TestService_DeleteCompetencyEvent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	})

	if err := svc.DeleteCompetencyEvent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetencyEvent failed: %v", err)
	}

	_, err := svc.GetCompetencyEventByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency event")
	}
}

// =========================================================================
// CompetencyEventTarget Service Tests (8.5)
// =========================================================================

func TestService_CreateCompetencyEventTarget(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	event, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	})

	req := CreateCompetencyEventTargetRequest{
		CompetencyEventID:  event.ID,
		OrganizationID:     createTestOrgID(),
		MissingSelf:        intPtr(0),
		MissingSuperior:    intPtr(1),
		MissingPeer:        intPtr(0),
		MissingSubordinate: intPtr(2),
	}

	resp, err := svc.CreateCompetencyEventTarget(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetencyEventTarget failed: %v", err)
	}
	if resp.CompetencyEventID != event.ID {
		t.Errorf("expected event ID '%s', got '%s'", event.ID, resp.CompetencyEventID)
	}
	if resp.MissingSuperior != 1 {
		t.Errorf("expected missing_superior 1, got %d", resp.MissingSuperior)
	}
}

func TestService_GetCompetencyEventTargetByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyEventTargetByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateCompetencyEventTarget(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	event, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	})
	created, _ := svc.CreateCompetencyEventTarget(ctx, CreateCompetencyEventTargetRequest{
		CompetencyEventID: event.ID,
		OrganizationID:    createTestOrgID(),
	})

	newMissing := 5
	updated, err := svc.UpdateCompetencyEventTarget(ctx, created.ID, UpdateCompetencyEventTargetRequest{
		MissingSelf: &newMissing,
	})
	if err != nil {
		t.Fatalf("UpdateCompetencyEventTarget failed: %v", err)
	}
	if updated.MissingSelf != 5 {
		t.Errorf("expected missing_self 5, got %d", updated.MissingSelf)
	}
}

func TestService_DeleteCompetencyEventTarget(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	event, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	})
	created, _ := svc.CreateCompetencyEventTarget(ctx, CreateCompetencyEventTargetRequest{
		CompetencyEventID: event.ID,
		OrganizationID:    createTestOrgID(),
	})

	if err := svc.DeleteCompetencyEventTarget(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetencyEventTarget failed: %v", err)
	}

	_, err := svc.GetCompetencyEventTargetByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency event target")
	}
}

// =========================================================================
// CompetencyScore Service Tests (8.6)
// =========================================================================

func TestService_CreateCompetencyScore(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCompetencyScoreRequest{
		OrganizationID:          createTestOrgID(),
		TechnicalGapPercentage:  10.5,
		ManagerialGapPercentage: 15.2,
		TotalGapPercentage:      12.0,
		TotalGradePercentage:    88.0,
	}

	resp, err := svc.CreateCompetencyScore(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetencyScore failed: %v", err)
	}
	if resp.TechnicalGapPercentage != 10.5 {
		t.Errorf("expected technical_gap 10.5, got %f", resp.TechnicalGapPercentage)
	}
	if resp.TotalGradePercentage != 88.0 {
		t.Errorf("expected total_grade 88.0, got %f", resp.TotalGradePercentage)
	}
}

func TestService_GetCompetencyScoreByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyScoreByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetCompetencyScoreByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyScoreByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent score")
	}
}

func TestService_UpdateCompetencyScore(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID:         createTestOrgID(),
		TotalGradePercentage:   75.0,
	})

	newGrade := 95.0
	updated, err := svc.UpdateCompetencyScore(ctx, created.ID, UpdateCompetencyScoreRequest{
		TotalGradePercentage: &newGrade,
	})
	if err != nil {
		t.Fatalf("UpdateCompetencyScore failed: %v", err)
	}
	if updated.TotalGradePercentage != 95.0 {
		t.Errorf("expected total_grade 95.0, got %f", updated.TotalGradePercentage)
	}
}

func TestService_DeleteCompetencyScore(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID: createTestOrgID(),
	})

	if err := svc.DeleteCompetencyScore(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetencyScore failed: %v", err)
	}

	_, err := svc.GetCompetencyScoreByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency score")
	}
}

// =========================================================================
// CompetencyScoreDetail Service Tests (8.7)
// =========================================================================

func TestService_CreateCompetencyScoreDetail(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	score, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID: createTestOrgID(),
	})
	comp, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Technical Skill"})

	req := CreateCompetencyScoreDetailRequest{
		CompetencyScoreID:     score.ID,
		CompetencyID:          comp.ID,
		Type:                  "technical",
		StandardLevel:         intPtr(4),
		StandardWeight:        70.0,
		EmployeeLevel:         intPtr(3),
		GapPercentage:         25.0,
		WeightedGapPercentage: 17.5,
	}

	resp, err := svc.CreateCompetencyScoreDetail(ctx, req)
	if err != nil {
		t.Fatalf("CreateCompetencyScoreDetail failed: %v", err)
	}
	if resp.Type != "technical" {
		t.Errorf("expected type 'technical', got '%s'", resp.Type)
	}
	if resp.StandardWeight != 70.0 {
		t.Errorf("expected standard_weight 70.0, got %f", resp.StandardWeight)
	}
}

func TestService_GetCompetencyScoreDetailByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetCompetencyScoreDetailByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateCompetencyScoreDetail(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	score, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID: createTestOrgID(),
	})
	comp, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Skill"})
	created, _ := svc.CreateCompetencyScoreDetail(ctx, CreateCompetencyScoreDetailRequest{
		CompetencyScoreID: score.ID,
		CompetencyID:      comp.ID,
		Type:              "technical",
	})

	newGap := 50.0
	updated, err := svc.UpdateCompetencyScoreDetail(ctx, created.ID, UpdateCompetencyScoreDetailRequest{
		GapPercentage: &newGap,
	})
	if err != nil {
		t.Fatalf("UpdateCompetencyScoreDetail failed: %v", err)
	}
	if updated.GapPercentage != 50.0 {
		t.Errorf("expected gap_percentage 50.0, got %f", updated.GapPercentage)
	}
}

func TestService_ListCompetencyScoreDetails(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	score, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID: createTestOrgID(),
	})
	comp1, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Skill A"})
	comp2, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Skill B"})

	svc.CreateCompetencyScoreDetail(ctx, CreateCompetencyScoreDetailRequest{
		CompetencyScoreID: score.ID,
		CompetencyID:      comp1.ID,
		Type:              "technical",
	})
	svc.CreateCompetencyScoreDetail(ctx, CreateCompetencyScoreDetailRequest{
		CompetencyScoreID: score.ID,
		CompetencyID:      comp2.ID,
		Type:              "managerial",
	})

	resp, err := svc.ListCompetencyScoreDetails(ctx, score.ID, 1, 10)
	if err != nil {
		t.Fatalf("ListCompetencyScoreDetails failed: %v", err)
	}
	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Data.([]CompetencyScoreDetailResponse)) != 2 {
		t.Errorf("expected 2 details, got %d", len(resp.Data.([]CompetencyScoreDetailResponse)))
	}
}

func TestService_DeleteCompetencyScoreDetail(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	score, _ := svc.CreateCompetencyScore(ctx, CreateCompetencyScoreRequest{
		OrganizationID: createTestOrgID(),
	})
	comp, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Skill"})
	created, _ := svc.CreateCompetencyScoreDetail(ctx, CreateCompetencyScoreDetailRequest{
		CompetencyScoreID: score.ID,
		CompetencyID:      comp.ID,
		Type:              "technical",
	})

	if err := svc.DeleteCompetencyScoreDetail(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetencyScoreDetail failed: %v", err)
	}

	_, err := svc.GetCompetencyScoreDetailByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency score detail")
	}
}
