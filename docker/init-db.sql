-- HRIS Platform - Database initialization script
-- Dijalankan oleh Docker Compose saat container postgres pertama kali start

-- Platform database sudah dibuat oleh POSTGRES_DB environment variable
-- Script ini hanya untuk setup tambahan

-- Buat superuser untuk tenant provisioning
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'hris_tenant_admin') THEN
        CREATE ROLE hris_tenant_admin WITH LOGIN PASSWORD 'tenant_admin_secret' SUPERUSER;
    END IF;
END
$$;

-- Buat tabel platform utama jika belum ada (migration akan handle sisanya)
CREATE TABLE IF NOT EXISTS public.schema_migrations (
    version bigint PRIMARY KEY,
    dirty boolean NOT NULL,
    applied_at timestamptz DEFAULT NOW()
);
