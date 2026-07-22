package jobmanagement

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Job Title CRUD Tests (9.1)
// =========================================================================

func TestRepository_CreateJobTitle(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	name := "Managerial Title"
	desc := "For managerial positions"
	t1 := &JobTitle{
		Name:         &name,
		Descriptions: &desc,
		Status:       int8Ptr(1),
	}

	if err := repo.CreateJobTitle(ctx, t1); err != nil {
		t.Fatalf("CreateJobTitle failed: %v", err)
	}
	if t1.ID == uuid.Nil {
		t.Error("expected job title ID to be generated")
	}
}

func TestRepository_FindJobTitleByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestJobTitle(ctx, repo)

	found, err := repo.FindJobTitleByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindJobTitleByID failed: %v", err)
	}
	if found.Name == nil || *found.Name != "Test Title" {
		t.Errorf("expected name 'Test Title', got %v", found.Name)
	}
	if found.Status == nil || *found.Status != 1 {
		t.Errorf("expected status 1, got %v", found.Status)
	}
}

func TestRepository_FindJobTitleByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindJobTitleByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent job title")
	}
}

func TestRepository_FindAllJobTitles_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Title %d", i+1)
		t1 := &JobTitle{Name: &name}
		if err := repo.CreateJobTitle(ctx, t1); err != nil {
			t.Fatalf("failed to create title %d: %v", i, err)
		}
	}

	titles, total, err := repo.FindAllJobTitles(ctx, 1, 2)
	if err != nil {
		t.Fatalf("FindAllJobTitles failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(titles) != 2 {
		t.Errorf("expected 2 titles on page 1, got %d", len(titles))
	}
}

func TestRepository_UpdateJobTitle(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestJobTitle(ctx, repo)
	newName := "Updated Title"
	created.Name = &newName

	if err := repo.UpdateJobTitle(ctx, created); err != nil {
		t.Fatalf("UpdateJobTitle failed: %v", err)
	}

	found, _ := repo.FindJobTitleByID(ctx, created.ID)
	if found.Name == nil || *found.Name != "Updated Title" {
		t.Errorf("expected name 'Updated Title', got %v", found.Name)
	}
}

func TestRepository_DeleteJobTitle(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestJobTitle(ctx, repo)

	if err := repo.DeleteJobTitle(ctx, created.ID); err != nil {
		t.Fatalf("DeleteJobTitle failed: %v", err)
	}

	_, err := repo.FindJobTitleByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting job title")
	}
}

// =========================================================================
// Job Title Sub CRUD Tests (9.2)
// =========================================================================

func TestRepository_JobTitleSubCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	title := createTestJobTitle(ctx, repo)

	// Create
	subName := "Sub Title A"
	sub := &JobTitleSub{
		JobManagementTitleID: &title.ID,
		Name:                &subName,
		Status:              int8Ptr(1),
	}
	if err := repo.CreateJobTitleSub(ctx, sub); err != nil {
		t.Fatalf("CreateJobTitleSub failed: %v", err)
	}
	if sub.ID == uuid.Nil {
		t.Error("expected sub ID to be generated")
	}

	// Find
	found, err := repo.FindJobTitleSubByID(ctx, sub.ID)
	if err != nil {
		t.Fatalf("FindJobTitleSubByID failed: %v", err)
	}
	if found.Name == nil || *found.Name != "Sub Title A" {
		t.Errorf("expected name 'Sub Title A', got %v", found.Name)
	}

	// List by title ID
	subs, err := repo.FindJobTitleSubsByTitleID(ctx, title.ID)
	if err != nil {
		t.Fatalf("FindJobTitleSubsByTitleID failed: %v", err)
	}
	if len(subs) != 1 {
		t.Errorf("expected 1 sub, got %d", len(subs))
	}

	// Update
	updatedName := "Sub Title B"
	found.Name = &updatedName
	if err := repo.UpdateJobTitleSub(ctx, found); err != nil {
		t.Fatalf("UpdateJobTitleSub failed: %v", err)
	}
	updated, _ := repo.FindJobTitleSubByID(ctx, sub.ID)
	if updated.Name == nil || *updated.Name != "Sub Title B" {
		t.Errorf("expected name 'Sub Title B', got %v", updated.Name)
	}

	// Delete
	if err := repo.DeleteJobTitleSub(ctx, sub.ID); err != nil {
		t.Fatalf("DeleteJobTitleSub failed: %v", err)
	}
	_, err = repo.FindJobTitleSubByID(ctx, sub.ID)
	if err == nil {
		t.Fatal("expected error after deleting sub")
	}
}

