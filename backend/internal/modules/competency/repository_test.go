package competency

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Competency Repository Tests (8.1)
// =========================================================================

func TestRepository_CreateCompetency(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	c := &Competency{
		Name: "Strategic Thinking",
		Field: strPtr("Managerial"),
		Cluster: strPtr("Core"),
	}

	if err := repo.CreateCompetency(ctx, c); err != nil {
		t.Fatalf("CreateCompetency failed: %v", err)
	}

	if c.ID == uuid.Nil {
		t.Error("expected competency ID to be generated")
	}
}

func TestRepository_FindCompetencyByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestCompetency(ctx, repo)

	found, err := repo.FindCompetencyByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindCompetencyByID failed: %v", err)
	}

	if found.Name != "Leadership" {
		t.Errorf("expected name 'Leadership', got '%s'", found.Name)
	}
}

func TestRepository_FindCompetencyByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindCompetencyByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent competency")
	}
}

func TestRepository_FindAllCompetencies_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		repo.CreateCompetency(ctx, &Competency{Name: "Competency"})
	}

	list, total, err := repo.FindAllCompetencies(ctx, 1, 2)
	if err != nil {
		t.Fatalf("FindAllCompetencies failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 competencies on page 1, got %d", len(list))
	}
}

func TestRepository_UpdateCompetency(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestCompetency(ctx, repo)

	created.Name = "Updated Leadership"
	if err := repo.UpdateCompetency(ctx, created); err != nil {
		t.Fatalf("UpdateCompetency failed: %v", err)
	}

	found, _ := repo.FindCompetencyByID(ctx, created.ID)
	if found.Name != "Updated Leadership" {
		t.Errorf("expected name 'Updated Leadership', got '%s'", found.Name)
	}
}

func TestRepository_DeleteCompetency(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestCompetency(ctx, repo)

	if err := repo.DeleteCompetency(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCompetency failed: %v", err)
	}

	_, err := repo.FindCompetencyByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency")
	}
}

// =========================================================================
// CompetenceValue Repository Tests (8.2 — legacy)
// =========================================================================

func TestRepository_CompetenceValueCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	v := &CompetenceValue{
		Type:  strPtr("score"),
		Level: intPtr(3),
		Name:  "Good",
		Point: intPtr(3),
	}
	if err := repo.CreateCompetenceValue(ctx, v); err != nil {
		t.Fatalf("CreateCompetenceValue failed: %v", err)
	}

	// Find
	found, err := repo.FindCompetenceValueByID(ctx, v.ID)
	if err != nil {
		t.Fatalf("FindCompetenceValueByID failed: %v", err)
	}
	if found.Name != "Good" {
		t.Errorf("expected name 'Good', got '%s'", found.Name)
	}

	// Update
	found.Name = "Excellent"
	found.Point = intPtr(5)
	repo.UpdateCompetenceValue(ctx, found)
	updated, _ := repo.FindCompetenceValueByID(ctx, v.ID)
	if updated.Name != "Excellent" {
		t.Errorf("expected name 'Excellent', got '%s'", updated.Name)
	}
	if *updated.Point != 5 {
		t.Errorf("expected point 5, got %d", *updated.Point)
	}

	// Delete
	repo.DeleteCompetenceValue(ctx, v.ID)
	_, err = repo.FindCompetenceValueByID(ctx, v.ID)
	if err == nil {
		t.Fatal("expected error after deleting competence value")
	}
}

// =========================================================================
// CompetencyValue Repository Tests (8.3 — structured)
// =========================================================================

