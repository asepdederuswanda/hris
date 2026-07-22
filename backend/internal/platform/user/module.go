package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const (
	ModuleName    = "Platform Users & Auth"
	ModuleSlug    = "platform-users"
	ModuleVersion = "1.0.0"
)

// NewModule membuat instance baru Platform User Module.
func NewModule(dbManager *database.Manager, authManager *auth.Manager, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, authManager, logger)
	handler := NewHandler(svc)

	return &userModule{
		handler: handler,
		service: svc,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type userModule struct {
	handler *Handler
	service *Service
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

func (m *userModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Platform user authentication, authorization, and user management",
		IsCore:      true,
		DependsOn:   []string{},
		Permissions: []string{
			"user.view",
			"user.create",
			"user.update",
			"user.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Platform Users",
				Icon:  "users-cog",
				Route: "/admin/platform/users",
			},
		},
	}
}

func (m *userModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler, m.authMW, m.rbacMW)
}

func (m *userModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PlatformUser{})
}

func (m *userModule) Seed(db *gorm.DB) error {
	return m.service.EnsureSeed(db)
}

func (m *userModule) Permissions() []string {
	return m.Info().Permissions
}
