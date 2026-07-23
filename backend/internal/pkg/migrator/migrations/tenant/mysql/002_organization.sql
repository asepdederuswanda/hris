-- =============================================================================
-- Tenant Migration: 002_organization
-- =============================================================================
-- Tabel untuk struktur organisasi tenant.
-- Setiap tenant memiliki database sendiri, sehingga tidak perlu company_id.
-- Semua primary key menggunakan CHAR(36) UUID.

-- ---------------------------------------------------------------------------
-- 2.1 Organization Summaries
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS organization_summaries (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(7) NOT NULL UNIQUE,
    decree_no   VARCHAR(20) NOT NULL UNIQUE,
    decree_date DATE NOT NULL,
    status      VARCHAR(20) NULL DEFAULT 'active',
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_orgsum_code (code),
    INDEX idx_orgsum_decree (decree_no),
    INDEX idx_orgsum_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 2.2 Organization Levels
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS organization_levels (
    id          CHAR(36) PRIMARY KEY,
    level_name  VARCHAR(30) NULL,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_orglevel_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 2.3 Zones
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS zones (
    id          CHAR(36) PRIMARY KEY,
    zone        VARCHAR(200) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_zones_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 2.4 Organizations (Tree Hierarchy)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS organizations (
    id                        CHAR(36) PRIMARY KEY,
    organization_summary_id   CHAR(36) NULL,
    code                      VARCHAR(2) NOT NULL,
    full_code                 VARCHAR(50) NOT NULL UNIQUE,
    nomenclature              VARCHAR(255) NOT NULL,
    description               VARCHAR(255) NULL,
    parent_id                 CHAR(36) NULL,
    zone_id                   CHAR(36) NULL,
    job_family_id             CHAR(36) NULL,
    grading_id                CHAR(36) NULL,
    level                     INT DEFAULT 0,
    sort_order                INT DEFAULT 0,
    created_by                CHAR(36) NULL,
    updated_by                CHAR(36) NULL,
    deleted_at                TIMESTAMP NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_orgs_parent (parent_id),
    INDEX idx_orgs_zone (zone_id),
    INDEX idx_orgs_job_family (job_family_id),
    INDEX idx_orgs_grading (grading_id),
    INDEX idx_orgs_full_code (full_code),
    INDEX idx_orgs_deleted_at (deleted_at),

    CONSTRAINT fk_orgs_summary FOREIGN KEY (organization_summary_id) REFERENCES organization_summaries(id) ON DELETE SET NULL,
    CONSTRAINT fk_orgs_parent  FOREIGN KEY (parent_id)               REFERENCES organizations(id)          ON DELETE SET NULL,
    CONSTRAINT fk_orgs_zone    FOREIGN KEY (zone_id)                 REFERENCES zones(id)                  ON DELETE SET NULL,
    CONSTRAINT fk_orgs_job_family FOREIGN KEY (job_family_id)        REFERENCES job_families(id)           ON DELETE SET NULL,
    CONSTRAINT fk_orgs_grading FOREIGN KEY (grading_id)              REFERENCES gradings(id)               ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 2.5 Job Family Competencies
-- NOTE: FK ke competencies(id) dihapus karena competencies ada di migration 008.
-- FK akan ditambahkan via ALTER TABLE di migration 008 (competency).
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_family_competencies (
    id            CHAR(36) PRIMARY KEY,
    job_family_id CHAR(36) NOT NULL,
    competency_id CHAR(36) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_jfc_family_competency (job_family_id, competency_id),
    INDEX idx_jfc_competency (competency_id),

    CONSTRAINT fk_jfc_job_family  FOREIGN KEY (job_family_id)  REFERENCES job_families(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 2.6 Positions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS positions (
    id                    CHAR(36) PRIMARY KEY,
    organization_id       CHAR(36) NOT NULL,
    job_family_id         CHAR(36) NULL,
    grading_id            CHAR(36) NULL,
    code                  VARCHAR(50) NOT NULL,
    title                 VARCHAR(200) NOT NULL,
    parent_position_id    CHAR(36) NULL,
    is_head               TINYINT(1) NOT NULL DEFAULT 0,
    headcount             INT NOT NULL DEFAULT 1,
    is_active             TINYINT(1) NOT NULL DEFAULT 1,
    effective_start_date  DATE NULL,
    effective_end_date    DATE NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_positions_org (organization_id),
    INDEX idx_positions_parent (parent_position_id),
    INDEX idx_positions_job_family (job_family_id),
    INDEX idx_positions_grading (grading_id),

    CONSTRAINT fk_positions_org       FOREIGN KEY (organization_id)    REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_positions_parent    FOREIGN KEY (parent_position_id) REFERENCES positions(id)    ON DELETE SET NULL,
    CONSTRAINT fk_positions_job_family FOREIGN KEY (job_family_id)     REFERENCES job_families(id)  ON DELETE SET NULL,
    CONSTRAINT fk_positions_grading   FOREIGN KEY (grading_id)         REFERENCES gradings(id)      ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
