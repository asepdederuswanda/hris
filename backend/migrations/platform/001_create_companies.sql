-- Migration: 001_create_companies
-- Database: Platform (MySQL)
-- Tabel untuk menyimpan data perusahaan/tenant

CREATE TABLE IF NOT EXISTS companies (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    npwp        VARCHAR(16) NULL,
    nib         VARCHAR(25) NULL,
    address     TEXT NULL,
    email       VARCHAR(255) NULL,
    phone       VARCHAR(20) NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    INDEX idx_companies_slug (slug),
    INDEX idx_companies_status (status),
    INDEX idx_companies_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
