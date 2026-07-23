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
    run_type                  VARCHAR(255) NOT NULL DEFAULT 'REGULAR',
    status                    VARCHAR(255) NOT NULL DEFAULT 'DRAFT',
    total_employees           INT NOT NULL DEFAULT 0,
    total_earning             DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_deduction           DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_employer_contribution DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_net                 DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_company_cost        DECIMAL(18, 2) NOT NULL DEFAULT 0,
    calculated_at             TIMESTAMP,
    reviewed_at               TIMESTAMP,
    approved_at               TIMESTAMP,
    locked_at                 TIMESTAMP,
    created_by                CHAR(36) NULL,
    updated_by                CHAR(36) NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_payroll_run_code UNIQUE (run_code),

    CONSTRAINT fk_payroll_run_period FOREIGN KEY (payroll_period_id) REFERENCES payroll_periods(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payroll_run_period ON payroll_runs (payroll_period_id);

CREATE INDEX IF NOT EXISTS idx_payroll_run_status ON payroll_runs (status);

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
    status                    VARCHAR(255) NOT NULL DEFAULT 'DRAFT',
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_payroll_run_emp UNIQUE (payroll_run_id, employee_id),

    CONSTRAINT fk_payroll_run_emp_run       FOREIGN KEY (payroll_run_id) REFERENCES payroll_runs(id)       ON DELETE CASCADE,
    CONSTRAINT fk_payroll_run_emp_employee  FOREIGN KEY (employee_id)    REFERENCES employees(id)          ON DELETE CASCADE,
    CONSTRAINT fk_payroll_run_emp_employment FOREIGN KEY (employment_id) REFERENCES employments(id)        ON DELETE SET NULL,
    CONSTRAINT fk_payroll_run_emp_position  FOREIGN KEY (position_id)    REFERENCES positions(id)          ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_payroll_run_emp_employee ON payroll_run_employees (employee_id);

CREATE INDEX IF NOT EXISTS idx_payroll_run_emp_employment ON payroll_run_employees (employment_id);

CREATE INDEX IF NOT EXISTS idx_payroll_run_emp_position ON payroll_run_employees (position_id);

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
    component_type          VARCHAR(255) NOT NULL,
    item_category           VARCHAR(255) NOT NULL DEFAULT 'EMPLOYEE_EARNING',
    paid_by                 VARCHAR(255) NOT NULL DEFAULT 'EMPLOYER',
    affects_gross_pay       SMALLINT NOT NULL DEFAULT 0,
    affects_net_pay         SMALLINT NOT NULL DEFAULT 0,
    affects_company_cost    SMALLINT NOT NULL DEFAULT 0,
    print_on_payslip        SMALLINT NOT NULL DEFAULT 1,
    amount                  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code           CHAR(3) NOT NULL DEFAULT 'IDR',
    source_group            VARCHAR(255) NOT NULL,
    source_table            VARCHAR(100) NULL,
    source_id               CHAR(36) NULL,
    source_type             VARCHAR(50) NULL,
    notes                   TEXT NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_payroll_items_run         FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)           ON DELETE CASCADE,
    CONSTRAINT fk_payroll_items_run_emp     FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_payroll_items_component   FOREIGN KEY (salary_component_id)     REFERENCES salary_components(id)      ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payroll_items_run_employee ON payroll_run_items (payroll_run_id, employee_id);

CREATE INDEX IF NOT EXISTS idx_payroll_items_run_emp ON payroll_run_items (payroll_run_employee_id);

CREATE INDEX IF NOT EXISTS idx_payroll_items_component ON payroll_run_items (salary_component_id);

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
    period_month              SMALLINT NOT NULL,
    period_code               VARCHAR(50) NOT NULL,
    employee_code             VARCHAR(50) NOT NULL,
    employee_name             VARCHAR(255) NOT NULL,
    position_title            VARCHAR(200) NULL,
    grading_name              VARCHAR(255) NULL,
    total_earning             DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_deduction           DECIMAL(18, 2) NOT NULL DEFAULT 0,
    net_amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    status                    VARCHAR(255) NOT NULL DEFAULT 'DRAFT',
    generated_at              TIMESTAMP,
    published_at              TIMESTAMP,
    cancelled_at              TIMESTAMP,
    created_by                CHAR(36) NULL,
    updated_by                CHAR(36) NULL,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_payslip_run_employee UNIQUE (payroll_run_id, employee_id),
CONSTRAINT uk_payslip_number UNIQUE (payslip_number),

    CONSTRAINT fk_payslip_run         FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)          ON DELETE CASCADE,
    CONSTRAINT fk_payslip_run_employee FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_payslip_employee    FOREIGN KEY (employee_id)             REFERENCES employees(id)             ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payslip_employee_period ON payroll_payslips (employee_id, period_year, period_month);

CREATE INDEX IF NOT EXISTS idx_payslip_status ON payroll_payslips (status);

CREATE INDEX IF NOT EXISTS idx_payslip_run_employee ON payroll_payslips (payroll_run_employee_id);

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
    has_npwp                        SMALLINT NOT NULL DEFAULT 1,
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
    status                          VARCHAR(255) NOT NULL DEFAULT 'CALCULATED',
    notes                           TEXT NULL,
    created_by                      CHAR(36) NULL,
    updated_by                      CHAR(36) NULL,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_pph21_log_run_employee UNIQUE (payroll_run_id, payroll_run_employee_id),

    CONSTRAINT fk_pph21_log_run            FOREIGN KEY (payroll_run_id)          REFERENCES payroll_runs(id)          ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_run_employee   FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_employee       FOREIGN KEY (employee_id)             REFERENCES employees(id)             ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_setting        FOREIGN KEY (pph21_setting_id)        REFERENCES pph21_settings(id)        ON DELETE CASCADE,
    CONSTRAINT fk_pph21_log_tax_profile    FOREIGN KEY (employee_tax_profile_id) REFERENCES employee_tax_profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pph21_log_run ON pph21_calculation_logs (payroll_run_id);

CREATE INDEX IF NOT EXISTS idx_pph21_log_employee ON pph21_calculation_logs (employee_id);

CREATE INDEX IF NOT EXISTS idx_pph21_log_setting ON pph21_calculation_logs (pph21_setting_id);

CREATE INDEX IF NOT EXISTS idx_pph21_log_tax_profile ON pph21_calculation_logs (employee_tax_profile_id);

-- ---------------------------------------------------------------------------
-- 7.6 Payroll Profile Change Logs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_profile_change_logs (
    id                  CHAR(36) PRIMARY KEY,
    employee_id         CHAR(36) NOT NULL,
    profile_table       VARCHAR(100) NOT NULL,
    profile_record_id   CHAR(36) NULL,
    action_type         VARCHAR(255) NOT NULL,
    reason              VARCHAR(255) NULL,
    notes               TEXT NULL,
    before_json         JSON NULL,
    after_json          JSON NULL,
    changed_by          CHAR(36) NULL,
    changed_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_payroll_profile_log_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payroll_profile_log_employee ON payroll_profile_change_logs (employee_id, changed_at);

CREATE INDEX IF NOT EXISTS idx_payroll_profile_log_profile ON payroll_profile_change_logs (profile_table, profile_record_id);
