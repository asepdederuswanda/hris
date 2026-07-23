// HRIS Platform - Go Modular Monolith
//
// Entry point untuk aplikasi HRIS platform.
// Melakukan inisialisasi shared infrastructure, mendaftarkan
// semua modul (platform & tenant), dan menjalankan HTTP server.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/authz"
	"github.com/inthros/hris-platform/internal/pkg/cache"
	"github.com/inthros/hris-platform/internal/pkg/config"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/logger"
	"github.com/inthros/hris-platform/internal/pkg/middleware"
	"github.com/inthros/hris-platform/internal/pkg/module"
	"github.com/inthros/hris-platform/internal/pkg/router"
	"github.com/inthros/hris-platform/internal/pkg/docs"
	"github.com/inthros/hris-platform/internal/pkg/migrator"

	// Platform modules
	"github.com/inthros/hris-platform/internal/platform/company"
	"github.com/inthros/hris-platform/internal/platform/license"
	"github.com/inthros/hris-platform/internal/platform/modulemgmt"
	"github.com/inthros/hris-platform/internal/platform/monitoring"
	"github.com/inthros/hris-platform/internal/platform/user"

	// Tenant modules
	"github.com/inthros/hris-platform/internal/modules/employee"
	"github.com/inthros/hris-platform/internal/modules/jobmanagement"
	"github.com/inthros/hris-platform/internal/modules/organization"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	migrateDown := flag.Bool("migrate-down", false, "Rollback all applied migrations and exit")
	migrateTo := flag.String("migrate-to", "", "Rollback migrations to specific version (exclusive) and exit")
	flag.Parse()

	// Validate flags: --migrate-down and --migrate-to are mutually exclusive
	if *migrateDown && *migrateTo != "" {
		fmt.Fprintln(os.Stderr, "ERROR: --migrate-down and --migrate-to are mutually exclusive")
		os.Exit(1)
	}

	// 1. Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize logger
	l := logger.New(cfg.Logger.Level, cfg.Logger.Format, "hris-platform")
	defer l.Sync()

	// 3. Initialize distributed cache (Redis + local in-memory + Pub/Sub invalidation)
	cacheManager, err := cache.New(cache.Config{
		RedisAddr:     cfg.Redis.RedisAddr(),
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
		DefaultTTL:    time.Duration(cfg.Cache.DefaultTTL) * time.Second,
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize cache", zap.Error(err))
	}
	defer cacheManager.Close()

	// 4. Initialize database manager (multi-tenant, multi-driver)
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
		LogLevel:                4, // Warn
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize database manager", zap.Error(err))
	}
	defer dbManager.CloseAll()

	// 5. Handle migration CLI commands (run and exit without starting server)
	if *migrateDown || *migrateTo != "" {
		cacheManager.Close()
		runMigrationCommand(l, dbManager, *migrateDown, *migrateTo)
		return
	}

	// 5. Initialize JWT auth manager
	authManager := auth.NewManager(auth.Config{
		Secret:          cfg.JWT.Secret,
		AccessTokenTTL:  time.Duration(cfg.JWT.AccessTokenTTL) * time.Minute,
		RefreshTokenTTL: time.Duration(cfg.JWT.RefreshTokenTTL) * time.Hour,
		Issuer:          cfg.JWT.Issuer,
	})

	// 6. Initialize module registry
	var platformModules []module.ModuleRegistration
	var tenantModules []module.ModuleRegistration

	// Create auth middleware once (reused across all platform modules)
	authMW := middleware.AuthJWT(authManager, l)

	// 6a. Register platform modules (ordered by priority)
	// Note: rbacMW diinisialisasi setelah SQL migrations karena membutuhkan
	// tabel RBAC di database. Tapi kita pass nil dulu dan update setelah init.
	platformModules = append(platformModules,
		module.ModuleRegistration{
			Module:   company.NewModule(dbManager, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 1,
		},
		module.ModuleRegistration{
			Module:   user.NewModule(dbManager, authManager, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 2,
		},
		module.ModuleRegistration{
			Module:   modulemgmt.NewModule(dbManager, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 3,
		},
		module.ModuleRegistration{
			Module:   license.NewModule(dbManager, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 4,
		},
	)

	// 6b. Register tenant modules
	tenantModules = append(tenantModules,
		module.ModuleRegistration{
			Module:   organization.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 1,
		},
		module.ModuleRegistration{
			Module:   employee.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 2,
		},
		module.ModuleRegistration{
			Module:   jobmanagement.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 3,
		},
	)

	// 7. Run SQL file migrations for platform modules
	l.Info("Running platform SQL migrations...")
	platformMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootPlatform)
	if err := platformMigrator.Up(); err != nil {
		l.Fatal("Platform SQL migration failed", zap.Error(err))
	}

	// 8. Run AutoMigrate for platform modules (sync GORM models to schema)
	l.Info("Running platform AutoMigrate...")
	for _, reg := range platformModules {
		if err := reg.Module.Migrate(dbManager.PlatformDB()); err != nil {
			l.Fatal("Platform AutoMigrate failed",
				zap.String("module", reg.Module.Info().Name),
				zap.Error(err),
			)
		}
		l.Info("Platform AutoMigrate completed",
			zap.String("module", reg.Module.Info().Name),
		)
	}

	// Note: Tenant migrations run during tenant provisioning,
	// not at startup. Each tenant gets its own database.
	l.Info("Running SQL seeders...")
	seederMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootSeeders)
	if err := seederMigrator.Up(); err != nil {
		l.Warn("SQL seeder warning", zap.Error(err))
	}

	// 10. Initialize database-backed RBAC enforcer (seeds defaults if empty)
	l.Info("Initializing database-backed RBAC enforcer...")
	rbacEnforcer, err := authz.NewEnforcerFromDB(dbManager.PlatformDB())
	if err != nil {
		l.Fatal("Failed to initialize RBAC enforcer", zap.Error(err))
	}
	rbacMW := authz.NewMiddleware(authz.MiddlewareConfig{Enforcer: rbacEnforcer})

	// 10a. Create RBAC service & handler for management API
	rbacLogger := l.Named("rbac")
	rbacRepo := authz.NewRepository(dbManager.PlatformDB())
	rbacService := authz.NewService(rbacRepo, rbacEnforcer, rbacLogger)
	rbacHandler := authz.NewHandler(rbacService)

	// 10b. Run module seeders for platform modules
	l.Info("Running platform module seeders...")
	for _, reg := range platformModules {
		if err := reg.Module.Seed(dbManager.PlatformDB()); err != nil {
			l.Warn("Platform seeder warning",
				zap.String("module", reg.Module.Info().Name),
				zap.Error(err),
			)
		}
	}

	// 11. Setup router and middleware
	r := router.Setup(
		router.Config{Mode: cfg.Server.Mode},
		middleware.AuthJWT(authManager, l),
		middleware.TenantRequired(),
		rbacMW,
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins:   cfg.CORS.AllowedOrigins,
			AllowedMethods:   cfg.CORS.AllowedMethods,
			AllowedHeaders:   cfg.CORS.AllowedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		}),
		middleware.RequestLogger(l),
		middleware.Recovery(l),
		platformModules,
		tenantModules,
	)

	// Register platform RBAC management routes (standalone)
	platformGroup := r.Group("/api/v1/platform")
	platformGroup.Use(authMW)
	platformGroup.Use(rbacMW)
	authz.RegisterRoutes(platformGroup, rbacHandler)

	// Register platform monitoring routes (standalone, no module interface needed)
	monitoringHandler := monitoring.NewHandler(dbManager, cacheManager, l)
	monitoring.RegisterRoutes(platformGroup, monitoringHandler, authMW, rbacMW)

	// Register Scalar API Documentation
	r.GET("/docs", docs.ScalarUIHandler())
	r.GET("/openapi.json", docs.OpenAPIHandler())

	// 12. Start server
	l.Info("Starting HRIS Platform server",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)
	l.Info("API Documentation available at",
		zap.String("url", "/docs"),
	)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		l.Fatal("Failed to start server", zap.Error(err))
	}
}

