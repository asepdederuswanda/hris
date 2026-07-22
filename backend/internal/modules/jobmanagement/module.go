package jobmanagement

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Job Management"
const ModuleSlug = "jobmanagement"
const ModuleVersion = "1.0.0"

// TenantDBFunc adalah fungsi resolver untuk mendapatkan koneksi
// tenant database berdasarkan context (company_id dari JWT claims).
type TenantDBFunc func(ctx context.Context) (*gorm.DB, error)

// NewTenantDBResolver membuat resolver yang mengambil company_id dari context.
func NewTenantDBResolver(dbManager *database.Manager) TenantDBFunc {
	return func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
}

func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	return &jmModule{
		handler: handler,
		logger:  logger,
	}
}

type jmModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *jmModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage job analysis including titles, values, objectives, responsibilities, competencies, and scoring",
		IsCore:      true,
		DependsOn:   []string{"organization"},
		Permissions: []string{
			"jobmanagement.view",
			"jobmanagement.create",
			"jobmanagement.update",
			"jobmanagement.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Job Management",
				Icon:  "briefcase",
				Route: "/admin/job-management",
				Children: []module.Menu{
					{Name: "Job Titles", Icon: "tag", Route: "/admin/job-management/titles"},
					{Name: "Job Values", Icon: "sliders", Route: "/admin/job-management/values"},
					{Name: "Job Objectives", Icon: "target", Route: "/admin/job-management/objectives"},
					{Name: "Job Identifications", Icon: "fingerprint", Route: "/admin/job-management/identifications"},
					{Name: "Responsibilities", Icon: "check-square", Route: "/admin/job-management/responsibilities"},
					{Name: "Authorities", Icon: "shield", Route: "/admin/job-management/authorities"},
					{Name: "Working Conditions", Icon: "activity", Route: "/admin/job-management/working-conditions"},
					{Name: "Job Scores", Icon: "bar-chart", Route: "/admin/job-management/scores"},
					{Name: "Competency Mapping", Icon: "grid", Route: "/admin/job-management/competencies"},
				},
			},
		},
	}
}

func (m *jmModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *jmModule) Migrate(db *gorm.DB) error {
	// SQL migration already handled by 009_job_management.sql migrator
	return nil
}

func (m *jmModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *jmModule) Permissions() []string {
	return m.Info().Permissions
}
