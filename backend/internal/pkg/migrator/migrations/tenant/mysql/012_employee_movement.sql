-- =============================================================================
-- Tenant Migration: 012_employee_movement
-- =============================================================================
-- Tabel untuk Employee Movement & Career Management
-- Mencakup: Promosi, Demosi, Mutasi/Rotasi, Perpanjangan Kontrak (PKWT),
--           Perubahan Status, Pensiun, dan Offboarding

-- ---------------------------------------------------------------------------
-- 12.1 Employee Movements (Riwayat Pergerakan Karyawan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_movements (
    id                         CHAR(36) PRIMARY KEY,
    employee_id                CHAR(36) NOT NULL,
    movement_type              VARCHAR(50) NOT NULL COMMENT 'promotion, demotion, mutation, contract_extension, status_change, retirement, offboarding, other',
    from_employment_id         CHAR(36) NULL,
    to_employment_id           CHAR(36) NULL,
    from_organization_id       CHAR(36) NULL,
    to_organization_id         CHAR(36) NULL,
    from_position_id           CHAR(36) NULL,
    to_position_id             CHAR(36) NULL,
    from_employment_status_id  CHAR(36) NULL,
    to_employment_status_id    CHAR(36) NULL,
    reason                     TEXT NULL,
    decision_letter_number     VARCHAR(50) NOT NULL,
    decision_letter_date       DATE NOT NULL,
    effective_date             DATE NOT NULL,
    status                     VARCHAR(20) DEFAULT 'draft' COMMENT 'draft, approved, executed, cancelled',
    notes                      TEXT NULL,
    approved_by                CHAR(36) NULL,
    approved_at                TIMESTAMP NULL,
    executed_by                CHAR(36) NULL,
    executed_at                TIMESTAMP NULL,
    created_by                 CHAR(36) NULL,
    updated_by                 CHAR(36) NULL,
    created_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_emp_mvmt_employee (employee_id),
    INDEX idx_emp_mvmt_type (movement_type),
    INDEX idx_emp_mvmt_status (status),
    INDEX idx_emp_mvmt_effective (effective_date),
    INDEX idx_emp_mvmt_from_org (from_organization_id),
    INDEX idx_emp_mvmt_to_org (to_organization_id),
    INDEX idx_emp_mvmt_from_pos (from_position_id),
    INDEX idx_emp_mvmt_to_pos (to_position_id),

    CONSTRAINT fk_empmvmt_employee    FOREIGN KEY (employee_id)                REFERENCES employees(id)            ON DELETE CASCADE,
    CONSTRAINT fk_empmvmt_from_empl   FOREIGN KEY (from_employment_id)         REFERENCES employments(id)          ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_empl     FOREIGN KEY (to_employment_id)           REFERENCES employments(id)          ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_org    FOREIGN KEY (from_organization_id)       REFERENCES organizations(id)        ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_org      FOREIGN KEY (to_organization_id)         REFERENCES organizations(id)        ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_pos    FOREIGN KEY (from_position_id)           REFERENCES positions(id)            ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_pos      FOREIGN KEY (to_position_id)             REFERENCES positions(id)            ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_empst  FOREIGN KEY (from_employment_status_id)  REFERENCES employment_statuses(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_empst    FOREIGN KEY (to_employment_status_id)    REFERENCES employment_statuses(id)  ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 12.2 Employee Contracts (PKWT & Perjanjian Kerja)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_contracts (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    contract_number       VARCHAR(50) NOT NULL,
    contract_type         VARCHAR(20) NOT NULL COMMENT 'pkwt, pkwtt, daily, other',
    start_date            DATE NOT NULL,
    end_date              DATE NULL,
    extension_count       INT DEFAULT 0,
    previous_contract_id  CHAR(36) NULL,
    decision_letter_number VARCHAR(50) NULL,
    notes                 TEXT NULL,
    document_url          VARCHAR(255) NULL,
    status                VARCHAR(20) DEFAULT 'active' COMMENT 'active, expired, extended, terminated',
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_emp_ctrct_employee (employee_id),
    INDEX idx_emp_ctrct_status (status),
    INDEX idx_emp_ctrct_type (contract_type),
    INDEX idx_emp_ctrct_end_date (end_date),

    CONSTRAINT fk_empctrct_employee       FOREIGN KEY (employee_id)               REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_empctrct_previous       FOREIGN KEY (previous_contract_id)      REFERENCES employee_contracts(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