// =========================================================================
// Job Value CRUD Tests (9.3)
// =========================================================================

func TestRepository_JobValueCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	v := &JobValue{
		Type: "education",
		Level: intPtr(2),
		Sort:  intPtr(1),
	}
	if err := repo.CreateJobValue(ctx, v); err != nil {
		t.Fatalf("CreateJobValue failed: %v", err)
	}

	// Find
	found, err := repo.FindJobValueByID(ctx, v.ID)
	if err != nil {
		t.Fatalf("FindJobValueByID failed: %v", err)
	}
	if found.Type != "education" {
		t.Errorf("expected type 'education', got '%s'", found.Type)
	}

	// Find by type
	byType, err := repo.FindJobValuesByType(ctx, "education")
	if err != nil {
		t.Fatalf("FindJobValuesByType failed: %v", err)
	}
	if len(byType) != 1 {
		t.Errorf("expected 1 value, got %d", len(byType))
	}

	// Update
	found.Type = "experience"
	if err := repo.UpdateJobValue(ctx, found); err != nil {
		t.Fatalf("UpdateJobValue failed: %v", err)
	}
	updated, _ := repo.FindJobValueByID(ctx, v.ID)
	if updated.Type != "experience" {
		t.Errorf("expected type 'experience', got '%s'", updated.Type)
	}

	// Delete
	if err := repo.DeleteJobValue(ctx, v.ID); err != nil {
		t.Fatalf("DeleteJobValue failed: %v", err)
	}
	_, err = repo.FindJobValueByID(ctx, v.ID)
	if err == nil {
		t.Fatal("expected error after deleting job value")
	}
}

func TestRepository_FindAllJobValues_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		v := &JobValue{
			Type: fmt.Sprintf("type_%d", i+1),
			Sort: intPtr(i + 1),
		}
		repo.CreateJobValue(ctx, v)
	}

	values, total, err := repo.FindAllJobValues(ctx, 1, 3)
	if err != nil {
		t.Fatalf("FindAllJobValues failed: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total 4, got %d", total)
	}
	if len(values) != 3 {
		t.Errorf("expected 3 values on page 1, got %d", len(values))
	}
}

// =========================================================================
// Job Objective CRUD Tests (9.4)
// =========================================================================

func TestRepository_JobObjectiveCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()

	// Create
	o := &JobObjective{
		OrganizationID: &orgID,
		Nomenclature:   "Strategic Planning",
		FullCode:       "STRAT-001",
		Objective:      strPtr("Plan and execute strategic initiatives"),
	}
	if err := repo.CreateJobObjective(ctx, o); err != nil {
		t.Fatalf("CreateJobObjective failed: %v", err)
	}

	// Find
	found, err := repo.FindJobObjectiveByID(ctx, o.ID)
	if err != nil {
		t.Fatalf("FindJobObjectiveByID failed: %v", err)
	}
	if found.Nomenclature != "Strategic Planning" {
		t.Errorf("expected 'Strategic Planning', got '%s'", found.Nomenclature)
	}

	// Update
	found.Nomenclature = "Tactical Planning"
	if err := repo.UpdateJobObjective(ctx, found); err != nil {
		t.Fatalf("UpdateJobObjective failed: %v", err)
	}
	updated, _ := repo.FindJobObjectiveByID(ctx, o.ID)
	if updated.Nomenclature != "Tactical Planning" {
		t.Errorf("expected 'Tactical Planning', got '%s'", updated.Nomenclature)
	}

	// Delete
	if err := repo.DeleteJobObjective(ctx, o.ID); err != nil {
		t.Fatalf("DeleteJobObjective failed: %v", err)
	}
	_, err = repo.FindJobObjectiveByID(ctx, o.ID)
	if err == nil {
		t.Fatal("expected error after deleting objective")
	}
}

