// Package database menyediakan manajemen koneksi database
// multi-tenant dengan support untuk PostgreSQL dan MySQL.
//
// Platform DB: database tunggal untuk manajemen platform
// (companies, modules, licenses, platform_users, dll).
//
// Tenant DB: satu database per company/tenant, diakses
// berdasarkan company_id dari JWT claims.
//
// Setiap koneksi bisa menggunakan driver PostgreSQL atau MySQL
// sesuai konfigurasi.
package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/inthros/hris-platform/internal/pkg/driver"
)

// Manager mengelola koneksi database untuk platform dan tenant.
type Manager struct {
	platformDB *gorm.DB
	tenants    map[string]*gorm.DB // company_id -> tenant db connection
	mu         sync.RWMutex
	cfg        *Config
	logger     *zap.Logger
	driver     driver.Type
}

// TenantConnection menyimpan detail koneksi tenant database.
type TenantConnection struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID string `gorm:"column:company_id;type:char(36);uniqueIndex;not null" json:"company_id"`
	Driver    string `gorm:"column:driver;type:varchar(20)" json:"driver"`
	Host      string `gorm:"column:host;type:varchar(255)" json:"host"`
	Port      int    `gorm:"column:port" json:"port"`
	DBName    string `gorm:"column:db_name;type:varchar(100)" json:"db_name"`
	Username  string `gorm:"column:username;type:varchar(100)" json:"username"`
	Password  string `gorm:"column:password;type:varchar(255)" json:"password"`
	SSLMode   string `gorm:"column:ssl_mode;type:varchar(20)" json:"sslmode"`
	IsActive  bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

// Config adalah konfigurasi untuk database manager.
type Config struct {
	Driver            string
	PlatformDSN       string
	PlatformHost      string
	PlatformPort      int
	PlatformUser      string
	PlatformPassword  string
	PlatformSSLMode   string
	TenantHost        string
	TenantPort        int
	TenantSuperUser   string
	TenantSuperPass   string
	TenantSSLMode     string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetimeMs int
	LogLevel          gormlogger.LogLevel
}

// NewManager membuat Manager baru dengan koneksi ke platform database.
func NewManager(cfg *Config, logger *zap.Logger) (*Manager, error) {
	platformDB, err := openGORM(cfg.Driver, cfg.PlatformDSN, cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to connect platform database: %w", err)
	}

	sqlDB, err := platformDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMs) * time.Millisecond)

	driver := driver.Parse(cfg.Driver)

	logger.Info("Connected to platform database",
		zap.String("driver", string(driver)),
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)

	return &Manager{
		platformDB: platformDB,
		tenants:    make(map[string]*gorm.DB),
		cfg:        cfg,
		logger:     logger,
		driver:     driver,
	}, nil
}

// PlatformDB mengembalikan koneksi ke platform database.
func (m *Manager) PlatformDB() *gorm.DB {
	return m.platformDB
}

// Driver mengembalikan tipe database driver yang digunakan.
func (m *Manager) Driver() string {
	return string(m.driver)
}

// TenantDB mengembalikan koneksi ke database tenant berdasarkan company_id.
// Koneksi di-cache setelah pertama kali dibuat.
func (m *Manager) TenantDB(companyID string) (*gorm.DB, error) {
	// Read lock first
	m.mu.RLock()
	db, exists := m.tenants[companyID]
	m.mu.RUnlock()

	if exists {
		return db, nil
	}

	return m.connectTenant(companyID)
}

// TenantDBFromContext mengambil company_id dari context (set by middleware)
// dan mengembalikan koneksi tenant database yang sesuai.
func (m *Manager) TenantDBFromContext(ctx context.Context) (*gorm.DB, error) {
	companyID, ok := ctx.Value("company_id").(string)
	if !ok || companyID == "" {
		return nil, fmt.Errorf("company_id not found in context")
	}
	return m.TenantDB(companyID)
}

func (m *Manager) connectTenant(companyID string) (*gorm.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check setelah write lock
	if db, exists := m.tenants[companyID]; exists {
		return db, nil
	}

	// Ambil detail koneksi dari platform database
	var conn TenantConnection
	if err := m.platformDB.
		Table("tenant_connections").
		Where("company_id = ? AND is_active = ?", companyID, true).
		First(&conn).Error; err != nil {
		return nil, fmt.Errorf("tenant connection not found for company %s: %w", companyID, err)
	}

	// Tentukan driver: prefer dari tenant connection, fallback ke default manager
	driver := conn.Driver
	if driver == "" {
		driver = string(m.driver)
	}

	dsn := buildDSN(driver, conn.Host, conn.Port, conn.Username, conn.Password, conn.DBName, conn.SSLMode)

	db, err := openGORM(driver, dsn, m.cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to connect tenant database for company %s: %w", companyID, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB for tenant: %w", err)
	}
	sqlDB.SetMaxOpenConns(m.cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.cfg.ConnMaxLifetimeMs) * time.Millisecond)

	m.tenants[companyID] = db
	m.logger.Info("Connected to tenant database",
		zap.String("company_id", companyID),
		zap.String("db_name", conn.DBName),
		zap.String("driver", driver),
	)

	return db, nil
}

