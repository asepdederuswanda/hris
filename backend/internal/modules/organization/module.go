package organization

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Organization Management"
const ModuleSlug = "organization"
const ModuleVersion = "1.0.0"

// TenantDBFunc adalah fungsi resolver untuk mendapatkan koneksi
// tenant database berdasarkan context (company_id dari JWT claims).
type TenantDBFunc func(ctx context.Context) (*gorm.DB, error)

// NewTenantDBResolver membuat resolver yang mengambil company_id dari context.
// Dipanggil oleh middleware TenantRequired yang sudah memastikan company_id ada.
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

	return &orgModule{
		handler: handler,
		logger:  logger,
	}
}

type orgModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *orgModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage organizational structure, hierarchy, and levels",
		IsCore:      true,
		DependsOn:   []string{},
		Permissions: []string{
			"organization.view",
			"organization.create",
			"organization.update",
			"organization.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Organization",
				Icon:  "building",
				Route: "/admin/organizations",
				Children: []module.Menu{
					{Name: "Organization Tree", Icon: "hierarchy", Route: "/admin/organizations"},
					{Name: "Zones", Icon: "map-pin", Route: "/admin/zones"},
					{Name: "Job Families", Icon: "briefcase", Route: "/admin/job-families"},
					{Name: "Positions", Icon: "users", Route: "/admin/positions"},
				},
			},
		},
	}
}

func (m *orgModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *orgModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Organization{})
}

func (m *orgModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *orgModule) Permissions() []string {
	return m.Info().Permissions
}