// =========================================================================
// Job Identification CRUD Tests (9.5)
// =========================================================================

func TestRepository_JobIdentificationCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	gradingID := uuid.New()

	i := &JobIdentification{
		OrganizationID: &orgID,
		Nomenclature:   "Senior Manager",
		FullCode:       "SM-001",
		GradingID:      gradingID,
	}
	if err := repo.CreateJobIdentification(ctx, i); err != nil {
		t.Fatalf("CreateJobIdentification failed: %v", err)
	}

	found, err := repo.FindJobIdentificationByID(ctx, i.ID)
	if err != nil {
		t.Fatalf("FindJobIdentificationByID failed: %v", err)
	}
	if found.Nomenclature != "Senior Manager" {
		t.Errorf("expected 'Senior Manager', got '%s'", found.Nomenclature)
	}

	// Update
	found.Nomenclature = "Junior Manager"
	repo.UpdateJobIdentification(ctx, found)
	updated, _ := repo.FindJobIdentificationByID(ctx, i.ID)
	if updated.Nomenclature != "Junior Manager" {
		t.Errorf("expected 'Junior Manager', got '%s'", updated.Nomenclature)
	}

	repo.DeleteJobIdentification(ctx, i.ID)
	_, err = repo.FindJobIdentificationByID(ctx, i.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Responsibility CRUD Tests (9.6)
// =========================================================================

func TestRepository_JobResponsibilityCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobResponsibility{
		Nomenclature: "Lead Projects",
		FullCode:     "RESP-001",
		MainTask:     strPtr("Lead cross-functional projects"),
		Activities:   strPtr("Planning, execution, reporting"),
		Outputs:      strPtr("Project deliverables"),
	}
	if err := repo.CreateJobResponsibility(ctx, r); err != nil {
		t.Fatalf("CreateJobResponsibility failed: %v", err)
	}

	found, _ := repo.FindJobResponsibilityByID(ctx, r.ID)
	if found.MainTask == nil || *found.MainTask != "Lead cross-functional projects" {
		t.Errorf("expected 'Lead cross-functional projects', got %v", found.MainTask)
	}

	// Update
	newTask := "Lead team initiatives"
	found.MainTask = &newTask
	repo.UpdateJobResponsibility(ctx, found)
	updated, _ := repo.FindJobResponsibilityByID(ctx, r.ID)
	if *updated.MainTask != "Lead team initiatives" {
		t.Errorf("expected 'Lead team initiatives', got '%s'", *updated.MainTask)
	}

	// Pagination test
	foundAll, total, err := repo.FindAllJobResponsibilities(ctx, 1, 20)
	if err != nil {
		t.Fatalf("FindAllJobResponsibilities failed: %v", err)
	}
	if total != 1 || len(foundAll) != 1 {
		t.Errorf("expected 1 responsibility, got total=%d len=%d", total, len(foundAll))
	}

	repo.DeleteJobResponsibility(ctx, r.ID)
	_, err = repo.FindJobResponsibilityByID(ctx, r.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Education Experience CRUD Tests (9.7)
// =========================================================================

func TestRepository_JobEducationExperienceCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	e := &JobEducationExperience{
		Nomenclature: "S1 Minimum",
		FullCode:     "EDU-001",
	}
	if err := repo.CreateJobEducationExperience(ctx, e); err != nil {
		t.Fatalf("CreateJobEducationExperience failed: %v", err)
	}

	found, _ := repo.FindJobEducationExperienceByID(ctx, e.ID)
	if found.Nomenclature != "S1 Minimum" {
		t.Errorf("expected 'S1 Minimum', got '%s'", found.Nomenclature)
	}

	found.Nomenclature = "S2 Minimum"
	repo.UpdateJobEducationExperience(ctx, found)

	// List
	all, total, _ := repo.FindAllJobEducationExperiences(ctx, 1, 20)
	if total != 1 || len(all) != 1 {
		t.Errorf("expected 1, got total=%d len=%d", total, len(all))
	}

	repo.DeleteJobEducationExperience(ctx, e.ID)
	_, err := repo.FindJobEducationExperienceByID(ctx, e.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job HR Authority CRUD Tests (9.8)
// =========================================================================

func TestRepository_JobHRAuthorityCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	a := &JobHRAuthority{
		Nomenclature: "Hire Staff",
		FullCode:     "HR-001",
		Description:  strPtr("Authority to hire staff up to manager level"),
	}
	if err := repo.CreateJobHRAuthority(ctx, a); err != nil {
		t.Fatalf("CreateJobHRAuthority failed: %v", err)
	}

	found, _ := repo.FindJobHRAuthorityByID(ctx, a.ID)
	if found.Nomenclature != "Hire Staff" {
		t.Errorf("expected 'Hire Staff', got '%s'", found.Nomenclature)
	}

	// Update
	desc := "Updated authority description"
	found.Description = &desc
	repo.UpdateJobHRAuthority(ctx, found)

	// List
	all, total, _ := repo.FindAllJobHRAuthorities(ctx, 1, 20)
	if total != 1 || len(all) != 1 {
		t.Errorf("expected 1, got total=%d len=%d", total, len(all))
	}

	repo.DeleteJobHRAuthority(ctx, a.ID)
	_, err := repo.FindJobHRAuthorityByID(ctx, a.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Operational Authority CRUD Tests (9.9)
// =========================================================================

func TestRepository_JobOperationalAuthorityCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	a := &JobOperationalAuthority{
		Nomenclature: "Approve Budget",
		FullCode:     "OP-001",
		Description:  strPtr("Authority to approve operational budget"),
	}
	if err := repo.CreateJobOperationalAuthority(ctx, a); err != nil {
		t.Fatalf("CreateJobOperationalAuthority failed: %v", err)
	}

	found, _ := repo.FindJobOperationalAuthorityByID(ctx, a.ID)
	if found.Nomenclature != "Approve Budget" {
		t.Errorf("expected 'Approve Budget', got '%s'", found.Nomenclature)
	}

	// Update
	found.Nomenclature = "Approve Large Budget"
	repo.UpdateJobOperationalAuthority(ctx, found)

	// List
	all, total, _ := repo.FindAllJobOperationalAuthorities(ctx, 1, 20)
	if total != 1 || len(all) != 1 {
		t.Errorf("expected 1, got total=%d len=%d", total, len(all))
	}

	repo.DeleteJobOperationalAuthority(ctx, a.ID)
	_, err := repo.FindJobOperationalAuthorityByID(ctx, a.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Working Activity CRUD Tests (9.10)
// =========================================================================

func TestRepository_JobWorkingActivityCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	a := &JobWorkingActivity{
		Nomenclature: "Daily Reporting",
		FullCode:     "ACT-001",
	}
	if err := repo.CreateJobWorkingActivity(ctx, a); err != nil {
		t.Fatalf("CreateJobWorkingActivity failed: %v", err)
	}

	found, _ := repo.FindJobWorkingActivityByID(ctx, a.ID)
	if found.Nomenclature != "Daily Reporting" {
		t.Errorf("expected 'Daily Reporting', got '%s'", found.Nomenclature)
	}

	found.Nomenclature = "Weekly Reporting"
	repo.UpdateJobWorkingActivity(ctx, found)
	repo.DeleteJobWorkingActivity(ctx, a.ID)
	_, err := repo.FindJobWorkingActivityByID(ctx, a.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Working Risk CRUD Tests (9.11)
// =========================================================================

func TestRepository_JobWorkingRiskCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobWorkingRisk{
		Nomenclature: "Work Accident",
		FullCode:     "RSK-001",
	}
	if err := repo.CreateJobWorkingRisk(ctx, r); err != nil {
		t.Fatalf("CreateJobWorkingRisk failed: %v", err)
	}

	found, _ := repo.FindJobWorkingRiskByID(ctx, r.ID)
	if found.Nomenclature != "Work Accident" {
		t.Errorf("expected 'Work Accident', got '%s'", found.Nomenclature)
	}

	// Update + Delete
	found.Nomenclature = "Health Hazard"
	repo.UpdateJobWorkingRisk(ctx, found)
	repo.DeleteJobWorkingRisk(ctx, r.ID)
	_, err := repo.FindJobWorkingRiskByID(ctx, r.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Relationship CRUD Tests (9.12)
// =========================================================================

func TestRepository_JobRelationshipCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobRelationship{
		Nomenclature: "Internal Coordination",
		FullCode:     "REL-001",
	}
	if err := repo.CreateJobRelationship(ctx, r); err != nil {
		t.Fatalf("CreateJobRelationship failed: %v", err)
	}

	found, _ := repo.FindJobRelationshipByID(ctx, r.ID)
	if found.Nomenclature != "Internal Coordination" {
		t.Errorf("expected 'Internal Coordination', got '%s'", found.Nomenclature)
	}

	// Full CRUD cycle
	found.Nomenclature = "External Coordination"
	repo.UpdateJobRelationship(ctx, found)
	updated, _ := repo.FindJobRelationshipByID(ctx, r.ID)
	if updated.Nomenclature != "External Coordination" {
		t.Errorf("expected 'External Coordination', got '%s'", updated.Nomenclature)
	}

	all, total, _ := repo.FindAllJobRelationships(ctx, 1, 20)
	if total != 1 || len(all) != 1 {
		t.Errorf("expected 1, got total=%d len=%d", total, len(all))
	}

	repo.DeleteJobRelationship(ctx, r.ID)
	_, err := repo.FindJobRelationshipByID(ctx, r.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Subordinate Control CRUD Tests (9.13)
// =========================================================================

func TestRepository_JobSubordinateControlCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	c := &JobSubordinateControl{
		Nomenclature: "Direct Reports",
		FullCode:     "SUB-001",
	}
	if err := repo.CreateJobSubordinateControl(ctx, c); err != nil {
		t.Fatalf("CreateJobSubordinateControl failed: %v", err)
	}

	found, _ := repo.FindJobSubordinateControlByID(ctx, c.ID)
	if found.Nomenclature != "Direct Reports" {
		t.Errorf("expected 'Direct Reports', got '%s'", found.Nomenclature)
	}

	repo.DeleteJobSubordinateControl(ctx, c.ID)
	_, err := repo.FindJobSubordinateControlByID(ctx, c.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Asset CRUD Tests (9.14)
// =========================================================================

func TestRepository_JobAssetCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	a := &JobAsset{
		Nomenclature: "Laptop & Equipment",
		FullCode:     "AST-001",
	}
	if err := repo.CreateJobAsset(ctx, a); err != nil {
		t.Fatalf("CreateJobAsset failed: %v", err)
	}

	found, _ := repo.FindJobAssetByID(ctx, a.ID)
	if found.Nomenclature != "Laptop & Equipment" {
		t.Errorf("expected 'Laptop & Equipment', got '%s'", found.Nomenclature)
	}

	repo.DeleteJobAsset(ctx, a.ID)
	_, err := repo.FindJobAssetByID(ctx, a.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Financial CRUD Tests (9.15)
// =========================================================================

func TestRepository_JobFinancialCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	f := &JobFinancial{
		Nomenclature: "Budget Authority",
		FullCode:     "FIN-001",
		IsAuthorized: true,
	}
	if err := repo.CreateJobFinancial(ctx, f); err != nil {
		t.Fatalf("CreateJobFinancial failed: %v", err)
	}

	found, _ := repo.FindJobFinancialByID(ctx, f.ID)
	if found.Nomenclature != "Budget Authority" {
		t.Errorf("expected 'Budget Authority', got '%s'", found.Nomenclature)
	}
	if found.IsAuthorized != true {
		t.Errorf("expected IsAuthorized true, got %v", found.IsAuthorized)
	}

	// Update
	found.IsAuthorized = false
	repo.UpdateJobFinancial(ctx, found)
	updated, _ := repo.FindJobFinancialByID(ctx, f.ID)
	if updated.IsAuthorized != false {
		t.Errorf("expected IsAuthorized false, got %v", updated.IsAuthorized)
	}

	repo.DeleteJobFinancial(ctx, f.ID)
	_, err := repo.FindJobFinancialByID(ctx, f.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Potency Competency CRUD Tests (9.16)
// =========================================================================

func TestRepository_JobPotencyCompetencyCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	c := &JobPotencyCompetency{
		Weight: float64Ptr(80.5),
	}
	if err := repo.CreateJobPotencyCompetency(ctx, c); err != nil {
		t.Fatalf("CreateJobPotencyCompetency failed: %v", err)
	}

	found, _ := repo.FindJobPotencyCompetencyByID(ctx, c.ID)
	if found.Weight == nil || *found.Weight != 80.5 {
		t.Errorf("expected weight 80.5, got %v", found.Weight)
	}

	// Update
	newWeight := 90.0
	found.Weight = &newWeight
	repo.UpdateJobPotencyCompetency(ctx, found)
	updated, _ := repo.FindJobPotencyCompetencyByID(ctx, c.ID)
	if updated.Weight == nil || *updated.Weight != 90.0 {
		t.Errorf("expected weight 90.0, got %v", updated.Weight)
	}

	repo.DeleteJobPotencyCompetency(ctx, c.ID)
	_, err := repo.FindJobPotencyCompetencyByID(ctx, c.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

// =========================================================================
// Job Score Tests (9.17)
// =========================================================================

func TestRepository_UpsertJobScore_Create(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	s := &JobScore{
		OrganizationID:            &orgID,
		JobValueWithFinancial:     1000,
		JobValueWithoutFinancial:  800,
		HasFinancialAuthority:     true,
	}
	if err := repo.UpsertJobScore(ctx, s); err != nil {
		t.Fatalf("UpsertJobScore (create) failed: %v", err)
	}
	if s.ID == uuid.Nil {
		t.Error("expected job score ID to be generated")
	}
}

func TestRepository_UpsertJobScore_Update(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	s := &JobScore{
		OrganizationID:        &orgID,
		JobValueWithoutFinancial: 500,
	}
	repo.UpsertJobScore(ctx, s)
	originalID := s.ID

	// Upsert again with same orgID -> should update
	s.JobValueWithFinancial = 1500
	if err := repo.UpsertJobScore(ctx, s); err != nil {
		t.Fatalf("UpsertJobScore (update) failed: %v", err)
	}
	if s.ID != originalID {
		t.Errorf("expected same ID after upsert, got %s vs %s", s.ID, originalID)
	}

	found, err := repo.FindJobScoreByOrganizationID(ctx, orgID)
	if err != nil {
		t.Fatalf("FindJobScoreByOrganizationID failed: %v", err)
	}
	if found.JobValueWithFinancial != 1500 {
		t.Errorf("expected JobValueWithFinancial 1500, got %d", found.JobValueWithFinancial)
	}
}

func TestRepository_FindJobScoreByOrganizationID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindJobScoreByOrganizationID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent job score")
	}
}

// =========================================================================
// Job Competency Group CRUD Tests (9.18)
// =========================================================================

func TestRepository_JobCompetencyGroupCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()

	// Create
	g := &JobCompetencyGroup{
		OrganizationID: &orgID,
		Category:       "technical",
		Weight:         70.0,
	}
	if err := repo.CreateJobCompetencyGroup(ctx, g); err != nil {
		t.Fatalf("CreateJobCompetencyGroup failed: %v", err)
	}

	// Find
	found, err := repo.FindJobCompetencyGroupByID(ctx, g.ID)
	if err != nil {
		t.Fatalf("FindJobCompetencyGroupByID failed: %v", err)
	}
	if found.Category != "technical" {
		t.Errorf("expected category 'technical', got '%s'", found.Category)
	}

	// Find by organization
	groups, err := repo.FindJobCompetencyGroupsByOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("FindJobCompetencyGroupsByOrganization failed: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}

	// Update
	found.Weight = 85.0
	if err := repo.UpdateJobCompetencyGroup(ctx, found); err != nil {
		t.Fatalf("UpdateJobCompetencyGroup failed: %v", err)
	}

	// Delete
	if err := repo.DeleteJobCompetencyGroup(ctx, g.ID); err != nil {
		t.Fatalf("DeleteJobCompetencyGroup failed: %v", err)
	}
	_, err = repo.FindJobCompetencyGroupByID(ctx, g.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
