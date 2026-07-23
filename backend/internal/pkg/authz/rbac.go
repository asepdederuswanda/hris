// Package authz menyediakan RBAC (Role-Based Access Control)
// untuk platform dan tenant endpoints menggunakan permission checking
// berbasis role dari JWT claims dengan role hierarchy.
//
// Role Hierarchy:
//   super_admin → akses penuh ke semua resource (platform + tenant)
//   company_admin → platform view-only + full tenant management
//   manager → tenant-level view/create/update (tanpa delete)
//   employee → tenant-level view-only
//
// Policy format: resource.action
// Contoh: company.create, user.view, module.activate, employee.view
//
// Inheritance: Jika suatu role tidak memiliki policy untuk resource/action tertentu,
// maka akan dicek ke parent role hingga ke super_admin atau sampai ketemu.
//
// Database-backed RBAC:
// Enforcer dapat diisi dari database platform via LoadFromDB().
// Saat startup, default policies di-seed ke database, dan enforcer
// memuat dari database. Admin dapat menambah/mengubah role & permission
// via API dan memanggil Reload() untuk sinkronisasi tanpa restart.
package authz

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role constants sesuai JWT claims.
type Role string

const (
	RoleSuperAdmin   Role = "super_admin"
	RoleCompanyAdmin Role = "company_admin"
	RoleManager      Role = "manager"
	RoleEmployee     Role = "employee"
)

// Decision hasil dari pemeriksaan permission.
type Decision string

const (
	DecisionAllow Decision = "allow"
	DecisionDeny  Decision = "deny"
)

// defaultPerm digunakan untuk mendefinisikan resource default dan actions-nya
// saat seeding data RBAC ke database.
type defaultPerm struct {
	resource string
	actions  []string
}

// Enforcer adalah RBAC enforcer yang memeriksa permission
// berdasarkan role dan resource-action yang diminta.
// Mendukung role hierarchy inheritance.
type Enforcer struct {
	// policies menyimpan aturan RBAC:
	// map[role]map[resource]allowedActions atau "*" untuk semua
	policies map[Role]map[string]string

	// hierarchy menyimpan parent dari setiap role untuk inheritance
	hierarchy map[Role]Role

	// db untuk database-backed RBAC (opsional)
	db *gorm.DB
}

// NewEnforcer membuat Enforcer baru dengan policy default (hardcoded).
// Untuk database-backed RBAC, gunakan NewEnforcerFromDB().
func NewEnforcer() *Enforcer {
	e := &Enforcer{
		policies:  make(map[Role]map[string]string),
		hierarchy: make(map[Role]Role),
	}
	e.loadDefaultHierarchy()
	e.loadDefaultPolicies()
	return e
}

// NewEnforcerFromDB membuat Enforcer dengan data dari database platform.
// Jika database kosong, akan di-seed dengan default policies.
// Enforcer ini support Reload() untuk sinkronasi tanpa restart.
func NewEnforcerFromDB(db *gorm.DB) (*Enforcer, error) {
	e := &Enforcer{
		policies:  make(map[Role]map[string]string),
		hierarchy: make(map[Role]Role),
		db:        db,
	}

	// Seed default roles & permissions jika database kosong
	if err := e.seedDefaults(db); err != nil {
		return nil, fmt.Errorf("failed to seed RBAC defaults: %w", err)
	}

	// Load dari database
	if err := e.loadFromDB(db); err != nil {
		return nil, fmt.Errorf("failed to load RBAC from database: %w", err)
	}

	return e, nil
}

// loadDefaultHierarchy memuat default role hierarchy untuk non-DB enforcer.
func (e *Enforcer) loadDefaultHierarchy() {
	e.hierarchy[RoleCompanyAdmin] = RoleSuperAdmin
	e.hierarchy[RoleManager] = RoleCompanyAdmin
	e.hierarchy[RoleEmployee] = RoleManager
}

