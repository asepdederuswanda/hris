package modulemgmt

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const (
	ModuleName    = "Module Management"
	ModuleSlug    = "module-management"
	ModuleVersion = "1.0.0"
)

// NewModule membuat instance baru Module Management Module.
func NewModule(dbManager *database.Manager, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, dbManager, logger)
	handler := NewHandler(svc)

	return &modulemgmtModule{
		handler: handler,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type modulemgmtModule struct {
	handler *Handler
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

func (m *modulemgmtModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage platform modules and their activation for companies",
		IsCore:      true,
		DependsOn:   []string{},
		Permissions: []string{
			"module.view",
			"module.create",
			"module.update",
			"module.activate",
			"module.deactivate",
		},
		Menus: []module.Menu{
			{
				Name:  "Modules",
				Icon:  "puzzle-piece",
				Route: "/admin/platform/modules",
			},
		},
	}
}

func (m *modulemgmtModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler, m.authMW, m.rbacMW)
}

func (m *modulemgmtModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PlatformModule{}, &CompanyModule{})
}

func (m *modulemgmtModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *modulemgmtModule) Permissions() []string {
	return m.Info().Permissions
}
