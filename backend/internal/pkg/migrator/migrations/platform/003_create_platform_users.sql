-- Migration: 003_create_platform_users
-- Database: Platform (Cross-Dialect)
-- Tabel untuk user platform (super admin, company admin)

CREATE TABLE IF NOT EXISTS platform_users (
    id              CHAR(36) PRIMARY KEY,
    company_id      CHAR(36) NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'company_admin',
    is_active       SMALLINT NOT NULL DEFAULT 1,
    last_login_at   TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pu_email ON platform_users (email);
CREATE INDEX idx_pu_role ON platform_users (role);
CREATE INDEX idx_pu_company ON platform_users (company_id);
CREATE INDEX idx_pu_active ON platform_users (is_active);

-- Note: Super admin seeder ada di migrations/seeders/001_seed_super_admin.sql
