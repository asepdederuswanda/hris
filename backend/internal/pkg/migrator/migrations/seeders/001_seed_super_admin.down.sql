-- Down Migration: 001_seed_super_admin
-- Rollback: Hapus user super admin default

DELETE FROM platform_users WHERE email = 'superadmin@hris-platform.com';
