// Package migrator menyediakan migration runner yang mengeksekusi
// file SQL embedded dan melacak migrasi yang sudah dijalankan.
//
// Format file:
//   - Up:   {prefix}_{nama}.sql          Contoh: 001_create_companies.sql
//   - Down: {prefix}_{nama}.down.sql     Contoh: 001_create_companies.down.sql
//
// Down file bersifat opsional. Jika tidak ada, migrasi tidak bisa di-rollback.
//
// Cara kerja:
//  1. Buat tabel schema_migrations (jika belum ada)
//  2. Scan direktori migrations/ untuk file *.sql
//  3. Pasangkan up file dengan down file berdasarkan prefix numerik
//  4. Up: Jalankan file yang belum tercatat di schema_migrations
//  5. Down: Rollback file yang sudah tercatat (reverse order)
package migrator

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migration merepresentasikan satu pasang migrasi (up + down).
type Migration struct {
	Version     string // "001", "002", dst
	Name        string // "create_companies"
	FilePath    string // "platform/001_create_companies.sql"
	Content     string // Isi file SQL (up)
	DownContent string // Isi file SQL untuk rollback (opsional, bisa kosong)
	HasDown     bool   // Apakah down file tersedia
}

// Migrator adalah engine migrasi SQL.
type Migrator struct {
	db            *gorm.DB
	logger        *zap.Logger
	migrationsFS  embed.FS
	rootDir       string // root di dalam embed FS, misal "migrations/platform"
}

// New membuat Migrator baru.
//
// Parameters:
//   - db: koneksi database GORM
//   - logger: logger
//   - migrationsFS: embedded filesystem yang berisi file SQL
//   - rootDir: path root di dalam embed FS (contoh: "migrations/platform")
func New(db *gorm.DB, logger *zap.Logger, migrationsFS embed.FS, rootDir string) *Migrator {
	return &Migrator{
		db:           db,
		logger:       logger,
		migrationsFS: migrationsFS,
		rootDir:      rootDir,
	}
}

// =============================================================================
// UP — Menjalankan migrasi
// =============================================================================

