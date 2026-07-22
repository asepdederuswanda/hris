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
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	OTEL     OTELConfig     `mapstructure:"otel"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

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
	MaxOpenConns      int    `mapstructure:"max_open_conns"`
	MaxIdleConns      int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeMs int    `mapstructure:"conn_max_lifetime_ms"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
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
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime_ms", 3600000)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

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
