package jobmanagement

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

// =========================================================================
// Job Title Service Tests (9.1)
// =========================================================================

func TestService_CreateJobTitle(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobTitleRequest{
		Name:         "Managerial",
		Descriptions: "For managerial positions",
		Status:       1,
	}

	resp, err := svc.CreateJobTitle(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobTitle failed: %v", err)
	}

	if resp.Name != "Managerial" {
		t.Errorf("expected name 'Managerial', got '%s'", resp.Name)
	}
	if resp.Status != 1 {
		t.Errorf("expected status 1, got %d", resp.Status)
	}
}

func TestService_CreateJobTitle_EmptyDescription(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobTitleRequest{
		Name: "Simple Title",
	}

	resp, err := svc.CreateJobTitle(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobTitle failed: %v", err)
	}
	if resp.Name != "Simple Title" {
		t.Errorf("expected name 'Simple Title', got '%s'", resp.Name)
	}
}

func TestService_GetJobTitleByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "Find Me"})

	found, err := svc.GetJobTitleByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetJobTitleByID failed: %v", err)
	}
	if found.Name != "Find Me" {
		t.Errorf("expected name 'Find Me', got '%s'", found.Name)
	}
}

func TestService_GetJobTitleByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetJobTitleByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetJobTitleByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetJobTitleByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent job title")
	}
}

func TestService_ListJobTitles_DefaultPagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.CreateJobTitle(ctx, CreateJobTitleRequest{
			Name: fmt.Sprintf("Title %d", i+1),
		})
	}

	resp, err := svc.ListJobTitles(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListJobTitles failed: %v", err)
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

func TestService_UpdateJobTitle(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "Before"})

	newName := "After"
	updated, err := svc.UpdateJobTitle(ctx, created.ID, UpdateJobTitleRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateJobTitle failed: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("expected name 'After', got '%s'", updated.Name)
	}
}

func TestService_DeleteJobTitle(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "To Delete"})

	if err := svc.DeleteJobTitle(ctx, created.ID); err != nil {
		t.Fatalf("DeleteJobTitle failed: %v", err)
	}

	_, err := svc.GetJobTitleByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting job title")
	}
}

// =========================================================================
// Job Title Sub Service Tests (9.2)
// =========================================================================

func TestService_CreateJobTitleSub(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	title, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "Main Title"})

	req := CreateJobTitleSubRequest{
		JobManagementTitleID: title.ID,
		Name:                 "Sub Title",
		Status:               1,
	}

	// Note: CreateJobTitleSub takes titleID as the parameter, not from the request body
	subResp, err := svc.CreateJobTitleSub(ctx, title.ID, req)
	if err != nil {
		t.Fatalf("CreateJobTitleSub failed: %v", err)
	}
	if subResp.Name != "Sub Title" {
		t.Errorf("expected name 'Sub Title', got '%s'", subResp.Name)
	}
}

func TestService_ListJobTitleSubs(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	title, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "Parent"})

	svc.CreateJobTitleSub(ctx, title.ID, CreateJobTitleSubRequest{
		JobManagementTitleID: title.ID,
		Name: "Sub A",
	})
	svc.CreateJobTitleSub(ctx, title.ID, CreateJobTitleSubRequest{
		JobManagementTitleID: title.ID,
		Name: "Sub B",
	})

	subs, err := svc.ListJobTitleSubs(ctx, title.ID)
	if err != nil {
		t.Fatalf("ListJobTitleSubs failed: %v", err)
	}
	if len(subs) != 2 {
		t.Errorf("expected 2 subs, got %d", len(subs))
	}
}

// =========================================================================
// Job Value Service Tests (9.3)
// =========================================================================

func TestService_CreateJobValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobValueRequest{
		Type:  "education",
		Level: intPtr(3),
		Sort:  intPtr(1),
	}

	resp, err := svc.CreateJobValue(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobValue failed: %v", err)
	}
	if resp.Type != "education" {
		t.Errorf("expected type 'education', got '%s'", resp.Type)
	}
}

func TestService_GetJobValueByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetJobValueByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdateJobValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobValue(ctx, CreateJobValueRequest{
		Type: "education",
		Sort: intPtr(1),
	})

	newType := "experience"
	updated, err := svc.UpdateJobValue(ctx, created.ID, UpdateJobValueRequest{
		Type: &newType,
	})
	if err != nil {
		t.Fatalf("UpdateJobValue failed: %v", err)
	}
	if updated.Type != "experience" {
		t.Errorf("expected type 'experience', got '%s'", updated.Type)
	}
}

func TestService_DeleteJobValue(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobValue(ctx, CreateJobValueRequest{Type: "delete-me"})

	if err := svc.DeleteJobValue(ctx, created.ID); err != nil {
		t.Fatalf("DeleteJobValue failed: %v", err)
	}
	_, err := svc.GetJobValueByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting job value")
	}
}