// ProvisionTenant membuat database tenant baru dan menyimpan konfigurasi koneksi.
// Langkah:
//  1. Generate nama database unik
//  2. Buat database via superuser connection
//  3. Simpan TenantConnection ke platform DB
//  4. Inisialisasi koneksi GORM ke tenant database
//
// Returns: tenant connection info dan error.
func (m *Manager) ProvisionTenant(companyID, dbName, dbUser, dbPassword, driverType string) (*TenantConnection, error) {
	// 1. Connect sebagai superuser untuk create database
	var superDSN string
	// Untuk create database, connect tanpa database tertentu
	switch driver.Parse(driverType) {
	case driver.MySQL:
		superDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			m.cfg.TenantSuperUser, m.cfg.TenantSuperPass,
			m.cfg.TenantHost, m.cfg.TenantPort,
		)
	default: // postgres
		superDSN = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
			m.cfg.TenantHost, m.cfg.TenantPort,
			m.cfg.TenantSuperUser, m.cfg.TenantSuperPass,
			m.cfg.TenantSSLMode,
		)
	}

	superDB, err := openGORM(driverType, superDSN, m.cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to connect as superuser: %w", err)
	}

	// 2. Create database
	switch driver.Parse(driverType) {
	case driver.MySQL:
		if err := superDB.Exec(fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
			dbName,
		)).Error; err != nil {
			return nil, fmt.Errorf("failed to create MySQL database: %w", err)
		}
	default: // postgres
		if err := superDB.Exec(fmt.Sprintf(
			"CREATE DATABASE \"%s\"",
			dbName,
		)).Error; err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL database: %w", err)
		}
	}

	m.logger.Info("Tenant database created",
		zap.String("company_id", companyID),
		zap.String("db_name", dbName),
		zap.String("driver", driverType),
	)

	// 3. Simpan TenantConnection
	conn := &TenantConnection{
		ID:        companyID, // reuse company_id sebagai ID (1-to-1)
		CompanyID: companyID,
		Driver:    driverType,
		Host:      m.cfg.TenantHost,
		Port:      m.cfg.TenantPort,
		DBName:    dbName,
		Username:  dbUser,
		Password:  dbPassword,
		SSLMode:   m.cfg.TenantSSLMode,
		IsActive:  true,
	}

	// Close superuser connection
	if sqlDB, err := superDB.DB(); err == nil {
		sqlDB.Close()
	}

	m.logger.Info("Tenant provisioned successfully",
		zap.String("company_id", companyID),
		zap.String("db_name", dbName),
	)

	return conn, nil
}

// SaveTenantConnection menyimpan atau mengupdate TenantConnection di platform DB.
func (m *Manager) SaveTenantConnection(conn *TenantConnection) error {
	return m.platformDB.Table("tenant_connections").Create(conn).Error
}

// FindTenantConnection mencari TenantConnection berdasarkan company_id.
func (m *Manager) FindTenantConnection(companyID string) (*TenantConnection, error) {
	var conn TenantConnection
	if err := m.platformDB.
		Table("tenant_connections").
		Where("company_id = ?", companyID).
		First(&conn).Error; err != nil {
		return nil, fmt.Errorf("tenant connection not found: %w", err)
	}
	return &conn, nil
}

// DeactivateTenantConnection menonaktifkan tenant connection (is_active = false)
// dan menutup koneksi yang di-cache.
func (m *Manager) DeactivateTenantConnection(companyID string) error {
	result := m.platformDB.
		Table("tenant_connections").
		Where("company_id = ?", companyID).
		Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("failed to deactivate tenant connection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant connection not found for company %s", companyID)
	}

	// Tutup koneksi cache agar koneksi baru tidak bisa dibuat
	m.CloseTenantConnection(companyID)

	m.logger.Info("Tenant connection deactivated",
		zap.String("company_id", companyID),
	)
	return nil
}

// ActivateTenantConnection mengaktifkan tenant connection (is_active = true).
func (m *Manager) ActivateTenantConnection(companyID string) error {
	result := m.platformDB.
		Table("tenant_connections").
		Where("company_id = ?", companyID).
		Update("is_active", true)
	if result.Error != nil {
		return fmt.Errorf("failed to activate tenant connection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant connection not found for company %s", companyID)
	}

	// Hapus cache koneksi lama supaya reconnect dengan is_active=true
	m.mu.Lock()
	delete(m.tenants, companyID)
	m.mu.Unlock()

	m.logger.Info("Tenant connection activated",
		zap.String("company_id", companyID),
	)
	return nil
}