// loadDefaultPolicies memuat default policies untuk non-DB enforcer.
func (e *Enforcer) loadDefaultPolicies() {
	e.AddPolicy(RoleSuperAdmin, "*", "*")

	// Company Admin: platform view-only + tenant full
	e.AddPolicy(RoleCompanyAdmin, "company", "view")
	e.AddPolicy(RoleCompanyAdmin, "user", "view")
	e.AddPolicy(RoleCompanyAdmin, "license", "view")
	e.AddPolicy(RoleCompanyAdmin, "organization", "*")
	e.AddPolicy(RoleCompanyAdmin, "employee", "*")
	e.AddPolicy(RoleCompanyAdmin, "attendance", "*")
	e.AddPolicy(RoleCompanyAdmin, "leave", "*")
	e.AddPolicy(RoleCompanyAdmin, "payroll", "*")
	e.AddPolicy(RoleCompanyAdmin, "competency", "*")
	e.AddPolicy(RoleCompanyAdmin, "jobmanagement", "*")
	e.AddPolicy(RoleCompanyAdmin, "employeemovement", "*")
	e.AddPolicy(RoleCompanyAdmin, "approval", "*")

	// Manager: view/create/update (no delete)
	e.AddPolicy(RoleManager, "organization", "view,create,update")
	e.AddPolicy(RoleManager, "employee", "view,create,update")
	e.AddPolicy(RoleManager, "attendance", "view")
	e.AddPolicy(RoleManager, "leave", "view,create")
	e.AddPolicy(RoleManager, "competency", "view,create,update")
	e.AddPolicy(RoleManager, "jobmanagement", "view,create,update")
	e.AddPolicy(RoleManager, "employeemovement", "view,create,update")
	e.AddPolicy(RoleManager, "approval", "view,create,update")

	// Employee: view-only
	e.AddPolicy(RoleEmployee, "organization", "view")
	e.AddPolicy(RoleEmployee, "employee", "view")
	e.AddPolicy(RoleEmployee, "attendance", "view")
	e.AddPolicy(RoleEmployee, "leave", "view")
	e.AddPolicy(RoleEmployee, "competency", "view")
	e.AddPolicy(RoleEmployee, "employeemovement", "view")
	e.AddPolicy(RoleEmployee, "payroll", "view")
}

// Reload menyegarkan semua policies dari database.
// Dipanggil setelah ada perubahan role/permission via API.
func (e *Enforcer) Reload() error {
	if e.db == nil {
		return fmt.Errorf("Enforcer not configured with database")
	}

	// Reset
	e.policies = make(map[Role]map[string]string)
	e.hierarchy = make(map[Role]Role)

	// Load ulang dari database
	return e.loadFromDB(e.db)
}

// loadFromDB memuat role hierarchy dan policies dari database.
func (e *Enforcer) loadFromDB(db *gorm.DB) error {
	// Load semua roles
	var roles []RbacRole
	if err := db.Preload("Permissions").Find(&roles).Error; err != nil {
		return fmt.Errorf("failed to load roles: %w", err)
	}

	if len(roles) == 0 {
		return fmt.Errorf("no roles found in database after seeding")
	}

	// Build slug → Role mapping and hierarchy
	roleSlugToRole := make(map[string]Role)
	roleByID := make(map[string]RbacRole)

	for _, r := range roles {
		roleSlugToRole[r.Slug] = Role(r.Slug)
		roleByID[r.ID.String()] = r
	}

	// Load hierarchy dari parent_id
	for _, r := range roles {
		if r.ParentID != nil {
			parentID := r.ParentID.String()
			if parent, ok := roleByID[parentID]; ok {
				e.hierarchy[Role(r.Slug)] = Role(parent.Slug)
			}
		}
	}

	// Load semua permissions
	var permissions []RbacPermission
	if err := db.Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to load permissions: %w", err)
	}

	// Build permission lookup: id → (resource, action)
	permMap := make(map[string]struct{ resource, action string })
	for _, p := range permissions {
		permMap[p.ID.String()] = struct{ resource, action string }{
			resource: p.Resource,
			action:   p.Action,
		}
	}

	// Build policies dari role_permissions
	for _, role := range roles {
		roleSlug := Role(role.Slug)
		if _, exists := e.policies[roleSlug]; !exists {
			e.policies[roleSlug] = make(map[string]string)
		}

		for _, rp := range role.Permissions {
			perm, ok := permMap[rp.PermissionID.String()]
			if !ok {
				continue
			}

			resource := perm.resource
			action := perm.action

			// Wildcard handling
			if resource == "*" {
				e.policies[roleSlug]["*"] = "*"
				continue
			}

			// Accumulate actions for same resource
			if existing, ok := e.policies[roleSlug][resource]; ok {
				if existing == "*" {
					continue // already wildcard
				}
				if action == "*" {
					e.policies[roleSlug][resource] = "*"
				} else if !strings.Contains(existing, action) {
					e.policies[roleSlug][resource] = existing + "," + action
				}
			} else {
				e.policies[roleSlug][resource] = action
			}
		}
	}

	return nil
}

