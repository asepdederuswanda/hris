-- =============================================================================
-- Tenant Migration: 003_employee
-- =============================================================================
-- Tabel untuk data kepegawaian tenant.

-- ---------------------------------------------------------------------------
-- 3.1 Employees (Core)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employees (
    id                CHAR(36) PRIMARY KEY,
    employee_id       VARCHAR(50) NOT NULL,
    nik               VARCHAR(16) NULL,
    family_id         VARCHAR(16) NULL,
    name              VARCHAR(255) NOT NULL,
    mother_name       VARCHAR(255) NULL,
    gender            VARCHAR(255) NULL,
    nationality_type  VARCHAR(255) NULL,
    nationality_id    CHAR(2) NULL,
    pob               VARCHAR(255) NULL,
    dob               DATE NULL,
    phone_number      VARCHAR(255) NULL,
    email             VARCHAR(255) NULL UNIQUE,
    linkedin          VARCHAR(255) NULL UNIQUE,
    ig                VARCHAR(255) NULL UNIQUE,
    profile_picture   VARCHAR(255) NULL,
    religion_id       CHAR(36) NULL,
    marital_status_id CHAR(36) NULL,
    status            VARCHAR(20) DEFAULT 'active',
    created_by        CHAR(36) NULL,
    updated_by        CHAR(36) NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_employees_religion       FOREIGN KEY (religion_id)       REFERENCES religions(id)       ON DELETE SET NULL,
    CONSTRAINT fk_employees_marital_status FOREIGN KEY (marital_status_id) REFERENCES marital_statuses(id) ON DELETE SET NULL,
    CONSTRAINT fk_employees_nationality    FOREIGN KEY (nationality_id)    REFERENCES countries(id)        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_employees_employee_id ON employees (employee_id);

CREATE INDEX IF NOT EXISTS idx_employees_nik ON employees (nik);

CREATE INDEX IF NOT EXISTS idx_employees_name ON employees (name);

CREATE INDEX IF NOT EXISTS idx_employees_religion ON employees (religion_id);

CREATE INDEX IF NOT EXISTS idx_employees_marital ON employees (marital_status_id);

CREATE INDEX IF NOT EXISTS idx_employees_nationality ON employees (nationality_id);

-- ---------------------------------------------------------------------------
-- 3.2 Employee Addresses
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_addresses (
    id            CHAR(36) PRIMARY KEY,
    employee_id   CHAR(36) NULL,
    type          VARCHAR(255) NULL,
    address       VARCHAR(255) NULL,
    province_id   CHAR(2) NULL,
    regency_id    CHAR(4) NULL,
    district_id   CHAR(6) NULL,
    village_id    CHAR(10) NULL,
    postal_code   VARCHAR(5) NULL,
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empaddr_employee  FOREIGN KEY (employee_id)  REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_empaddr_province  FOREIGN KEY (province_id)  REFERENCES provinces(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_regency   FOREIGN KEY (regency_id)   REFERENCES regencies(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_district  FOREIGN KEY (district_id)  REFERENCES districts(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_village   FOREIGN KEY (village_id)   REFERENCES villages(id)   ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_empaddr_employee ON employee_addresses (employee_id);

CREATE INDEX IF NOT EXISTS idx_empaddr_province ON employee_addresses (province_id);

CREATE INDEX IF NOT EXISTS idx_empaddr_regency ON employee_addresses (regency_id);

CREATE INDEX IF NOT EXISTS idx_empaddr_district ON employee_addresses (district_id);

CREATE INDEX IF NOT EXISTS idx_empaddr_village ON employee_addresses (village_id);

-- ---------------------------------------------------------------------------
-- 3.3 Emergency Contacts
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS emergency_contacts (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NULL,
    name                  VARCHAR(255) NOT NULL,
    relationship_type_id  CHAR(36) NULL,
    phone_number          VARCHAR(50) NOT NULL,
    address               VARCHAR(255) NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_emcontact_employee    FOREIGN KEY (employee_id)          REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_emcontact_relation    FOREIGN KEY (relationship_type_id) REFERENCES relationship_types(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_emcontact_employee ON emergency_contacts (employee_id);

CREATE INDEX IF NOT EXISTS idx_emcontact_relationship ON emergency_contacts (relationship_type_id);

-- ---------------------------------------------------------------------------
-- 3.4 Employee Families
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_families (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NULL,
    nik                   VARCHAR(16) NULL,
    name                  VARCHAR(255) NOT NULL,
    dob                   DATE NULL,
    relationship_type_id  CHAR(36) NULL,
    education_id          CHAR(36) NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empfam_employee    FOREIGN KEY (employee_id)          REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_empfam_relation    FOREIGN KEY (relationship_type_id) REFERENCES relationship_types(id) ON DELETE SET NULL,
    CONSTRAINT fk_empfam_education   FOREIGN KEY (education_id)         REFERENCES educations(id)        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_empfam_employee ON employee_families (employee_id);

CREATE INDEX IF NOT EXISTS idx_empfam_relationship ON employee_families (relationship_type_id);

CREATE INDEX IF NOT EXISTS idx_empfam_education ON employee_families (education_id);

-- ---------------------------------------------------------------------------
-- 3.5 Employee Educations
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_educations (
    id              CHAR(36) PRIMARY KEY,
    employee_id     CHAR(36) NULL,
    education_id    CHAR(36) NULL,
    name            VARCHAR(255) NOT NULL,
    major           VARCHAR(255) NULL,
    graduation_year INTEGER NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empedu_employee  FOREIGN KEY (employee_id)  REFERENCES employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_empedu_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_empedu_employee ON employee_educations (employee_id);

CREATE INDEX IF NOT EXISTS idx_empedu_education ON employee_educations (education_id);

-- ---------------------------------------------------------------------------
-- 3.6 Employee Experiences
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_experiences (
    id          CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NULL,
    company     VARCHAR(255) NOT NULL,
    position    VARCHAR(255) NULL,
    start_year  INTEGER NULL,
    end_year    INTEGER NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empexp_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_empexp_employee ON employee_experiences (employee_id);

-- ---------------------------------------------------------------------------
-- 3.7 Employee Documents
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_documents (
    id          CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NULL,
    name        VARCHAR(255) NOT NULL,
    file        VARCHAR(255) NOT NULL,
    note        VARCHAR(255) NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empdoc_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_empdoc_employee ON employee_documents (employee_id);

-- ---------------------------------------------------------------------------
-- 3.8 Employee Insurances
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_insurances (
    id          CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NULL,
    category    VARCHAR(255) NULL,
    number      VARCHAR(100) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(100) NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_empins_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_empins_employee ON employee_insurances (employee_id);

-- ---------------------------------------------------------------------------
-- 3.9 Employments (Riwayat Jabatan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employments (
    id                       CHAR(36) PRIMARY KEY,
    employee_id              CHAR(36) NULL,
    organization_id          CHAR(36) NULL,
    position_id              CHAR(36) NULL,
    employment_status_id     CHAR(36) NULL,
    decision_letter_number   VARCHAR(50) NOT NULL,
    decision_letter_date     DATE NOT NULL,
    effective_date           DATE NOT NULL,
    effective_end_date       DATE NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_employments_employee FOREIGN KEY (employee_id)          REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_employments_org      FOREIGN KEY (organization_id)     REFERENCES organizations(id)      ON DELETE SET NULL,
    CONSTRAINT fk_employments_position FOREIGN KEY (position_id)         REFERENCES positions(id)          ON DELETE SET NULL,
    CONSTRAINT fk_employments_status   FOREIGN KEY (employment_status_id) REFERENCES employment_statuses(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_employments_employee ON employments (employee_id);

CREATE INDEX IF NOT EXISTS idx_employments_org ON employments (organization_id);

CREATE INDEX IF NOT EXISTS idx_employments_position ON employments (position_id);

CREATE INDEX IF NOT EXISTS idx_employments_status ON employments (employment_status_id);
