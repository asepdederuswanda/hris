-- =============================================================================
-- Tenant Migration: 001_master_data
-- =============================================================================
-- Tabel master data yang digunakan oleh semua modul tenant.
-- Setiap tenant memiliki database sendiri, sehingga tidak perlu company_id.
-- Semua primary key menggunakan CHAR(36) UUID.

-- ---------------------------------------------------------------------------
-- 1.1 Religions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS religions (
    id          CHAR(36) PRIMARY KEY,
    religion    VARCHAR(200) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.2 Educations
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS educations (
    id          CHAR(36) PRIMARY KEY,
    education   VARCHAR(200) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.3 Marital Statuses
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS marital_statuses (
    id              CHAR(36) PRIMARY KEY,
    marital_status  VARCHAR(100) NOT NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.4 Countries (reference data, no UUID needed)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS countries (
    id          CHAR(2) PRIMARY KEY,
    code        VARCHAR(2) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_countries_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.5 Provinces (Wilayah Administrasi Indonesia)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS provinces (
    id          CHAR(2) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.6 Regencies / Cities (Kabupaten/Kota)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS regencies (
    id          CHAR(4) PRIMARY KEY,
    province_id CHAR(2) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_regencies_province (province_id),
    CONSTRAINT fk_regencies_province FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.7 Districts (Kecamatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS districts (
    id          CHAR(6) PRIMARY KEY,
    regency_id  CHAR(4) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_districts_regency (regency_id),
    CONSTRAINT fk_districts_regency FOREIGN KEY (regency_id) REFERENCES regencies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.8 Villages (Kelurahan/Desa)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS villages (
    id          CHAR(10) PRIMARY KEY,
    district_id CHAR(6) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_villages_district (district_id),
    CONSTRAINT fk_villages_district FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.9 Relationship Types
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS relationship_types (
    id          CHAR(36) PRIMARY KEY,
    slug        VARCHAR(200) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.10 Employment Statuses
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employment_statuses (
    id              CHAR(36) PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    has_duration    TINYINT(1) NOT NULL DEFAULT 0,
    duration        INT NULL,
    duration_type   ENUM('days', 'months', 'years') NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.11 Gradings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gradings (
    id            CHAR(36) PRIMARY KEY,
    grading_name  VARCHAR(30) NULL,
    status        TINYINT NULL,
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.12 Job Families
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_families (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_job_families_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.13 PPh21 PTKP (Penghasilan Tidak Kena Pajak)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ptkps (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    ptkp        BIGINT NOT NULL,
    `group`     CHAR(1) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 1.14 Tarif Efektif Rata-rata (TER)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ters (
    id          CHAR(36) PRIMARY KEY,
    `group`     CHAR(1) NOT NULL,
    bruto_min   BIGINT NULL,
    bruto_max   BIGINT NULL,
    rate        DECIMAL(10, 2) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
