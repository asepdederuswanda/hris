-- Migration: 004_create_modules
-- Database: Platform (MySQL)
-- Tabel untuk mendaftarkan modul yang tersedia di platform

CREATE TABLE IF NOT EXISTS modules (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    version     VARCHAR(20) NOT NULL,
    description TEXT NULL,
    is_core     TINYINT(1) NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    INDEX idx_modules_slug (slug),
    INDEX idx_modules_core (is_core),
    INDEX idx_modules_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