// =========================================================================
// Job Objective Service Tests (9.4)
// =========================================================================

func TestService_CreateJobObjective(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobObjectiveRequest{
		OrganizationID: createTestOrgID(),
		Nomenclature:   "Lead Strategy",
		FullCode:       "OBJ-001",
		Objective:      "Lead corporate strategy development",
	}

	resp, err := svc.CreateJobObjective(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobObjective failed: %v", err)
	}
	if resp.Nomenclature != "Lead Strategy" {
		t.Errorf("expected 'Lead Strategy', got '%s'", resp.Nomenclature)
	}
	if resp.FullCode != "OBJ-001" {
		t.Errorf("expected 'OBJ-001', got '%s'", resp.FullCode)
	}
}

func TestService_ListJobObjectives_Pagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		svc.CreateJobObjective(ctx, CreateJobObjectiveRequest{
			OrganizationID: createTestOrgID(),
			Nomenclature:   fmt.Sprintf("Objective %d", i+1),
			FullCode:       fmt.Sprintf("OBJ-%03d", i+1),
		})
	}

	resp, err := svc.ListJobObjectives(ctx, 1, 2)
	if err != nil {
		t.Fatalf("ListJobObjectives failed: %v", err)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
	if len(resp.Data.([]JobObjectiveResponse)) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(resp.Data.([]JobObjectiveResponse)))
	}
}

// =========================================================================
// Job Identification Service Tests (9.5)
// =========================================================================

func TestService_CreateJobIdentification(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobIdentificationRequest{
		OrganizationID: createTestOrgID(),
		Nomenclature:   "Senior Manager",
		FullCode:       "SM-001",
		GradingID:      uuid.New().String(),
	}

	resp, err := svc.CreateJobIdentification(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobIdentification failed: %v", err)
	}
	if resp.Nomenclature != "Senior Manager" {
		t.Errorf("expected 'Senior Manager', got '%s'", resp.Nomenclature)
	}
}

func TestService_UpdateJobIdentification(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobIdentification(ctx, CreateJobIdentificationRequest{
		OrganizationID: createTestOrgID(),
		Nomenclature:   "Manager",
		FullCode:       "MGR-001",
		GradingID:      uuid.New().String(),
	})

	newNomenclature := "Director"
	updated, err := svc.UpdateJobIdentification(ctx, created.ID, UpdateJobIdentificationRequest{
		Nomenclature: &newNomenclature,
	})
	if err != nil {
		t.Fatalf("UpdateJobIdentification failed: %v", err)
	}
	if updated.Nomenclature != "Director" {
		t.Errorf("expected 'Director', got '%s'", updated.Nomenclature)
	}
}

// =========================================================================
// Job Responsibility Service Tests (9.6)
// =========================================================================

func TestService_CreateJobResponsibility(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobResponsibilityRequest{
		OrganizationID:    createTestOrgID(),
		Nomenclature:      "Project Management",
		FullCode:          "RESP-001",
		MainTask:          "Manage projects",
		Activities:        "Planning, execution",
		Outputs:           "Project reports",
		SuccessIndicators: "On-time delivery",
	}

	resp, err := svc.CreateJobResponsibility(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobResponsibility failed: %v", err)
	}
	if resp.Nomenclature != "Project Management" {
		t.Errorf("expected 'Project Management', got '%s'", resp.Nomenclature)
	}
	if resp.MainTask != "Manage projects" {
		t.Errorf("expected 'Manage projects', got '%s'", resp.MainTask)
	}
}

// =========================================================================
// Job HRAuthority Service Tests (9.8)
// =========================================================================

func TestService_CreateJobHRAuthority(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateJobHRAuthorityRequest{
		OrganizationID: createTestOrgID(),
		Nomenclature:   "Hiring Authority",
		FullCode:       "HR-001",
		Description:    "Can hire up to manager level",
	}

	resp, err := svc.CreateJobHRAuthority(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobHRAuthority failed: %v", err)
	}
	if resp.Nomenclature != "Hiring Authority" {
		t.Errorf("expected 'Hiring Authority', got '%s'", resp.Nomenclature)
	}
	if resp.Description != "Can hire up to manager level" {
		t.Errorf("expected description mismatch, got '%s'", resp.Description)
	}
}

// =========================================================================
// Job Score Service Tests (9.17)
// =========================================================================

func TestService_UpsertJobScore_Create(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	orgID := createTestOrgID()
	req := UpdateJobScoreRequest{
		JobValueWithFinancial:    uint64Ptr(2000),
		JobValueWithoutFinancial: uint64Ptr(1500),
		HasFinancialAuthority:    boolPtr(true),
	}

	resp, err := svc.UpsertJobScore(ctx, orgID, req)
	if err != nil {
		t.Fatalf("UpsertJobScore failed: %v", err)
	}
	if resp.JobValueWithFinancial != 2000 {
		t.Errorf("expected 2000, got %d", resp.JobValueWithFinancial)
	}
	if resp.HasFinancialAuthority != true {
		t.Errorf("expected true, got %v", resp.HasFinancialAuthority)
	}
}

