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
    is_location_required    SMALLINT NOT NULL DEFAULT 0,
    is_face_required        SMALLINT NOT NULL DEFAULT 0,
    is_overtime_enabled     SMALLINT NOT NULL DEFAULT 0,
    max_distance_meter      INT NULL,
    late_tolerance_minutes  INT NULL,
    overtime_min_minutes    INT NULL,
    created_by              CHAR(36) NULL,
    updated_by              CHAR(36) NULL,
    deleted_by              CHAR(36) NULL,
    deleted_at              TIMESTAMP,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_att_setting_deleted_at ON attendance_company_settings (deleted_at);

-- ---------------------------------------------------------------------------
-- 4.2 Attendance Company Shifts
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_company_shifts (
    id                CHAR(36) PRIMARY KEY,
    shift_name        VARCHAR(255) NOT NULL,
    check_in_time     TIME NOT NULL,
    check_out_time    TIME NOT NULL,
    is_cross_midnight SMALLINT NOT NULL DEFAULT 0,
    created_by        CHAR(36) NULL,
    updated_by        CHAR(36) NULL,
    deleted_at        TIMESTAMP,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_att_shift_deleted_at ON attendance_company_shifts (deleted_at);

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
    is_day_off          SMALLINT NULL,
    created_by          CHAR(36) NULL,
    updated_by          CHAR(36) NULL,
    deleted_at          TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_emp_shift_date UNIQUE (employee_id, attendance_shift_id, effective_date_from),

    CONSTRAINT fk_empshift_employee FOREIGN KEY (employee_id)         REFERENCES employees(id)              ON DELETE CASCADE,
    CONSTRAINT fk_empshift_shift    FOREIGN KEY (attendance_shift_id) REFERENCES attendance_company_shifts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_empshift_shift ON attendance_employee_shifts (attendance_shift_id);

CREATE INDEX IF NOT EXISTS idx_empshift_deleted_at ON attendance_employee_shifts (deleted_at);

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
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_attloc_deleted_at ON attendance_locations (deleted_at);

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
    last_seen_at  TIMESTAMP,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_device_uuid UNIQUE (device_uuid)
);

