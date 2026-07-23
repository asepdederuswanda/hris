-- Migration: 002_create_tenant_connections
-- Database: Platform (Cross-Dialect)
-- Tabel untuk menyimpan konfigurasi koneksi database per tenant

CREATE TABLE IF NOT EXISTS tenant_connections (
    id          CHAR(36) PRIMARY KEY,
    company_id  CHAR(36) NOT NULL UNIQUE,
    driver      VARCHAR(20) NOT NULL DEFAULT 'postgres',
    host        VARCHAR(255) NOT NULL DEFAULT 'localhost',
    port        INT NOT NULL DEFAULT 5432,
    db_name     VARCHAR(100) NOT NULL,
    username    VARCHAR(100) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    ssl_mode    VARCHAR(20) NOT NULL DEFAULT 'require',
    is_active   SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    CONSTRAINT fk_tc_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE INDEX idx_tc_company ON tenant_connections (company_id);
CREATE INDEX idx_tc_active ON tenant_connections (is_active);
CREATE INDEX idx_tc_deleted_at ON tenant_connections (deleted_at);