func TestRepository_CompetencyValueCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	v := &CompetencyValue{
		Type:  "technical",
		Name:  "Analytical Thinking",
		Slug:  "analytical-thinking",
		Level: 4,
		Code:  strPtr("T-004"),
	}
	if err := repo.CreateCompetencyValue(ctx, v); err != nil {
		t.Fatalf("CreateCompetencyValue failed: %v", err)
	}

	// Find
	found, err := repo.FindCompetencyValueByID(ctx, v.ID)
	if err != nil {
		t.Fatalf("FindCompetencyValueByID failed: %v", err)
	}
	if found.Name != "Analytical Thinking" {
		t.Errorf("expected name 'Analytical Thinking', got '%s'", found.Name)
	}
	if found.Slug != "analytical-thinking" {
		t.Errorf("expected slug 'analytical-thinking', got '%s'", found.Slug)
	}

	// Update
	found.Level = 5
	repo.UpdateCompetencyValue(ctx, found)
	updated, _ := repo.FindCompetencyValueByID(ctx, v.ID)
	if updated.Level != 5 {
		t.Errorf("expected level 5, got %d", updated.Level)
	}

	// List
	repo.CreateCompetencyValue(ctx, &CompetencyValue{
		Type:  "managerial",
		Name:  "Team Leadership",
		Slug:  "team-leadership",
		Level: 3,
	})
	list, total, _ := repo.FindAllCompetencyValues(ctx, 1, 10)
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}

	// Delete
	repo.DeleteCompetencyValue(ctx, v.ID)
	_, err = repo.FindCompetencyValueByID(ctx, v.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency value")
	}
}

// =========================================================================
// CompetencyEvent Repository Tests (8.4)
// =========================================================================

func TestRepository_CompetencyEventCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	e := &CompetencyEvent{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
		Status:     "active",
	}
	if err := repo.CreateCompetencyEvent(ctx, e); err != nil {
		t.Fatalf("CreateCompetencyEvent failed: %v", err)
	}

	// Find
	found, err := repo.FindCompetencyEventByID(ctx, e.ID)
	if err != nil {
		t.Fatalf("FindCompetencyEventByID failed: %v", err)
	}
	if found.Type != "manual" {
		t.Errorf("expected type 'manual', got '%s'", found.Type)
	}

	// Update
	found.Status = "closed"
	repo.UpdateCompetencyEvent(ctx, found)
	updated, _ := repo.FindCompetencyEventByID(ctx, e.ID)
	if updated.Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", updated.Status)
	}

	// Delete
	repo.DeleteCompetencyEvent(ctx, e.ID)
	_, err = repo.FindCompetencyEventByID(ctx, e.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency event")
	}
}

// =========================================================================
// CompetencyEventTarget Repository Tests (8.5)
// =========================================================================

func TestRepository_CompetencyEventTargetCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create parent event
	event := &CompetencyEvent{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
	}
	repo.CreateCompetencyEvent(ctx, event)

	// Create target
	target := &CompetencyEventTarget{
		CompetencyEventID:  event.ID,
		OrganizationID:     uuid.New(),
		MissingSuperior:    1,
		MissingSubordinate: 2,
	}
	if err := repo.CreateCompetencyEventTarget(ctx, target); err != nil {
		t.Fatalf("CreateCompetencyEventTarget failed: %v", err)
	}

	// Find with preload
	found, err := repo.FindCompetencyEventTargetByID(ctx, target.ID)
	if err != nil {
		t.Fatalf("FindCompetencyEventTargetByID failed: %v", err)
	}
	if found.OrganizationID != target.OrganizationID {
		t.Errorf("expected organization ID '%s', got '%s'", target.OrganizationID, found.OrganizationID)
	}
	if found.CompetencyEvent == nil {
		t.Error("expected CompetencyEvent relation to be preloaded")
	}

	// Update
	found.MissingSelf = 3
	repo.UpdateCompetencyEventTarget(ctx, found)
	updated, _ := repo.FindCompetencyEventTargetByID(ctx, target.ID)
	if updated.MissingSelf != 3 {
		t.Errorf("expected missing_self 3, got %d", updated.MissingSelf)
	}

	// Delete
	repo.DeleteCompetencyEventTarget(ctx, target.ID)
	_, err = repo.FindCompetencyEventTargetByID(ctx, target.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency event target")
	}
}

