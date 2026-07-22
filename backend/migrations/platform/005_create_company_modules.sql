-- Migration: 005_create_company_modules
-- Database: Platform (MySQL)
-- Tabel relasi many-to-many antara company dan module (module yang diaktifkan)

CREATE TABLE IF NOT EXISTS company_modules (
    company_id   CHAR(36) NOT NULL,
    module_id    CHAR(36) NOT NULL,
    enabled      TINYINT(1) NOT NULL DEFAULT 1,
    activated_at TIMESTAMP NULL,

    PRIMARY KEY (company_id, module_id),

    INDEX idx_cm_company (company_id),
    INDEX idx_cm_module (module_id),
    INDEX idx_cm_enabled (enabled),

    CONSTRAINT fk_cm_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_cm_module  FOREIGN KEY (module_id)  REFERENCES modules(id)   ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
