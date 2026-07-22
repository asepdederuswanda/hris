# HRIS Platform

**Human Resource Information System — Enterprise Grade**

Modular monolith backend untuk platform HRIS enterprise dengan arsitektur multi-tenant, dibangun menggunakan **Go** + **Gin** + **GORM**.

---

## 📋 Daftar Isi

- [Arsitektur](#arsitektur)
- [Tech Stack](#tech-stack)
- [Struktur Proyek](#struktur-proyek)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Konfigurasi](#konfigurasi)
- [Environment Variables](#environment-variables)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Docker](#docker)
- [Testing](#testing)
- [Module SDK](#module-sdk)
- [Database Migration](#database-migration)
- [Tenant Provisioning](#tenant-provisioning-end-to-end-verified-)
- [Roadmap](#roadmap)
- [Dokumentasi Lainnya](#dokumentasi-lainnya)

---

## 🏗️ Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Modular Monolith                       │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐ │
│  │  Platform    │  │   Shared    │  │   Tenant Modules     │ │
│  │  Management  │  │   Kernel    │  │                      │ │
│  │              │  │             │  │  ┌──────────────────┐│ │
│  │ • Company    │  │ • Config    │  │  │ Organization     ││ │
│  │ • Module     │  │ • Database  │  │  │ Employee         ││ │
│  │ • License    │  │ • Auth/JWT  │  │  │ Attendance       ││ │
│  │ • User       │  │ • Middleware│  │  │ Leave            ││ │
│  │ • Feature    │  │ • Router   │  │  │ Payroll          ││ │
│  │   Flag       │  │ • Logger   │  │  │ Competency       ││ │
│  │              │  │ • Module   │  │  │ Job Management   ││ │
│  │              │  │   SDK      │  │  │ Approval         ││ │
│  └─────────────┘  └─────────────┘  │  └──────────────────┘│ │
│                                     └──────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
         │                    │                      │
         ▼                    ▼                      ▼
  ┌─────────────┐   ┌──────────────┐      ┌──────────────────┐
  │ Platform DB  │   │  Tenant 1 DB │      │   Tenant 2 DB    │
  │ (PostgreSQL) │   │  (PostgreSQL │      │   (MySQL)        │
  │  or MySQL)   │   │   or MySQL)  │      │    or MySQL)     │
  └─────────────┘   └──────────────┘      └──────────────────┘
```

### Design Principles

| Prinsip | Deskripsi |
|---|---|
| **Modular Monolith** | Satu binary Go, multiple internal modules dengan batasan tegas |
| **Multi-Tenant** | Database per tenant — isolasi data total antar company |
| **Module SDK** | Setiap modul mengikuti kontrak interface yang sama |
| **Multi-Driver DB** | Support PostgreSQL & MySQL, bisa berbeda per tenant |
| **API-First** | REST API dengan JWT authentication |
| **Security by Design** | JWT + RBAC + Audit Log + Rate Limiting |

---

## ⚙️ Tech Stack

| Layer | Teknologi |
|---|---|
| **Bahasa** | Go 1.22+ |
| **HTTP Framework** | Gin |
| **ORM** | GORM (PostgreSQL & MySQL) |
| **Database** | PostgreSQL 16+ / MySQL 8+ |
| **Cache** | Redis 7+ |
| **Queue** | Asynq (Redis-based) |
| **Auth** | JWT (golang-jwt) |
| **Config** | Viper + godotenv |
| **Logging** | Zap (structured) |
| **Tracing** | OpenTelemetry (optional) |
| **Container** | Docker + Docker Compose |

---

## 📁 Struktur Proyek

```
hris-platform/
├── backend/                          # Go backend
│   ├── cmd/
│   │   ├── server/main.go            # Entry point API server
│   │   └── installer/main.go         # CLI installer (provisioning)
│   ├── internal/
│   │   ├── platform/                 # Platform Management Module
│   │   │   └── company/              #   Company CRUD + TenantConnection
│   │   ├── modules/                  # Tenant Modules
│   │   │   ├── organization/         #   Organization tree CRUD
│   │   │   ├── employee/             #   Employee CRUD + 8 sub-modules
│   │   │   └── jobmanagement/        #   Job Management (18 entities)
│   │   └── pkg/                      # Shared Kernel
│   │       ├── config/               # Viper configuration loader
│   │       ├── database/             # Multi-tenant DB manager
│   │       ├── driver/               # Shared DB driver type
│   │       ├── auth/                 # JWT generation & validation
│   │       ├── middleware/           # Gin middleware
│   │       ├── router/               # Router setup & module registration
│   │       ├── logger/               # Zap logger
│   │       └── module/               # Module SDK interface
│   ├── config/
│   │   └── config.yaml               # Base configuration
│   └── internal/
│       └── pkg/
│           └── migrator/             # SQL Migration Engine
│               ├── migrator.go       #  Core: Up(), Down(), DownTo()
│               ├── embed.go          #  //go:embed migrations/
│               └── migrations/
│                   ├── platform/     #  Platform DDL (6 files)
│                   ├── seeders/      #  Seed data (1 file)
│                   └── tenant/       #  Tenant template (future)
│   ├── docker/
│   │   └── Dockerfile
│   ├── .env.example                  # Environment template
│   ├── Makefile
│   ├── go.mod
│   └── go.sum
├── frontend/                         # Vue 3 + TypeScript (future)
├── docker/
│   └── docker-compose.yml            # Full infra compose
└── docs/
    ├── platform-architecture-design.md
    └── analisis-blueprint-vs-existing.md
```

---

## ✅ Prerequisites

- **Go** 1.22+
- **PostgreSQL** 16+ atau **MySQL** 8+
- **Redis** 7+
- **Docker** & Docker Compose (optional, untuk development environment)

---

## 🚀 Quick Start

### 1. Clone & Setup

```bash
git clone <repo-url> hris-platform
cd hris-platform/backend

# Copy environment
cp .env.example .env

# Edit .env sesuai environment lokal
nano .env

# Download dependencies
go mod download
```

> ⚠️ **Jangan commit file `.env` ke repository!** File ini berisi secret/kredensial lokal. Sudah seharusnya di-*gitignore*.

### 2. Database Setup

Buat database platform:

```bash
# PostgreSQL
createdb hris_platform

# atau via psql
psql -U postgres -c "CREATE DATABASE hris_platform;"
```

### 3. Run Server

```bash
# Development mode
make run

# atau langsung
go run ./cmd/server --config ./config/config.yaml
```

Server akan berjalan di `http://localhost:8080`.

### 4. Verify

```bash
curl http://localhost:8080/healthz
# {"status":"ok","service":"hris-platform"}
```

---

## 🔧 Konfigurasi

### Priority Loading

Konfigurasi di-load dengan urutan prioritas (low → high):

```
1. Default values (hardcoded)
2. config/config.yaml
3. .env file (via godotenv)
4. OS environment variables (HRIS_* prefix)
```

### Configuration File

```yaml
# backend/config/config.yaml
server:
  port: "8080"
  mode: "debug"               # debug | release | test

database:
  driver: "postgres"           # postgres | mysql
  platform_host: "localhost"
  platform_port: 5432
  platform_db: "hris_platform"
  # ... lihat config.yaml untuk lengkapnya
```

### Berganti ke MySQL

```yaml
database:
  driver: "mysql"
  platform_host: "localhost"
  platform_port: 3306
```

---

## 🌐 Environment Variables

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp backend/.env.example backend/.env
```

| Variable | Default | Deskripsi |
|---|---|---|
| `HRIS_SERVER_PORT` | `8080` | Port server |
| `HRIS_SERVER_MODE` | `debug` | Mode: debug/release/test |
| `HRIS_DATABASE_DRIVER` | `postgres` | Driver: postgres/mysql |
| `HRIS_DATABASE_PLATFORM_HOST` | `localhost` | Host platform DB |
| `HRIS_DATABASE_PLATFORM_PORT` | `5432` | Port platform DB |
| `HRIS_DATABASE_PLATFORM_DB` | `hris_platform` | Nama platform DB |
| `HRIS_REDIS_HOST` | `localhost` | Host Redis |
| `HRIS_REDIS_PORT` | `6379` | Port Redis |
| `HRIS_JWT_SECRET` | - | Secret key untuk JWT |
| `HRIS_JWT_ACCESS_TOKEN_TTL` | `15` | Access token expiry (menit) |
| `HRIS_JWT_REFRESH_TOKEN_TTL` | `24` | Refresh token expiry (jam) |
| `HRIS_LOGGER_LEVEL` | `info` | Log level |
| `HRIS_CORS_ALLOWED_ORIGINS` | `*` | CORS origins |

> Lihat [`.env.example`](backend/.env.example) untuk daftar lengkap.

---

## 📖 API Documentation

### Interactive Docs (Scalar UI)

Dokumentasi API interaktif tersedia setelah server berjalan:

```
http://localhost:8080/docs          # Scalar UI (interaktif)
http://localhost:8080/openapi.json  # OpenAPI 3.0 spec (JSON)
```

### Health Check

```bash
GET  /healthz          # Server health
GET  /readyz           # Readiness check
```

---

### 🔒 RBAC (Role-Based Access Control)

Platform menggunakan RBAC dengan **4 role** dalam hierarki:

```
super_admin  (full access — platform + tenant)
  └── company_admin  (platform view + tenant full)
        └── manager  (tenant view/create/update)
              └── employee  (tenant view-only)
```

**Inheritance behavior:**
- Jika suatu role **tidak memiliki** policy untuk resource → inherit dari parent
- Jika suatu role **memiliki** policy (dengan daftar action terbatas) → action yang tidak tercantum = intentional deny
- Contoh: `employee` punya `organization: view` → `create` akan **denied**, tidak inherit dari `manager`

#### Platform Resources

| Role | Resource | Actions |
|------|----------|--------|
| **super_admin** | `*` (semua) | `*` (semua) |
| **company_admin** | `company`, `user`, `license` | `view` |

**Detail Permission per Endpoint (Company):**

| Permission | super_admin | company_admin | Deskripsi |
|-----------|:-----------:|:-------------:|-----------|
| `company.view` | ✅ | ✅ | List & detail companies |
| `company.create` | ✅ | ❌ | Create company + provision |
| `company.update` | ✅ | ❌ | Update company info |
| `company.delete` | ✅ | ❌ | Soft delete (deactivate) |
| `company.suspend` | ✅ | ❌ | Suspend + deactivate connection |
| `company.activate` | ✅ | ❌ | Activate + reactivate connection |
| `company.terminate` | ✅ | ❌ | Terminate + drop database |
| `company.backup` | ✅ | ❌ | Backup tenant (Phase 2) |
| `company.restore` | ✅ | ❌ | Restore tenant (Phase 2) |

#### Tenant Resources

| Role | Organization | Employee | Attendance | Leave | Payroll | Approval |
|------|:-----------:|:--------:|:----------:|:-----:|:-------:|:--------:|
| **super_admin** | ✅ * | ✅ * | ✅ * | ✅ * | ✅ * | ✅ * |
| **company_admin** | ✅ * | ✅ * | ✅ * | ✅ * | ✅ * | ✅ * |
| **manager** | ✅ V/C/U | ✅ V/C/U | ✅ V | ✅ V/C | ❌ | ✅ V/C/U |
| **employee** | ✅ V | ✅ V | ✅ V | ✅ V | ✅ V | ❌ |

> * = full access, V = view, C = create, U = update

**Permission format:** `resource.action` (contoh: `company.create`, `user.view`)

Endpoint yang tidak memiliki akses akan mengembalikan **403 Forbidden**:

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to perform this action",
    "details": {
      "role": "company_admin",
      "resource": "company",
      "action": "create"
    }
  }
}
```

---

### Platform API (`/api/v1/platform`)

> 💡 **Auth required**: ✅ = JWT + RBAC, 🔓 = Public

#### 🔐 Auth

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `POST` | `/login` | Login platform admin | 🔓 Public | - |
| `POST` | `/refresh` | Refresh access token | 🔓 Public | - |

#### 👥 Users

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/users` | List platform users | ✅ | `user.view` (super_admin, company_admin) |
| `POST` | `/users` | Create platform user | ✅ | `user.create` (super_admin only) |
| `GET` | `/users/:id` | Get user detail | ✅ | `user.view` (super_admin, company_admin) |
| `PUT` | `/users/:id` | Update user | ✅ | `user.update` (super_admin only) |

#### 🏢 Companies — Tenant Lifecycle Management

| Method | Endpoint | Deskripsi | Auth | RBAC | Cleanup? |
|---|---|---|---|---|---|
| `POST` | `/companies` | Create company + provision tenant DB 🔄 | ✅ | `company.create` (super_admin only) | ❌ |
| `GET` | `/companies` | List all companies | ✅ | `company.view` (super_admin, company_admin) | ❌ |
| `GET` | `/companies/:id` | Get company detail | ✅ | `company.view` (super_admin, company_admin) | ❌ |
| `PUT` | `/companies/:id` | Update company info | ✅ | `company.update` (super_admin only) | ❌ |
| `DELETE` | `/companies/:id` | **Soft delete** — deactivate connection + deleted_at | ✅ | `company.delete` (super_admin only) | ✅ Deactivate |
| `POST` | `/companies/:id/suspend` | **Suspend** — deactivate connection, set status suspended | ✅ | `company.suspend` (super_admin only) | ✅ Deactivate |
| `POST` | `/companies/:id/activate` | **Activate** — reactivate connection, set status active | ✅ | `company.activate` (super_admin only) | ✅ Reactivate |
| `POST` | `/companies/:id/terminate` | **Terminate** — drop database + remove connection, set status terminated | ✅ | `company.terminate` (super_admin only) | ✅ Drop DB |
| `POST` | `/companies/:id/backup` | Trigger tenant backup (Phase 2) | ✅ | `company.backup` (super_admin only) | ❌ |
| `POST` | `/companies/:id/restore` | Trigger tenant restore (Phase 2) | ✅ | `company.restore` (super_admin only) | ❌ |

**Lifecycle Actions vs Database Cleanup:**

| Action | Company Status | Tenant Connection | Tenant Database |
|--------|---------------|-------------------|-----------------|
| `Create` | `active` | `is_active = true` | ✅ Created + migrated (106 tables) |
| `Suspend` | `suspended` | `is_active = false` + cache cleared | ✅ Data preserved |
| `Activate` | `active` | `is_active = true` + cache cleared | ✅ Reconnected |
| `Soft Delete` | (hidden via `deleted_at`) | `is_active = false` + cache cleared | ✅ Data preserved |
| `Terminate` | `terminated` | ❌ Removed from DB | ❌ Dropped |

#### 🧩 Modules

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/modules` | List all registered modules | ✅ | `module.view` (super_admin only) |
| `POST` | `/modules` | Register new module | ✅ | `module.create` (super_admin only) |
| `GET` | `/modules/:id` | Get module detail | ✅ | `module.view` (super_admin only) |
| `PUT` | `/modules/:id` | Update module | ✅ | `module.update` (super_admin only) |
| `GET` | `/modules/:id/companies` | List companies using module | ✅ | `module.view` (super_admin only) |
| `POST` | `/modules/:id/activate` | Activate module for company | ✅ | `module.activate` (super_admin only) |
| `POST` | `/modules/:id/deactivate` | Deactivate module for company | ✅ | `module.deactivate` (super_admin only) |

#### 🔑 Licenses

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/licenses` | List all licenses | ✅ | `license.view` (super_admin, company_admin) |
| `POST` | `/licenses` | Create license for company | ✅ | `license.create` (super_admin only) |
| `GET` | `/licenses/:id` | Get license detail | ✅ | `license.view` (super_admin, company_admin) |
| `PUT` | `/licenses/:id` | Update license | ✅ | `license.update` (super_admin only) |

#### 📊 Monitoring

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/monitoring/health` | Platform health + DB status | ✅ | `monitoring.view` (super_admin only) |
| `GET` | `/monitoring/tenants` | List tenant connections health | ✅ | `monitoring.view` (super_admin only) |
| `GET` | `/monitoring/tenants/:id` | Get tenant health detail | ✅ | `monitoring.view` (super_admin only) |

### Tenant API (`/api/v1/tenant`)

Semua endpoint tenant memerlukan **JWT Bearer Token** di header:
```
Authorization: Bearer <access_token>
```

| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/organizations` | List organizations |
| `POST` | `/organizations` | Create organization |
| `GET` | `/organizations/:id` | Get organization |
| `PUT` | `/organizations/:id` | Update organization |
| `DELETE` | `/organizations/:id` | Delete organization |
| `GET` | `/organizations?tree=true` | Get organization tree |
| `GET` | `/employees` | List employees (pagination) |
| `POST` | `/employees` | Create employee |
| `GET` | `/employees/:id` | Get employee with all sub-modules |
| `PUT` | `/employees/:id` | Update employee |
| `DELETE` | `/employees/:id` | Delete employee (hard delete) |
| `POST` | `/employees/:id/addresses` | Create address (type: MAIN/DOMICILE) |
| `PUT` | `/employees/:id/addresses/:addressId` | Update address |
| `DELETE` | `/employees/:id/addresses/:addressId` | Delete address |
| `POST` | `/employees/:id/emergency-contacts` | Create emergency contact |
| `PUT` | `/employees/:id/emergency-contacts/:contactId` | Update emergency contact |
| `DELETE` | `/employees/:id/emergency-contacts/:contactId` | Delete emergency contact |
| `POST` | `/employees/:id/families` | Create family member |
| `PUT` | `/employees/:id/families/:familyId` | Update family member |
| `DELETE` | `/employees/:id/families/:familyId` | Delete family member |
| `POST` | `/employees/:id/educations` | Create education record |
| `PUT` | `/employees/:id/educations/:educationId` | Update education |
| `DELETE` | `/employees/:id/educations/:educationId` | Delete education |
| `POST` | `/employees/:id/experiences` | Create work experience |
| `PUT` | `/employees/:id/experiences/:experienceId` | Update experience |
| `DELETE` | `/employees/:id/experiences/:experienceId` | Delete experience |
| `POST` | `/employees/:id/documents` | Create document |
| `PUT` | `/employees/:id/documents/:documentId` | Update document |
| `DELETE` | `/employees/:id/documents/:documentId` | Delete document |
| `POST` | `/employees/:id/insurances` | Create insurance (BPJS) |
| `PUT` | `/employees/:id/insurances/:insuranceId` | Update insurance |
| `DELETE` | `/employees/:id/insurances/:insuranceId` | Delete insurance |
| `POST` | `/employees/:id/employments` | Create employment record |
| `PUT` | `/employees/:id/employments/:employmentId` | Update employment |
| `DELETE` | `/employees/:id/employments/:employmentId` | Delete employment |

**Job Management — Titles (9.1)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/titles` | List job titles |
| `POST` | `/job-management/titles` | Create job title |
| `GET` | `/job-management/titles/:id` | Get job title with subs |
| `PUT` | `/job-management/titles/:id` | Update job title |
| `DELETE` | `/job-management/titles/:id` | Delete job title |
| `POST` | `/job-management/titles/:id/subs` | Create job title sub |
| `GET` | `/job-management/titles/:id/subs` | List subs under a title |
| `GET` | `/job-management/titles/:id/subs/:subId` | Get job title sub |
| `PUT` | `/job-management/titles/:id/subs/:subId` | Update job title sub |
| `DELETE` | `/job-management/titles/:id/subs/:subId` | Delete job title sub |

**Job Management — Values (9.3)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/values` | List job values |
| `POST` | `/job-management/values` | Create job value |
| `GET` | `/job-management/values/:id` | Get job value |
| `PUT` | `/job-management/values/:id` | Update job value |
| `DELETE` | `/job-management/values/:id` | Delete job value |

**Job Management — Objectives (9.4)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/objectives` | List job objectives |
| `POST` | `/job-management/objectives` | Create job objective |
| `GET` | `/job-management/objectives/:id` | Get job objective |
| `PUT` | `/job-management/objectives/:id` | Update job objective |
| `DELETE` | `/job-management/objectives/:id` | Delete job objective |

**Job Management — Identifications (9.5)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/identifications` | List job identifications |
| `POST` | `/job-management/identifications` | Create job identification |
| `GET` | `/job-management/identifications/:id` | Get job identification |
| `PUT` | `/job-management/identifications/:id` | Update job identification |
| `DELETE` | `/job-management/identifications/:id` | Delete job identification |

**Job Management — Responsibilities (9.6)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/responsibilities` | List job responsibilities |
| `POST` | `/job-management/responsibilities` | Create job responsibility |
| `GET` | `/job-management/responsibilities/:id` | Get job responsibility |
| `PUT` | `/job-management/responsibilities/:id` | Update job responsibility |
| `DELETE` | `/job-management/responsibilities/:id` | Delete job responsibility |

**Job Management — Education Experiences (9.7)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/education-experiences` | List education experiences |
| `POST` | `/job-management/education-experiences` | Create education experience |
| `GET` | `/job-management/education-experiences/:id` | Get education experience |
| `PUT` | `/job-management/education-experiences/:id` | Update education experience |
| `DELETE` | `/job-management/education-experiences/:id` | Delete education experience |

**Job Management — HR Authorities (9.8)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/hr-authorities` | List HR authorities |
| `POST` | `/job-management/hr-authorities` | Create HR authority |
| `GET` | `/job-management/hr-authorities/:id` | Get HR authority |
| `PUT` | `/job-management/hr-authorities/:id` | Update HR authority |
| `DELETE` | `/job-management/hr-authorities/:id` | Delete HR authority |

**Job Management — Operational Authorities (9.9)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/operational-authorities` | List operational authorities |
| `POST` | `/job-management/operational-authorities` | Create operational authority |
| `GET` | `/job-management/operational-authorities/:id` | Get operational authority |
| `PUT` | `/job-management/operational-authorities/:id` | Update operational authority |
| `DELETE` | `/job-management/operational-authorities/:id` | Delete operational authority |

**Job Management — Working Activities (9.10)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/working-activities` | List working activities |
| `POST` | `/job-management/working-activities` | Create working activity |
| `GET` | `/job-management/working-activities/:id` | Get working activity |
| `PUT` | `/job-management/working-activities/:id` | Update working activity |
| `DELETE` | `/job-management/working-activities/:id` | Delete working activity |

**Job Management — Working Risks (9.11)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/working-risks` | List working risks |
| `POST` | `/job-management/working-risks` | Create working risk |
| `GET` | `/job-management/working-risks/:id` | Get working risk |
| `PUT` | `/job-management/working-risks/:id` | Update working risk |
| `DELETE` | `/job-management/working-risks/:id` | Delete working risk |

**Job Management — Relationships (9.12)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/relationships` | List job relationships |
| `POST` | `/job-management/relationships` | Create job relationship |
| `GET` | `/job-management/relationships/:id` | Get job relationship |
| `PUT` | `/job-management/relationships/:id` | Update job relationship |
| `DELETE` | `/job-management/relationships/:id` | Delete job relationship |

**Job Management — Subordinate Controls (9.13)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/subordinate-controls` | List subordinate controls |
| `POST` | `/job-management/subordinate-controls` | Create subordinate control |
| `GET` | `/job-management/subordinate-controls/:id` | Get subordinate control |
| `PUT` | `/job-management/subordinate-controls/:id` | Update subordinate control |
| `DELETE` | `/job-management/subordinate-controls/:id` | Delete subordinate control |

**Job Management — Assets (9.14)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/assets` | List job assets |
| `POST` | `/job-management/assets` | Create job asset |
| `GET` | `/job-management/assets/:id` | Get job asset |
| `PUT` | `/job-management/assets/:id` | Update job asset |
| `DELETE` | `/job-management/assets/:id` | Delete job asset |

**Job Management — Financials (9.15)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/financials` | List job financials |
| `POST` | `/job-management/financials` | Create job financial |
| `GET` | `/job-management/financials/:id` | Get job financial |
| `PUT` | `/job-management/financials/:id` | Update job financial |
| `DELETE` | `/job-management/financials/:id` | Delete job financial |

**Job Management — Potency Competencies (9.16)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/potency-competencies` | List potency competencies |
| `POST` | `/job-management/potency-competencies` | Create potency competency |
| `GET` | `/job-management/potency-competencies/:id` | Get potency competency |
| `PUT` | `/job-management/potency-competencies/:id` | Update potency competency |
| `DELETE` | `/job-management/potency-competencies/:id` | Delete potency competency |

**Job Management — Scores (9.17)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/scores` | List all scores |
| `GET` | `/job-management/scores/org/:orgId` | Get score by organization |
| `PUT` | `/job-management/scores/org/:orgId` | Upsert score for organization |

**Job Management — Competency Groups (9.18)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/job-management/competency-groups` | Create competency group |
| `GET` | `/job-management/competency-groups` | List competency groups |
| `GET` | `/job-management/competency-groups/:id` | Get competency group |
| `PUT` | `/job-management/competency-groups/:id` | Update competency group |
| `DELETE` | `/job-management/competency-groups/:id` | Delete competency group |

### Response Format

```json
// Success
{
    "success": true,
    "data": { ... },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 100
    }
}

// Error
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Email is required"
    }
}
```

---

## 💻 Development

### Makefile Commands

Jalankan dari direktori `backend/`:

```bash
make build            # Build binary
make run              # Run server
make run-hot          # Run with hot reload (air)
make test             # Run all tests
make lint             # Run linter
make vet              # Run go vet
make coverage         # Run tests with coverage report
make docker           # Build Docker image
make tidy             # Tidy dependencies
make clean            # Clean build artifacts
make help             # Show all commands
```

### Menambah Module Baru

Setiap modul harus mengimplementasikan `module.Module` interface:

```go
type Module interface {
    Info() ModuleInfo
    RegisterRoutes(router *gin.RouterGroup)
    Migrate(db *gorm.DB) error
    Seed(db *gorm.DB) error
    Permissions() []string
}
```

Template struktur module:

```text
internal/modules/{nama_module}/
├── handler.go
├── service.go
├── repository.go
├── model.go
├── dto.go
├── routes.go
└── module.go
```

---

## 🐳 Docker

### Development Environment

```bash
# Start semua services
make docker-compose-up

# Services:
# - PostgreSQL :5432
# - Redis      :6379
# - API       :8080
# - Asynqmon  :8081

# Stop
make docker-compose-down
```

### Build Docker Image

```bash
make docker
docker run -p 8080:8080 hris-platform:latest
```

---

## 🧪 Testing

```bash
# Semua test
make test

# Dengan coverage
make coverage     # hasil: coverage.html

# Short test (tanpa integration)
make test-short
```

---

## 🧩 Module SDK

Setiap modul di HRIS Platform harus mengikuti kontrak Module SDK:

1. **ModuleInfo** — identitas modul (nama, slug, versi, dependensi)
2. **RegisterRoutes** — daftarkan semua endpoint HTTP
3. **Migrate** — migrasi database modul
4. **Seed** — data awal / seeder
5. **Permissions** — daftar permission yang dibutuhkan

```yaml
# manifest.yaml (contoh)
name: "Organization Management"
slug: "organization"
version: "1.0.0"
depends_on: []
permissions:
  - "organization.view"
  - "organization.create"
menus:
  - name: "Organization"
    icon: "building"
    route: "/admin/organizations"
```

---

## 🔄 Database Migration

Platform menggunakan **Go SQL Migration Runner** (`internal/pkg/migrator/`) dengan dua mode:

### SQL Migrator (file `.sql` embedded)

File SQL migration di-embed ke binary Go via `//go:embed` — tidak perlu file eksternal saat runtime.

**File Convention:**
```text
001_create_companies.sql          ← Up migration (DDL)
001_create_companies.down.sql     ← Down migration (rollback) — optional
```

**Startup Flow:**
```text
[1] SQL Migrator (platform/*.sql)  → CREATE TABLE, indexes, FKs
[2] GORM AutoMigrate               → Sync Go struct columns
[3] SQL Seeders (seeders/*.sql)    → INSERT initial data
[4] Module Seeders                 → Business logic seed data
```

### CLI Usage

```bash
# Normal startup (run all pending migrations + seeders)
go run ./cmd/server --config ./config/config.yaml

# Rollback ALL applied migrations (down files required)
go run ./cmd/server --config ./config/config.yaml --migrate-down

# Rollback to specific version (exclusive)
# Contoh: applied [001,002,003,004,005,006]
# --migrate-to 004 → rollback 006, 005 (004 tetap applied)
go run ./cmd/server --config ./config/config.yaml --migrate-to 004
```

> ⚠️ **--migrate-down** dan **--migrate-to** bersifat mutually exclusive.
> Setelah selesai, server akan exit (tidak start).

### Tracking Table

Setiap migration yang sukses dicatat di tabel `schema_migrations`:

```sql
CREATE TABLE schema_migrations (
    version     VARCHAR(14) PRIMARY KEY,   -- '001', '002'
    name        VARCHAR(255) NOT NULL,      -- 'create_companies'
    applied_at  TIMESTAMP NOT NULL,         -- waktu eksekusi
    checksum    VARCHAR(64) NOT NULL,
    file_path   VARCHAR(500) NOT NULL
);
```

### Migration Files

- `internal/pkg/migrator/migrations/platform/` — 6 platform DDL files
- `internal/pkg/migrator/migrations/seeders/` — 1 seeder file
- `internal/pkg/migrator/migrations/tenant/` — Tenant template (future)

### Transaction Safety

- **Up**: SQL content + insert `schema_migrations` record dalam 1 database transaction
- **Down**: Down SQL + delete `schema_migrations` record dalam 1 database transaction
- Jika salah satu gagal → entire transaction rollback

---

## 🏗️ Tenant Provisioning (End-to-End Verified ✅)

Platform memiliki **Tenant Provisioning Engine** yang secara otomatis membuat database dan menjalankan migrasi saat company baru dibuat.

### Provisioning Flow

```
POST /api/v1/platform/companies
    ↓
1. Generate slug dari company name
2. Cek duplikasi slug (tambahkan UUID suffix jika perlu)
3. Save company ke platform DB (status: active)
4. provisionTenant():
   ├── a. Generate db_name = hris_{slug}
   ├── b. Connect sebagai superuser (root@localhost)
   ├── c. Buat database tenant (CREATE DATABASE IF NOT EXISTS)
   ├── d. Simpan TenantConnection ke platform DB (ID = companyID)
   ├── e. Connect ke tenant DB via GORM
   └── f. Jalankan 11 tenant SQL migrations (106 tables)
5. Jika provisioning berhasil → company status: active
6. Jika provisioning gagal → company status: suspended (data tetap tersimpan)
```

### Tenant Migration Files (11 files → 106 tables)

| File | Isi |
|------|-----|
| `001_master_data.sql` | Master tables (religions, educations, countries, provinces, dll) |
| `002_organization.sql` | Organization structure, zones, job families, positions |
| `003_employee.sql` | Employees, employments, families, educations, documents |
| `004_attendance.sql` | Attendance settings, shifts, events, overtime |
| `005_leave.sql` | Leave types, requests, accrual policies, company holidays |
| `006_payroll_structure.sql` | Salary components, grades, payroll periods |
| `007_approval.sql` | Approval flows, steps, instances, tasks |
| `008_competency.sql` | Competencies, values, events, scores + FK dari migration 002 |
| `009_job_management.sql` | Job titling, values, objectives, responsibilities |
| `010_permissions.sql` | Roles, permissions, model_has_roles/ permissions |
| `011_pph21.sql` | PPh21 settings, tax brackets, PTKP rates |

> **Catatan:** Total 106 tabel termasuk `schema_migrations` (auto-created oleh migrator engine).

### Daftar Lengkap 106 Tabel Tenant

**Approval (5):**
`approval_actions`, `approval_flow_steps`, `approval_flows`, `approval_instances`, `approval_tasks`

**Attendance (10):**
`attendance_company_settings`, `attendance_company_shifts`, `attendance_device_captures`,
`attendance_employee_shifts`, `attendance_events`, `attendance_exempt_positions`,
`attendance_face_captures`, `attendance_locations`, `attendance_overtime_requests`, `attendance_sessions`

**BPJS (2):**
`bpjs_rate_components`, `bpjs_settings`

**Competency (7):**
`competence_values`, `competencies`, `competency_event_targets`, `competency_events`,
`competency_score_details`, `competency_scores`, `competency_values`

**Employee (14):**
`emergency_contacts`, `employee_addresses`, `employee_bank_profiles`, `employee_bpjs_profiles`,
`employee_documents`, `employee_educations`, `employee_experiences`, `employee_families`,
`employee_insurances`, `employee_leave_balances`, `employee_payroll_profiles`, `employee_tax_profiles`,
`employees`, `employments`

**Job Management (20):**
`job_families`, `job_family_competencies`, `job_management_assets`, `job_management_competency_groups`,
`job_management_education_experiences`, `job_management_financials`, `job_management_hr_authorities`,
`job_management_identifications`, `job_management_objectives`, `job_management_operational_authorities`,
`job_management_potency_competencies`, `job_management_relationships`, `job_management_responsibilities`,
`job_management_scores`, `job_management_subordinate_controls`, `job_management_title_subs`,
`job_management_titles`, `job_management_values`, `job_management_working_activities`,
`job_management_working_risks`

**Leave (6):**
`company_holidays`, `leave_accrual_policies`, `leave_reasons`, `leave_request_details`, `leave_requests`, `leave_types`

**Master Data (12):**
`countries`, `districts`, `document_templates`, `educations`, `employment_statuses`,
`gradings`, `marital_statuses`, `provinces`, `regencies`, `relationship_types`, `religions`, `villages`

**Organization (5):**
`organization_levels`, `organization_summaries`, `organizations`, `positions`, `zones`

**Payroll (16):**
`payroll_payslips`, `payroll_periods`, `payroll_profile_change_logs`, `payroll_run_employees`,
`payroll_run_items`, `payroll_runs`, `pph21_calculation_logs`, `pph21_ptkp_rates`,
`pph21_settings`, `pph21_tax_brackets`, `ptkps`, `salary_change_logs`,
`salary_components`, `salary_employee_adjustments`, `salary_employee_components`, `salary_grade_components`

**Settings & Permissions (7):**
`feature_permission`, `features`, `model_has_permissions`, `model_has_roles`,
`permissions`, `role_has_permissions`, `roles`

**System (1):**
`schema_migrations`

**Tax (1):**
`ters`

### Provisioning Test Results ✅

**Test Date:** 22 Juli 2026  
**Environment:** Development — Laragon MySQL 8.0  
**Driver:** MySQL (`multiStatements=true`, `parseTime=True`)

| Item | Status | Detail |
|------|--------|--------|
| Company status | ✅ **active** | API mengembalikan `status: "active"` |
| Tenant database | ✅ Created | `hris_final-provision-test` |
| Tenant connection | ✅ Saved | Record di `tenant_connections` tersimpan |
| Migrations | ✅ **11 files** | 001 → 011 sukses semua |
| Total tables | ✅ **106 tables** | Setiap migrasi menciptakan tabel sesuai DDL |
| Server log | ✅ Clean | "Tenant provisioning completed successfully" |

#### API Test Response

```json
{
  "success": true,
  "data": {
    "id": "0b77f721-cee5-46f0-bc75-45219fc2316d",
    "name": "Final Provision Test",
    "slug": "final-provision-test",
    "email": "final@test.com",
    "phone": "021777777",
    "status": "active"
  }
}
```

### Issues Resolved During Development

| Issue | Root Cause | Fix |
|-------|-----------|-----|
| `TenantConnection` ID field | Struct Go tanpa `ID` → INSERT gagal (PK CHAR(36) tanpa default) | Tambah field ID — reuse companyID sebagai PK |
| `ssl_mode` column mismatch | GORM tag `sslmode` (tanpa underscore) tapi DDL `ssl_mode` | GORM tag: `sslmode` → `ssl_mode` |
| Multi-statement SQL gagal | MySQL driver default blokir multi-statement | Tambah `multiStatements=true` ke DSN |
| Access denied tenant DB | Kredensial pakai nama DB sebagai user (user tidak ada) | Development: gunakan `root`/`""` |
| FK dependency cross-file | `002_organization.sql` punya FK ke tabel `competencies` (migration 008) | Pindah FK ke migration 008 via ALTER TABLE |

---

## ✅ Completed Work

### 📄 Documentation

| # | Item | File |
|---|------|------|
| ✅ | Analisis blueprint v3 vs existing Laravel app | `docs/analisis-blueprint-vs-existing.md` |
| ✅ | Platform architecture design (modular monolith, multi-tenant) | `docs/platform-architecture-design.md` |
| ✅ | Environment variables template | `backend/.env.example` |
| ✅ | Build & development Makefile | `backend/Makefile` |
| ✅ | README utama proyek | `README.md` |

### ⚙️ Shared Kernel (Backend Infrastructure)

| # | Package | Files | Deskripsi |
|---|---------|-------|-----------|
| ✅ | `internal/pkg/config/` | `config.go` | Viper configuration loader (YAML + .env + env vars) |
| ✅ | `internal/pkg/database/` | `manager.go` | Multi-tenant DB connection manager with caching |
| ✅ | `internal/pkg/driver/` | `driver.go` | Shared DB driver type (PostgreSQL / MySQL) |
| ✅ | `internal/pkg/auth/` | `jwt.go` | JWT token generation & validation (access + refresh) |
| ✅ | `internal/pkg/middleware/` | `auth.go`, `cors.go`, `logger.go`, `recovery.go`, `tenant.go` | Gin middleware stack (auth, CORS, logging, recovery, tenant resolver) |
| ✅ | `internal/pkg/router/` | `router.go` | Router setup & module registration |
| ✅ | `internal/pkg/logger/` | `logger.go` | Zap structured logger |
| ✅ | `internal/pkg/module/` | `module.go` | Module SDK interface & registration |

### 🏢 Platform Module — Company

| # | Item | Files |
|---|------|-------|
| ✅ | Company model (UUID, status lifecycle) | `internal/platform/company/model.go` |
| ✅ | Company CRUD (Create, Read, Update, Delete) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Company lifecycle (Suspend, Activate, Backup, Restore) | `handler.go`, `service.go` |
| ✅ | Company request/response DTOs | `dto.go` |
| ✅ | Company routes registration | `routes.go` |
| ✅ | TenantConnection model (multi-tenant DB config) | `tenant_connection.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |

### 👤 Platform Module — Users & Auth

| # | Item | Files |
|---|------|-------|
| ✅ | PlatformUser model (UUID, bcrypt password, roles) | `internal/platform/user/model.go` |
| ✅ | JWT Authentication (Login, Refresh Token) | `service.go`, `handler.go` |
| ✅ | User CRUD (Create, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs | `dto.go` |
| ✅ | Routes registration | `routes.go` |
| ✅ | Auto-seed super admin (development) | `module.go` (via `EnsureSeed`) |

### 🧩 Platform Module — Module Management

| # | Item | Files |
|---|------|-------|
| ✅ | PlatformModule model (registered modules) | `internal/platform/modulemgmt/model.go` |
| ✅ | CompanyModule model (company-module association) | `internal/platform/modulemgmt/model.go` |
| ✅ | Module CRUD (Register, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Module activation/deactivation for company | `handler.go`, `service.go` |
| ✅ | Company-module listing | `repository.go` (via JOIN query) |
| ✅ | Routes registration | `routes.go` |

### 🔑 Platform Module — License Management

| # | Item | Files |
|---|------|-------|
| ✅ | License model (plan types, limits, dates) | `internal/platform/license/model.go` |
| ✅ | License CRUD (Create, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Plan-based auto calculation (max employees, modules) | `service.go` |
| ✅ | License key generation (via xid) | `service.go` |
| ✅ | Routes registration | `routes.go` |

### 📊 Platform Module — Monitoring

| # | Item | Files |
|---|------|-------|
| ✅ | Platform health check (DB connection status) | `internal/platform/monitoring/handler.go` |
| ✅ | Tenant connection health listing | `internal/platform/monitoring/handler.go` |
| ✅ | Tenant connection health detail | `internal/platform/monitoring/handler.go` |
| ✅ | Routes registration | `internal/platform/monitoring/routes.go` |

### 🏛️ Tenant Module — Organization

| # | Item | Files |
|---|------|-------|
| ✅ | Organization model (tree hierarchy with parent_id) | `internal/modules/organization/model.go` |
| ✅ | Organization CRUD + Tree view | `handler.go`, `service.go`, `repository.go` |
| ✅ | Organization request/response DTOs | `dto.go` |
| ✅ | Organization routes registration | `routes.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |

### 📦 API Documentation

| # | Item | File |
|---|------|------|
| ✅ | OpenAPI 3.0 JSON specification (**80+ endpoints**) | `internal/pkg/docs/openapi.json` |
| ✅ | Scalar UI served at `/docs` (interactive documentation) | `internal/pkg/docs/scalar.go` |
| ✅ | OpenAPI spec served at `/openapi.json` | `internal/pkg/docs/scalar.go` |

### 🐳 Infrastructure

| # | Item | File |
|---|------|------|
| ✅ | Multi-stage Dockerfile | `backend/docker/Dockerfile` |
| ✅ | Docker Compose (PostgreSQL, Redis, API, Asynqmon) | `docker/docker-compose.yml` |
| ✅ | PostgreSQL init script | `docker/init-db.sql` |
| ✅ | CLI Installer stub for tenant provisioning | `backend/cmd/installer/main.go` |
| ✅ | API server entry point | `backend/cmd/server/main.go` |
| ✅ | Server config YAML | `backend/config/config.yaml` |
| ✅ | Go module dependencies | `go.mod`, `go.sum` |

### 🗄️ Database Support

| # | Item | Detail |
|---|------|--------|
| ✅ | PostgreSQL driver | `gorm.io/driver/postgres v1.5.9` |
| ✅ | MySQL driver | `gorm.io/driver/mysql v1.6.0` |
| ✅ | Multi-driver DSN generation | PostgreSQL & MySQL format in config |
| ✅ | Per-tenant driver configuration | `driver` field in `tenant_connections` table |
| ✅ | Shared driver type | `internal/pkg/driver/driver.go` |

### 📦 Go Dependencies

```
github.com/gin-gonic/gin v1.10.0          # HTTP framework
github.com/golang-jwt/jwt/v5 v5.2.1       # JWT auth
github.com/google/uuid v1.6.0             # UUID generation
github.com/gosimple/slug v1.15.0          # URL slug
github.com/joho/godotenv v1.5.1           # .env file loader
github.com/spf13/viper v1.19.0            # Configuration
go.uber.org/zap v1.27.0                   # Structured logging
gorm.io/driver/mysql v1.6.0               # MySQL driver
gorm.io/driver/postgres v1.5.9            # PostgreSQL driver
gorm.io/gorm v1.30.0                      # ORM
```

### ✅ Tenant Provisioning

| # | Item | Detail |
|---|------|--------|
| ✅ | Provisioning Engine | Database creation + TenantConnection save |
| ✅ | Tenant SQL Migrations | 11 migration files → 106 tables |
| ✅ | Multi-statement MySQL support | `multiStatements=true` di DSN |
| ✅ | Error handling / graceful failure | Company status = `suspended` jika provisioning gagal |
| ✅ | End-to-end test | Company active ✅, 106 tables ✅, MySQL |

### ✅ Tenant Lifecycle Management

| # | Item | Detail |
|---|------|--------|
| ✅ | `Suspend` endpoint | Deactivate connection + status suspended |
| ✅ | `Activate` endpoint | Reactivate connection + status active |
| ✅ | `Terminate` endpoint | Drop database + remove connection + status terminated |
| ✅ | Soft delete cleanup | `DELETE /:id` now deactivates connection too |
| ✅ | Cached connection cleanup | `CloseTenantConnection` removes cached GORM connections |
| ✅ | End-to-end test | Suspend ✅ → Activate ✅ → Terminate ✅ |

### 🏛️ Tenant Module — Employee ✅

| # | Item | Files |
|---|------|-------|
| ✅ | Employee model (9 GORM models with UUID hooks) | `internal/modules/employee/model.go` |
| ✅ | Employee CRUD (Create, Read, Update, Delete) | `handler.go`, `service.go`, `repository.go` |
| ✅ | 8 sub-modules (Addresses, Emergency Contacts, Families, Educations, Experiences, Documents, Insurances, Employments) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs with validation (oneof, required, email, max) | `dto.go` |
| ✅ | Paginated List response | `service.go` |
| ✅ | Routes registration (30+ nested endpoints) | `routes.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |
| ✅ | AutoMigrate for 9 models during tenant provisioning | `module.go` |
| ✅ | End-to-end tested: Create, List, Get, Update, Delete + all sub-modules | API verified |

### 🏛️ Tenant Module — Job Management ✅

| # | Item | Files |
|---|------|-------|
| ✅ | Job Management model (18 GORM entities with UUID hooks) | `internal/modules/jobmanagement/model.go` |
| ✅ | Full CRUD for all 18 entities (Titles, Subs, Values, Objectives, Identifications, Responsibilities, Education Experiences, HR Authorities, Operational Authorities, Working Activities, Working Risks, Relationships, Subordinate Controls, Assets, Financials, Potency Competencies, Scores, Competency Groups) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs with validation for all 18 entities | `dto.go` |
| ✅ | Paginated List responses for list endpoints | `service.go` |
| ✅ | Routes registration (36+ endpoints) | `routes.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |
| ✅ | SQLite-integrated unit tests (74 tests: 22 repository + 29 service + 11 handler + 3 route + 9 integration) | `*_test.go` |
| ✅ | OpenAPI 3.0 documentation (35 schemas + 36 endpoints) | `internal/pkg/docs/openapi.json` |
| ✅ | RBAC permission: `jobmanagement.*` for company_admin | `internal/pkg/authz/rbac.go` |

### Build Status

```bash
$ go vet ./...    # ✅ Lulus
$ go build ./...  # ✅ Berhasil
```

---

## 🗺️ Roadmap

### Phase 1: Foundation ✅
- ✅ Platform Architecture Design
- ✅ Core Shared Packages (config, database, auth, middleware, router)
- ✅ Platform Module — Company Management (CRUD + TenantConnection)
- ✅ Tenant Module — Organization (tree CRUD)
- ✅ Docker & CI/CD setup
- ✅ Multi-database driver support (PostgreSQL + MySQL)
- ✅ Environment configuration (.env)
- ✅ Tenant Provisioning Engine (11 migrations → 106 tables, end-to-end verified)

### Phase 2: Tenant Core Modules ✅
- ✅ Employee Module (full CRUD + 8 sub-modules: addresses, emergency contacts, families, educations, experiences, documents, insurances, employments)

### Phase 3: Payroll & Complex
- ✅ **Job Management Module** (16+ models, 18 GORM entities, full CRUD + scoring) -- [Selesai 22 Juli 2026]
- ⬜ Payroll Module (calculator, PPh21, BPJS)
- ⬜ Competency Module

### Phase 4: Operations
- ⬜ Attendance & Leave Modules
- ⬜ Approval Module (multi-step workflow)
- ⬜ Frontend Integration (Vue 3 + PrimeVue)

### Phase 5: Polish
- ⬜ E2E Testing (Playwright)
- ⬜ Performance Optimization

---

## 📚 Dokumentasi Lainnya

| Dokumen | Deskripsi |
|---|---|
| [`docs/platform-architecture-design.md`](docs/platform-architecture-design.md) | Architecture design document lengkap |
| [`docs/analisis-blueprint-vs-existing.md`](docs/analisis-blueprint-vs-existing.md) | Analisis blueprint vs existing Laravel app |
| [`backend/.env.example`](backend/.env.example) | Template environment variables |
| [`backend/Makefile`](backend/Makefile) | Build & development commands |

---

## 🛡️ Security

- **JWT Authentication** — Access + Refresh token
- **Casbin RBAC** — Role-based access control
- **Multi-Tenant Isolation** — Database per tenant
- **Input Validation** — via go-playground/validator
- **SQL Injection Prevention** — GORM parameterized queries
- **Audit Log** — Semua mutasi tercatat

---

## 📄 Lisensi

Proprietary — All rights reserved.

---

*Dokumen ini diperbarui pada: 22 Juli 2026*
