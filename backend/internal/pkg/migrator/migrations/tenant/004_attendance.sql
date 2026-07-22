-- =============================================================================
-- Tenant Migration: 004_attendance
-- =============================================================================
-- Tabel untuk modul absensi/kehadiran tenant.

-- ---------------------------------------------------------------------------
-- 4.1 Attendance Company Settings
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_company_settings (
    id                      CHAR(36) PRIMARY KEY,
    latitude                DECIMAL(10, 8) NULL,
    longitude               DECIMAL(11, 8) NULL,
    is_location_required    TINYINT(1) NOT NULL DEFAULT 0,
    is_face_required        TINYINT(1) NOT NULL DEFAULT 0,
    is_overtime_enabled     TINYINT(1) NOT NULL DEFAULT 0,
    max_distance_meter      INT NULL,
    late_tolerance_minutes  INT NULL,
    overtime_min_minutes    INT NULL,
    created_by              CHAR(36) NULL,
    updated_by              CHAR(36) NULL,
    deleted_by              CHAR(36) NULL,
    deleted_at              TIMESTAMP NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_att_setting_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.2 Attendance Company Shifts
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_company_shifts (
    id                CHAR(36) PRIMARY KEY,
    shift_name        VARCHAR(255) NOT NULL,
    check_in_time     TIME NOT NULL,
    check_out_time    TIME NOT NULL,
    is_cross_midnight TINYINT(1) NOT NULL DEFAULT 0,
    created_by        CHAR(36) NULL,
    updated_by        CHAR(36) NULL,
    deleted_at        TIMESTAMP NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_att_shift_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.3 Attendance Employee Shifts
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_employee_shifts (
    id                  CHAR(36) PRIMARY KEY,
    employee_id         CHAR(36) NOT NULL,
    attendance_shift_id CHAR(36) NOT NULL,
    effective_date_from DATE NOT NULL,
    effective_date_to   DATE NULL,
    days_of_week_mask   INT NULL,
    is_day_off          TINYINT(1) NULL,
    created_by          CHAR(36) NULL,
    updated_by          CHAR(36) NULL,
    deleted_at          TIMESTAMP NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_emp_shift_date (employee_id, attendance_shift_id, effective_date_from),
    INDEX idx_empshift_shift (attendance_shift_id),
    INDEX idx_empshift_deleted_at (deleted_at),

    CONSTRAINT fk_empshift_employee FOREIGN KEY (employee_id)         REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_empshift_shift    FOREIGN KEY (attendance_shift_id) REFERENCES attendance_company_shifts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.4 Attendance Locations
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_locations (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    latitude    DECIMAL(10, 7) NOT NULL,
    longitude   DECIMAL(10, 7) NOT NULL,
    radius_m    INT NOT NULL DEFAULT 50,
    created_by  CHAR(36) NULL,
    updated_by  CHAR(36) NULL,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_attloc_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.5 Attendance Device Captures
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_device_captures (
    id            CHAR(36) PRIMARY KEY,
    device_uuid   VARCHAR(100) NOT NULL,
    device_type   VARCHAR(20) NOT NULL,
    os_version    VARCHAR(50) NULL,
    model         VARCHAR(100) NULL,
    app_version   VARCHAR(50) NULL,
    last_seen_at  TIMESTAMP NULL,
    deleted_at    TIMESTAMP NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_device_uuid (device_uuid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.6 Attendance Face Captures
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_face_captures (
    id              CHAR(36) PRIMARY KEY,
    employee_id     CHAR(36) NOT NULL,
    captured_at     DATETIME NOT NULL,
    image_url       TEXT NOT NULL,
    image_sha256    CHAR(64) NOT NULL,
    liveness_score  DECIMAL(6, 3) NULL,
    match_score     DECIMAL(6, 3) NULL,
    verified        TINYINT(1) NULL DEFAULT 0,
    provider        VARCHAR(50) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_face_employee_time (employee_id, captured_at),

    CONSTRAINT fk_face_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.7 Attendance Events
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_events (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    overtime_request_id   CHAR(36) NULL,
    event_type            ENUM('CHECKIN', 'CHECKOUT') NOT NULL,
    event_time_utc        DATETIME NOT NULL,
    event_time_local      DATETIME NOT NULL,
    device_id             CHAR(36) NULL,
    latitude              DECIMAL(10, 7) NOT NULL,
    longitude             DECIMAL(10, 7) NOT NULL,
    accuracy_m            INT NULL,
    location_provider     ENUM('GPS', 'NETWORK', 'WIFI', 'MANUAL') NULL,
    validated_location_id CHAR(36) NULL,
    distance_m            INT NULL,
    is_in_geofence        TINYINT(1) NOT NULL DEFAULT 0,
    face_capture_id       CHAR(36) NULL,
    validation_status     ENUM('PENDING', 'VALID', 'INVALID', 'OVERRIDDEN') NOT NULL DEFAULT 'PENDING',
    validation_note       VARCHAR(255) NULL,
    deleted_at            TIMESTAMP NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_att_event_employee (employee_id),
    INDEX idx_att_event_device (device_id),
    INDEX idx_att_event_face (face_capture_id),
    INDEX idx_att_event_overtime (overtime_request_id),
    INDEX idx_att_event_location (validated_location_id),
    INDEX idx_att_event_deleted_at (deleted_at),

    CONSTRAINT fk_att_event_employee  FOREIGN KEY (employee_id)    REFERENCES employees(id)                ON DELETE CASCADE,
    CONSTRAINT fk_att_event_device    FOREIGN KEY (device_id)      REFERENCES attendance_device_captures(id) ON DELETE SET NULL,
    CONSTRAINT fk_att_event_face      FOREIGN KEY (face_capture_id) REFERENCES attendance_face_captures(id)  ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.8 Attendance Sessions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_sessions (
    id                           CHAR(36) PRIMARY KEY,
    employee_id                  CHAR(36) NOT NULL,
    work_date                    DATE NOT NULL,
    shift_id                     CHAR(36) NULL,
    is_overtime_day              TINYINT(1) NULL DEFAULT 0,
    overtime_request_id          CHAR(36) NULL,
    approved_overtime_start_local DATETIME(6) NULL,
    approved_overtime_end_local   DATETIME(6) NULL,
    leave_request_id             CHAR(36) NULL,
    leave_fraction               DECIMAL(4, 2) NULL,
    planned_start_local          DATETIME(6) NULL,
    planned_end_local            DATETIME(6) NULL,
    checkin_event_id             CHAR(36) NULL,
    checkout_event_id            CHAR(36) NULL,
    status                       ENUM('OPEN', 'CLOSED', 'MISSING_CHECKIN', 'MISSING_CHECKOUT', 'ABSENT', 'DAY_OFF', 'EXEMPT', 'LEAVE') NOT NULL DEFAULT 'OPEN',
    lateness_minutes             INT NOT NULL DEFAULT 0,
    early_leave_minutes          INT NOT NULL DEFAULT 0,
    work_minutes                 INT NOT NULL DEFAULT 0,
    break_minutes                INT NOT NULL DEFAULT 0,
    overtime_minutes             INT NOT NULL DEFAULT 0,
    deleted_at                   TIMESTAMP NULL,
    created_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_session_emp_date (employee_id, work_date),
    INDEX idx_att_session_shift (shift_id),
    INDEX idx_att_session_checkin (checkin_event_id),
    INDEX idx_att_session_checkout (checkout_event_id),
    INDEX idx_att_session_overtime (overtime_request_id),
    INDEX idx_att_session_leave (leave_request_id),
    INDEX idx_att_session_deleted_at (deleted_at),

    CONSTRAINT fk_att_session_employee  FOREIGN KEY (employee_id)       REFERENCES employees(id)             ON DELETE CASCADE,
    CONSTRAINT fk_att_session_shift     FOREIGN KEY (shift_id)          REFERENCES attendance_company_shifts(id) ON DELETE SET NULL,
    CONSTRAINT fk_att_session_checkin   FOREIGN KEY (checkin_event_id)  REFERENCES attendance_events(id)     ON DELETE SET NULL,
    CONSTRAINT fk_att_session_checkout  FOREIGN KEY (checkout_event_id) REFERENCES attendance_events(id)     ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.9 Attendance Overtime Requests
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_overtime_requests (
    id                   CHAR(36) PRIMARY KEY,
    employee_id          CHAR(36) NOT NULL,
    work_date            DATE NOT NULL,
    start_time_local     DATETIME(6) NOT NULL,
    end_time_local       DATETIME(6) NOT NULL,
    requested_minutes    INT NOT NULL,
    reason               VARCHAR(255) NULL,
    status               ENUM('SUBMITTED', 'APPROVED', 'REJECTED') NOT NULL DEFAULT 'SUBMITTED',
    approved_by          CHAR(36) NULL,
    approved_at          DATETIME(6) NULL,
    approval_note        VARCHAR(255) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_att_overtime_employee (employee_id),
    INDEX idx_att_overtime_status (status),
    INDEX idx_att_overtime_date (work_date),

    CONSTRAINT fk_att_overtime_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 4.10 Attendance Exempt Positions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_exempt_positions (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    is_exempt       TINYINT(1) NOT NULL DEFAULT 1,
    note            VARCHAR(255) NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    deleted_at      TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_exempt_org (organization_id),
    INDEX idx_att_exempt_deleted_at (deleted_at),

    CONSTRAINT fk_att_exempt_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
