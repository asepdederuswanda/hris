// CLI Installer untuk HRIS Platform
//
// Menyediakan perintah CLI untuk:
//   - Provision tenant baru (create database, run migrations, seed)
//   - Backup & restore tenant
//   - Health check tenant
//
// Usage: go run ./cmd/installer --help
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/inthros/hris-platform/internal/pkg/config"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/logger"
	"github.com/inthros/hris-platform/internal/pkg/migrator"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")

	provisionCmd := flag.NewFlagSet("provision", flag.ExitOnError)
	provisionCmd.String("config", "", "Path to configuration file")
	companyID := provisionCmd.String("company", "", "Company ID to provision")
	dbName := provisionCmd.String("db-name", "", "Tenant database name (optional, generated from company slug if empty)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load config for all commands
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Init logger
	l := logger.New(cfg.Logger.Level, cfg.Logger.Format, "hris-installer")
	defer l.Sync()

	// Init database manager
	dbManager, err := database.NewManager(&database.Config{
		Driver:                  cfg.Database.Driver,
		PlatformDSN:             cfg.Database.PlatformDSN(),
		PlatformHost:            cfg.Database.PlatformHost,
		PlatformPort:            cfg.Database.PlatformPort,
		PlatformUser:            cfg.Database.PlatformUser,
		PlatformPassword:        cfg.Database.PlatformPassword,
		PlatformSSLMode:         cfg.Database.PlatformSSLMode,
		TenantHost:              cfg.Database.TenantHost,
		TenantPort:              cfg.Database.TenantPort,
		TenantSuperUser:         cfg.Database.TenantSuperUser,
		TenantSuperPass:         cfg.Database.TenantSuperPass,
		TenantSSLMode:           cfg.Database.TenantSSLMode,
		MaxOpenConns:            cfg.Database.MaxOpenConns,
		MaxIdleConns:            cfg.Database.MaxIdleConns,
		ConnMaxLifetimeMs:       cfg.Database.ConnMaxLifetimeMs,
		TenantMaxOpenConns:      cfg.Database.TenantMaxOpenConns,
		TenantMaxIdleConns:      cfg.Database.TenantMaxIdleConns,
		TenantConnMaxLifetimeMs: cfg.Database.TenantConnMaxLifetimeMs,
		TenantConnMaxIdleTimeMs: cfg.Database.TenantConnMaxIdleTimeMs,
		LogLevel:                4,
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize database manager", zap.Error(err))
	}
	defer dbManager.CloseAll()

	switch os.Args[1] {
	case "provision":
		provisionCmd.Parse(os.Args[2:])
		handleProvision(l, dbManager, *companyID, *dbName)

	case "migrate":
		if len(os.Args) < 3 {
			log.Fatal("Usage: installer migrate --company=<id>")
		}
		migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
		migrateCompanyID := migrateCmd.String("company", "", "Company ID to migrate")
		migrateCmd.Parse(os.Args[2:])
		handleTenantMigrate(l, dbManager, *migrateCompanyID)

	case "encrypt-passwords":
		handleEncryptPasswords(l, dbManager)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("HRIS Platform CLI Installer")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  installer [--config=<path>] provision --company=<id> [--db-name=<name>]")
	fmt.Println("  installer [--config=<path>] migrate --company=<id>")
	fmt.Println("  installer [--config=<path>] encrypt-passwords")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  provision          Provision a new tenant (create database + run migrations)")
	fmt.Println("  migrate            Run pending tenant migrations for an existing company")
	fmt.Println("  encrypt-passwords  Encrypt legacy plaintext passwords in tenant_connections")
	fmt.Println("")
	fmt.Println("Environment:")
	fmt.Println("  HRIS_ENCRYPTION_KEY    Required for encrypt-passwords (32-byte hex, 64 chars)")
}

func handleProvision(l *zap.Logger, dbManager *database.Manager, companyID, dbName string) {
	if companyID == "" {
		log.Fatal("company is required")
	}

	l.Info("Starting tenant provisioning from CLI",
		zap.String("company_id", companyID),
	)

	if dbName == "" {
		dbName = fmt.Sprintf("hris_tenant_%s", companyID[:8])
	}

	// 1. Create database
	conn, err := dbManager.ProvisionTenant(companyID, dbName, "root", "", dbManager.Driver())
	if err != nil {
		l.Fatal("Failed to create tenant database", zap.Error(err))
	}

	// 2. Simpan TenantConnection
	if err := dbManager.SaveTenantConnection(conn); err != nil {
		l.Fatal("Failed to save tenant connection", zap.Error(err))
	}
	l.Info("Tenant connection saved")

	// 3. Dapatkan koneksi ke tenant DB
	tenantDB, err := dbManager.TenantDB(companyID)
	if err != nil {
		l.Fatal("Failed to connect to tenant database", zap.Error(err))
	}

	// 4. Jalankan tenant migrations (pilih dialect sesuai driver)
	l.Info("Running tenant SQL migrations...")
	tenantRoot := migrator.TenantRootPath(dbManager.Driver())
	tenantMigrator := migrator.New(tenantDB, l, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		l.Fatal("Tenant migration failed", zap.Error(err))
	}

	l.Info("Tenant provisioning completed successfully",
		zap.String("company_id", companyID),
		zap.String("db_name", dbName),
	)
}

func handleEncryptPasswords(l *zap.Logger, dbManager *database.Manager) {
	l.Info("Starting legacy password encryption...")
	l.Warn("Ensure HRIS_ENCRYPTION_KEY environment variable is set before running this command")
	l.Warn("Passwords encrypted with a different key will NOT be recoverable!")

	encrypted, errCount, err := dbManager.EncryptLegacyPasswords()
	if err != nil {
		l.Fatal("Failed to encrypt legacy passwords", zap.Error(err))
	}

	l.Info("Legacy password encryption completed",
		zap.Int("passwords_encrypted", encrypted),
		zap.Int("errors", errCount),
	)

	if encrypted == 0 && errCount == 0 {
		l.Info("No legacy plaintext passwords found — all credentials are already encrypted")
	}
}

func handleTenantMigrate(l *zap.Logger, dbManager *database.Manager, companyID string) {
	if companyID == "" {
		log.Fatal("company is required")
	}

	l.Info("Running tenant migration upgrade",
		zap.String("company_id", companyID),
	)

	tenantDB, err := dbManager.TenantDB(companyID)
	if err != nil {
		l.Fatal("Failed to connect to tenant database", zap.Error(err))
	}

	tenantRoot := migrator.TenantRootPath(dbManager.Driver())
	tenantMigrator := migrator.New(tenantDB, l, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		l.Fatal("Tenant migration failed", zap.Error(err))
	}

	l.Info("Tenant migration completed successfully")
}
