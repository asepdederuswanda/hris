-- =============================================================================
-- Tenant Migration: 011_settings
-- =============================================================================
-- Tabel untuk pengaturan tenant: fitur, role, permission, hari libur, template dokumen.

-- ---------------------------------------------------------------------------
-- 11.1 Features (Fitur yang tersedia)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS features (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL,
    "group"     VARCHAR(255) NULL,
    description VARCHAR(255) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ---------------------------------------------------------------------------
-- 11.2 Permissions (Izin)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS permissions (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    guard_name  VARCHAR(255) NOT NULL DEFAULT 'web',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_permission_name_guard UNIQUE (name, guard_name)
);

-- ---------------------------------------------------------------------------
-- 11.3 Roles (Peran)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS roles (
    id              CHAR(36) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    guard_name      VARCHAR(255) NOT NULL DEFAULT 'web',
    description     VARCHAR(255) NULL,
    is_default      SMALLINT NULL DEFAULT 0,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_role_name_guard UNIQUE (name, guard_name)
);

CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles (deleted_at);

-- ---------------------------------------------------------------------------
-- 11.4 Feature Permission (Relasi fitur ↔ permission)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS feature_permission (
    id              CHAR(36) PRIMARY KEY,
    feature_id      CHAR(36) NOT NULL,
    permission_id   CHAR(36) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_fp_feature    FOREIGN KEY (feature_id)    REFERENCES features(id)    ON DELETE CASCADE,
    CONSTRAINT fk_fp_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fp_feature ON feature_permission (feature_id);

CREATE INDEX IF NOT EXISTS idx_fp_permission ON feature_permission (permission_id);

-- ---------------------------------------------------------------------------
-- 11.5 Model Has Roles (User ↔ Role)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS model_has_roles (
    role_id     CHAR(36) NOT NULL,
    model_type  VARCHAR(255) NOT NULL,
    model_id    CHAR(36) NOT NULL,

    PRIMARY KEY (role_id, model_id, model_type),

    CONSTRAINT fk_mhr_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_mhr_model ON model_has_roles (model_id, model_type);

-- ---------------------------------------------------------------------------
-- 11.6 Model Has Permissions (User ↔ Permission langsung)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS model_has_permissions (
    permission_id CHAR(36) NOT NULL,
    model_type    VARCHAR(255) NOT NULL,
    model_id      CHAR(36) NOT NULL,

    PRIMARY KEY (permission_id, model_id, model_type),

    CONSTRAINT fk_mhp_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_mhp_model ON model_has_permissions (model_id, model_type);

-- ---------------------------------------------------------------------------
-- 11.7 Role Has Permissions (Role ↔ Permission)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS role_has_permissions (
    permission_id CHAR(36) NOT NULL,
    role_id       CHAR(36) NOT NULL,

    PRIMARY KEY (permission_id, role_id),

    CONSTRAINT fk_rhp_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_rhp_role       FOREIGN KEY (role_id)       REFERENCES roles(id)       ON DELETE CASCADE
);

-- ---------------------------------------------------------------------------
-- 11.8 Company Holidays (Hari Libur Perusahaan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS company_holidays (
    id            CHAR(36) PRIMARY KEY,
    holiday_date  DATE NOT NULL,
    name          VARCHAR(200) NOT NULL,
    description   TEXT NULL,
    is_active     SMALLINT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_company_holiday UNIQUE (holiday_date)
);

-- ---------------------------------------------------------------------------
-- 11.9 Document Templates (Template Dokumen)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS document_templates (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    type        VARCHAR(255) NOT NULL,
    content     TEXT NULL,
    is_active   SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
