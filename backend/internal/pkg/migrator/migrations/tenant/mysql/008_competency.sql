-- =============================================================================
-- Tenant Migration: 008_competency
-- =============================================================================
-- Tabel untuk modul kompetensi tenant.

-- ---------------------------------------------------------------------------
-- 8.1 Competencies (Master kompetensi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competencies (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    field       VARCHAR(255) NULL,
    cluster     VARCHAR(255) NULL,
    definition  TEXT NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.2 Competence Values (Nilai kompetensi — legacy style)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competence_values (
    id          CHAR(36) PRIMARY KEY,
    type        ENUM('Pendidikan','Pengalaman Kerja','Psychological','Technical','Managerial',
                     'Communicating & Influencing Skill','Problem Solving & Decision Making',
                     'Jumlah Uang','Wewenang','Dampak pada Hasil Akhir (Memiliki Wewenang Keuangan)',
                     'Dampak pada Hasil Akhir (Tidak Memiliki Wewenang Keuangan)',
                     'Nilai Asset','Wewenang Asset','Total Bawahan','Lingkup Hubungan Kerja',
                     'Frekuensi Hubungan Kerja','Aktifitas Fisik','Lingkungan Kerja',
                     'Resiko/Bahaya') NULL,
    level       INT NULL,
    name        VARCHAR(255) NOT NULL,
    point       INT NULL,
    description VARCHAR(255) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.3 Competency Values (Nilai kompetensi — structured)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_values (
    id          CHAR(36) PRIMARY KEY,
    type        VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    level       TINYINT UNSIGNED NOT NULL,
    code        VARCHAR(255) NULL,
    description TEXT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_compval_type_level (type, level),
    INDEX idx_compval_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.4 Competency Events (Periode penilaian kompetensi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_events (
    id            CHAR(36) PRIMARY KEY,
    type          VARCHAR(20) NOT NULL COMMENT 'auto | manual',
    period_type   VARCHAR(20) NOT NULL COMMENT 'annual | semester | quarter',
    period_year   SMALLINT UNSIGNED NOT NULL,
    period_number TINYINT UNSIGNED NULL COMMENT 'semester 1-2 / quarter 1-4',
    status        VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'draft | active | closed',
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_comp_event_period (type, period_type, period_year, period_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.5 Competency Event Targets
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_event_targets (
    id                  CHAR(36) PRIMARY KEY,
    competency_event_id CHAR(36) NOT NULL,
    organization_id     CHAR(36) NOT NULL,
    employee_id         CHAR(36) NULL,
    missing_self        TINYINT(1) NOT NULL DEFAULT 0,
    missing_superior    TINYINT(1) NOT NULL DEFAULT 0,
    missing_peer        TINYINT(1) NOT NULL DEFAULT 0,
    missing_subordinate TINYINT(1) NOT NULL DEFAULT 0,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_comp_event_target (competency_event_id, organization_id),
    INDEX idx_comp_target_org (organization_id),
    INDEX idx_comp_target_employee (employee_id),

    CONSTRAINT fk_comp_target_event      FOREIGN KEY (competency_event_id) REFERENCES competency_events(id)  ON DELETE CASCADE,
    CONSTRAINT fk_comp_target_org        FOREIGN KEY (organization_id)     REFERENCES organizations(id)      ON DELETE CASCADE,
    CONSTRAINT fk_comp_target_employee   FOREIGN KEY (employee_id)         REFERENCES employees(id)          ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.6 Competency Scores (Skor penilaian per organisasi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_scores (
    id                        CHAR(36) PRIMARY KEY,
    organization_id           CHAR(36) NOT NULL,
    employee_id               CHAR(36) NULL,
    technical_gap_percentage  DECIMAL(6, 2) NOT NULL DEFAULT 0,
    managerial_gap_percentage DECIMAL(6, 2) NOT NULL DEFAULT 0,
    total_gap_percentage      DECIMAL(6, 2) NOT NULL DEFAULT 0,
    total_grade_percentage    DECIMAL(6, 2) NOT NULL DEFAULT 0,
    competency_event_id       CHAR(36) NULL,
    assessed_at               TIMESTAMP NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_comp_score_org (organization_id),
    INDEX idx_comp_score_employee (employee_id),
    INDEX idx_comp_score_event (competency_event_id),

    CONSTRAINT fk_comp_score_org     FOREIGN KEY (organization_id)   REFERENCES organizations(id)     ON DELETE CASCADE,
    CONSTRAINT fk_comp_score_employee FOREIGN KEY (employee_id)      REFERENCES employees(id)         ON DELETE SET NULL,
    CONSTRAINT fk_comp_score_event   FOREIGN KEY (competency_event_id) REFERENCES competency_events(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 8.7 Competency Score Details
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_score_details (
    id                    CHAR(36) PRIMARY KEY,
    competency_score_id   CHAR(36) NOT NULL,
    competency_id         CHAR(36) NOT NULL,
    type                  ENUM('technical', 'managerial') NOT NULL,
    standard_level        TINYINT UNSIGNED NULL,
    standard_weight       DECIMAL(6, 2) NOT NULL DEFAULT 0,
    employee_level        TINYINT UNSIGNED NULL,
    gap_percentage        DECIMAL(6, 2) NOT NULL DEFAULT 0,
    weighted_gap_percentage DECIMAL(6, 2) NOT NULL DEFAULT 0,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_comp_score_detail (competency_score_id, competency_id),
    INDEX idx_comp_detail_competency (competency_id),

    CONSTRAINT fk_comp_detail_score      FOREIGN KEY (competency_score_id) REFERENCES competency_scores(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_detail_competency FOREIGN KEY (competency_id)       REFERENCES competencies(id)       ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- FK dari migration 002: job_family_competencies → competencies
-- Ditambahkan di sini karena competencies table baru ada setelah migration 008.
-- ---------------------------------------------------------------------------
ALTER TABLE job_family_competencies
    ADD CONSTRAINT fk_jfc_competency
    FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE;