-- ---------------------------------------------------------------------------
-- 4.6 Attendance Face Captures
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_face_captures (
    id              CHAR(36) PRIMARY KEY,
    employee_id     CHAR(36) NOT NULL,
    captured_at     TIMESTAMP NOT NULL,
    image_url       TEXT NOT NULL,
    image_sha256    CHAR(64) NOT NULL,
    liveness_score  DECIMAL(6, 3) NULL,
    match_score     DECIMAL(6, 3) NULL,
    verified        SMALLINT NULL DEFAULT 0,
    provider        VARCHAR(50) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_face_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_face_employee_time ON attendance_face_captures (employee_id, captured_at);

-- ---------------------------------------------------------------------------
-- 4.7 Attendance Events
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_events (
    id                    CHAR(36) PRIMARY KEY,
    employee_id           CHAR(36) NOT NULL,
    overtime_request_id   CHAR(36) NULL,
    event_type            VARCHAR(255) NOT NULL,
    event_time_utc        TIMESTAMP NOT NULL,
    event_time_local      TIMESTAMP NOT NULL,
    device_id             CHAR(36) NULL,
    latitude              DECIMAL(10, 7) NOT NULL,
    longitude             DECIMAL(10, 7) NOT NULL,
    accuracy_m            INT NULL,
    location_provider     VARCHAR(255) NULL,
    validated_location_id CHAR(36) NULL,
    distance_m            INT NULL,
    is_in_geofence        SMALLINT NOT NULL DEFAULT 0,
    face_capture_id       CHAR(36) NULL,
    validation_status     VARCHAR(255) NOT NULL DEFAULT 'PENDING',
    validation_note       VARCHAR(255) NULL,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_att_event_employee  FOREIGN KEY (employee_id)    REFERENCES employees(id)                ON DELETE CASCADE,
    CONSTRAINT fk_att_event_device    FOREIGN KEY (device_id)      REFERENCES attendance_device_captures(id) ON DELETE SET NULL,
    CONSTRAINT fk_att_event_face      FOREIGN KEY (face_capture_id) REFERENCES attendance_face_captures(id)  ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_att_event_employee ON attendance_events (employee_id);

CREATE INDEX IF NOT EXISTS idx_att_event_device ON attendance_events (device_id);

CREATE INDEX IF NOT EXISTS idx_att_event_face ON attendance_events (face_capture_id);

CREATE INDEX IF NOT EXISTS idx_att_event_overtime ON attendance_events (overtime_request_id);

CREATE INDEX IF NOT EXISTS idx_att_event_location ON attendance_events (validated_location_id);

CREATE INDEX IF NOT EXISTS idx_att_event_deleted_at ON attendance_events (deleted_at);

-- ---------------------------------------------------------------------------
-- 4.8 Attendance Sessions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_sessions (
    id                           CHAR(36) PRIMARY KEY,
    employee_id                  CHAR(36) NOT NULL,
    work_date                    DATE NOT NULL,
    shift_id                     CHAR(36) NULL,
    is_overtime_day              SMALLINT NULL DEFAULT 0,
    overtime_request_id          CHAR(36) NULL,
    approved_overtime_start_local TIMESTAMP(6) NULL,
    approved_overtime_end_local   TIMESTAMP(6) NULL,
    leave_request_id             CHAR(36) NULL,
    leave_fraction               DECIMAL(4, 2) NULL,
    planned_start_local          TIMESTAMP(6) NULL,
    planned_end_local            TIMESTAMP(6) NULL,
    checkin_event_id             CHAR(36) NULL,
    checkout_event_id            CHAR(36) NULL,
    status                       VARCHAR(255) NOT NULL DEFAULT 'OPEN',
    lateness_minutes             INT NOT NULL DEFAULT 0,
    early_leave_minutes          INT NOT NULL DEFAULT 0,
    work_minutes                 INT NOT NULL DEFAULT 0,
    break_minutes                INT NOT NULL DEFAULT 0,
    overtime_minutes             INT NOT NULL DEFAULT 0,
    deleted_at                   TIMESTAMP,
    created_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_session_emp_date UNIQUE (employee_id, work_date),

    CONSTRAINT fk_att_session_employee  FOREIGN KEY (employee_id)       REFERENCES employees(id)             ON DELETE CASCADE,
    CONSTRAINT fk_att_session_shift     FOREIGN KEY (shift_id)          REFERENCES attendance_company_shifts(id) ON DELETE SET NULL,
    CONSTRAINT fk_att_session_checkin   FOREIGN KEY (checkin_event_id)  REFERENCES attendance_events(id)     ON DELETE SET NULL,
    CONSTRAINT fk_att_session_checkout  FOREIGN KEY (checkout_event_id) REFERENCES attendance_events(id)     ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_att_session_shift ON attendance_sessions (shift_id);

CREATE INDEX IF NOT EXISTS idx_att_session_checkin ON attendance_sessions (checkin_event_id);

CREATE INDEX IF NOT EXISTS idx_att_session_checkout ON attendance_sessions (checkout_event_id);

CREATE INDEX IF NOT EXISTS idx_att_session_overtime ON attendance_sessions (overtime_request_id);

CREATE INDEX IF NOT EXISTS idx_att_session_leave ON attendance_sessions (leave_request_id);

CREATE INDEX IF NOT EXISTS idx_att_session_deleted_at ON attendance_sessions (deleted_at);

-- ---------------------------------------------------------------------------
-- 4.9 Attendance Overtime Requests
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_overtime_requests (
    id                   CHAR(36) PRIMARY KEY,
    employee_id          CHAR(36) NOT NULL,
    work_date            DATE NOT NULL,
    start_time_local     TIMESTAMP(6) NOT NULL,
    end_time_local       TIMESTAMP(6) NOT NULL,
    requested_minutes    INT NOT NULL,
    reason               VARCHAR(255) NULL,
    status               VARCHAR(255) NOT NULL DEFAULT 'SUBMITTED',
    approved_by          CHAR(36) NULL,
    approved_at          TIMESTAMP(6) NULL,
    approval_note        VARCHAR(255) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_att_overtime_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_att_overtime_employee ON attendance_overtime_requests (employee_id);

CREATE INDEX IF NOT EXISTS idx_att_overtime_status ON attendance_overtime_requests (status);

CREATE INDEX IF NOT EXISTS idx_att_overtime_date ON attendance_overtime_requests (work_date);

-- ---------------------------------------------------------------------------
-- 4.10 Attendance Exempt Positions
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendance_exempt_positions (
    id              CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    is_exempt       SMALLINT NOT NULL DEFAULT 1,
    note            VARCHAR(255) NULL,
    created_by      CHAR(36) NULL,
    updated_by      CHAR(36) NULL,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_exempt_org UNIQUE (organization_id),

    CONSTRAINT fk_att_exempt_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_att_exempt_deleted_at ON attendance_exempt_positions (deleted_at);
