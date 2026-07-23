// Package migrator — integration test against real MySQL and PostgreSQL databases.
//
// Prerequisites:
//   - MySQL running on localhost:3306, root without password
//   - PostgreSQL running on localhost:5432, postgres:password
//   - Database hris_migrate_test must exist on both (created by script)
//
// Run: cd backend && go test -v -run TestMigratorIntegration ./internal/pkg/migrator/... -timeout 120s
//
// or for a specific dialect:
//   go test -v -run TestMigratorIntegration/MySQL ./internal/pkg/migrator/... -timeout 60s
//   go test -v -run TestMigratorIntegration/PostgreSQL ./internal/pkg/migrator/... -timeout 60s
package migrator

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Daftar tabel yang diharapkan dari platform migrations.
var expectedPlatformTables = []string{
	"companies",
	"tenant_connections",
	"platform_users",
	"modules",
	"company_modules",
	"licenses",
}

// expectedTenantKeyTables adalah tabel kunci yang harus ada setelah tenant migrations.
// Dipilih dari berbagai migration (001-011) untuk memastikan semuanya tereksekusi.
// Nama tabel diambil dari file SQL yang sebenarnya.
var expectedTenantKeyTables = []string{
	"religions",              // 001_master_data
	"positions",              // 002_organization
	"employees",              // 003_employee
	"attendance_sessions",    // 004_attendance
	"leave_types",            // 005_leave
	"salary_components",      // 006_payroll_structure
	"payroll_payslips",       // 007_payroll_run
	"competencies",           // 008_competency
	"job_management_titles",  // 009_job_management
	"approval_flows",         // 010_approval
	"permissions",            // 011_settings
}

// =========================================================================
// TestTenantRootPath — unit test untuk fungsi pemilih dialect
// =========================================================================

func TestTenantRootPath(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		expected string
	}{
		{"mysql", "mysql", "migrations/tenant/mysql"},
		{"postgres", "postgres", "migrations/tenant/postgres"},
		{"empty_defaults_to_mysql", "", "migrations/tenant/mysql"},
		{"unknown_defaults_to_mysql", "sqlite", "migrations/tenant/mysql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TenantRootPath(tt.driver)
			if got != tt.expected {
				t.Errorf("TenantRootPath(%q) = %q, want %q", tt.driver, got, tt.expected)
			}
		})
	}
}

// =========================================================================
// Integration Test — runs against real databases
// =========================================================================

// mysqlDSN adalah DSN untuk koneksi ke MySQL test database.
const mysqlDSN = "root:@tcp(localhost:3306)/hris_migrate_test?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

// postgresDSN adalah DSN untuk koneksi ke PostgreSQL test database.
const postgresDSN = "host=localhost port=5432 user=postgres password=password dbname=hris_migrate_test sslmode=disable"

// testLogger adalah zap.Logger untuk test.
func testLogger(t *testing.T) *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	return logger
}