// seedDefaults memasukkan role dan permission default ke database
// jika tabel masih kosong. Hanya berjalan sekali saat pertama kali.
func (e *Enforcer) seedDefaults(db *gorm.DB) error {
	var count int64

	// Cek apakah sudah ada data
	if err := db.Model(&RbacRole{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // sudah ada data, skip seed
	}

	// Buat default roles
	superAdminUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	companyAdminUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	managerUUID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	employeeUUID := uuid.MustParse("00000000-0000-0000-0000-000000000004")

	roles := []RbacRole{
		{ID: superAdminUUID, Name: "Super Admin", Slug: string(RoleSuperAdmin), Description: strPtr("Full access to all resources"), IsSystem: true},
		{ID: companyAdminUUID, Name: "Company Admin", Slug: string(RoleCompanyAdmin), Description: strPtr("Platform view + tenant full access"), ParentID: &superAdminUUID, IsSystem: true},
		{ID: managerUUID, Name: "Manager", Slug: string(RoleManager), Description: strPtr("Tenant view/create/update"), ParentID: &companyAdminUUID, IsSystem: true},
		{ID: employeeUUID, Name: "Employee", Slug: string(RoleEmployee), Description: strPtr("Tenant view-only"), ParentID: &managerUUID, IsSystem: true},
	}

	for _, r := range roles {
		if err := db.Create(&r).Error; err != nil {
			return fmt.Errorf("failed to seed role %s: %w", r.Slug, err)
		}
	}

	// Buat permissions dari hardcoded defaults
	allResources := []defaultPerm{
		// Platform resources
		{"company", []string{"view", "create", "update", "delete", "suspend", "activate", "terminate", "backup", "restore"}},
		{"user", []string{"view", "create", "update"}},
		{"module", []string{"view", "create", "update", "activate", "deactivate"}},
		{"license", []string{"view", "create", "update"}},
		{"monitoring", []string{"view"}},
		// Tenant resources
		{"organization", []string{"view", "create", "update", "delete"}},
		{"employee", []string{"view", "create", "update", "delete"}},
		{"attendance", []string{"view", "create", "update", "delete"}},
		{"leave", []string{"view", "create", "update", "delete"}},
		{"payroll", []string{"view", "create", "update", "delete"}},
		{"competency", []string{"view", "create", "update", "delete"}},
		{"jobmanagement", []string{"view", "create", "update", "delete"}},
		{"employeemovement", []string{"view", "create", "update", "delete"}},
		{"approval", []string{"view", "create", "update", "delete"}},
	}

	// Simpan permission IDs untuk mapping role_permissions nanti
	permIDs := make(map[string]map[string]uuid.UUID) // resource → action → id

	for _, r := range allResources {
		if _, ok := permIDs[r.resource]; !ok {
			permIDs[r.resource] = make(map[string]uuid.UUID)
		}
		for _, action := range r.actions {
			perm := RbacPermission{
				Resource:    r.resource,
				Action:      action,
				Description: strPtr(fmt.Sprintf("Can %s %s", action, r.resource)),
				IsSystem:    true,
			}
			if err := db.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to seed permission %s.%s: %w", r.resource, action, err)
			}
			permIDs[r.resource][action] = perm.ID
		}
	}

	// Helper untuk menambah role_permissions
	addPerms := func(roleID uuid.UUID, perms map[string][]string) error {
		for resource, actions := range perms {
			if resource == "*" {
				// Wildcard: tambahkan semua permission
				for _, r := range allResources {
					for _, action := range r.actions {
						if id, ok := permIDs[r.resource][action]; ok {
							rp := RbacRolePermission{RoleID: roleID, PermissionID: id}
							if err := db.Create(&rp).Error; err != nil {
								return err
							}
						}
					}
				}
			} else {
				for _, action := range actions {
					if action == "*" {
						// Semua action untuk resource ini
						for _, a := range allPermActions(resource, allResources) {
							if id, ok := permIDs[resource][a]; ok {
								rp := RbacRolePermission{RoleID: roleID, PermissionID: id}
								if err := db.Create(&rp).Error; err != nil {
									return err
								}
							}
						}
					} else {
						if id, ok := permIDs[resource][action]; ok {
							rp := RbacRolePermission{RoleID: roleID, PermissionID: id}
							if err := db.Create(&rp).Error; err != nil {
								return err
							}
						}
					}
				}
			}
		}
		return nil
	}

	// Super Admin: semua permission
	if err := addPerms(superAdminUUID, map[string][]string{"*": {"*"}}); err != nil {
		return err
	}

	// Company Admin: platform view-only + tenant full
	companyAdminPerms := map[string][]string{
		"company":          {"view"},
		"user":             {"view"},
		"license":          {"view"},
		"organization":     {"*"},
		"employee":         {"*"},
		"attendance":       {"*"},
		"leave":            {"*"},
		"payroll":          {"*"},
		"competency":       {"*"},
		"jobmanagement":    {"*"},
		"employeemovement": {"*"},
		"approval":         {"*"},
	}
	if err := addPerms(companyAdminUUID, companyAdminPerms); err != nil {
		return err
	}

	// Manager: view/create/update (no delete)
	managerPerms := map[string][]string{
		"organization":     {"view", "create", "update"},
		"employee":         {"view", "create", "update"},
		"attendance":       {"view"},
		"leave":            {"view", "create"},
		"competency":       {"view", "create", "update"},
		"jobmanagement":    {"view", "create", "update"},
		"employeemovement": {"view", "create", "update"},
		"approval":         {"view", "create", "update"},
	}
	if err := addPerms(managerUUID, managerPerms); err != nil {
		return err
	}

	// Employee: view-only
	employeePerms := map[string][]string{
		"organization":     {"view"},
		"employee":         {"view"},
		"attendance":       {"view"},
		"leave":            {"view"},
		"payroll":          {"view"},
		"competency":       {"view"},
		"employeemovement": {"view"},
	}
	if err := addPerms(employeeUUID, employeePerms); err != nil {
		return err
	}

	return nil
}