// Up menjalankan semua migrasi yang belum dijalankan.
func (m *Migrator) Up() error {
	// 1. Buat schema_migrations table jika belum ada
	if err := m.ensureTrackingTable(); err != nil {
		return fmt.Errorf("failed to create tracking table: %w", err)
	}

	// 2. Scan dan parse file migrasi
	migrations, err := m.scanMigrations()
	if err != nil {
		return fmt.Errorf("failed to scan migrations: %w", err)
	}

	if len(migrations) == 0 {
		m.logger.Info("No pending migrations found",
			zap.String("dir", m.rootDir))
		return nil
	}

	// 3. Dapatkan daftar migrasi yang sudah dijalankan
	applied, err := m.getAppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedSet := make(map[string]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	// 4. Jalankan migrasi yang pending
	pendingCount := 0
	for _, mig := range migrations {
		if appliedSet[mig.Version] {
			m.logger.Debug("Migration already applied, skipping",
				zap.String("version", mig.Version),
				zap.String("name", mig.Name))
			continue
		}

		m.logger.Info("Running migration",
			zap.String("version", mig.Version),
			zap.String("name", mig.Name),
			zap.String("file", mig.FilePath))

		if err := m.executeUpMigration(mig); err != nil {
			return fmt.Errorf("migration %s (%s) failed: %w",
				mig.Version, mig.Name, err)
		}

		pendingCount++
	}

	if pendingCount == 0 {
		m.logger.Info("All migrations already applied",
			zap.String("dir", m.rootDir))
	} else {
		m.logger.Info("Migrations completed",
			zap.Int("applied", pendingCount),
			zap.String("dir", m.rootDir))
	}

	return nil
}

// =============================================================================
// DOWN — Rollback migrasi
// =============================================================================

// Down melakukan rollback SEMUA migrasi yang sudah dijalankan (reverse order).
func (m *Migrator) Down() error {
	return m.DownTo("")
}

// DownTo melakukan rollback migrasi hingga (tidak termasuk) versi target.
// Rollback dilakukan dalam reverse order (dari versi tertinggi ke terendah).
//
// Contoh: applied = [001, 002, 003, 004, 005, 006]
//   - DownTo("004") → rollback 006, 005 (004 dan sebelumnya tetap applied)
//   - DownTo("")    → rollback semua (006, 005, 004, 003, 002, 001)
func (m *Migrator) DownTo(targetVersion string) error {
	// 1. Pastikan tracking table ada
	if err := m.ensureTrackingTable(); err != nil {
		return fmt.Errorf("failed to create tracking table: %w", err)
	}

	// 2. Dapatkan daftar versi yang sudah dijalankan (descending)
	applied, err := m.getAppliedVersionsDesc()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		m.logger.Info("No migrations to rollback",
			zap.String("dir", m.rootDir))
		return nil
	}

	// 3. Scan migrations untuk mendapatkan down content
	migrations, err := m.scanMigrations()
	if err != nil {
		return fmt.Errorf("failed to scan migrations: %w", err)
	}

	// Build lookup map by version
	migMap := make(map[string]Migration)
	for _, mig := range migrations {
		migMap[mig.Version] = mig
	}

	// 4. Tentukan mana yang akan di-rollback
	rollbackVersions := make([]string, 0)
	for _, v := range applied {
		if targetVersion == "" || v > targetVersion {
			rollbackVersions = append(rollbackVersions, v)
		}
	}

	if len(rollbackVersions) == 0 {
		m.logger.Info("No migrations to rollback (already at target version)",
			zap.String("target", targetVersion),
			zap.String("dir", m.rootDir))
		return nil
	}

	// 5. Rollback satu per satu (sudah descending order dari getAppliedVersionsDesc)
	rolledBack := 0
	for _, v := range rollbackVersions {
		mig, exists := migMap[v]
		if !exists {
			m.logger.Warn("Migration file not found, skipping rollback",
				zap.String("version", v))
			continue
		}

		if !mig.HasDown {
			m.logger.Warn("No down file for migration, skipping rollback",
				zap.String("version", mig.Version),
				zap.String("name", mig.Name))
			continue
		}

		m.logger.Info("Rolling back migration",
			zap.String("version", mig.Version),
			zap.String("name", mig.Name))

		if err := m.executeDownMigration(mig); err != nil {
			return fmt.Errorf("rollback %s (%s) failed: %w",
				mig.Version, mig.Name, err)
		}

		rolledBack++
	}

	if rolledBack == 0 {
		m.logger.Warn("No migrations were rolled back (missing down files)",
			zap.String("dir", m.rootDir))
	} else {
		m.logger.Info("Rollback completed",
			zap.Int("rolled_back", rolledBack),
			zap.String("dir", m.rootDir))
	}

	return nil
}

// =============================================================================
// MIGRATION EXECUTION
// =============================================================================

// executeUpMigration menjalankan satu file migrasi up dalam transaction.
func (m *Migrator) executeUpMigration(mig Migration) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Jalankan SQL up
		if err := tx.Exec(mig.Content).Error; err != nil {
			return fmt.Errorf("execute SQL failed: %w", err)
		}

		// Catat ke tabel tracking
		record := schemaMigration{
			Version:   mig.Version,
			Name:      mig.Name,
			AppliedAt: time.Now(),
			Checksum:  fmt.Sprintf("%d", len(mig.Content)),
			FilePath:  mig.FilePath,
		}

		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})
}

// executeDownMigration menjalankan satu file migrasi down dalam transaction.
func (m *Migrator) executeDownMigration(mig Migration) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Jalankan SQL down
		if err := tx.Exec(mig.DownContent).Error; err != nil {
			return fmt.Errorf("execute down SQL failed: %w", err)
		}

		// Hapus dari tabel tracking
		if err := tx.Where("version = ?", mig.Version).Delete(&schemaMigration{}).Error; err != nil {
			return fmt.Errorf("failed to remove migration record: %w", err)
		}

		return nil
	})
}

// =============================================================================
// TRACKING TABLE
// =============================================================================

// schemaMigration model untuk tabel tracking.
type schemaMigration struct {
	Version   string    `gorm:"primaryKey;column:version"`
	Name      string    `gorm:"column:name"`
	AppliedAt time.Time `gorm:"column:applied_at"`
	Checksum  string    `gorm:"column:checksum"`
	FilePath  string    `gorm:"column:file_path"`
}

func (schemaMigration) TableName() string {
	return "schema_migrations"
}

