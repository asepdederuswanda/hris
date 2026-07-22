// Package module mendefinisikan kontrak Module SDK yang wajib
// diimplementasikan oleh setiap modul (platform & tenant).
package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TargetDB mendefinisikan jenis database target untuk sebuah modul.
type TargetDB string

const (
	TargetPlatform TargetDB = "platform"
	TargetTenant   TargetDB = "tenant"
)

// Menu merepresentasikan item menu sidebar untuk sebuah modul.
type Menu struct {
	Name     string `yaml:"name" json:"name"`
	Icon     string `yaml:"icon" json:"icon"`
	Route    string `yaml:"route" json:"route"`
	Parent   string `yaml:"parent,omitempty" json:"parent,omitempty"`
	Children []Menu `yaml:"children,omitempty" json:"children,omitempty"`
}

// ModuleInfo berisi metadata identitas sebuah modul.
type ModuleInfo struct {
	Name        string   `yaml:"name" json:"name"`
	Slug        string   `yaml:"slug" json:"slug"`
	Version     string   `yaml:"version" json:"version"`
	Description string   `yaml:"description" json:"description"`
	IsCore      bool     `yaml:"is_core" json:"is_core"`
	DependsOn   []string `yaml:"depends_on" json:"depends_on"`
	Permissions []string `yaml:"permissions" json:"permissions"`
	Menus       []Menu   `yaml:"menus" json:"menus"`
}

// Module adalah kontrak SDK yang wajib diimplementasikan setiap modul.
//
// Setiap modul bertanggung jawab atas:
//  1. Identitas diri (Info)
//  2. Pendaftaran route API (RegisterRoutes)
//  3. Migrasi database sendiri (Migrate)
//  4. Data awal / seeder (Seed)
//  5. Daftar permission (Permissions)
type Module interface {
	// Info mengembalikan metadata modul (nama, slug, versi, dll).
	Info() ModuleInfo

	// RegisterRoutes mendaftarkan semua endpoint HTTP modul
	// ke dalam router group yang disediakan.
	RegisterRoutes(router *gin.RouterGroup)

	// Migrate menjalankan migration database untuk modul ini.
	// Parameter db adalah koneksi ke database yang sesuai
	// (platform DB untuk modul platform, tenant DB untuk modul tenant).
	Migrate(db *gorm.DB) error

	// Seed menjalankan seeder data awal untuk modul ini.
	Seed(db *gorm.DB) error

	// Permissions mengembalikan daftar permission yang
	// dibutuhkan/disediakan oleh modul ini.
	Permissions() []string
}

// ModuleRegistration digunakan untuk mendaftarkan modul beserta
// tipe database yang digunakannya (platform atau tenant).
type ModuleRegistration struct {
	Module   Module
	TargetDB TargetDB
	Priority int // Urutan inisialisasi (semakin kecil semakin awal)
}