// allPermActions mengembalikan semua action yang terdaftar untuk suatu resource.
func allPermActions(resource string, resources []defaultPerm) []string {
	for _, r := range resources {
		if r.resource == resource {
			return r.actions
		}
	}
	return nil
}

// strPtr helper.
func strPtr(s string) *string {
	return &s
}

// Check memeriksa apakah role tertentu diizinkan mengakses
// resource dengan action tertentu. Mendukung role hierarchy:
// jika role tidak memiliki policy, akan dicek ke parent role.
//
// Parameters:
//   - role: role user (dari JWT claims)
//   - resource: resource yang akan diakses (contoh: "company", "employee", "organization")
//   - action: action yang dilakukan (contoh: "view", "create", "update", "delete")
//
// Returns:
//   - DecisionAllow jika diizinkan
//   - DecisionDeny jika tidak diizinkan
func (e *Enforcer) Check(role, resource, action string) Decision {
	currentRole := Role(role)

	// Traverse role hierarchy: current → parent → grandparent → ...
	for currentRole != "" {
		rolePolicies, exists := e.policies[currentRole]
		if !exists {
			// Role tidak dikenal, lanjut ke parent
			currentRole = e.hierarchy[currentRole]
			continue
		}

		// Check wildcard resource (super admin)
		if allowedActions, ok := rolePolicies["*"]; ok {
			if allowedActions == "*" {
				return DecisionAllow
			}
		}

		// Check resource-specific policies
		allowedActions, exists := rolePolicies[resource]
		if !exists {
			// Resource tidak didefinisikan di role ini, cek parent
			currentRole = e.hierarchy[currentRole]
			continue
		}

		// Check wildcard action
		if allowedActions == "*" {
			return DecisionAllow
		}

		// Check specific action
		actions := strings.Split(allowedActions, ",")
		for _, a := range actions {
			if strings.TrimSpace(a) == action {
				return DecisionAllow
			}
		}

		// Resource ada di policy role ini tapi action tidak cocok.
		// Karena role ini sudah memiliki explicit policy untuk resource tsb
		// (dengan daftar action yang terbatas), missing action = intentional deny.
		// Jangan traverse ke parent — langsung return Deny.
		return DecisionDeny
	}

	return DecisionDeny
}