// ensureTrackingTable membuat tabel schema_migrations jika belum ada.
// Menggunakan SQL cross-dialect (tanpa ENGINE/CHARSET MySQL-specific)
// agar kompatibel dengan PostgreSQL dan MySQL.
func (m *Migrator) ensureTrackingTable() error {
	sql := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version     VARCHAR(14) PRIMARY KEY,
		name        VARCHAR(255) NOT NULL,
		applied_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		checksum    VARCHAR(64) NOT NULL DEFAULT '',
		file_path   VARCHAR(500) NOT NULL DEFAULT ''
	);`

	return m.db.Exec(sql).Error
}

// getAppliedVersions mengembalikan daftar version yang sudah dijalankan (ASC).
func (m *Migrator) getAppliedVersions() ([]string, error) {
	var versions []string
	err := m.db.Model(&schemaMigration{}).
		Select("version").
		Order("version ASC").
		Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// getAppliedVersionsDesc mengembalikan daftar version yang sudah dijalankan (DESC).
func (m *Migrator) getAppliedVersionsDesc() ([]string, error) {
	var versions []string
	err := m.db.Model(&schemaMigration{}).
		Select("version").
		Order("version DESC").
		Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// =============================================================================
// FILE SCANNING
// =============================================================================

// scanMigrations membaca semua file .sql dan .down.sql dari embedded FS,
// memasangkan up dengan down, dan mengembalikan daftar migrasi terurut.
func (m *Migrator) scanMigrations() ([]Migration, error) {
	type upFile struct {
		Version  string
		Name     string
		FilePath string
		Content  string
	}

	type downFile struct {
		Version string
		Content string
	}

	var upFiles []upFile
	downFiles := make(map[string]string) // version → content

	err := fs.WalkDir(m.migrationsFS, m.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := strings.ToLower(d.Name())

		// Skip non-SQL files
		if !strings.HasSuffix(name, ".sql") {
			return nil
		}

		// === DOWN FILE ===
		if strings.HasSuffix(name, ".down.sql") {
			content, err := m.migrationsFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", path, err)
			}

			// Strip .down.sql → "001_create_companies"
			base := strings.TrimSuffix(d.Name(), ".down.sql")
			version, _ := parseFilenameFromBase(base)
			if version == "" {
				m.logger.Warn("Skipping down file without numeric prefix",
					zap.String("file", d.Name()))
				return nil
			}

			downFiles[version] = string(content)
			return nil
		}

		// === UP FILE ===
		// Skip .down.sql sudah ditangani di atas, jadi ini hanya file .sql biasa
		content, err := m.migrationsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		filename := filepath.Base(path)
		version, name := parseFilename(filename)

		if version == "" {
			m.logger.Warn("Skipping migration file without numeric prefix",
				zap.String("file", filename))
			return nil
		}

		// Hitung relative path dari rootDir
		relPath := path
		if strings.HasPrefix(relPath, m.rootDir) {
			relPath = strings.TrimPrefix(relPath, m.rootDir+"/")
		}

		upFiles = append(upFiles, upFile{
			Version:  version,
			Name:     name,
			FilePath: relPath,
			Content:  string(content),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Urutkan up files berdasarkan version (numerik)
	sort.Slice(upFiles, func(i, j int) bool {
		return upFiles[i].Version < upFiles[j].Version
	})

	// Build Migration list dengan down content
	migrations := make([]Migration, 0, len(upFiles))
	for _, uf := range upFiles {
		downContent, hasDown := downFiles[uf.Version]
		migrations = append(migrations, Migration{
			Version:     uf.Version,
			Name:        uf.Name,
			FilePath:    uf.FilePath,
			Content:     uf.Content,
			DownContent: downContent,
			HasDown:     hasDown,
		})
	}

	return migrations, nil
}

// parseFilename mengekstrak version dan nama dari filename up.
// Contoh: "001_create_companies.sql" → version="001", name="create_companies"
func parseFilename(filename string) (version, name string) {
	base := strings.TrimSuffix(filename, ".sql")
	return parseFilenameFromBase(base)
}

// parseFilenameFromBase mengekstrak version dan nama dari base name (tanpa extension).
// Contoh: "001_create_companies" → version="001", name="create_companies"
func parseFilenameFromBase(base string) (version, name string) {
	// Cari separator pertama (underscore setelah prefix numerik)
	parts := strings.SplitN(base, "_", 2)
	if len(parts) < 2 {
		return "", ""
	}

	version = parts[0]
	name = parts[1]

	// Validasi version numerik
	for _, c := range version {
		if c < '0' || c > '9' {
			return "", ""
		}
	}

	return version, name
}
