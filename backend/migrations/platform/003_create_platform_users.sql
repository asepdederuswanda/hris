-- Migration: 003_create_platform_users
-- Database: Platform (MySQL)
-- Tabel untuk user platform (super admin, company admin)

CREATE TABLE IF NOT EXISTS platform_users (
    id              CHAR(36) PRIMARY KEY,
    company_id      CHAR(36) NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'company_admin' COMMENT 'super_admin | company_admin',
    is_active       TINYINT(1) NOT NULL DEFAULT 1,
    last_login_at   TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_pu_email (email),
    INDEX idx_pu_role (role),
    INDEX idx_pu_company (company_id),
    INDEX idx_pu_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Note: Super admin seeder ada di migrations/seeders/001_seed_super_admin.sql
