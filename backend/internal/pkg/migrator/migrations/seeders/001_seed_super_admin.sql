-- =============================================================================
-- Super Admin Seeder
-- =============================================================================
-- Jalankan setelah migration selesai untuk membuat user super admin default.
-- Password: admin123 (bcrypt hash)
-- Email: superadmin@hris-platform.com

INSERT IGNORE INTO platform_users (id, email, password_hash, name, role, is_active)
SELECT UUID(), 'superadmin@hris-platform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Super Admin', 'super_admin', 1
WHERE NOT EXISTS (SELECT 1 FROM platform_users WHERE email = 'superadmin@hris-platform.com');