// AddPolicy menambahkan aturan RBAC baru secara dinamis.
func (e *Enforcer) AddPolicy(role Role, resource, actions string) {
	if _, exists := e.policies[role]; !exists {
		e.policies[role] = make(map[string]string)
	}
	e.policies[role][resource] = actions
}

// MustCheck sama seperti Check tetapi mengembalikan error
// jika tidak diizinkan, cocok untuk digunakan di middleware.
func (e *Enforcer) MustCheck(role, resource, action string) error {
	if e.Check(role, resource, action) == DecisionDeny {
		return fmt.Errorf("RBAC: role '%s' not allowed to %s %s",
			role, action, resource)
	}
	return nil
}

// ResourceFromPath mengekstrak nama resource dari path endpoint.
// Mendukung platform dan tenant paths.
//
// Contoh:
//   - /api/v1/platform/companies -> "company"
//   - /api/v1/platform/users -> "user"
//   - /api/v1/tenant/organizations -> "organization"
//   - /api/v1/tenant/employees -> "employee"
func ResourceFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Path pattern: api/v1/{domain}/{resource}/...
	// domain = "platform" atau "tenant"
	for i, part := range parts {
		if (part == "platform" || part == "tenant") && i+1 < len(parts) {
			resource := parts[i+1]
			return singularize(resource)
		}
	}
	return ""
}

// ActionFromMethod mengonversi HTTP method ke action name.
func ActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "view"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return strings.ToLower(method)
	}
}

// singularize mengubah plural ke singular secara sederhana.
func singularize(s string) string {
	// Handle irregular plurals
	irregular := map[string]string{
		"companies":                  "company",
		"licenses":                   "license",
		"modules":                    "module",
		"users":                      "user",
		"monitoring":                 "monitoring",
		"tenants":                    "tenant",
		"organizations":              "organization",
		"employees":                  "employee",
		"attendances":                "attendance",
		"competencies":               "competency",
		"job-management":             "jobmanagement",
		"employee-movements":         "employeemovement",
	}
	if singular, ok := irregular[s]; ok {
		return singular
	}
	// Simple rules: hapus "ies" → "y"
	if strings.HasSuffix(s, "ies") && len(s) > 3 {
		return s[:len(s)-3] + "y"
	}
	// Simple rules: hapus "s" di akhir
	if strings.HasSuffix(s, "s") && len(s) > 1 {
		return s[:len(s)-1]
	}
	return s
}