func TestService_UpsertJobScore_Update(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	orgID := createTestOrgID()
	svc.UpsertJobScore(ctx, orgID, UpdateJobScoreRequest{
		JobValueWithFinancial: uint64Ptr(1000),
	})

	// Upsert again
	resp, err := svc.UpsertJobScore(ctx, orgID, UpdateJobScoreRequest{
		JobValueWithFinancial: uint64Ptr(3000),
	})
	if err != nil {
		t.Fatalf("UpsertJobScore (update) failed: %v", err)
	}
	if resp.JobValueWithFinancial != 3000 {
		t.Errorf("expected 3000, got %d", resp.JobValueWithFinancial)
	}
}

func TestService_GetJobScoreByOrganization(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	orgID := createTestOrgID()
	svc.UpsertJobScore(ctx, orgID, UpdateJobScoreRequest{
		JobValueWithFinancial: uint64Ptr(5000),
	})

	found, err := svc.GetJobScoreByOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("GetJobScoreByOrganization failed: %v", err)
	}
	if found.JobValueWithFinancial != 5000 {
		t.Errorf("expected 5000, got %d", found.JobValueWithFinancial)
	}
}

func TestService_GetJobScoreByOrganization_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetJobScoreByOrganization(ctx, "nonexistent-org")
	if err == nil {
		t.Fatal("expected error for non-existent job score")
	}
}

// =========================================================================
// Job Competency Group Service Tests (9.18)
// =========================================================================

func TestService_CreateJobCompetencyGroup(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	orgID := createTestOrgID()
	req := CreateJobCompetencyGroupRequest{
		OrganizationID: orgID,
		Category:       "technical",
		Weight:         70.0,
	}

	resp, err := svc.CreateJobCompetencyGroup(ctx, req)
	if err != nil {
		t.Fatalf("CreateJobCompetencyGroup failed: %v", err)
	}
	if resp.Category != "technical" {
		t.Errorf("expected 'technical', got '%s'", resp.Category)
	}
	if resp.Weight != 70.0 {
		t.Errorf("expected weight 70.0, got %f", resp.Weight)
	}
}

func TestService_ListJobCompetencyGroups(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	orgID := createTestOrgID()
	svc.CreateJobCompetencyGroup(ctx, CreateJobCompetencyGroupRequest{
		OrganizationID: orgID,
		Category:       "technical",
		Weight:         60.0,
	})
	svc.CreateJobCompetencyGroup(ctx, CreateJobCompetencyGroupRequest{
		OrganizationID: orgID,
		Category:       "managerial",
		Weight:         40.0,
	})

	groups, err := svc.ListJobCompetencyGroups(ctx, orgID)
	if err != nil {
		t.Fatalf("ListJobCompetencyGroups failed: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestService_UpdateJobCompetencyGroup(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateJobCompetencyGroup(ctx, CreateJobCompetencyGroupRequest{
		OrganizationID: createTestOrgID(),
		Category:       "technical",
		Weight:         50.0,
	})

	newWeight := 80.0
	updated, err := svc.UpdateJobCompetencyGroup(ctx, created.ID, UpdateJobCompetencyGroupRequest{
		Weight: &newWeight,
	})
	if err != nil {
		t.Fatalf("UpdateJobCompetencyGroup failed: %v", err)
	}
	if updated.Weight != 80.0 {
		t.Errorf("expected weight 80.0, got %f", updated.Weight)
	}
}

// =========================================================================
// Nested Entity Tests: JobValue with JobTitleSub relation
// =========================================================================

func TestService_CreateJobValue_WithTitleSub(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	title, _ := svc.CreateJobTitle(ctx, CreateJobTitleRequest{Name: "Functional"})
	sub, _ := svc.CreateJobTitleSub(ctx, title.ID, CreateJobTitleSubRequest{
		JobManagementTitleID: title.ID,
		Name: "Sub Functional",
	})

	valueReq := CreateJobValueRequest{
		JobManagementTitleSubID: &sub.ID,
		Type:                   "education",
		Sort:                   intPtr(1),
	}
	resp, err := svc.CreateJobValue(ctx, valueReq)
	if err != nil {
		t.Fatalf("CreateJobValue with title sub failed: %v", err)
	}
	if resp.Type != "education" {
		t.Errorf("expected 'education', got '%s'", resp.Type)
	}
	// The sub name should have been copied
	if resp.JobManagementTitleSubName != "Sub Functional" {
		t.Errorf("expected sub name 'Sub Functional', got '%s'", resp.JobManagementTitleSubName)
	}
}
