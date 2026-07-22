# HRIS Enterprise Blueprint

## Roadmap

1.  Existing Code Analysis
2.  Existing System Documentation
3.  Backend Implementation
4.  API Documentation
5.  Backend Testing
6.  Frontend Implementation
7.  Frontend Testing

## Sprint Planning

  Sprint   Target
  -------- --------------------------------------------
  1        Project setup, monorepo, CI/CD, Docker
  2        Existing code analysis & inventory
  3        Business process & ERD documentation
  4        Core backend (config, auth, module loader)
  5-10     Master & HR modules
  11-14    Payroll, Approval, Reporting
  15-18    Frontend implementation
  19       Integration & UAT
  20       Production readiness

## Milestones

-   M1: Existing system documented
-   M2: Core platform selesai
-   M3: Semua API tersedia
-   M4: Frontend selesai
-   M5: UAT selesai
-   M6: Go Live

## Definition of Done

-   Code review selesai
-   Unit test ≥80%
-   Swagger diperbarui
-   Migration tersedia
-   Tidak ada critical vulnerability
-   Lulus QA

## Coding Standard

### Backend

-   gofmt, golangci-lint
-   Clean Architecture + Modular Monolith
-   Dependency Injection
-   Context-aware request \### Frontend
-   TypeScript strict
-   ESLint + Prettier
-   Composition API
-   Feature based structure

## Git Strategy

-   main
-   develop
-   feature/\*
-   release/\*
-   hotfix/\*

## Architecture

-   Modular Monolith
-   Multi Company
-   Module enable/disable
-   Redis
-   MySQL
-   Queue Asynq

## Database Convention

-   snake_case
-   ULID/UUIDv7 atau BIGINT sesuai kebutuhan
-   audit columns
-   foreign key & index
-   company_id pada seluruh tabel bisnis

## API Guideline

-   REST
-   OpenAPI 3
-   Versioning /api/v1
-   Consistent response
-   Pagination
-   Filtering
-   Sorting

## UI Guideline

-   PrimeVue
-   Tailwind CSS
-   Responsive
-   Dark mode ready

## Security Checklist

-   JWT
-   RBAC
-   Rate limiting
-   Audit log
-   Input validation
-   OWASP ASVS review

## Performance Target

-   API \<200 ms (normal)
-   P95 \<500 ms
-   Redis cache
-   Server-side pagination

## Deployment

-   Docker
-   GitHub Actions
-   Staging
-   Production
-   Rollback

## Backup & DR

-   Daily DB backup
-   Weekly full backup
-   Restore test bulanan

## Observability

-   Structured logging
-   Prometheus
-   Grafana
-   OpenTelemetry

## Testing Strategy

-   Unit
-   Integration
-   API
-   E2E (Playwright)
-   UAT

## Module Checklist

-   Entity
-   DTO
-   Repository
-   Service
-   Handler
-   Route
-   Migration
-   Seeder
-   Permission
-   API Docs
-   Unit Test
-   Frontend
-   E2E

## Plugin SDK

Setiap modul wajib menyediakan: - module manifest - route registration -
menu registration - permission registration - migration - seeder -
configuration
