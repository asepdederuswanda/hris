-- =============================================================================
-- Tenant Migration: 009_job_management
-- =============================================================================
-- Tabel untuk modul job management (analisis jabatan) tenant.

-- ---------------------------------------------------------------------------
-- 9.1 Job Management Titles (Jenis Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_titles (
    id            CHAR(36) PRIMARY KEY,
    name          VARCHAR(100) NULL,
    descriptions  TEXT NULL,
    status        TINYINT NULL,
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.2 Job Management Title Subs (Sub Jenis Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_title_subs (
    id                         CHAR(36) PRIMARY KEY,
    job_management_title_id    CHAR(36) NULL,
    job_management_title_name  VARCHAR(100) NULL,
    name                       VARCHAR(100) NULL,
    descriptions               TEXT NULL,
    status                     TINYINT NULL,
    created_by                 CHAR(36) NULL,
    updated_by                 CHAR(36) NULL,
    created_at                 TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.3 Job Management Values (Nilai Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_values (
    id                             CHAR(36) PRIMARY KEY,
    job_management_title_sub_id    CHAR(36) NULL,
    job_management_title_sub_name  VARCHAR(100) NULL,
    type                           VARCHAR(255) NOT NULL,
    level                          INT NULL,
    descriptions                   TEXT NULL,
    note                           TEXT NULL,
    sort                           INT NULL,
    created_by                     CHAR(36) NULL,
    updated_by                     CHAR(36) NULL,
    created_at                     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.4 Job Management Objectives (Tujuan Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_objectives (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NULL,
    nomenclature    VARCHAR(50) NOT NULL,
    full_code       VARCHAR(20) NOT NULL,
    objective       TEXT NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmo_organization (organization_id),
    CONSTRAINT fk_jmo_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.5 Job Management Identifications (Identitas Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_identifications (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NULL,
    nomenclature    VARCHAR(50) NOT NULL,
    full_code       VARCHAR(20) NOT NULL,
    grading_id      CHAR(36) NOT NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmi_organization (organization_id),
    INDEX idx_jmi_grading (grading_id),

    CONSTRAINT fk_jmi_org     FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    CONSTRAINT fk_jmi_grading FOREIGN KEY (grading_id)       REFERENCES gradings(id)      ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.6 Job Management Responsibilities (Tanggung Jawab)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_responsibilities (
    id                CHAR(36) PRIMARY KEY,
    organization_id   CHAR(36) NULL,
    nomenclature      VARCHAR(50) NOT NULL,
    full_code         VARCHAR(20) NOT NULL,
    main_task         TEXT NULL,
    activities        TEXT NULL,
    outputs           TEXT NULL,
    success_indicators TEXT NULL,
    created_by        CHAR(36) NULL,
    updated_by        CHAR(36) NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmr_organization (organization_id),
    CONSTRAINT fk_jmr_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.7 Job Management Education Experiences (Pendidikan & Pengalaman)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_education_experiences (
    id                                    CHAR(36) PRIMARY KEY,
    organization_id                       CHAR(36) NULL,
    nomenclature                          VARCHAR(50) NOT NULL,
    full_code                             VARCHAR(20) NOT NULL,
    job_management_value_education_id     CHAR(36) NULL,
    job_management_value_experience_id    CHAR(36) NULL,
    created_by                            CHAR(36) NULL,
    updated_by                            CHAR(36) NULL,
    created_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmee_organization (organization_id),
    CONSTRAINT fk_jmee_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.8 Job Management HR Authorities (Kewenangan SDM)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_hr_authorities (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NULL,
    nomenclature    VARCHAR(50) NOT NULL,
    full_code       VARCHAR(20) NOT NULL,
    description     TEXT NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmha_organization (organization_id),
    CONSTRAINT fk_jmha_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.9 Job Management Operational Authorities (Kewenangan Operasional)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_operational_authorities (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NULL,
    nomenclature    VARCHAR(50) NOT NULL,
    full_code       VARCHAR(20) NOT NULL,
    description     TEXT NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmoa_organization (organization_id),
    CONSTRAINT fk_jmoa_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.10 Job Management Working Activities (Aktivitas Kerja)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_working_activities (
    id                       CHAR(36) PRIMARY KEY,
    organization_id          CHAR(36) NULL,
    nomenclature             VARCHAR(50) NOT NULL,
    full_code                VARCHAR(20) NOT NULL,
    job_management_value_id  CHAR(36) NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmwa_organization (organization_id),
    CONSTRAINT fk_jmwa_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.11 Job Management Working Risks (Risiko Kerja)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_working_risks (
    id                                CHAR(36) PRIMARY KEY,
    organization_id                   CHAR(36) NULL,
    nomenclature                      VARCHAR(50) NOT NULL,
    full_code                         VARCHAR(20) NOT NULL,
    job_management_value_environment_id CHAR(36) NULL,
    job_management_value_hazard_id    CHAR(36) NULL,
    created_by                        CHAR(36) NULL,
    updated_by                        CHAR(36) NULL,
    created_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmwr_organization (organization_id),
    CONSTRAINT fk_jmwr_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.12 Job Management Relationships (Hubungan Kerja)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_relationships (
    id                                    CHAR(36) PRIMARY KEY,
    organization_id                       CHAR(36) NULL,
    nomenclature                          VARCHAR(50) NOT NULL,
    full_code                             VARCHAR(20) NOT NULL,
    job_management_value_relationship_id  CHAR(36) NULL,
    job_management_value_frequency_id     CHAR(36) NULL,
    created_by                            CHAR(36) NULL,
    updated_by                            CHAR(36) NULL,
    created_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmrel_organization (organization_id),
    CONSTRAINT fk_jmrel_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.13 Job Management Subordinate Controls (Bawahan yang Dikendalikan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_subordinate_controls (
    id                       CHAR(36) PRIMARY KEY,
    organization_id          CHAR(36) NULL,
    nomenclature             VARCHAR(50) NOT NULL,
    full_code                VARCHAR(20) NOT NULL,
    job_management_value_id  CHAR(36) NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmsc_organization (organization_id),
    CONSTRAINT fk_jmsc_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.14 Job Management Assets (Aset Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_assets (
    id                              CHAR(36) PRIMARY KEY,
    organization_id                 CHAR(36) NULL,
    nomenclature                    VARCHAR(50) NOT NULL,
    full_code                       VARCHAR(20) NOT NULL,
    job_management_value_asset_id   CHAR(36) NULL,
    job_management_value_authority_id CHAR(36) NULL,
    created_by                      CHAR(36) NULL,
    updated_by                      CHAR(36) NULL,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmast_organization (organization_id),
    CONSTRAINT fk_jmast_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.15 Job Management Financials (Keuangan Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_financials (
    id                              CHAR(36) PRIMARY KEY,
    organization_id                 CHAR(36) NULL,
    nomenclature                    VARCHAR(50) NOT NULL,
    full_code                       VARCHAR(20) NOT NULL,
    is_authorized                   TINYINT(1) NOT NULL DEFAULT 0,
    job_management_value_cash_id    CHAR(36) NULL,
    job_management_value_authority_id CHAR(36) NULL,
    job_management_value_impact_id  CHAR(36) NULL,
    created_by                      CHAR(36) NULL,
    updated_by                      CHAR(36) NULL,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmfin_organization (organization_id),
    CONSTRAINT fk_jmfin_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.16 Job Management Potency Competencies (Kompetensi Potensi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_potency_competencies (
    id                       CHAR(36) PRIMARY KEY,
    organization_id          CHAR(36) NULL,
    job_management_value_id  CHAR(36) NULL,
    competency_id            CHAR(36) NULL,
    weight                   DECIMAL(8, 2) NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmpc_organization (organization_id),
    INDEX idx_jmpc_competency (competency_id),

    CONSTRAINT fk_jmpc_org        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    CONSTRAINT fk_jmpc_competency FOREIGN KEY (competency_id)   REFERENCES competencies(id)   ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.17 Job Management Scores (Skor Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_scores (
    id                         CHAR(36) PRIMARY KEY,
    organization_id            CHAR(36) NOT NULL UNIQUE,
    job_value_with_financial   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    job_value_without_financial BIGINT UNSIGNED NOT NULL DEFAULT 0,
    has_financial_authority    TINYINT(1) NOT NULL DEFAULT 0,
    components                 JSON NULL,
    sub_component_points       JSON NULL,
    calculated_at              TIMESTAMP NULL,
    created_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_jmscore_org (organization_id),

    CONSTRAINT fk_jmscore_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 9.18 Job Management Competency Groups (Bobot Kompetensi per Organisasi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_competency_groups (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    category        ENUM('technical', 'managerial') NOT NULL,
    weight          DECIMAL(8, 2) NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_jmcg_org_category (organization_id, category),

    CONSTRAINT fk_jmcg_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
