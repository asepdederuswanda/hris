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
package authz

import (
	"fmt"
	"strings"
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

// Enforcer adalah RBAC enforcer yang memeriksa permission
// berdasarkan role dan resource-action yang diminta.
// Mendukung role hierarchy inheritance.
type Enforcer struct {
	// policies menyimpan aturan RBAC:
	// map[role]map[resource]allowedActions atau "*" untuk semua
	policies map[Role]map[string]string

	// hierarchy menyimpan parent dari setiap role untuk inheritance
	hierarchy map[Role]Role
}

// NewEnforcer membuat Enforcer baru dengan policy default.
func NewEnforcer() *Enforcer {
	e := &Enforcer{
		policies:  make(map[Role]map[string]string),
		hierarchy: make(map[Role]Role),
	}
	e.loadDefaultHierarchy()
	e.loadDefaultPolicies()
	return e
}

// loadDefaultHierarchy memuat role hierarchy:
//
//	super_admin (root — tidak punya parent)
//	  └── company_admin
//	        └── manager
//	              └── employee
func (e *Enforcer) loadDefaultHierarchy() {
	e.hierarchy[RoleCompanyAdmin] = RoleSuperAdmin
	e.hierarchy[RoleManager] = RoleCompanyAdmin
	e.hierarchy[RoleEmployee] = RoleManager
	// RoleSuperAdmin tidak punya parent (root)
}

// loadDefaultPolicies memuat aturan RBAC default untuk semua role.
//
// super_admin:
//   - Platform & Tenant: semua resource & action
//
// company_admin:
//   - Platform: view-only (company, user, license)
//   - Tenant: full management (org, employee, attendance, leave, payroll, dll)
//
// manager:
//   - Platform: none
//   - Tenant: view + create + update (tanpa delete)
//
// employee:
//   - Platform: none
//   - Tenant: view-only
func (e *Enforcer) loadDefaultPolicies() {
	// ========================================
	// SUPER ADMIN — akses ke semua resource & action
	// ========================================
	e.policies[RoleSuperAdmin] = map[string]string{
		"*": "*", // wildcard: semua resource & action
	}

	// ========================================
	// COMPANY ADMIN — platform view-only + tenant full
	// ========================================
	e.policies[RoleCompanyAdmin] = map[string]string{
		// Platform resources (view-only)
		"company":      "view",
		"user":         "view",
		"license":      "view",

		// Tenant resources (full access)
		"organization": "*",
		"employee":     "*",
		"attendance":   "*",
		"leave":        "*",
		"payroll":      "*",
		"competency":   "*",
		"jobmanagement": "*",
		"approval":     "*",
	}

	// ========================================
	// MANAGER — tenant view + create + update (no delete)
	// ========================================
	e.policies[RoleManager] = map[string]string{
		"organization": "view,create,update",
		"employee":     "view,create,update",
		"attendance":   "view",
		"leave":        "view,create",
		"approval":     "view,create,update",
	}

	// ========================================
	// EMPLOYEE — tenant view-only
	// ========================================
	e.policies[RoleEmployee] = map[string]string{
		"organization": "view",
		"employee":     "view",
		"attendance":   "view",
		"leave":        "view",
		"payroll":      "view",
	}
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
		"job-management": "jobmanagement",
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
