-- Down Migration: 005_leave
DROP TABLE IF EXISTS employee_leave_balances;
DROP TABLE IF EXISTS leave_request_details;
DROP TABLE IF EXISTS leave_requests;
DROP TABLE IF EXISTS leave_reasons;
DROP TABLE IF EXISTS leave_accrual_policies;
DROP TABLE IF EXISTS leave_types;
