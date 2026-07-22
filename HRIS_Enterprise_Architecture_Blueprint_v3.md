# HRIS Enterprise Architecture Blueprint v3

## 1. Vision & Goals

Membangun platform HRIS Enterprise berbasis Go + Vue dengan arsitektur
Modular Monolith, Multi-Tenant (Database per Company), dan Platform
Management.

## 2. Product Scope

-   Platform Admin
-   Multi Company
-   Module Management
-   License Management
-   Feature Flag
-   HRIS Modules

## 3. Roadmap

1.  Existing Code Analysis
2.  Existing Documentation
3.  Platform Architecture Design
4.  Module Management Design
5.  Backend Implementation
6.  API Documentation
7.  Backend Testing
8.  Frontend Implementation
9.  Frontend Testing
10. UAT & Go Live

## 4. Multi-Tenant Architecture

### Platform Database

-   companies
-   company_connections
-   modules
-   module_dependencies
-   company_modules
-   licenses
-   feature_flags
-   platform_users
-   audit_logs

### Tenant Database (1 Company = 1 Database)

-   employees
-   attendance
-   leave
-   payroll
-   overtime
-   recruitment
-   performance
-   training
-   asset
-   approval
-   reports
-   settings

Catatan: - Tidak menggunakan company_id pada tabel bisnis tenant. -
Semua tenant memiliki struktur database yang sama.

## 5. Tenant Provisioning Engine

Saat company dibuat: 1. Validasi lisensi. 2. Membuat database. 3.
Membuat database user. 4. Menyimpan koneksi tenant. 5. Menjalankan core
migration. 6. Menjalankan migration modul aktif. 7. Menjalankan seeder.
8. Membuat Super Admin. 9. Mengaktifkan modul. 10. Audit log. 11. Health
check tenant.

## 6. Platform Management

-   Company Management
-   Module Management
-   License Management
-   Feature Flag
-   System Configuration
-   Tenant Monitoring
-   Tenant Backup & Restore

## 7. Module SDK

Setiap modul wajib memiliki: - manifest.yaml - routes - handlers -
services - repositories - entities - dto - permissions - menus -
migrations - seeders - configs - tests - api docs

## 8. Backend Architecture

-   Go
-   Gin
-   GORM
-   Redis
-   Asynq
-   Casbin
-   Zap
-   OpenTelemetry

## 9. Frontend Architecture

-   Vue 3
-   TypeScript
-   Vite
-   PrimeVue
-   Tailwind CSS
-   Pinia
-   TanStack Query

## 10. Monorepo

backend/ platform/ tenant/ shared/ installer/ frontend/ platform-admin/
tenant/ shared/ docs/ docker/

## 11. Security

-   JWT
-   RBAC
-   Audit Log
-   Rate Limiting
-   MFA Ready
-   OWASP ASVS

## 12. CI/CD

-   GitHub Actions
-   Docker
-   Auto Test
-   Auto Build
-   Auto Deploy

## 13. Testing

-   Unit
-   Integration
-   API
-   E2E (Playwright)
-   UAT
-   Performance

## 14. Disaster Recovery

-   Backup per tenant
-   Restore per tenant
-   Point-in-time recovery
-   Disaster recovery procedure

## 15. Future Roadmap

-   Plugin Marketplace
-   AI Assistant
-   BI Dashboard
-   Mobile Apps
-   Public API
-   Kubernetes Deployment
