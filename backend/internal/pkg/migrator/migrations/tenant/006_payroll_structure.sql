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
    component_type              ENUM('EARNING', 'DEDUCTION', 'EMPLOYER_CONTRIBUTION', 'INFORMATION') NOT NULL,
    calculation_type            ENUM('FIXED', 'MANUAL') NOT NULL DEFAULT 'FIXED',
    is_taxable                  TINYINT(1) NOT NULL DEFAULT 1,
    is_bpjs_base                TINYINT(1) NOT NULL DEFAULT 0,
    is_recurring                TINYINT(1) NOT NULL DEFAULT 1,
    is_proratable               TINYINT(1) NOT NULL DEFAULT 1,
    print_on_salary_structure   TINYINT(1) NOT NULL DEFAULT 1,
    display_order               INT NOT NULL DEFAULT 100,
    status                      ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_salary_comp_code (code),
    INDEX idx_salary_comp_type (component_type),
    INDEX idx_salary_comp_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    is_mandatory          TINYINT(1) NOT NULL DEFAULT 1,
    is_default            TINYINT(1) NOT NULL DEFAULT 1,
    status                ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                 TEXT NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_grade_comp_start (grading_id, salary_component_id, effective_start_date),
    INDEX idx_grade_comp_lookup (grading_id, effective_start_date, effective_end_date, status),
    INDEX idx_grade_comp_component (salary_component_id),

    CONSTRAINT fk_grade_comp_grading   FOREIGN KEY (grading_id)          REFERENCES gradings(id)          ON DELETE SET NULL,
    CONSTRAINT fk_grade_comp_component FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    source_type           ENUM('MANUAL', 'POSITION', 'GRADE', 'ORGANIZATION', 'COMPANY_POLICY', 'MIGRATION', 'IMPORT') NOT NULL DEFAULT 'MANUAL',
    source_ref_id         CHAR(36) NULL,
    effective_start_date  DATE NOT NULL,
    effective_end_date    DATE NULL,
    notes                 TEXT NULL,
    status                ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_emp_comp_start (employee_id, salary_component_id, effective_start_date),
    INDEX idx_emp_comp_effective (employee_id, effective_start_date, effective_end_date, status),
    INDEX idx_emp_comp_component (salary_component_id),
    INDEX idx_emp_comp_employment (employment_id),
    INDEX idx_emp_comp_position (position_id),
    INDEX idx_emp_comp_grading (grading_id),
    INDEX idx_emp_comp_source (source_type, source_ref_id),

    CONSTRAINT fk_emp_comp_employee   FOREIGN KEY (employee_id)         REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_emp_comp_employment FOREIGN KEY (employment_id)       REFERENCES employments(id)       ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_position   FOREIGN KEY (position_id)         REFERENCES positions(id)         ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_grading    FOREIGN KEY (grading_id)          REFERENCES gradings(id)          ON DELETE SET NULL,
    CONSTRAINT fk_emp_comp_component  FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.4 Salary Change Logs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS salary_change_logs (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_salary_component_id CHAR(36) NULL,
    salary_component_id         CHAR(36) NULL,
    action_type                 ENUM('CREATE', 'UPDATE', 'END', 'REACTIVATE', 'DELETE') NOT NULL,
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
    changed_at                  TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_salary_change_employee (employee_id, changed_at),
    INDEX idx_salary_change_component (salary_component_id),
    INDEX idx_salary_change_action (action_type),

    CONSTRAINT fk_salary_change_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    period_month          TINYINT NOT NULL,
    amount                DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code         CHAR(3) NOT NULL DEFAULT 'IDR',
    source_type           ENUM('MANUAL', 'IMPORT', 'CORRECTION', 'POLICY', 'MIGRATION') NOT NULL DEFAULT 'MANUAL',
    reason                VARCHAR(255) NULL,
    notes                 TEXT NULL,
    status                ENUM('DRAFT', 'APPROVED', 'CANCELLED', 'APPLIED') NOT NULL DEFAULT 'DRAFT',
    approved_by           CHAR(36) NULL,
    approved_at           TIMESTAMP NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_adj_employee_period (employee_id, period_year, period_month, status),
    INDEX idx_adj_component (salary_component_id),
    INDEX idx_adj_employment (employment_id),
    INDEX idx_adj_position (position_id),

    CONSTRAINT fk_adj_employee   FOREIGN KEY (employee_id)         REFERENCES employees(id)         ON DELETE CASCADE,
    CONSTRAINT fk_adj_employment FOREIGN KEY (employment_id)       REFERENCES employments(id)       ON DELETE SET NULL,
    CONSTRAINT fk_adj_position   FOREIGN KEY (position_id)         REFERENCES positions(id)         ON DELETE SET NULL,
    CONSTRAINT fk_adj_component  FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.6 Payroll Periods
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_periods (
    id            CHAR(36) PRIMARY KEY,
    period_code   VARCHAR(50) NOT NULL,
    period_year   INT NOT NULL,
    period_month  TINYINT NOT NULL,
    start_date    DATE NOT NULL,
    end_date      DATE NOT NULL,
    as_of_date    DATE NOT NULL,
    status        ENUM('OPEN', 'CLOSED') NOT NULL DEFAULT 'OPEN',
    created_by    CHAR(36) NULL,
    updated_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_period_code (period_code),
    UNIQUE KEY uk_period_year_month (period_year, period_month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.7 Employee Payroll Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_payroll_profiles (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    employment_id         CHAR(36) NULL,
    payroll_group_code    VARCHAR(50) NOT NULL DEFAULT 'MONTHLY',
    payroll_frequency     ENUM('MONTHLY', 'WEEKLY', 'DAILY') NOT NULL DEFAULT 'MONTHLY',
    payment_method        ENUM('BANK_TRANSFER', 'CASH', 'CHEQUE') NOT NULL DEFAULT 'BANK_TRANSFER',
    salary_currency       CHAR(3) NOT NULL DEFAULT 'IDR',
    is_payroll_active     TINYINT(1) NOT NULL DEFAULT 1,
    effective_start_date  DATE NOT NULL,
    effective_end_date    DATE NULL,
    status                ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                 TEXT NULL,
    created_by            CHAR(36) NULL,
    updated_by            CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_payroll_profile_start (employee_id, effective_start_date),
    INDEX idx_payroll_profile_employee_effective (employee_id, effective_start_date, effective_end_date, status),
    INDEX idx_payroll_profile_group (payroll_group_code, status),
    INDEX idx_payroll_profile_employment (employment_id),

    CONSTRAINT fk_payroll_profile_employee   FOREIGN KEY (employee_id)   REFERENCES employees(id)   ON DELETE CASCADE,
    CONSTRAINT fk_payroll_profile_employment FOREIGN KEY (employment_id) REFERENCES employments(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    is_primary                  TINYINT(1) NOT NULL DEFAULT 1,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bank_profile_start (employee_id, effective_start_date),
    INDEX idx_bank_profile_effective (employee_id, effective_start_date, effective_end_date, status),
    INDEX idx_bank_profile_primary (employee_id, is_primary, status),
    INDEX idx_bank_profile_payroll (employee_payroll_profile_id),

    CONSTRAINT fk_bank_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_bank_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.9 Employee BPJS Profiles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_bpjs_profiles (
    id                          CHAR(36) PRIMARY KEY,
    employee_id                 CHAR(36) NOT NULL,
    employee_payroll_profile_id CHAR(36) NOT NULL,
    bpjs_health_active          TINYINT(1) NOT NULL DEFAULT 0,
    bpjs_health_no              VARCHAR(50) NULL,
    bpjs_health_registered_name VARCHAR(255) NULL,
    bpjs_tk_active              TINYINT(1) NOT NULL DEFAULT 0,
    bpjs_tk_no                  VARCHAR(50) NULL,
    bpjs_tk_registered_name     VARCHAR(255) NULL,
    jkk_risk_class              ENUM('VERY_LOW', 'LOW', 'MEDIUM', 'HIGH', 'VERY_HIGH') NOT NULL DEFAULT 'LOW',
    pension_active              TINYINT(1) NOT NULL DEFAULT 1,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bpjs_profile_start (employee_id, effective_start_date),
    INDEX idx_bpjs_profile_effective (employee_id, effective_start_date, effective_end_date, status),
    INDEX idx_bpjs_profile_payroll (employee_payroll_profile_id),

    CONSTRAINT fk_bpjs_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_bpjs_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    tax_method                  ENUM('GROSS', 'GROSS_UP', 'NETT') NOT NULL DEFAULT 'GROSS',
    is_taxable                  TINYINT(1) NOT NULL DEFAULT 1,
    has_npwp                    TINYINT(1) NOT NULL DEFAULT 0,
    effective_start_date        DATE NOT NULL,
    effective_end_date          DATE NULL,
    status                      ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                       TEXT NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tax_profile_start (employee_id, effective_start_date),
    INDEX idx_tax_profile_effective (employee_id, effective_start_date, effective_end_date, status),
    INDEX idx_tax_profile_payroll (employee_payroll_profile_id),

    CONSTRAINT fk_tax_profile_employee FOREIGN KEY (employee_id)                 REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_tax_profile_payroll  FOREIGN KEY (employee_payroll_profile_id) REFERENCES employee_payroll_profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.11 BPJS Settings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS bpjs_settings (
    id                       CHAR(36) PRIMARY KEY,
    setting_code             VARCHAR(50) NOT NULL,
    setting_name             VARCHAR(150) NOT NULL,
    base_source              ENUM('BPJS_BASE_COMPONENTS', 'BASIC_SALARY', 'GROSS_EARNING') NOT NULL DEFAULT 'BPJS_BASE_COMPONENTS',
    health_max_base_amount   DECIMAL(18, 2) NULL,
    pension_max_base_amount  DECIMAL(18, 2) NULL,
    default_jkk_risk_class   ENUM('VERY_LOW', 'LOW', 'MEDIUM', 'HIGH', 'VERY_HIGH') NOT NULL DEFAULT 'LOW',
    rounding_mode            ENUM('NONE', 'ROUND', 'CEIL', 'FLOOR') NOT NULL DEFAULT 'ROUND',
    effective_start_date     DATE NOT NULL,
    effective_end_date       DATE NULL,
    status                   ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                    TEXT NULL,
    created_by               CHAR(36) NULL,
    updated_by               CHAR(36) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bpjs_setting_code (setting_code),
    INDEX idx_bpjs_setting_effective (effective_start_date, effective_end_date, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.12 BPJS Rate Components
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS bpjs_rate_components (
    id                      CHAR(36) PRIMARY KEY,
    bpjs_setting_id         CHAR(36) NOT NULL,
    rate_code               VARCHAR(80) NOT NULL,
    rate_name               VARCHAR(180) NOT NULL,
    bpjs_program            ENUM('HEALTH', 'JHT', 'JP', 'JKK', 'JKM', 'JKP') NOT NULL,
    paid_by                 ENUM('EMPLOYEE', 'EMPLOYER') NOT NULL,
    salary_component_id     CHAR(36) NULL,
    rate_percent            DECIMAL(8, 4) NOT NULL DEFAULT 0,
    fixed_amount            DECIMAL(18, 2) NULL,
    min_base_amount         DECIMAL(18, 2) NULL,
    max_base_amount         DECIMAL(18, 2) NULL,
    jkk_risk_class          ENUM('VERY_LOW', 'LOW', 'MEDIUM', 'HIGH', 'VERY_HIGH') NULL,
    is_employee_deduction   TINYINT(1) NOT NULL DEFAULT 0,
    is_employer_contribution TINYINT(1) NOT NULL DEFAULT 0,
    generate_to_payroll_item TINYINT(1) NOT NULL DEFAULT 1,
    print_on_payslip        TINYINT(1) NOT NULL DEFAULT 1,
    display_order           INT NOT NULL DEFAULT 0,
    effective_start_date    DATE NOT NULL,
    effective_end_date      DATE NULL,
    status                  ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                   TEXT NULL,
    created_by              CHAR(36) NULL,
    updated_by              CHAR(36) NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_bpjs_rate_code_start (rate_code, effective_start_date),
    INDEX idx_bpjs_rate_setting (bpjs_setting_id),
    INDEX idx_bpjs_rate_program (bpjs_program, paid_by, status),
    INDEX idx_bpjs_rate_effective (effective_start_date, effective_end_date, status),
    INDEX idx_bpjs_rate_component (salary_component_id),

    CONSTRAINT fk_bpjs_rate_setting   FOREIGN KEY (bpjs_setting_id)     REFERENCES bpjs_settings(id)     ON DELETE CASCADE,
    CONSTRAINT fk_bpjs_rate_component FOREIGN KEY (salary_component_id) REFERENCES salary_components(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 6.13 PPh21 Settings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pph21_settings (
    id                                CHAR(36) PRIMARY KEY,
    setting_code                      VARCHAR(50) NOT NULL,
    setting_name                      VARCHAR(150) NOT NULL,
    calculation_method                ENUM('REGULAR_GROSS_ANNUALIZED') NOT NULL DEFAULT 'REGULAR_GROSS_ANNUALIZED',
    default_tax_method                ENUM('GROSS', 'GROSS_UP', 'NETT') NOT NULL DEFAULT 'GROSS',
    pph21_component_id                CHAR(36) NOT NULL,
    occupational_expense_rate_percent DECIMAL(8, 4) NOT NULL DEFAULT 5,
    occupational_expense_max_monthly  DECIMAL(18, 2) NOT NULL DEFAULT 500000,
    occupational_expense_max_yearly   DECIMAL(18, 2) NOT NULL DEFAULT 6000000,
    deduct_bpjs_health_employee       TINYINT(1) NOT NULL DEFAULT 0,
    deduct_bpjs_jht_employee          TINYINT(1) NOT NULL DEFAULT 1,
    deduct_bpjs_jp_employee           TINYINT(1) NOT NULL DEFAULT 1,
    annualization_months              TINYINT NOT NULL DEFAULT 12,
    pkp_rounding_unit                 DECIMAL(18, 2) NOT NULL DEFAULT 1000,
    non_npwp_multiplier_percent       DECIMAL(8, 4) NOT NULL DEFAULT 100,
    rounding_mode                     ENUM('NONE', 'ROUND', 'CEIL', 'FLOOR') NOT NULL DEFAULT 'ROUND',
    effective_start_date              DATE NOT NULL,
    effective_end_date                DATE NULL,
    status                            ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    notes                             TEXT NULL,
    created_by                        CHAR(36) NULL,
    updated_by                        CHAR(36) NULL,
    created_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_pph21_setting_code (setting_code),
    INDEX idx_pph21_setting_effective (effective_start_date, effective_end_date, status),
    INDEX idx_pph21_setting_component (pph21_component_id),

    CONSTRAINT fk_pph21_setting_component FOREIGN KEY (pph21_component_id) REFERENCES salary_components(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    status               ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_by           CHAR(36) NULL,
    updated_by           CHAR(36) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_pph21_ptkp_status_start (ptkp_status, effective_start_date),
    INDEX idx_pph21_ptkp_effective (effective_start_date, effective_end_date, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    status               ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_by           CHAR(36) NULL,
    updated_by           CHAR(36) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_pph21_bracket_order_start (bracket_order, effective_start_date),
    INDEX idx_pph21_bracket_effective (effective_start_date, effective_end_date, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
