-- Down Migration: 004_attendance
DROP TABLE IF EXISTS attendance_exempt_positions;
DROP TABLE IF EXISTS attendance_overtime_requests;
DROP TABLE IF EXISTS attendance_sessions;
DROP TABLE IF EXISTS attendance_events;
DROP TABLE IF EXISTS attendance_face_captures;
DROP TABLE IF EXISTS attendance_device_captures;
DROP TABLE IF EXISTS attendance_locations;
DROP TABLE IF EXISTS attendance_employee_shifts;
DROP TABLE IF EXISTS attendance_company_shifts;
DROP TABLE IF EXISTS attendance_company_settings;
