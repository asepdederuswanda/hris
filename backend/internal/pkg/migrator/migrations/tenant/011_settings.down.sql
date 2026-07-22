-- Down Migration: 011_settings
DROP TABLE IF EXISTS document_templates;
DROP TABLE IF EXISTS company_holidays;
DROP TABLE IF EXISTS role_has_permissions;
DROP TABLE IF EXISTS model_has_permissions;
DROP TABLE IF EXISTS model_has_roles;
DROP TABLE IF EXISTS feature_permission;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS features;
