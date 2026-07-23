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

## 3. Module Completion Status

### 3.1 Modul Inti (Completed ✅)

| Modul | Status | Detail |
|-------|:------:|--------|
| **Platform & Tenant Management** | ✅ Completed | Provisioning DB multi-tenant, isolasi database, switching context tenant, lifecycle management (Suspend/Activate/Terminate) |
| **Organization Management** | ✅ Completed | Multi-Company Architecture, Dynamic Department Hierarchy (Adjacency List), Location & Geofencing Zones, Organization Summary |
| **Employee Management** | ✅ Completed | Data personal, kontak, alamat, keluarga, pendidikan, dokumen, riwayat kerja, rekening/pajak, 8 sub-modules |
| **Job Management** | ✅ Completed | 18 GORM entities: Titles, Subs, Values, Objectives, Identifications, Responsibilities, Education Experiences, HR/Operational Authorities, Working Activities/Risks, Relationships, Subordinate Controls, Assets, Financials, Potency Competencies, Scores, Competency Groups |
| **Competency Management** | 🗄️ DB Schema Only | Tabel database tersedia (DDL migration 008_competency.sql — 7 tabel), Go module belum diimplementasikan |

### 3.2 Modul Operasional & Siklus Karier (Planned 🗓️)

| Modul | Prioritas | Scope |
|-------|:---------:|-------|
| **Organization History, Versioning & Cloning** | 🟢 High | Change Capture, Full Structure Cloning DRAFT, Version Audit Trail |
| **Employee Movement & Career Management** | 🔴 High | Promosi/Demosi, Perpanjangan Kontrak (PKWT), Pensiun & Offboarding/PHK |
| **Time & Attendance** | 🔴 High | Presensi, penjadwalan shift, lembur (overtime), kalkulasi keterlambatan |
| **Leave & Time Off** | 🔴 High | Pengajuan cuti/sakit/izin, kuota tahunan, multi-level approval |
| **Payroll & Compensation Engine** | 🔴 High | Kalkulasi gaji, tunjangan/potongan, PPh 21, BPJS, slip gaji digital |
| **Performance Management** | 🟡 Medium | KPI, OKR, review 360 terintegrasi Job Management & Competency |
| **Reimbursement & Claim** | 🟡 Medium | Klaim kesehatan & operasional dinas |
| **Recruitment & Onboarding (ATS)** | 🟡 Medium | Kandidat, alur seleksi, otomatisasi onboarding |

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
Semua tenant memiliki struktur database yang sama (106 tables).

## 5. Tenant Provisioning Engine

Saat company dibuat: 1. Validasi lisensi. 2. Membuat database. 3.
Membuat database user. 4. Menyimpan koneksi tenant. 5. Menjalankan core
migration (11 files → 106 tables). 6. Menjalankan migration modul aktif.
7. Menjalankan seeder. 8. Membuat Super Admin. 9. Mengaktifkan modul. 10.
Audit log. 11. Health check tenant.

### End-to-End Test Results ✅

| Item | Status | Detail |
|------|:------:|--------|
| Company status | ✅ active | API mengembalikan `status: "active"` |
| Tenant database | ✅ Created | Database tenant ter-create di MySQL |
| Tenant connection | ✅ Saved | Record di `tenant_connections` tersimpan |
| Migrations | ✅ 11 files | 001 → 011 sukses semua |
| Total tables | ✅ 106 tables | Setiap migrasi menciptakan tabel sesuai DDL |

## 6. Gap Analysis & Technical Recommendations

Berdasarkan tinjauan arsitektur, berikut area kritis untuk production readiness:

### 6.1 Tenant Lifecycle & Resource Cleanup ✅
**Masalah:** Penanganan status tenant berisiko menyisakan koneksi TCP/GORM.
**Status:** ✅ **Sudah diimplementasikan** — `CloseTenantConnection(companyID)` sudah ada dan dipanggil di seluruh lifecycle (Suspend/Activate/Terminate).

