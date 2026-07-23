-- Migration: 006_create_licenses
-- Database: Platform (Cross-Dialect)
-- Tabel untuk menyimpan lisensi per company

CREATE TABLE IF NOT EXISTS licenses (
    id             CHAR(36) PRIMARY KEY,
    company_id     CHAR(36) NOT NULL,
    license_key    VARCHAR(100) NOT NULL UNIQUE,
    plan_type      VARCHAR(50) NOT NULL,
    max_employees  INT NOT NULL DEFAULT 0,
    max_modules    INT NOT NULL DEFAULT 0,
    start_date     DATE NOT NULL,
    end_date       DATE NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP NULL,

    CONSTRAINT fk_licenses_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE INDEX idx_licenses_company ON licenses (company_id);
CREATE INDEX idx_licenses_key ON licenses (license_key);
CREATE INDEX idx_licenses_plan ON licenses (plan_type);
CREATE INDEX idx_licenses_status ON licenses (status);
CREATE INDEX idx_licenses_deleted_at ON licenses (deleted_at);
