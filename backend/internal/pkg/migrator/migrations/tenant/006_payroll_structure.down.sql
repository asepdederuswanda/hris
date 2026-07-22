-- Down Migration: 006_payroll_structure
DROP TABLE IF EXISTS pph21_tax_brackets;
DROP TABLE IF EXISTS pph21_ptkp_rates;
DROP TABLE IF EXISTS pph21_settings;
DROP TABLE IF EXISTS bpjs_rate_components;
DROP TABLE IF EXISTS bpjs_settings;
DROP TABLE IF EXISTS employee_tax_profiles;
DROP TABLE IF EXISTS employee_bpjs_profiles;
DROP TABLE IF EXISTS employee_bank_profiles;
DROP TABLE IF EXISTS employee_payroll_profiles;
DROP TABLE IF EXISTS payroll_periods;
DROP TABLE IF EXISTS salary_employee_adjustments;
DROP TABLE IF EXISTS salary_change_logs;
DROP TABLE IF EXISTS salary_employee_components;
DROP TABLE IF EXISTS salary_grade_components;
DROP TABLE IF EXISTS salary_components;