// =========================================================================
// CompetencyScore Repository Tests (8.6)
// =========================================================================

func TestRepository_CompetencyScoreCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	s := &CompetencyScore{
		OrganizationID:          uuid.New(),
		TechnicalGapPercentage:  10.5,
		ManagerialGapPercentage: 15.2,
		TotalGapPercentage:      12.0,
		TotalGradePercentage:    88.0,
	}
	if err := repo.CreateCompetencyScore(ctx, s); err != nil {
		t.Fatalf("CreateCompetencyScore failed: %v", err)
	}

	// Find with preload (Details)
	found, err := repo.FindCompetencyScoreByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("FindCompetencyScoreByID failed: %v", err)
	}
	if found.OrganizationID != s.OrganizationID {
		t.Errorf("expected org ID '%s', got '%s'", s.OrganizationID, found.OrganizationID)
	}
	if found.Details == nil {
		t.Error("expected Details slice to be initialized")
	}

	// Update
	found.TotalGradePercentage = 95.0
	repo.UpdateCompetencyScore(ctx, found)
	updated, _ := repo.FindCompetencyScoreByID(ctx, s.ID)
	if updated.TotalGradePercentage != 95.0 {
		t.Errorf("expected total_grade 95.0, got %f", updated.TotalGradePercentage)
	}

	// List
	repo.CreateCompetencyScore(ctx, &CompetencyScore{
		OrganizationID: uuid.New(),
	})
	list, total, _ := repo.FindAllCompetencyScores(ctx, 1, 10)
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}

	// Delete
	repo.DeleteCompetencyScore(ctx, s.ID)
	_, err = repo.FindCompetencyScoreByID(ctx, s.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency score")
	}
}

// =========================================================================
// CompetencyScoreDetail Repository Tests (8.7)
// =========================================================================

func TestRepository_CompetencyScoreDetailCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create parent score and competency
	score := createTestCompetencyScore(ctx, repo)
	comp := createTestCompetency(ctx, repo)

	// Create detail
	d := &CompetencyScoreDetail{
		CompetencyScoreID:     score.ID,
		CompetencyID:          comp.ID,
		Type:                  "technical",
		StandardLevel:         intPtr(4),
		StandardWeight:        70.0,
		EmployeeLevel:         intPtr(3),
		GapPercentage:         25.0,
		WeightedGapPercentage: 17.5,
	}
	if err := repo.CreateCompetencyScoreDetail(ctx, d); err != nil {
		t.Fatalf("CreateCompetencyScoreDetail failed: %v", err)
	}

	// Find
	found, err := repo.FindCompetencyScoreDetailByID(ctx, d.ID)
	if err != nil {
		t.Fatalf("FindCompetencyScoreDetailByID failed: %v", err)
	}
	if found.Type != "technical" {
		t.Errorf("expected type 'technical', got '%s'", found.Type)
	}
	if found.StandardWeight != 70.0 {
		t.Errorf("expected standard_weight 70.0, got %f", found.StandardWeight)
	}

	// Update
	found.GapPercentage = 50.0
	repo.UpdateCompetencyScoreDetail(ctx, found)
	updated, _ := repo.FindCompetencyScoreDetailByID(ctx, d.ID)
	if updated.GapPercentage != 50.0 {
		t.Errorf("expected gap_percentage 50.0, got %f", updated.GapPercentage)
	}

	// List by score ID
	repo.CreateCompetencyScoreDetail(ctx, &CompetencyScoreDetail{
		CompetencyScoreID: score.ID,
		CompetencyID:      comp.ID,
		Type:              "managerial",
	})
	list, total, _ := repo.FindAllCompetencyScoreDetails(ctx, score.ID, 1, 10)
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 details, got %d", len(list))
	}

	// Delete
	repo.DeleteCompetencyScoreDetail(ctx, d.ID)
	_, err = repo.FindCompetencyScoreDetailByID(ctx, d.ID)
	if err == nil {
		t.Fatal("expected error after deleting competency score detail")
	}
}