### 6.2 Keamanan Kredensial Database Tenant 🔴 High
**Masalah:** Kredensial koneksi DB tenant disimpan dalam bentuk plain text.
**Solusi:** Enkripsi AES-256-GCM pada kolom `TenantConnection.Password` via `HRIS_ENCRYPTION_KEY`.

### 6.3 Optimasi Connection Pooling 🟡 Medium
**Masalah:** `SetMaxOpenConns(100)` per tenant berisiko habiskan kuota koneksi.
**Solusi:** Pool 10–20 koneksi per tenant + PgBouncer di production.

### 6.4 Penanganan Dialek SQL Migrasi 🔴 High
**Masalah:** Perbedaan sintaksis DDL PostgreSQL vs MySQL pada file .sql mentah.
**Solusi:** Pisahkan migrasi per driver atau gunakan Goose/Atlas.

### 6.5 Sinkronisasi Cache Terdistribusi 🟡 Medium
**Masalah:** Update cache Redis perlu dikonsumsi konsisten di multi-node.
**Solusi:** Redis Pub/Sub untuk invalidasi cache real-time.

### Priority Matrix

| Area | Komponen | Prioritas | Action Item Utama |
| :--- | :--- | :---: | :--- |
| **Security** | Tenant Credentials | 🔴 High | Enkripsi AES-256-GCM pada `tenant_connections.password` |
| **Database** | SQL Dialect | 🔴 High | Pemisahan DDL PostgreSQL dan MySQL |
| **Performance** | Connection Pool | 🟡 Medium | Pool 10-20 per tenant + PgBouncer |
| **Architecture** | Cache Sync | 🟡 Medium | Redis Pub/Sub untuk invalidasi cache |

## 7. Platform Management

-   Company Management
-   Module Management
-   License Management
-   Feature Flag
-   System Configuration
-   Tenant Monitoring
-   Tenant Backup & Restore

## 8. Module SDK

Setiap modul wajib memiliki: - manifest.yaml - routes - handlers -
services - repositories - entities - dto - permissions - menus -
migrations - seeders - configs - tests - api docs

## 9. Backend Architecture

-   Go 1.22+
-   Gin
-   GORM (PostgreSQL & MySQL)
-   Redis + Asynq
-   Custom RBAC (role-based permission dengan hierarchy)
-   Zap (structured logging)
-   OpenTelemetry (optional)

## 10. Frontend Architecture

-   Vue 3
-   TypeScript
-   Vite
-   PrimeVue
-   Tailwind CSS
-   Pinia
-   TanStack Query

## 11. Monorepo

```text
hris-platform/
├── backend/
│   ├── cmd/
│   │   ├── server/          # API server entry point
│   │   └── installer/       # CLI installer (provisioning)
│   ├── internal/
│   │   ├── platform/        # Platform Management modules
│   │   ├── modules/         # Tenant modules
│   │   └── pkg/             # Shared Kernel
│   ├── config/
│   └── docker/
├── frontend/                # Vue 3 (future)
├── docker/
└── docs/
```

## 12. Security

-   JWT (access + refresh token)
-   RBAC (4 roles: super_admin, company_admin, manager, employee)
-   Audit Log
-   Rate Limiting
-   MFA Ready
-   OWASP ASVS

## 13. CI/CD

-   GitHub Actions
-   Docker
-   Auto Test
-   Auto Build
-   Auto Deploy

## 14. Testing

-   Unit (Go testing + testify)
-   Integration (SQLite in-memory)
-   API
-   E2E (Playwright)
-   UAT
-   Performance

## 15. Disaster Recovery

-   Backup per tenant
-   Restore per tenant
-   Point-in-time recovery
-   Disaster recovery procedure

## 16. Future Roadmap

-   Plugin Marketplace
-   AI Assistant
-   BI Dashboard
-   Mobile Apps
-   Public API
-   Kubernetes Deployment
