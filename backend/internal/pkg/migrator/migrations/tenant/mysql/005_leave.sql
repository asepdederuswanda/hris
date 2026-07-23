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
    is_paid               TINYINT(1) NOT NULL DEFAULT 1,
    requires_attachment   TINYINT(1) NOT NULL DEFAULT 0,
    allow_half_day        TINYINT(1) NOT NULL DEFAULT 1,
    default_quota_days    INT NULL,
    quota_period          ENUM('YEAR', 'MONTH', 'NONE') NOT NULL DEFAULT 'YEAR',
    counts_against_quota  TINYINT(1) NOT NULL DEFAULT 1,
    allow_hourly          TINYINT(1) NOT NULL DEFAULT 0,
    is_active             TINYINT(1) NOT NULL DEFAULT 1,
    deleted_at            TIMESTAMP NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_leave_type_name (name),
    INDEX idx_leave_type_active (is_active),
    INDEX idx_leave_type_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_accrual_policy (leave_type_id, effective_from),
    INDEX idx_accrual_leave_type (leave_type_id),

    CONSTRAINT fk_accrual_leave_type FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 5.3 Leave Reasons
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_reasons (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    is_other    TINYINT(1) NOT NULL DEFAULT 0,
    sort_order  INT NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_leave_reason_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 5.4 Leave Requests
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_requests (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    leave_type_id         CHAR(36) NOT NULL,
    request_start_date    DATE NOT NULL,
    request_end_date      DATE NOT NULL,
    duration_mode         ENUM('FULL_DAY', 'HALF_DAY_AM', 'HALF_DAY_PM', 'HOURLY') NOT NULL DEFAULT 'FULL_DAY',
    requested_days        DECIMAL(6, 2) NOT NULL DEFAULT 0.00,
    leave_reason_id       CHAR(36) NULL,
    leave_reason_note     VARCHAR(255) NULL,
    attachment_url        TEXT NULL,
    status                ENUM('DRAFT', 'SUBMITTED', 'PENDING_APPROVAL', 'APPROVED_FINAL', 'REJECTED_FINAL', 'CANCELLED') NOT NULL DEFAULT 'SUBMITTED',
    supervisor_id         CHAR(36) NULL,
    supervisor_action_at  DATETIME(6) NULL,
    supervisor_note       VARCHAR(255) NULL,
    hr_id                 CHAR(36) NULL,
    hr_action_at          DATETIME(6) NULL,
    hr_note               VARCHAR(255) NULL,
    approval_instance_id  CHAR(36) NULL,
    start_time            TIME NULL,
    end_time              TIME NULL,
    submitted_at          DATETIME(6) NULL,
    approved_at           DATETIME(6) NULL,
    rejected_at           DATETIME(6) NULL,
    cancelled_at          DATETIME(6) NULL,
    deleted_at            TIMESTAMP NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_leave_req_employee_date (employee_id, request_start_date, request_end_date),
    INDEX idx_leave_req_employee_created (employee_id, created_at),
    INDEX idx_leave_req_type (leave_type_id),
    INDEX idx_leave_req_reason (leave_reason_id),
    INDEX idx_leave_req_status (status),
    INDEX idx_leave_req_approval_instance (approval_instance_id),
    INDEX idx_leave_req_deleted_at (deleted_at),

    CONSTRAINT fk_leave_req_employee   FOREIGN KEY (employee_id)   REFERENCES employees(id)   ON DELETE CASCADE,
    CONSTRAINT fk_leave_req_type       FOREIGN KEY (leave_type_id)  REFERENCES leave_types(id)  ON DELETE CASCADE,
    CONSTRAINT fk_leave_req_reason     FOREIGN KEY (leave_reason_id) REFERENCES leave_reasons(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 5.5 Leave Request Details
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_request_details (
    id                    CHAR(36) PRIMARY KEY,
    leave_request_id      CHAR(36) NOT NULL,
    employee_id           CHAR(36) NOT NULL,
    leave_date            DATE NOT NULL,
    day_fraction          DECIMAL(4, 2) NOT NULL DEFAULT 1.00,
    is_paid               TINYINT(1) NOT NULL DEFAULT 1,
    approval_instance_id  CHAR(36) NULL,
    deleted_at            TIMESTAMP NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_leave_req_day (leave_request_id, leave_date),
    INDEX idx_leave_detail_emp_date (employee_id, leave_date),
    INDEX idx_leave_detail_request (leave_request_id),
    INDEX idx_leave_detail_deleted_at (deleted_at),

    CONSTRAINT fk_leave_detail_request FOREIGN KEY (leave_request_id) REFERENCES leave_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_leave_detail_employee FOREIGN KEY (employee_id)     REFERENCES employees(id)      ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    last_adjustment_at      DATETIME(6) NULL,
    created_at              TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_leave_balance (employee_id, leave_type_id, period_year),
    INDEX idx_leave_balance_type (leave_type_id),

    CONSTRAINT fk_leave_balance_employee FOREIGN KEY (employee_id)  REFERENCES employees(id)  ON DELETE CASCADE,
    CONSTRAINT fk_leave_balance_type     FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