// runMigrationCommand mengeksekusi perintah migration CLI dan exit.
// Digunakan untuk --migrate-down dan --migrate-to flags.
func runMigrationCommand(l *zap.Logger, dbManager *database.Manager, down bool, to string) {
	l.Info("Migration command detected, running in CLI mode")

	platformMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootPlatform)
	seederMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootSeeders)

	if down {
		// Rollback all: seeders first (reverse), then platform
		l.Info("Rolling back all seeders...")
		if err := seederMigrator.Down(); err != nil {
			l.Fatal("Seeder rollback failed", zap.Error(err))
		}

		l.Info("Rolling back all platform migrations...")
		if err := platformMigrator.Down(); err != nil {
			l.Fatal("Platform migration rollback failed", zap.Error(err))
		}
	} else if to != "" {
		// Rollback to specific version
		l.Info("Rolling back platform migrations to version",
			zap.String("target", to))
		if err := platformMigrator.DownTo(to); err != nil {
			l.Fatal("Platform migration rollback failed", zap.Error(err))
		}

		l.Info("Rolling back seeders to version",
			zap.String("target", to))
		if err := seederMigrator.DownTo(to); err != nil {
			l.Warn("Seeder rollback warning", zap.Error(err))
		}
	}

	l.Info("Migration command completed successfully")
}
