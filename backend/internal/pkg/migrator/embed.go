// Package migrator menyediakan migration runner yang mengeksekusi
// file SQL embedded dan melacak migrasi yang sudah dijalankan.
package migrator

import "embed"

// MigrationsFS berisi semua file SQL migration yang di-embed ke binary.
// Subdirectory:
//   - migrations/platform/ → file DDL untuk platform tables
//   - migrations/seeders/  → file DML untuk seed data
//
//go:embed migrations
var MigrationsFS embed.FS

// RootPlatform adalah root directory untuk platform migrations di dalam embed FS.
const RootPlatform = "migrations/platform"

// RootSeeders adalah root directory untuk seeders di dalam embed FS.
const RootSeeders = "migrations/seeders"

// RootTenantMySQL adalah root directory untuk tenant migrations MySQL.
const RootTenantMySQL = "migrations/tenant/mysql"

// RootTenantPostgres adalah root directory untuk tenant migrations PostgreSQL.
const RootTenantPostgres = "migrations/tenant/postgres"

// RootTenant — backward compatibility alias. Seluruh kode baru harus
// menggunakan TenantRootPath(driver) untuk memilih dialect MySQL atau PostgreSQL.
const RootTenant = RootTenantMySQL

// TenantRootPath mengembalikan root path migrations sesuai driver database.
// Parameters:
//   - driver: "mysql" atau "postgres"
// Returns: path ke direktori migrations yang sesuai
func TenantRootPath(driver string) string {
	switch driver {
	case "mysql":
		return RootTenantMySQL
	case "postgres":
		return RootTenantPostgres
	default:
		// Fallback ke mysql untuk backward compatibility
		return RootTenantMySQL
	}
}
