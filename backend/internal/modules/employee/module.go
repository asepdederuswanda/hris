package employee

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Employee Management"
const ModuleSlug = "employee"
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

	return &empModule{
		handler: handler,
		logger:  logger,
	}
}

type empModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *empModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage employee data including addresses, families, educations, documents, and employment history",
		IsCore:      true,
		DependsOn:   []string{"organization"},
		Permissions: []string{
			"employee.view",
			"employee.create",
			"employee.update",
			"employee.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Employee",
				Icon:  "users",
				Route: "/admin/employees",
				Children: []module.Menu{
					{Name: "Employee List", Icon: "list", Route: "/admin/employees"},
					{Name: "Employee Create", Icon: "user-plus", Route: "/admin/employees/create"},
				},
			},
		},
	}
}

func (m *empModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *empModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Employee{},
		&EmployeeAddress{},
		&EmergencyContact{},
		&EmployeeFamily{},
		&EmployeeEducation{},
		&EmployeeExperience{},
		&EmployeeDocument{},
		&EmployeeInsurance{},
		&Employment{},
	)
}

func (m *empModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *empModule) Permissions() []string {
	return m.Info().Permissions
}
