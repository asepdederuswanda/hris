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
    movement_type              VARCHAR(50) NOT NULL,
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
    status                     VARCHAR(20) DEFAULT 'draft',
    notes                      TEXT NULL,
    approved_by                CHAR(36) NULL,
    approved_at                TIMESTAMP NULL,
    executed_by                CHAR(36) NULL,
    executed_at                TIMESTAMP NULL,
    created_by                 CHAR(36) NULL,
    updated_by                 CHAR(36) NULL,
    created_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_empmvmt_employee    FOREIGN KEY (employee_id)                REFERENCES employees(id)            ON DELETE CASCADE,
    CONSTRAINT fk_empmvmt_from_empl   FOREIGN KEY (from_employment_id)         REFERENCES employments(id)          ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_empl     FOREIGN KEY (to_employment_id)           REFERENCES employments(id)          ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_org    FOREIGN KEY (from_organization_id)       REFERENCES organizations(id)        ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_org      FOREIGN KEY (to_organization_id)         REFERENCES organizations(id)        ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_pos    FOREIGN KEY (from_position_id)           REFERENCES positions(id)            ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_pos      FOREIGN KEY (to_position_id)             REFERENCES positions(id)            ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_from_empst  FOREIGN KEY (from_employment_status_id)  REFERENCES employment_statuses(id)  ON DELETE SET NULL,
    CONSTRAINT fk_empmvmt_to_empst    FOREIGN KEY (to_employment_status_id)    REFERENCES employment_statuses(id)  ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_emp_mvmt_employee ON employee_movements (employee_id);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_type ON employee_movements (movement_type);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_status ON employee_movements (status);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_effective ON employee_movements (effective_date);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_from_org ON employee_movements (from_organization_id);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_to_org ON employee_movements (to_organization_id);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_from_pos ON employee_movements (from_position_id);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_to_pos ON employee_movements (to_position_id);

-- ---------------------------------------------------------------------------
-- 12.2 Employee Contracts (PKWT & Perjanjian Kerja)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_contracts (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    contract_number       VARCHAR(50) NOT NULL,
    contract_type         VARCHAR(20) NOT NULL,
    start_date            DATE NOT NULL,
    end_date              DATE NULL,
    extension_count       INT DEFAULT 0,
    previous_contract_id  CHAR(36) NULL,
    decision_letter_number VARCHAR(50) NULL,
    notes                 TEXT NULL,
    document_url          VARCHAR(255) NULL,
    status                VARCHAR(20) DEFAULT 'active',
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_empctrct_employee       FOREIGN KEY (employee_id)               REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_empctrct_previous       FOREIGN KEY (previous_contract_id)      REFERENCES employee_contracts(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_emp_ctrct_employee ON employee_contracts (employee_id);
CREATE INDEX IF NOT EXISTS idx_emp_ctrct_status ON employee_contracts (status);
CREATE INDEX IF NOT EXISTS idx_emp_ctrct_type ON employee_contracts (contract_type);
CREATE INDEX IF NOT EXISTS idx_emp_ctrct_end_date ON employee_contracts (end_date);
