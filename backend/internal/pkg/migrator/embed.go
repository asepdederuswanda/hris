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

// RootTenant adalah root directory untuk tenant migrations di dalam embed FS.
// Berisi 11 migration files (001-011) untuk semua modul tenant.
const RootTenant = "migrations/tenant"
