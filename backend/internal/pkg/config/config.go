// Package config menyediakan konfigurasi aplikasi terpusat
// menggunakan Viper dengan dukungan file, environment variables,
// dan default values.
package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/inthros/hris-platform/internal/pkg/driver"
)

// Config adalah struktur utama konfigurasi aplikasi.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Cache    CacheConfig    `mapstructure:"cache"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	OTEL     OTELConfig     `mapstructure:"otel"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// Catatan: Kunci enkripsi AES-256-GCM dibaca langsung dari environment variable
// HRIS_ENCRYPTION_KEY oleh package internal/pkg/crypto, bukan dari config.
// Lihat crypto.go untuk detail format key.

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Mode         string `mapstructure:"mode"` // debug, release, test
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	// Driver database yang digunakan: "postgres" atau "mysql"
	Driver string `mapstructure:"driver"`

	PlatformHost     string `mapstructure:"platform_host"`
	PlatformPort     int    `mapstructure:"platform_port"`
	PlatformDB       string `mapstructure:"platform_db"`
	PlatformUser     string `mapstructure:"platform_user"`
	PlatformPassword string `mapstructure:"platform_password"`
	PlatformSSLMode  string `mapstructure:"platform_sslmode"`

	TenantHost        string `mapstructure:"tenant_host"`
	TenantPort        int    `mapstructure:"tenant_port"`
	TenantSuperUser   string `mapstructure:"tenant_super_user"`
	TenantSuperPass   string `mapstructure:"tenant_super_password"`
	TenantSSLMode     string `mapstructure:"tenant_sslmode"`

	// Platform connection pool (single DB, moderate pool)
	MaxOpenConns      int `mapstructure:"max_open_conns"`
	MaxIdleConns      int `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeMs int `mapstructure:"conn_max_lifetime_ms"`

	// Tenant connection pool (per-tenant, smaller pool to avoid connection storm)
	// Set lebih rendah dari platform pool karena ada N tenant.
	// Rekomendasi: 5-15 open, 2-5 idle, max lifetime 30-60 menit.
	TenantMaxOpenConns      int `mapstructure:"tenant_max_open_conns"`
	TenantMaxIdleConns      int `mapstructure:"tenant_max_idle_conns"`
	TenantConnMaxLifetimeMs int `mapstructure:"tenant_conn_max_lifetime_ms"`
	// MaxIdleTime: koneksi idle akan ditutup setelah durasi ini.
	// Berguna untuk mengurangi koneksi yang tidak terpakai saat traffic sepi.
	TenantConnMaxIdleTimeMs int `mapstructure:"tenant_conn_max_idle_time_ms"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// CacheConfig untuk distributed cache dengan Redis Pub/Sub.
type CacheConfig struct {
	// Default TTL untuk cache items (detik). Default: 300 (5 menit).
	DefaultTTL int `mapstructure:"default_ttl"`
	// KeyPrefix untuk namespacing cache keys. Default: "hris:cache".
	KeyPrefix string `mapstructure:"key_prefix"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`  // menit
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"` // jam
	Issuer          string `mapstructure:"issuer"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, console
	OutputPath string `mapstructure:"output_path"` // stdout atau file path
}

type OTELConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	ServiceName       string  `mapstructure:"service_name"`
	CollectorEndpoint string  `mapstructure:"collector_endpoint"`
	SampleRate        float64 `mapstructure:"sample_rate"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// Load membaca konfigurasi dari file dan environment variables.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 0. Load .env file ke OS environment (optional, untuk local development)
	// godotenv akan load file .env dan memasukkan variabel ke os.Environ()
	// sehingga Viper AutomaticEnv() bisa membacanya dengan prefix HRIS_
	// .env file bersifat opsional, skip jika tidak ditemukan
	_ = godotenv.Load()

	// Default values
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)

	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.platform_host", "localhost")
	v.SetDefault("database.platform_port", 5432)
	v.SetDefault("database.platform_db", "hris_platform")
	v.SetDefault("database.platform_user", "hris")
	v.SetDefault("database.platform_sslmode", "disable")

	v.SetDefault("database.tenant_host", "localhost")
	v.SetDefault("database.tenant_port", 5432)
	v.SetDefault("database.tenant_super_user", "postgres")
	v.SetDefault("database.tenant_sslmode", "disable")
	v.SetDefault("database.max_open_conns", 10)           // Platform: 10 cukup untuk single DB
	v.SetDefault("database.max_idle_conns", 5)            // Platform: 5 idle
	v.SetDefault("database.conn_max_lifetime_ms", 3600000) // Platform: 1 jam

	v.SetDefault("database.tenant_max_open_conns", 10)              // Per tenant: 10 (dengan 50 tenant = 500 total)
	v.SetDefault("database.tenant_max_idle_conns", 3)               // Per tenant: 3 idle
	v.SetDefault("database.tenant_conn_max_lifetime_ms", 1800000)   // Per tenant: 30 menit
	v.SetDefault("database.tenant_conn_max_idle_time_ms", 300000)   // Per tenant: 5 menit idle → close

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	v.SetDefault("cache.default_ttl", 300) // 5 menit
	v.SetDefault("cache.key_prefix", "hris:cache")

	v.SetDefault("jwt.access_token_ttl", 15)  // 15 menit
	v.SetDefault("jwt.refresh_token_ttl", 24) // 24 jam
	v.SetDefault("jwt.issuer", "hris-platform")

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.output_path", "stdout")

	v.SetDefault("otel.enabled", false)
	v.SetDefault("otel.service_name", "hris-platform")
	v.SetDefault("otel.sample_rate", 0.1)

	v.SetDefault("cors.allowed_origins", []string{"*"})
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.max_age", 86400)

	// Configuration file (YAML)
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}
	v.ReadInConfig()

	// Environment variables (highest priority)
	// Termasuk variabel dari .env file yang sudah di-load oleh godotenv.Load()
	v.SetEnvPrefix("HRIS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// PlatformDSN mengembalikan connection string untuk platform database
// sesuai driver yang dikonfigurasi.
func (c *DatabaseConfig) PlatformDSN() string {
	switch driver.Parse(c.Driver) {
	case driver.MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.PlatformUser, c.PlatformPassword,
			c.PlatformHost, c.PlatformPort,
			c.PlatformDB,
		)
	default: // postgres
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.PlatformHost, c.PlatformPort,
			c.PlatformUser, c.PlatformPassword,
			c.PlatformDB, c.PlatformSSLMode,
		)
	}
}

// TenantDSN mengembalikan connection string untuk tenant database
// sesuai driver yang dikonfigurasi.
func (c *DatabaseConfig) TenantDSN(dbName, dbUser, dbPassword string) string {
	switch driver.Parse(c.Driver) {
	case driver.MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword,
			c.TenantHost, c.TenantPort,
			dbName,
		)
	default: // postgres
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.TenantHost, c.TenantPort,
			dbUser, dbPassword,
			dbName, c.TenantSSLMode,
		)
	}
}

// SuperuserDSN mengembalikan connection string dengan user super
// (digunakan untuk provisioning tenant database).
func (c *DatabaseConfig) SuperuserDSN() string {
	switch driver.Parse(c.Driver) {
	case driver.MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			c.TenantSuperUser, c.TenantSuperPass,
			c.TenantHost, c.TenantPort,
		)
	default: // postgres
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
			c.TenantHost, c.TenantPort,
			c.TenantSuperUser, c.TenantSuperPass,
			c.TenantSSLMode,
		)
	}
}

// RedisAddr mengembalikan alamat Redis dalam format host:port.
func (c *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
