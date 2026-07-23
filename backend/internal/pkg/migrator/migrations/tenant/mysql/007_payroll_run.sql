-- =============================================================================
-- Tenant Migration: 007_payroll_run
-- =============================================================================
-- Tabel untuk proses payroll run tenant.

-- ---------------------------------------------------------------------------
-- 7.1 Payroll Runs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_runs (
    id                        CHAR(36) PRIMARY KEY,
    payroll_period_id         CHAR(36) NOT NULL,
    run_code                  VARCHAR(50) NOT NULL,
    run_type                  ENUM('REGULAR', 'OFF_CYCLE', 'THR', 'BONUS') NOT NULL DEFAULT 'REGULAR',
    status                    ENUM('DRAFT', 'CALCULATED', 'REVIEWED', 'APPROVED', 'LOCKED', 'CANCELLED') NOT NULL DEFAULT 'DRAFT',
    total_employees           INT NOT NULL DEFAULT 0,
    total_earning             DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_deduction           DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_employer_contribution DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_net                 DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_company_cost        DECIMAL(18, 2) NOT NULL DEFAULT 0,
    calculated_at             TIMESTAMP NULL,
    reviewed_at               TIMESTAMP NULL,
    approved_at               TIMESTAMP NULL,
    locked_at                 TIMESTAMP NULL,
    created_by                CHAR(36) NULL,
    updated_by                CHAR(36) NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_payroll_run_code (run_code),
    INDEX idx_payroll_run_period (payroll_period_id),
    INDEX idx_payroll_run_status (status),

    CONSTRAINT fk_payroll_run_period FOREIGN KEY (payroll_period_id) REFERENCES payroll_periods(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 7.2 Payroll Run Employees
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_run_employees (
    id                        CHAR(36) PRIMARY KEY,
    payroll_run_id            CHAR(36) NOT NULL,
    employee_id               CHAR(36) NOT NULL,
    employment_id             CHAR(36) NULL,
    position_id               CHAR(36) NULL,
    grading_id                CHAR(36) NULL,
    employee_code             VARCHAR(50) NOT NULL,
    employee_name             VARCHAR(255) NOT NULL,
    position_title            VARCHAR(200) NULL,
    grading_name              VARCHAR(255) NULL,
    total_earning             DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_deduction           DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_employer_contribution DECIMAL(18, 2) NOT NULL DEFAULT 0,
    net_amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_company_cost        DECIMAL(18, 2) NOT NULL DEFAULT 0,
    status                    ENUM('DRAFT', 'CALCULATED', 'EXCLUDED') NOT NULL DEFAULT 'DRAFT',
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_payroll_run_emp (payroll_run_id, employee_id),
    INDEX idx_payroll_run_emp_employee (employee_id),
    INDEX idx_payroll_run_emp_employment (employment_id),
    INDEX idx_payroll_run_emp_position (position_id),

    CONSTRAINT fk_payroll_run_emp_run       FOREIGN KEY (payroll_run_id) REFERENCES payroll_runs(id)       ON DELETE CASCADE,
    CONSTRAINT fk_payroll_run_emp_employee  FOREIGN KEY (employee_id)    REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_payroll_run_emp_employment FOREIGN KEY (employment_id) REFERENCES employments(id)        ON DELETE SET NULL,
    CONSTRAINT fk_payroll_run_emp_position  FOREIGN KEY (position_id)    REFERENCES positions(id)          ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 7.3 Payroll Run Items (detail per komponen gaji)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_run_items (
    id                      CHAR(36) PRIMARY KEY,
    payroll_run_id          CHAR(36) NOT NULL,
    payroll_run_employee_id CHAR(36) NOT NULL,
    employee_id             CHAR(36) NOT NULL,
    salary_component_id     CHAR(36) NOT NULL,
    component_code          VARCHAR(50) NOT NULL,
    component_name          VARCHAR(150) NOT NULL,
    component_type          ENUM('EARNING', 'DEDUCTION', 'EMPLOYER_CONTRIBUTION', 'INFORMATION') NOT NULL,
    item_category           ENUM('EMPLOYEE_EARNING', 'EMPLOYEE_DEDUCTION', 'EMPLOYER_CONTRIBUTION', 'INFORMATION') NOT NULL DEFAULT 'EMPLOYEE_EARNING',
    paid_by                 ENUM('EMPLOYEE', 'EMPLOYER', 'NONE') NOT NULL DEFAULT 'EMPLOYER',
    affects_gross_pay       TINYINT(1) NOT NULL DEFAULT 0,
    affects_net_pay         TINYINT(1) NOT NULL DEFAULT 0,
    affects_company_cost    TINYINT(1) NOT NULL DEFAULT 0,
    print_on_payslip        TINYINT(1) NOT NULL DEFAULT 1,
    amount                  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code           CHAR(3) NOT NULL DEFAULT 'IDR',
    source_group            ENUM('STRUCTURE', 'ADJUSTMENT', 'STATUTORY') NOT NULL,
    source_table            VARCHAR(100) NULL,
    source_id               CHAR(36) NULL,
    source_type             VARCHAR(50) NULL,
    notes                   TEXT NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_payroll_items_run_employee (payroll_run_id, employee_id),
    INDEX idx_payroll_items_run_emp (payroll_run_employee_id),
    INDEX idx_payroll_items_component (salary_component_id),

    CONSTRAINT fk_payroll_items_run         FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)           ON DELETE CASCADE,
    CONSTRAINT fk_payroll_items_run_emp     FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_payroll_items_component   FOREIGN KEY (salary_component_id)     REFERENCES salary_components(id)      ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 7.4 Payroll Payslips
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_payslips (
    id                        CHAR(36) PRIMARY KEY,
    payroll_run_id            CHAR(36) NOT NULL,
    payroll_run_employee_id   CHAR(36) NOT NULL,
    employee_id               CHAR(36) NOT NULL,
    payslip_number            VARCHAR(100) NOT NULL,
    period_year               INT NOT NULL,
    period_month              TINYINT NOT NULL,
    period_code               VARCHAR(50) NOT NULL,
    employee_code             VARCHAR(50) NOT NULL,
    employee_name             VARCHAR(255) NOT NULL,
    position_title            VARCHAR(200) NULL,
    grading_name              VARCHAR(255) NULL,
    total_earning             DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_deduction           DECIMAL(18, 2) NOT NULL DEFAULT 0,
    net_amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    status                    ENUM('DRAFT', 'PUBLISHED', 'CANCELLED') NOT NULL DEFAULT 'DRAFT',
    generated_at              TIMESTAMP NULL,
    published_at              TIMESTAMP NULL,
    cancelled_at              TIMESTAMP NULL,
    created_by                CHAR(36) NULL,
    updated_by                CHAR(36) NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_payslip_run_employee (payroll_run_id, employee_id),
    UNIQUE KEY uk_payslip_number (payslip_number),
    INDEX idx_payslip_employee_period (employee_id, period_year, period_month),
    INDEX idx_payslip_status (status),
    INDEX idx_payslip_run_employee (payroll_run_employee_id),

    CONSTRAINT fk_payslip_run         FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)          ON DELETE CASCADE,
    CONSTRAINT fk_payslip_run_employee FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_payslip_employee    FOREIGN KEY (employee_id)             REFERENCES employees(id)             ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 7.5 PPh21 Calculation Logs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pph21_calculation_logs (
    id                              CHAR(36) PRIMARY KEY,
    payroll_run_id                  CHAR(36) NOT NULL,
    payroll_run_employee_id         CHAR(36) NOT NULL,
    employee_id                     CHAR(36) NOT NULL,
    pph21_setting_id                CHAR(36) NOT NULL,
    employee_tax_profile_id         CHAR(36) NOT NULL,
    payroll_run_item_id             CHAR(36) NULL,
    tax_method                      VARCHAR(20) NOT NULL,
    ptkp_status                     VARCHAR(20) NOT NULL,
    has_npwp                        TINYINT(1) NOT NULL DEFAULT 1,
    gross_monthly                   DECIMAL(18, 2) NOT NULL DEFAULT 0,
    occupational_expense_monthly    DECIMAL(18, 2) NOT NULL DEFAULT 0,
    bpjs_tax_deductible_monthly     DECIMAL(18, 2) NOT NULL DEFAULT 0,
    net_monthly                     DECIMAL(18, 2) NOT NULL DEFAULT 0,
    net_annualized                  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    ptkp_annual                     DECIMAL(18, 2) NOT NULL DEFAULT 0,
    pkp_annual                      DECIMAL(18, 2) NOT NULL DEFAULT 0,
    annual_tax_before_npwp_mult     DECIMAL(18, 2) NOT NULL DEFAULT 0,
    non_npwp_multiplier_percent     DECIMAL(8, 4) NOT NULL DEFAULT 100,
    annual_tax_after_npwp_mult      DECIMAL(18, 2) NOT NULL DEFAULT 0,
    pph21_monthly                   DECIMAL(18, 2) NOT NULL DEFAULT 0,
    formula_json                    JSON NULL,
    status                          ENUM('CALCULATED', 'SKIPPED', 'ERROR') NOT NULL DEFAULT 'CALCULATED',
    notes                           TEXT NULL,
    created_by                      CHAR(36) NULL,
    updated_by                      CHAR(36) NULL,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_pph21_log_run_employee (payroll_run_id, payroll_run_employee_id),
    INDEX idx_pph21_log_run (payroll_run_id),
    INDEX idx_pph21_log_employee (employee_id),
    INDEX idx_pph21_log_setting (pph21_setting_id),
    INDEX idx_pph21_log_tax_profile (employee_tax_profile_id),

    CONSTRAINT fk_pph21_log_run            FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)          ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_run_employee   FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_employee       FOREIGN KEY (employee_id)             REFERENCES employees(id)             ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_setting        FOREIGN KEY (pph21_setting_id)        REFERENCES pph21_settings(id)        ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_tax_profile    FOREIGN KEY (employee_tax_profile_id) REFERENCES employee_tax_profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 7.6 Payroll Profile Change Logs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_profile_change_logs (
    id                  CHAR(36) PRIMARY KEY,
    employee_id         CHAR(36) NOT NULL,
    profile_table       VARCHAR(100) NOT NULL,
    profile_record_id   CHAR(36) NULL,
    action_type         ENUM('CREATE', 'UPDATE', 'END', 'REACTIVATE', 'DELETE') NOT NULL,
    reason              VARCHAR(255) NULL,
    notes               TEXT NULL,
    before_json         JSON NULL,
    after_json          JSON NULL,
    changed_by          CHAR(36) NULL,
    changed_at          TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_payroll_profile_log_employee (employee_id, changed_at),
    INDEX idx_payroll_profile_log_profile (profile_table, profile_record_id),

    CONSTRAINT fk_payroll_profile_log_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
