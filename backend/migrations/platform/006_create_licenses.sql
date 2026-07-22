-- Migration: 006_create_licenses
-- Database: Platform (MySQL)
-- Tabel untuk menyimpan lisensi per company

CREATE TABLE IF NOT EXISTS licenses (
    id             CHAR(36) PRIMARY KEY,
    company_id     CHAR(36) NOT NULL,
    license_key    VARCHAR(100) NOT NULL UNIQUE,
    plan_type      VARCHAR(50) NOT NULL COMMENT 'free | basic | pro | enterprise',
    max_employees  INT NOT NULL DEFAULT 0 COMMENT '0 = unlimited',
    max_modules    INT NOT NULL DEFAULT 0 COMMENT '0 = unlimited',
    start_date     DATE NOT NULL,
    end_date       DATE NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'active | expired | suspended | cancelled',
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP NULL,

    INDEX idx_licenses_company (company_id),
    INDEX idx_licenses_key (license_key),
    INDEX idx_licenses_plan (plan_type),
    INDEX idx_licenses_status (status),
    INDEX idx_licenses_deleted_at (deleted_at),

    CONSTRAINT fk_licenses_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
