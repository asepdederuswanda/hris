package employeemovement

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Employee Movement & Career Management"
const ModuleSlug = "employeemovement"
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

	return &empMovementModule{
		handler: handler,
		logger:  logger,
	}
}

type empMovementModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *empMovementModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage employee career movements (promotions, demotions, mutations, contract extensions, retirements, offboarding) and employment contracts (PKWT/PKWTT)",
		IsCore:      true,
		DependsOn:   []string{"employee", "organization"},
		Permissions: []string{
			"employeemovement.view",
			"employeemovement.create",
			"employeemovement.update",
			"employeemovement.delete",
			"employeemovement.approve",
			"employeemovement.execute",
		},
		Menus: []module.Menu{
			{
				Name:  "Career",
				Icon:  "trending-up",
				Route: "/admin/career",
				Children: []module.Menu{
					{Name: "Movements", Icon: "shuffle", Route: "/admin/career/movements"},
					{Name: "Contracts", Icon: "file-text", Route: "/admin/career/contracts"},
				},
			},
		},
	}
}

func (m *empMovementModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *empMovementModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&EmployeeMovement{},
		&EmployeeContract{},
	)
}

func (m *empMovementModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *empMovementModule) Permissions() []string {
	return m.Info().Permissions
}
