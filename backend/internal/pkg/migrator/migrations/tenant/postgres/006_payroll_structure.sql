-- =============================================================================
-- Tenant Migration: 006_payroll_structure
-- =============================================================================
-- Tabel untuk struktur payroll tenant (salary components, profiles, BPJS, PPh21).

-- ---------------------------------------------------------------------------
-- 6.1 Salary Components
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_components (
    id                          CHAR(36) PRIMARY KEY,
    code                        VARCHAR(50) NOT NULL,
    name                        VARCHAR(150) NOT NULL,
    description                 TEXT NULL,
    component_type              VARCHAR(255) NOT NULL,
    calculation_type            VARCHAR(255) NOT NULL DEFAULT 'FIXED',
    is_taxable                  SMALLINT NOT NULL DEFAULT 1,
    is_bpjs_base                SMALLINT NOT NULL DEFAULT 0,
    is_recurring                SMALLINT NOT NULL DEFAULT 1,
    is_proratable               SMALLINT NOT NULL DEFAULT 1,
    print_on_salary_structure   SMALLINT NOT NULL DEFAULT 1,
    display_order               INT NOT NULL DEFAULT 100,
    status                      VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_salary_comp_code UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_salary_comp_type ON salary_components (component_type);

CREATE INDEX IF NOT EXISTS idx_salary_comp_status ON salary_components (status);

-- ---------------------------------------------------------------------------
-- 6.2 Salary Grade Components
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_grade_components (
    id                    CHAR(36) PRIMARY KEY,
    grading_id            CHAR(36) NULL,
    salary_component_id   CHAR(36) NOT NULL,
    amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code         CHAR(3) NOT NULL DEFAULT 'IDR',
    effective_start_date  DATE NOT NULL,
    effective_end_date    DATE NULL,
    is_mandatory          SMALLINT NOT NULL DEFAULT 1,
    is_default            SMALLINT NOT NULL DEFAULT 1,
    status                VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                 TEXT NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_grade_comp_start UNIQUE (grading_id, salary_component_id, effective_start_date),

    CONSTRAINT fk_grade_comp_grading   FOREIGN KEY (grading_id)          REFERENCES gradings(id)          ON DELETE SET NULL,
    CONSTRAINT fk_grade_comp_component FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_grade_comp_lookup ON salary_grade_components (grading_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_grade_comp_component ON salary_grade_components (salary_component_id);

-- ---------------------------------------------------------------------------
-- 6.3 Employee Salary Components
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_employee_components (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    employment_id         CHAR(36) NULL,
    position_id           CHAR(36) NULL,
    grading_id            CHAR(36) NULL,
    salary_component_id   CHAR(36) NOT NULL,
    amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code         CHAR(3) NOT NULL DEFAULT 'IDR',
    source_type           VARCHAR(255) NOT NULL DEFAULT 'MANUAL',
    source_ref_id         CHAR(36) NULL,
    effective_start_date  DATE NOT NULL,
    effective_end_date    DATE NULL,
    notes                 TEXT NULL,
    status                VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_emp_comp_start UNIQUE (employee_id, salary_component_id, effective_start_date),

    CONSTRAINT fk_emp_comp_employee   FOREIGN KEY (employee_id)         REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_emp_comp_employment FOREIGN KEY (employment_id)       REFERENCES employments(id)       ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_position   FOREIGN KEY (position_id)         REFERENCES positions(id)         ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_grading    FOREIGN KEY (grading_id)          REFERENCES gradings(id)          ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_component  FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_emp_comp_effective ON salary_employee_components (employee_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_emp_comp_component ON salary_employee_components (salary_component_id);

CREATE INDEX IF NOT EXISTS idx_emp_comp_employment ON salary_employee_components (employment_id);

CREATE INDEX IF NOT EXISTS idx_emp_comp_position ON salary_employee_components (position_id);

CREATE INDEX IF NOT EXISTS idx_emp_comp_grading ON salary_employee_components (grading_id);

CREATE INDEX IF NOT EXISTS idx_emp_comp_source ON salary_employee_components (source_type, source_ref_id);

-- ---------------------------------------------------------------------------
-- 6.4 Salary Change Logs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_change_logs (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_salary_component_id CHAR(36) NULL,
    salary_component_id         CHAR(36) NULL,
    action_type                 VARCHAR(255) NOT NULL,
    old_amount                  DECIMAL(18, 2) NULL,
    new_amount                  DECIMAL(18, 2) NULL,
    old_effective_start_date    DATE NULL,
    new_effective_start_date    DATE NULL,
    old_effective_end_date      DATE NULL,
    new_effective_end_date      DATE NULL,
    reason                      VARCHAR(255) NULL,
    notes                       TEXT NULL,
    before_json                 JSON NULL,
    after_json                  JSON NULL,
    changed_by                  CHAR(36) NULL,
    changed_at                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_salary_change_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_salary_change_employee ON salary_change_logs (employee_id, changed_at);

CREATE INDEX IF NOT EXISTS idx_salary_change_component ON salary_change_logs (salary_component_id);

CREATE INDEX IF NOT EXISTS idx_salary_change_action ON salary_change_logs (action_type);

-- ---------------------------------------------------------------------------
-- 6.5 Salary Employee Adjustments (One-time adjustments per period)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_employee_adjustments (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    employment_id         CHAR(36) NULL,
    position_id           CHAR(36) NULL,
    salary_component_id   CHAR(36) NOT NULL,
    period_year           INT NOT NULL,
    period_month          SMALLINT NOT NULL,
    amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code         CHAR(3) NOT NULL DEFAULT 'IDR',
    source_type           VARCHAR(255) NOT NULL DEFAULT 'MANUAL',
    reason                VARCHAR(255) NULL,
    notes                 TEXT NULL,
    status                VARCHAR(255) NOT NULL DEFAULT 'DRAFT',
    approved_by           CHAR(36) NULL,
    approved_at           TIMESTAMP,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_adj_employee   FOREIGN KEY (employee_id)         REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_adj_employment FOREIGN KEY (employment_id)       REFERENCES employments(id)       ON DELETE SET NULL,
    CONSTRAINT fk_adj_position   FOREIGN KEY (position_id)         REFERENCES positions(id)         ON DELETE SET NULL,
    CONSTRAINT fk_adj_component  FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_adj_employee_period ON salary_employee_adjustments (employee_id, period_year, period_month, status);

CREATE INDEX IF NOT EXISTS idx_adj_component ON salary_employee_adjustments (salary_component_id);

CREATE INDEX IF NOT EXISTS idx_adj_employment ON salary_employee_adjustments (employment_id);

CREATE INDEX IF NOT EXISTS idx_adj_position ON salary_employee_adjustments (position_id);

-- ---------------------------------------------------------------------------
-- 6.6 Payroll Periods
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_periods (
    id            CHAR(36) PRIMARY KEY,
    period_code   VARCHAR(50) NOT NULL,
    period_year   INT NOT NULL,
    period_month  SMALLINT NOT NULL,
    start_date    DATE NOT NULL,
    end_date      DATE NOT NULL,
    as_of_date    DATE NOT NULL,
    status        VARCHAR(255) NOT NULL DEFAULT 'OPEN',
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_period_code UNIQUE (period_code),
CONSTRAINT uk_period_year_month UNIQUE (period_year, period_month)
);

-- ---------------------------------------------------------------------------
-- 6.7 Employee Payroll Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_payroll_profiles (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    employment_id         CHAR(36) NULL,
    payroll_group_code    VARCHAR(50) NOT NULL DEFAULT 'MONTHLY',
    payroll_frequency     VARCHAR(255) NOT NULL DEFAULT 'MONTHLY',
    payment_method        VARCHAR(255) NOT NULL DEFAULT 'BANK_TRANSFER',
    salary_currency       CHAR(3) NOT NULL DEFAULT 'IDR',
    is_payroll_active     SMALLINT NOT NULL DEFAULT 1,
    effective_start_date  DATE NOT NULL,
    effective_end_date    DATE NULL,
    status                VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                 TEXT NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_payroll_profile_start UNIQUE (employee_id, effective_start_date),

    CONSTRAINT fk_payroll_profile_employee   FOREIGN KEY (employee_id)   REFERENCES employees(id)   ON DELETE CASCADE,
    CONSTRAINT fk_payroll_profile_employment FOREIGN KEY (employment_id) REFERENCES employments(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_payroll_profile_employee_effective ON employee_payroll_profiles (employee_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_payroll_profile_group ON employee_payroll_profiles (payroll_group_code, status);

CREATE INDEX IF NOT EXISTS idx_payroll_profile_employment ON employee_payroll_profiles (employment_id);

-- ---------------------------------------------------------------------------
-- 6.8 Employee Bank Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_bank_profiles (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_payroll_profile_id CHAR(36) NOT NULL,
    bank_code                   VARCHAR(50) NULL,
    bank_name                   VARCHAR(150) NOT NULL,
    bank_branch                 VARCHAR(150) NULL,
    bank_account_number         VARCHAR(100) NOT NULL,
    bank_account_holder_name    VARCHAR(255) NOT NULL,
    is_primary                  SMALLINT NOT NULL DEFAULT 1,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_bank_profile_start UNIQUE (employee_id, effective_start_date),

    CONSTRAINT fk_bank_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_bank_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_bank_profile_effective ON employee_bank_profiles (employee_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_bank_profile_primary ON employee_bank_profiles (employee_id, is_primary, status);

CREATE INDEX IF NOT EXISTS idx_bank_profile_payroll ON employee_bank_profiles (employee_payroll_profile_id);

-- ---------------------------------------------------------------------------
-- 6.9 Employee BPJS Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_bpjs_profiles (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_payroll_profile_id CHAR(36) NOT NULL,
    bpjs_health_active          SMALLINT NOT NULL DEFAULT 0,
    bpjs_health_no              VARCHAR(50) NULL,
    bpjs_health_registered_name VARCHAR(255) NULL,
    bpjs_tk_active              SMALLINT NOT NULL DEFAULT 0,
    bpjs_tk_no                  VARCHAR(50) NULL,
    bpjs_tk_registered_name     VARCHAR(255) NULL,
    jkk_risk_class              VARCHAR(255) NOT NULL DEFAULT 'LOW',
    pension_active              SMALLINT NOT NULL DEFAULT 1,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_bpjs_profile_start UNIQUE (employee_id, effective_start_date),

    CONSTRAINT fk_bpjs_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_bpjs_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_bpjs_profile_effective ON employee_bpjs_profiles (employee_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_bpjs_profile_payroll ON employee_bpjs_profiles (employee_payroll_profile_id);

-- ---------------------------------------------------------------------------
-- 6.10 Employee Tax Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_tax_profiles (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_payroll_profile_id CHAR(36) NOT NULL,
    npwp                        VARCHAR(50) NULL,
    npwp_registered_name        VARCHAR(255) NULL,
    ptkp_status                 VARCHAR(20) NULL,
    tax_method                  VARCHAR(255) NOT NULL DEFAULT 'GROSS',
    is_taxable                  SMALLINT NOT NULL DEFAULT 1,
    has_npwp                    SMALLINT NOT NULL DEFAULT 0,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_tax_profile_start UNIQUE (employee_id, effective_start_date),

    CONSTRAINT fk_tax_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_tax_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tax_profile_effective ON employee_tax_profiles (employee_id, effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_tax_profile_payroll ON employee_tax_profiles (employee_payroll_profile_id);

-- ---------------------------------------------------------------------------
-- 6.11 BPJS Settings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS bpjs_settings (
    id                       CHAR(36) PRIMARY KEY,
    setting_code             VARCHAR(50) NOT NULL,
    setting_name             VARCHAR(150) NOT NULL,
    base_source              VARCHAR(255) NOT NULL DEFAULT 'BPJS_BASE_COMPONENTS',
    health_max_base_amount   DECIMAL(18, 2) NULL,
    pension_max_base_amount  DECIMAL(18, 2) NULL,
    default_jkk_risk_class   VARCHAR(255) NOT NULL DEFAULT 'LOW',
    rounding_mode            VARCHAR(255) NOT NULL DEFAULT 'ROUND',
    effective_start_date     DATE NOT NULL,
    effective_end_date       DATE NULL,
    status                   VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                    TEXT NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_bpjs_setting_code UNIQUE (setting_code)
);

CREATE INDEX IF NOT EXISTS idx_bpjs_setting_effective ON bpjs_settings (effective_start_date, effective_end_date, status);

-- ---------------------------------------------------------------------------
-- 6.12 BPJS Rate Components
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS bpjs_rate_components (
    id                      CHAR(36) PRIMARY KEY,
    bpjs_setting_id         CHAR(36) NOT NULL,
    rate_code               VARCHAR(80) NOT NULL,
    rate_name               VARCHAR(180) NOT NULL,
    bpjs_program            VARCHAR(255) NOT NULL,
    paid_by                 VARCHAR(255) NOT NULL,
    salary_component_id     CHAR(36) NULL,
    rate_percent            DECIMAL(8, 4) NOT NULL DEFAULT 0,
    fixed_amount            DECIMAL(18, 2) NULL,
    min_base_amount         DECIMAL(18, 2) NULL,
    max_base_amount         DECIMAL(18, 2) NULL,
    jkk_risk_class          VARCHAR(255) NULL,
    is_employee_deduction   SMALLINT NOT NULL DEFAULT 0,
    is_employer_contribution SMALLINT NOT NULL DEFAULT 0,
    generate_to_payroll_item SMALLINT NOT NULL DEFAULT 1,
    print_on_payslip        SMALLINT NOT NULL DEFAULT 1,
    display_order           INT NOT NULL DEFAULT 0,
    effective_start_date    DATE NOT NULL,
    effective_end_date      DATE NULL,
    status                  VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                   TEXT NULL,
    created_by              CHAR(36) NULL,
    updated_by              CHAR(36) NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_bpjs_rate_code_start UNIQUE (rate_code, effective_start_date),

    CONSTRAINT fk_bpjs_rate_setting   FOREIGN KEY (bpjs_setting_id)     REFERENCES bpjs_settings(id)     ON DELETE CASCADE,
    CONSTRAINT fk_bpjs_rate_component FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_bpjs_rate_setting ON bpjs_rate_components (bpjs_setting_id);

CREATE INDEX IF NOT EXISTS idx_bpjs_rate_program ON bpjs_rate_components (bpjs_program, paid_by, status);

CREATE INDEX IF NOT EXISTS idx_bpjs_rate_effective ON bpjs_rate_components (effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_bpjs_rate_component ON bpjs_rate_components (salary_component_id);

-- ---------------------------------------------------------------------------
-- 6.13 PPh21 Settings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pph21_settings (
    id                                CHAR(36) PRIMARY KEY,
    setting_code                      VARCHAR(50) NOT NULL,
    setting_name                      VARCHAR(150) NOT NULL,
    calculation_method                VARCHAR(255) NOT NULL DEFAULT 'REGULAR_GROSS_ANNUALIZED',
    default_tax_method                VARCHAR(255) NOT NULL DEFAULT 'GROSS',
    pph21_component_id                CHAR(36) NOT NULL,
    occupational_expense_rate_percent DECIMAL(8, 4) NOT NULL DEFAULT 5,
    occupational_expense_max_monthly  DECIMAL(18, 2) NOT NULL DEFAULT 500000,
    occupational_expense_max_yearly   DECIMAL(18, 2) NOT NULL DEFAULT 6000000,
    deduct_bpjs_health_employee       SMALLINT NOT NULL DEFAULT 0,
    deduct_bpjs_jht_employee          SMALLINT NOT NULL DEFAULT 1,
    deduct_bpjs_jp_employee           SMALLINT NOT NULL DEFAULT 1,
    annualization_months              SMALLINT NOT NULL DEFAULT 12,
    pkp_rounding_unit                 DECIMAL(18, 2) NOT NULL DEFAULT 1000,
    non_npwp_multiplier_percent       DECIMAL(8, 4) NOT NULL DEFAULT 100,
    rounding_mode                     VARCHAR(255) NOT NULL DEFAULT 'ROUND',
    effective_start_date              DATE NOT NULL,
    effective_end_date                DATE NULL,
    status                            VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    notes                             TEXT NULL,
    created_by                        CHAR(36) NULL,
    updated_by                        CHAR(36) NULL,
    created_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_pph21_setting_code UNIQUE (setting_code),

    CONSTRAINT fk_pph21_setting_component FOREIGN KEY (pph21_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pph21_setting_effective ON pph21_settings (effective_start_date, effective_end_date, status);

CREATE INDEX IF NOT EXISTS idx_pph21_setting_component ON pph21_settings (pph21_component_id);

-- ---------------------------------------------------------------------------
-- 6.14 PPh21 PTKP Rates
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pph21_ptkp_rates (
    id                   CHAR(36) PRIMARY KEY,
    ptkp_status          VARCHAR(20) NOT NULL,
    description          VARCHAR(255) NULL,
    annual_amount        DECIMAL(18, 2) NOT NULL DEFAULT 0,
    effective_start_date DATE NOT NULL,
    effective_end_date   DATE NULL,
    status               VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_by           CHAR(36) NULL,
    updated_by           CHAR(36) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_pph21_ptkp_status_start UNIQUE (ptkp_status, effective_start_date)
);

CREATE INDEX IF NOT EXISTS idx_pph21_ptkp_effective ON pph21_ptkp_rates (effective_start_date, effective_end_date, status);

-- ---------------------------------------------------------------------------
-- 6.15 PPh21 Tax Brackets
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pph21_tax_brackets (
    id                   CHAR(36) PRIMARY KEY,
    bracket_order        INT NOT NULL,
    lower_bound          DECIMAL(18, 2) NOT NULL DEFAULT 0,
    upper_bound          DECIMAL(18, 2) NULL,
    rate_percent         DECIMAL(8, 4) NOT NULL DEFAULT 0,
    effective_start_date DATE NOT NULL,
    effective_end_date   DATE NULL,
    status               VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_by           CHAR(36) NULL,
    updated_by           CHAR(36) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_pph21_bracket_order_start UNIQUE (bracket_order, effective_start_date)
);

CREATE INDEX IF NOT EXISTS idx_pph21_bracket_effective ON pph21_tax_brackets (effective_start_date, effective_end_date, status);
