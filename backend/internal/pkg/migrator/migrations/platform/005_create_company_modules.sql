-- Migration: 005_create_company_modules
-- Database: Platform (Cross-Dialect)
-- Tabel relasi many-to-many antara company dan module (module yang diaktifkan)

CREATE TABLE IF NOT EXISTS company_modules (
    company_id   CHAR(36) NOT NULL,
    module_id    CHAR(36) NOT NULL,
    enabled      SMALLINT NOT NULL DEFAULT 1,
    activated_at TIMESTAMP NULL,

    PRIMARY KEY (company_id, module_id),

    CONSTRAINT fk_cm_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_cm_module  FOREIGN KEY (module_id)  REFERENCES modules(id)   ON DELETE CASCADE
);

CREATE INDEX idx_cm_company ON company_modules (company_id);
CREATE INDEX idx_cm_module ON company_modules (module_id);
CREATE INDEX idx_cm_enabled ON company_modules (enabled);