// openDB membuka koneksi database dengan driver tertentu.
// Menerima testing.TB sehingga bisa dipakai oleh *testing.T dan *testing.B.
func openDB(t testing.TB, driver, dsn string) *gorm.DB {
	t.Helper()

	var dialector gorm.Dialector
	switch strings.ToLower(driver) {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		t.Fatalf("unsupported driver: %s", driver)
		return nil
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		t.Fatalf("failed to connect to %s: %v", driver, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB for %s: %v", driver, err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(30 * time.Second)

	return db
}

// tableExists mengecek apakah tabel ada di database.
func tableExists(t testing.TB, db *gorm.DB, tableName, driver string) bool {
	t.Helper()

	var count int64
	var sql string

	switch driver {
	case "mysql":
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	case "postgres":
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_catalog = current_database() AND table_schema = 'public' AND table_name = ?"
	default:
		t.Fatalf("unsupported driver: %s", driver)
	}

	if err := db.Raw(sql, tableName).Scan(&count).Error; err != nil {
		t.Logf("Warning: failed to check table %s: %v", tableName, err)
		return false
	}

	return count > 0
}

// TestMigratorIntegration menjalankan migrasi terhadap database nyata.
func TestMigratorIntegration(t *testing.T) {
	// Skip jika bukan integration test
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("MySQL", func(t *testing.T) {
		testDialect(t, "mysql", mysqlDSN)
	})

	t.Run("PostgreSQL", func(t *testing.T) {
		testDialect(t, "postgres", postgresDSN)
	})
}

// testDialect menjalankan test migrasi untuk satu dialect.
func testDialect(t *testing.T, driver, dsn string) {
	logger := testLogger(t)
	db := openDB(t, driver, dsn)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// ============================
	// TEST 1: Platform Migrations
	// ============================
	t.Run("Platform", func(t *testing.T) {
		m := New(db, logger, MigrationsFS, RootPlatform)

		// Run Up
		if err := m.Up(); err != nil {
			t.Fatalf("Platform migrations UP failed: %v", err)
		}

		// Verify tables exist
		for _, table := range expectedPlatformTables {
			if !tableExists(t, db, table, driver) {
				t.Errorf("Platform migration: expected table %q not found", table)
			} else {
				t.Logf("  ✓ Table %q exists", table)
			}
		}

		// Verify tracking table schema_migrations has records
		var count int64
		db.Model(&schemaMigration{}).Count(&count)
		if count == 0 {
			t.Error("schema_migrations table is empty after platform migrations")
		} else {
			t.Logf("  ✓ schema_migrations has %d records", count)
		}

		// Run Up again — should be idempotent
		if err := m.Up(); err != nil {
			t.Errorf("Platform migrations UP (2nd run, idempotent check) failed: %v", err)
		}

		// Verify count didn't change
		var count2 int64
		db.Model(&schemaMigration{}).Count(&count2)
		if count2 != count {
			t.Errorf("schema_migrations count changed after 2nd Up: %d -> %d", count, count2)
		}

		// Run Down
		if err := m.Down(); err != nil {
			t.Fatalf("Platform migrations DOWN failed: %v", err)
		}

		t.Logf("  ✅ Platform migrations for %s completed successfully", driver)
	})

	// ============================
	// TEST 2: Tenant Migrations
	// ============================
	// Catatan: Platform dan Tenant migrations berbagi tabel schema_migrations
	// dengan nomor versi yang berbeda (platform: 001-006, tenant: 001-011).
	// Setelah platform test, versi 001-006 sudah tercatat. Agar tenant test
	// tidak skip versi 001-006, kita bersihkan schema_migrations dulu.
	t.Run("Tenant", func(t *testing.T) {
		// Reset schema_migrations untuk menghindari overflow versi
		if err := db.Exec("DELETE FROM schema_migrations WHERE 1=1").Error; err != nil {
			t.Fatalf("Failed to clear schema_migrations: %v", err)
		}

		tenantRoot := TenantRootPath(driver)
		m := New(db, logger, MigrationsFS, tenantRoot)

		t.Logf("Using tenant root: %s", tenantRoot)

		// Run Up
		if err := m.Up(); err != nil {
			t.Fatalf("Tenant migrations UP for %s failed: %v", driver, err)
		}

		// Verify key tables exist (strict check)
		var missingTables int
		for _, table := range expectedTenantKeyTables {
			if !tableExists(t, db, table, driver) {
				t.Errorf("Tenant migration: expected table %q NOT found", table)
				missingTables++
			} else {
				t.Logf("  ✓ Table %q exists", table)
			}
		}

		if missingTables > 0 {
			t.Fatalf("%d expected tenant tables are missing", missingTables)
		}

		// Verify tracking table
		var count int64
		db.Model(&schemaMigration{}).Count(&count)
		if count == 0 {
			t.Error("schema_migrations is empty after tenant migrations")
		} else {
			t.Logf("  ✓ schema_migrations has %d records after tenant migrations", count)
		}

		// Run Down — rollback semua
		if err := m.Down(); err != nil {
			t.Fatalf("Tenant migrations DOWN for %s failed: %v", driver, err)
		}

		t.Logf("  ✅ Tenant migrations for %s completed successfully", driver)
	})
}

// =========================================================================
// TestEmbeddedFiles — verifikasi bahwa semua file migration ter-embed
// =========================================================================

func TestEmbeddedFiles(t *testing.T) {
	// Check MySQL tenant directory
	mysqlFiles, err := MigrationsFS.ReadDir("migrations/tenant/mysql")
	if err != nil {
		t.Fatalf("Failed to read MySQL tenant migrations dir: %v", err)
	}
	t.Logf("MySQL tenant files: %d", len(mysqlFiles))

	// Check PostgreSQL tenant directory
	pgFiles, err := MigrationsFS.ReadDir("migrations/tenant/postgres")
	if err != nil {
		t.Fatalf("Failed to read PostgreSQL tenant migrations dir: %v", err)
	}
	t.Logf("PostgreSQL tenant files: %d", len(pgFiles))

	// Verify same count
	if len(mysqlFiles) != len(pgFiles) {
		t.Errorf("MySQL and PostgreSQL tenant dirs have different file counts: %d vs %d",
			len(mysqlFiles), len(pgFiles))
	}

	// Build sets of file names
	mysqlNames := make(map[string]bool)
	pgNames := make(map[string]bool)

	for _, f := range mysqlFiles {
		if !f.IsDir() {
			mysqlNames[f.Name()] = true
		}
	}
	for _, f := range pgFiles {
		if !f.IsDir() {
			pgNames[f.Name()] = true
		}
	}

	// Check MySQL has all PG files
	for name := range pgNames {
		if !mysqlNames[name] {
			t.Errorf("MySQL missing file that PostgreSQL has: %s", name)
		}
	}

	// Check PG has all MySQL files
	for name := range mysqlNames {
		if !pgNames[name] {
			t.Errorf("PostgreSQL missing file that MySQL has: %s", name)
		}
	}

	// Verify platform files
	platformFiles, err := MigrationsFS.ReadDir("migrations/platform")
	if err != nil {
		t.Fatalf("Failed to read platform migrations dir: %v", err)
	}
	t.Logf("Platform files: %d", len(platformFiles))

	// Verify seeder files
	seederFiles, err := MigrationsFS.ReadDir("migrations/seeders")
	if err != nil {
		t.Fatalf("Failed to read seeders dir: %v", err)
	}
	t.Logf("Seeder files: %d", len(seederFiles))

	// Print content summary of a key PostgreSQL file to verify dialect
	pg001Content, err := MigrationsFS.ReadFile("migrations/tenant/postgres/001_master_data.sql")
	if err != nil {
		t.Fatalf("Failed to read PG 001_master_data.sql: %v", err)
	}
	content := string(pg001Content)

	// PostgreSQL should NOT have ENGINE clause
	if strings.Contains(content, "ENGINE=") {
		t.Error("PostgreSQL migration file contains MySQL-only ENGINE clause")
	}

	// PostgreSQL should NOT have backtick identifiers
	if strings.Contains(content, "`") {
		t.Error("PostgreSQL migration file contains backtick identifiers (MySQL-only)")
	}

	t.Log("  ✅ All embedded file checks passed")
}

// =========================================================================
// Benchmark — mengukur performa migrasi (kedua dialect)
// =========================================================================

func BenchmarkMigrator(b *testing.B) {
	// Skip benchmark di CI
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	drivers := []struct {
		name   string
		driver string
		dsn    string
	}{
		{"MySQL", "mysql", mysqlDSN},
		{"PostgreSQL", "postgres", postgresDSN},
	}

	for _, d := range drivers {
		b.Run(fmt.Sprintf("Platform_%s", d.name), func(b *testing.B) {
			benchDialect(b, d.dsn, d.driver, RootPlatform)
		})
		b.Run(fmt.Sprintf("Tenant_%s", d.name), func(b *testing.B) {
			benchDialect(b, d.dsn, d.driver, TenantRootPath(d.driver))
		})
	}
}

// benchDialect menjalankan benchmark up+down migration.
// schema_migrations di-drop di setiap iterasi agar migrasi benar-benar
// dijalankan (bukan skip karena sudah applied).
func benchDialect(b *testing.B, dsn, driver, root string) {
	logger, _ := zap.NewDevelopment()

	for i := 0; i < b.N; i++ {
		func() {
			db := openDB(b, driver, dsn)
			sqlDB, _ := db.DB()
			defer sqlDB.Close()

			// Reset schema_migrations agar migrasi benar-benar dijalankan
			db.Exec("DROP TABLE IF EXISTS schema_migrations")

			m := New(db, logger, MigrationsFS, root)
			if err := m.Up(); err != nil {
				b.Fatalf("Up failed: %v", err)
			}
			if err := m.Down(); err != nil {
				b.Fatalf("Down failed: %v", err)
			}
		}()
	}
}
