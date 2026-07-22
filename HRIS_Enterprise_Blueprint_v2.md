# HRIS Enterprise Blueprint (Updated)

## Architecture Update

### Multi-Tenant Strategy

**Database per Company (Database-per-Tenant)**

    Platform Database
    ├── companies
    ├── company_connections
    ├── modules
    ├── module_dependencies
    ├── company_modules
    ├── feature_flags
    ├── licenses
    └── platform_users

    Tenant Database (1 database/company)
    ├── employees
    ├── attendance
    ├── leave
    ├── payroll
    ├── asset
    ├── approval
    ├── training
    └── settings

> Tidak menggunakan `company_id` pada tabel bisnis tenant karena setiap
> database hanya dimiliki satu perusahaan.

## Phase 2.5 - Platform Management

### Company Management

-   Create Company
-   Update Company
-   Suspend/Activate Company
-   Delete Company
-   Company Branding
-   Company Domain
-   Database Connection Management

### Module Management

-   Module Registry
-   Module Loader
-   Module Dependency
-   Dynamic Menu
-   Dynamic Route
-   Dynamic Permission
-   Dynamic Migration
-   Dynamic Seeder

### License Management

-   Trial
-   Basic
-   Professional
-   Enterprise
-   Expiration
-   License Validation

### Feature Flag

-   Enable/Disable Feature per Company
-   Beta Feature
-   Experimental Feature

## Tenant Provisioning Engine

Saat company dibuat sistem otomatis:

1.  Membuat database tenant.
2.  Membuat user database.
3.  Menyimpan koneksi ke Platform Database.
4.  Menjalankan migration core.
5.  Menjalankan migration modul aktif.
6.  Menjalankan seeder.
7.  Membuat Super Admin tenant.
8.  Mengaktifkan lisensi.
9.  Mengaktifkan module.
10. Menulis audit log.

## Platform Database

-   companies
-   company_connections
-   modules
-   module_dependencies
-   company_modules
-   licenses
-   feature_flags
-   platform_users
-   audit_logs

## Monorepo

    hris-platform/
    ├── backend/
    │   ├── platform/
    │   ├── tenant/
    │   ├── shared/
    │   └── installer/
    ├── frontend/
    │   ├── platform-admin/
    │   ├── tenant/
    │   └── shared/
    ├── docs/
    └── docker/

## Deliverables Tambahan

-   Platform Management
-   Tenant Provisioning Engine
-   Database-per-Tenant Architecture
-   Module SDK
-   License SDK
-   Feature Flag Engine
-   Company Installer
-   Migration Runner
-   Tenant Backup & Restore Guide
-   Disaster Recovery Guide
