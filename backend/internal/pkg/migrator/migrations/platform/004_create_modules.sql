-- Migration: 004_create_modules
-- Database: Platform (Cross-Dialect)
-- Tabel untuk mendaftarkan modul yang tersedia di platform

CREATE TABLE IF NOT EXISTS modules (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    version     VARCHAR(20) NOT NULL,
    description TEXT NULL,
    is_core     SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);

CREATE INDEX idx_modules_slug ON modules (slug);
CREATE INDEX idx_modules_core ON modules (is_core);
CREATE INDEX idx_modules_deleted_at ON modules (deleted_at);
