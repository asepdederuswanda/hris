package competency

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Competency Management"
const ModuleSlug = "competency"
const ModuleVersion = "1.0.0"

type TenantDBFunc func(ctx context.Context) (*gorm.DB, error)

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

	return &compModule{
		handler: handler,
		logger:  logger,
	}
}

type compModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *compModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage competency dictionaries, values, events, scoring, and gap analysis",
		IsCore:      true,
		DependsOn:   []string{"organization", "employee"},
		Permissions: []string{
			"competency.view",
			"competency.create",
			"competency.update",
			"competency.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Competency",
				Icon:  "award",
				Route: "/admin/competency",
				Children: []module.Menu{
					{Name: "Competencies", Icon: "book", Route: "/admin/competency/competencies"},
					{Name: "Values", Icon: "sliders", Route: "/admin/competency/values"},
					{Name: "Events", Icon: "calendar", Route: "/admin/competency/events"},
					{Name: "Scores", Icon: "bar-chart", Route: "/admin/competency/scores"},
				},
			},
		},
	}
}

func (m *compModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *compModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Competency{},
		&CompetenceValue{},
		&CompetencyValue{},
		&CompetencyEvent{},
		&CompetencyEventTarget{},
		&CompetencyScore{},
		&CompetencyScoreDetail{},
	)
}

func (m *compModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *compModule) Permissions() []string {
	return m.Info().Permissions
}