// RemoveTenantConnection menghapus record TenantConnection dari platform DB.
func (m *Manager) RemoveTenantConnection(companyID string) error {
	// Tutup koneksi cache terlebih dahulu
	m.CloseTenantConnection(companyID)

	result := m.platformDB.
		Table("tenant_connections").
		Where("company_id = ?", companyID).
		Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to remove tenant connection: %w", result.Error)
	}

	m.logger.Info("Tenant connection removed",
		zap.String("company_id", companyID),
	)
	return nil
}

// CloseTenantConnection menutup koneksi tenant yang di-cache.
func (m *Manager) CloseTenantConnection(companyID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, exists := m.tenants[companyID]; exists {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		delete(m.tenants, companyID)
		m.logger.Info("Tenant connection closed",
			zap.String("company_id", companyID),
		)
	}
}

// DropTenantDB menghapus database tenant melalui superuser connection.
func (m *Manager) DropTenantDB(companyID string) error {
	// Cari TenantConnection untuk mendapatkan nama database
	conn, err := m.FindTenantConnection(companyID)
	if err != nil {
		return fmt.Errorf("cannot find tenant connection: %w", err)
	}

	// Tutup koneksi cache sebelum drop database
	m.CloseTenantConnection(companyID)

	// Build superuser DSN (tanpa database tertentu)
	driverType := conn.Driver
	if driverType == "" {
		driverType = string(m.driver)
	}

	var superDSN string
	switch driver.Parse(driverType) {
	case driver.MySQL:
		superDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			m.cfg.TenantSuperUser, m.cfg.TenantSuperPass,
			m.cfg.TenantHost, m.cfg.TenantPort,
		)
	default: // postgres
		superDSN = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
			m.cfg.TenantHost, m.cfg.TenantPort,
			m.cfg.TenantSuperUser, m.cfg.TenantSuperPass,
			m.cfg.TenantSSLMode,
		)
	}

	superDB, err := openGORM(driverType, superDSN, m.cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to connect as superuser: %w", err)
	}
	defer func() {
		if sqlDB, err := superDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// DROP DATABASE
	switch driver.Parse(driverType) {
	case driver.MySQL:
		if err := superDB.Exec(fmt.Sprintf(
			"DROP DATABASE IF EXISTS `%s`",
			conn.DBName,
		)).Error; err != nil {
			return fmt.Errorf("failed to drop MySQL database: %w", err)
		}
	default: // postgres
		// PostgreSQL: terminate connections first then drop
		superDB.Exec(fmt.Sprintf(
			"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s'",
			conn.DBName,
		))
		if err := superDB.Exec(fmt.Sprintf(
			"DROP DATABASE IF EXISTS \"%s\"",
			conn.DBName,
		)).Error; err != nil {
			return fmt.Errorf("failed to drop PostgreSQL database: %w", err)
		}
	}

	m.logger.Info("Tenant database dropped",
		zap.String("company_id", companyID),
		zap.String("db_name", conn.DBName),
	)
	return nil
}

// CloseAll menutup semua koneksi database.
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close platform DB
	if sqlDB, err := m.platformDB.DB(); err == nil {
		sqlDB.Close()
	}

	// Close all tenant connections
	for companyID, db := range m.tenants {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		delete(m.tenants, companyID)
	}

	return nil
}

// HealthCheck melakukan ping ke semua koneksi database.
func (m *Manager) HealthCheck() map[string]error {
	results := make(map[string]error)

	// Platform DB
	if sqlDB, err := m.platformDB.DB(); err == nil {
		results["platform"] = sqlDB.Ping()
	} else {
		results["platform"] = err
	}

	// Tenant connections
	m.mu.RLock()
	for companyID, db := range m.tenants {
		if sqlDB, err := db.DB(); err == nil {
			results[fmt.Sprintf("tenant:%s", companyID)] = sqlDB.Ping()
		}
	}
	m.mu.RUnlock()

	return results
}

// ========================================================================
// Helper functions
// ========================================================================

// openGORM membuka koneksi GORM dengan driver yang sesuai.
func openGORM(drv, dsn string, logLevel gormlogger.LogLevel) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch driver.Parse(drv) {
	case driver.MySQL:
		dialector = mysql.Open(dsn)
	default: // postgres
		dialector = postgres.Open(dsn)
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
}

// buildDSN membuat connection string sesuai driver.
// Tenant DSN menggunakan multiStatements=true agar bisa menjalankan
// file SQL migration yang berisi multiple statements (CREATE TABLE, INSERT, dll).
func buildDSN(drv, host string, port int, user, password, dbName, sslMode string) string {
	switch driver.Parse(drv) {
	case driver.MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
			user, password, host, port, dbName,
		)
	default: // postgres
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbName, sslMode,
		)
	}
}
