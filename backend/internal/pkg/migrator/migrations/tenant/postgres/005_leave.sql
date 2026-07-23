-- =============================================================================
-- Tenant Migration: 005_leave
-- =============================================================================
-- Tabel untuk modul cuti/izin tenant.

-- ---------------------------------------------------------------------------
-- 5.1 Leave Types
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_types (
    id                    CHAR(36) PRIMARY KEY,
    code                  VARCHAR(50) NOT NULL DEFAULT '',
    name                  VARCHAR(50) NOT NULL,
    description           VARCHAR(200) NOT NULL,
    is_paid               SMALLINT NOT NULL DEFAULT 1,
    requires_attachment   SMALLINT NOT NULL DEFAULT 0,
    allow_half_day        SMALLINT NOT NULL DEFAULT 1,
    default_quota_days    INT NULL,
    quota_period          VARCHAR(255) NOT NULL DEFAULT 'INTEGER',
    counts_against_quota  SMALLINT NOT NULL DEFAULT 1,
    allow_hourly          SMALLINT NOT NULL DEFAULT 0,
    is_active             SMALLINT NOT NULL DEFAULT 1,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_leave_type_name UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_leave_type_active ON leave_types (is_active);

CREATE INDEX IF NOT EXISTS idx_leave_type_deleted_at ON leave_types (deleted_at);

-- ---------------------------------------------------------------------------
-- 5.2 Leave Accrual Policies
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_accrual_policies (
    id                CHAR(36) PRIMARY KEY,
    leave_type_id     CHAR(36) NOT NULL,
    base_quota_days   DECIMAL(6, 2) NOT NULL,
    extra_every_years INT NOT NULL DEFAULT 2,
    extra_days        DECIMAL(6, 2) NOT NULL DEFAULT 1.00,
    max_extra_days    DECIMAL(6, 2) NULL,
    effective_from    DATE NOT NULL,
    effective_to      DATE NULL,
    deleted_at        INT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_accrual_policy UNIQUE (leave_type_id, effective_from),

    CONSTRAINT fk_accrual_leave_type FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_accrual_leave_type ON leave_accrual_policies (leave_type_id);

-- ---------------------------------------------------------------------------
-- 5.3 Leave Reasons
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_reasons (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    is_other    SMALLINT NOT NULL DEFAULT 0,
    sort_order  INT NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_leave_reason_deleted_at ON leave_reasons (deleted_at);

-- ---------------------------------------------------------------------------
-- 5.4 Leave Requests
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_requests (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    leave_type_id         CHAR(36) NOT NULL,
    request_start_date    DATE NOT NULL,
    request_end_date      DATE NOT NULL,
    duration_mode         VARCHAR(255) NOT NULL DEFAULT 'FULL_DAY',
    requested_days        DECIMAL(6, 2) NOT NULL DEFAULT 0.00,
    leave_reason_id       CHAR(36) NULL,
    leave_reason_note     VARCHAR(255) NULL,
    attachment_url        TEXT NULL,
    status                VARCHAR(255) NOT NULL DEFAULT 'SUBMITTED',
    supervisor_id         CHAR(36) NULL,
    supervisor_action_at  TIMESTAMP(6) NULL,
    supervisor_note       VARCHAR(255) NULL,
    hr_id                 CHAR(36) NULL,
    hr_action_at          TIMESTAMP(6) NULL,
    hr_note               VARCHAR(255) NULL,
    approval_instance_id  CHAR(36) NULL,
    start_time            TIME NULL,
    end_time              TIME NULL,
    submitted_at          TIMESTAMP(6) NULL,
    approved_at           TIMESTAMP(6) NULL,
    rejected_at           TIMESTAMP(6) NULL,
    cancelled_at          TIMESTAMP(6) NULL,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_leave_req_employee   FOREIGN KEY (employee_id)   REFERENCES employees(id)   ON DELETE CASCADE,
    CONSTRAINT fk_leave_req_type       FOREIGN KEY (leave_type_id)  REFERENCES leave_types(id)  ON DELETE CASCADE,
    CONSTRAINT fk_leave_req_reason     FOREIGN KEY (leave_reason_id) REFERENCES leave_reasons(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_leave_req_employee_date ON leave_requests (employee_id, request_start_date, request_end_date);

CREATE INDEX IF NOT EXISTS idx_leave_req_employee_created ON leave_requests (employee_id, created_at);

CREATE INDEX IF NOT EXISTS idx_leave_req_type ON leave_requests (leave_type_id);

CREATE INDEX IF NOT EXISTS idx_leave_req_reason ON leave_requests (leave_reason_id);

CREATE INDEX IF NOT EXISTS idx_leave_req_status ON leave_requests (status);

CREATE INDEX IF NOT EXISTS idx_leave_req_approval_instance ON leave_requests (approval_instance_id);

CREATE INDEX IF NOT EXISTS idx_leave_req_deleted_at ON leave_requests (deleted_at);

-- ---------------------------------------------------------------------------
-- 5.5 Leave Request Details
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_request_details (
    id                    CHAR(36) PRIMARY KEY,
    leave_request_id      CHAR(36) NOT NULL,
    employee_id           CHAR(36) NOT NULL,
    leave_date            DATE NOT NULL,
    day_fraction          DECIMAL(4, 2) NOT NULL DEFAULT 1.00,
    is_paid               SMALLINT NOT NULL DEFAULT 1,
    approval_instance_id  CHAR(36) NULL,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_leave_req_day UNIQUE (leave_request_id, leave_date),

    CONSTRAINT fk_leave_detail_request FOREIGN KEY (leave_request_id) REFERENCES leave_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_leave_detail_employee FOREIGN KEY (employee_id)     REFERENCES employees(id)      ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_leave_detail_emp_date ON leave_request_details (employee_id, leave_date);

CREATE INDEX IF NOT EXISTS idx_leave_detail_request ON leave_request_details (leave_request_id);

CREATE INDEX IF NOT EXISTS idx_leave_detail_deleted_at ON leave_request_details (deleted_at);

-- ---------------------------------------------------------------------------
-- 5.6 Employee Leave Balances
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_leave_balances (
    id                      CHAR(36) PRIMARY KEY,
    employee_id             CHAR(36) NOT NULL,
    leave_type_id           CHAR(36) NOT NULL,
    period_year             INT NOT NULL,
    quota_days              DECIMAL(6, 2) NOT NULL DEFAULT 0,
    used_days               DECIMAL(6, 2) NOT NULL DEFAULT 0,
    remaining_days          DECIMAL(6, 2) NOT NULL DEFAULT 0,
    last_adjustment_ref     VARCHAR(50) NULL,
    last_adjustment_ref_id  CHAR(36) NULL,
    last_adjustment_at      TIMESTAMP(6) NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_leave_balance UNIQUE (employee_id, leave_type_id, period_year),

    CONSTRAINT fk_leave_balance_employee FOREIGN KEY (employee_id)  REFERENCES employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_leave_balance_type     FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_leave_balance_type ON employee_leave_balances (leave_type_id);
