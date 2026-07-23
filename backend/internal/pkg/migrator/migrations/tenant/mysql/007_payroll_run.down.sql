-- Down Migration: 007_payroll_run
DROP TABLE IF EXISTS payroll_profile_change_logs;
DROP TABLE IF EXISTS pph21_calculation_logs;
DROP TABLE IF EXISTS payroll_payslips;
DROP TABLE IF EXISTS payroll_run_items;
DROP TABLE IF EXISTS payroll_run_employees;
DROP TABLE IF EXISTS payroll_runs;
