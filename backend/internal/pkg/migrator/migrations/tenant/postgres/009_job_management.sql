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
    status        SMALLINT NULL,
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ---------------------------------------------------------------------------
-- 9.2 Job Management Title Subs (Sub Jenis Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_title_subs (
    id                         CHAR(36) PRIMARY KEY,
    job_management_title_id    CHAR(36) NULL,
    job_management_title_name  VARCHAR(100) NULL,
    name                       VARCHAR(100) NULL,
    descriptions               TEXT NULL,
    status                     SMALLINT NULL,
    created_by                 CHAR(36) NULL,
    updated_by                 CHAR(36) NULL,
    created_at                 TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
    updated_at                     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmo_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmo_organization ON job_management_objectives (organization_id);

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
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_jmi_org     FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    CONSTRAINT fk_jmi_grading FOREIGN KEY (grading_id)       REFERENCES gradings(id)      ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_jmi_organization ON job_management_identifications (organization_id);

CREATE INDEX IF NOT EXISTS idx_jmi_grading ON job_management_identifications (grading_id);

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
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmr_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmr_organization ON job_management_responsibilities (organization_id);

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
    updated_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmee_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmee_organization ON job_management_education_experiences (organization_id);

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
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmha_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmha_organization ON job_management_hr_authorities (organization_id);

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
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmoa_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmoa_organization ON job_management_operational_authorities (organization_id);

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
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmwa_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmwa_organization ON job_management_working_activities (organization_id);

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
    updated_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmwr_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmwr_organization ON job_management_working_risks (organization_id);

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
    updated_at                            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmrel_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmrel_organization ON job_management_relationships (organization_id);

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
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmsc_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmsc_organization ON job_management_subordinate_controls (organization_id);

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
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmast_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmast_organization ON job_management_assets (organization_id);

-- ---------------------------------------------------------------------------
-- 9.15 Job Management Financials (Keuangan Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_financials (
    id                              CHAR(36) PRIMARY KEY,
    organization_id                 CHAR(36) NULL,
    nomenclature                    VARCHAR(50) NOT NULL,
    full_code                       VARCHAR(20) NOT NULL,
    is_authorized                   SMALLINT NOT NULL DEFAULT 0,
    job_management_value_cash_id    CHAR(36) NULL,
    job_management_value_authority_id CHAR(36) NULL,
    job_management_value_impact_id  CHAR(36) NULL,
    created_by                      CHAR(36) NULL,
    updated_by                      CHAR(36) NULL,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_jmfin_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmfin_organization ON job_management_financials (organization_id);

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
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_jmpc_org        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    CONSTRAINT fk_jmpc_competency FOREIGN KEY (competency_id)   REFERENCES competencies(id)   ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_jmpc_organization ON job_management_potency_competencies (organization_id);

CREATE INDEX IF NOT EXISTS idx_jmpc_competency ON job_management_potency_competencies (competency_id);

-- ---------------------------------------------------------------------------
-- 9.17 Job Management Scores (Skor Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_scores (
    id                         CHAR(36) PRIMARY KEY,
    organization_id            CHAR(36) NOT NULL UNIQUE,
    job_value_with_financial   BIGINT NOT NULL DEFAULT 0,
    job_value_without_financial BIGINT NOT NULL DEFAULT 0,
    has_financial_authority    SMALLINT NOT NULL DEFAULT 0,
    components                 JSON NULL,
    sub_component_points       JSON NULL,
    calculated_at              TIMESTAMP,
    created_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_jmscore_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_jmscore_org ON job_management_scores (organization_id);

-- ---------------------------------------------------------------------------
-- 9.18 Job Management Competency Groups (Bobot Kompetensi per Organisasi)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_management_competency_groups (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    category        VARCHAR(255) NOT NULL,
    weight          DECIMAL(8, 2) NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_jmcg_org_category UNIQUE (organization_id, category),

    CONSTRAINT fk_jmcg_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);
