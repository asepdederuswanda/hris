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
    gender            ENUM('M','F') NULL,
    nationality_type  ENUM('WNI','WNA') NULL,
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
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_employees_employee_id (employee_id),
    INDEX idx_employees_nik (nik),
    INDEX idx_employees_name (name),
    INDEX idx_employees_religion (religion_id),
    INDEX idx_employees_marital (marital_status_id),
    INDEX idx_employees_nationality (nationality_id),

    CONSTRAINT fk_employees_religion       FOREIGN KEY (religion_id)       REFERENCES religions(id)       ON DELETE SET NULL,
    CONSTRAINT fk_employees_marital_status FOREIGN KEY (marital_status_id) REFERENCES marital_statuses(id) ON DELETE SET NULL,
    CONSTRAINT fk_employees_nationality    FOREIGN KEY (nationality_id)    REFERENCES countries(id)        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 3.2 Employee Addresses
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_addresses (
    id            CHAR(36) PRIMARY KEY,
    employee_id   CHAR(36) NULL,
    type          ENUM('MAIN','DOMICILE') NULL,
    address       VARCHAR(255) NULL,
    province_id   CHAR(2) NULL,
    regency_id    CHAR(4) NULL,
    district_id   CHAR(6) NULL,
    village_id    CHAR(10) NULL,
    postal_code   VARCHAR(5) NULL,
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empaddr_employee (employee_id),
    INDEX idx_empaddr_province (province_id),
    INDEX idx_empaddr_regency (regency_id),
    INDEX idx_empaddr_district (district_id),
    INDEX idx_empaddr_village (village_id),

    CONSTRAINT fk_empaddr_employee  FOREIGN KEY (employee_id)  REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_empaddr_province  FOREIGN KEY (province_id)  REFERENCES provinces(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_regency   FOREIGN KEY (regency_id)   REFERENCES regencies(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_district  FOREIGN KEY (district_id)  REFERENCES districts(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empaddr_village   FOREIGN KEY (village_id)   REFERENCES villages(id)   ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_emcontact_employee (employee_id),
    INDEX idx_emcontact_relationship (relationship_type_id),

    CONSTRAINT fk_emcontact_employee    FOREIGN KEY (employee_id)          REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_emcontact_relation    FOREIGN KEY (relationship_type_id) REFERENCES relationship_types(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empfam_employee (employee_id),
    INDEX idx_empfam_relationship (relationship_type_id),
    INDEX idx_empfam_education (education_id),

    CONSTRAINT fk_empfam_employee    FOREIGN KEY (employee_id)          REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_empfam_relation    FOREIGN KEY (relationship_type_id) REFERENCES relationship_types(id) ON DELETE SET NULL,
    CONSTRAINT fk_empfam_education   FOREIGN KEY (education_id)         REFERENCES educations(id)        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 3.5 Employee Educations
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_educations (
    id              CHAR(36) PRIMARY KEY,
    employee_id     CHAR(36) NULL,
    education_id    CHAR(36) NULL,
    name            VARCHAR(255) NOT NULL,
    major           VARCHAR(255) NULL,
    graduation_year YEAR NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empedu_employee (employee_id),
    INDEX idx_empedu_education (education_id),

    CONSTRAINT fk_empedu_employee  FOREIGN KEY (employee_id)  REFERENCES employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_empedu_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 3.6 Employee Experiences
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_experiences (
    id          CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NULL,
    company     VARCHAR(255) NOT NULL,
    position    VARCHAR(255) NULL,
    start_year  YEAR NULL,
    end_year    YEAR NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empexp_employee (employee_id),

    CONSTRAINT fk_empexp_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empdoc_employee (employee_id),

    CONSTRAINT fk_empdoc_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 3.8 Employee Insurances
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_insurances (
    id          CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NULL,
    category    ENUM('BPJS','Non BPJS') NULL,
    number      VARCHAR(100) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(100) NULL,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_empins_employee (employee_id),

    CONSTRAINT fk_empins_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_employments_employee (employee_id),
    INDEX idx_employments_org (organization_id),
    INDEX idx_employments_position (position_id),
    INDEX idx_employments_status (employment_status_id),

    CONSTRAINT fk_employments_employee FOREIGN KEY (employee_id)          REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_employments_org      FOREIGN KEY (organization_id)     REFERENCES organizations(id)      ON DELETE SET NULL,
    CONSTRAINT fk_employments_position FOREIGN KEY (position_id)         REFERENCES positions(id)          ON DELETE SET NULL,
    CONSTRAINT fk_employments_status   FOREIGN KEY (employment_status_id) REFERENCES employment_statuses(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
