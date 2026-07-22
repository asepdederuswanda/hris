# Analisis Perbandingan: HRIS Enterprise Architecture Blueprint v3 vs Existing Application (inthros-web)

**Dokumen:** Analisis Gap & Kesesuaian Arsitektur
**Tanggal:** 22 Juli 2026
**Versi:** 1.0

---

## Daftar Isi

1. [Executive Summary](#1-executive-summary)
2. [Perbandingan Stack Teknologi](#2-perbandingan-stack-teknologi)
3. [Multi-Tenant Architecture](#3-multi-tenant-architecture)
4. [Domain & Modul Coverage](#4-domain--modul-coverage)
5. [Platform Management Features](#5-platform-management-features)
6. [Arsitektur & Alur Data](#6-arsitektur--alur-data)
7. [Struktur Kode & Pola](#7-struktur-kode--pola)
8. [Fitur Existing yang Tidak Ada di Blueprint](#8-fitur-existing-yang-tidak-ada-di-blueprint)
9. [CI/CD & Deployment](#9-cicd--deployment)
10. [Gap Analysis & Prioritas](#10-gap-analysis--prioritas)
11. [Strengths Existing](#11-strengths-existing)
12. [Rekomendasi Strategis](#12-rekomendasi-strategis)
13. [Lampiran: Struktur File Key](#13-lampiran-struktur-file-key)

---

## 1. Executive Summary

Dokumen ini menganalisis kesesuaian antara **HRIS Enterprise Architecture Blueprint v3** (target arsitektur ideal) dengan **aplikasi existing inthros-web** yang sudah berjalan.

Referensi utama analisis:
- [`HRIS_Enterprise_Architecture_Blueprint_v3.md`](../HRIS_Enterprise_Architecture_Blueprint_v3.md)
- [`AI_CODING_RULES.md`](../inthros-web/AI_CODING_RULES.md) — Standar coding project existing
- Seluruh direktori [`inthros-web/`](../inthros-web/) — Source code aplikasi existing

### Ringkasan Project (Existing inthros-web)

| Metrik | Jumlah |
|---|---|
| **Controllers** | ~50+ |
| **Models (Eloquent)** | ~80+ |
| **Services** | ~40+ |
| **Form Requests** | ~60+ |
| **Migrations** | ~80 files |
| **Seeders** | ~20 files |
| **Stored Procedures** | 11 files |
| **Vue Pages** | ~30+ pages |
| **Vue Components** | ~60+ components |
| **Routes** | ~800+ lines di `web.php` |

### Temuan Utama

| Aspek | Status | Severity |
|---|---|---|
| **Stack backend berbeda total** (Laravel PHP vs Go) | 🔴 Gap Kritis | High |
| **Arsitektur frontend-backend** (Inertia monolith vs REST API) | 🔴 Gap Kritis | High |
| **Multi-tenant model** (single-DB dengan company_id vs DB-per-company) | 🔴 Gap Kritis | High |
| **Domain coverage** | 🟢 Existing lebih luas | - |
| **Payroll & tax Indonesia** | 🟢 Existing sudah mature | - |
| **Module Management & License** | 🟡 Belum ada | Medium |
| **Testing & Observability** | 🟡 Masih minimal | Medium |

---

## 2. Perbandingan Stack Teknologi

### 2.1 Tabel Perbandingan

| Aspek | Blueprint v3 (Target) | Existing (inthros-web) | Gap |
|---|---|---|---|
| **Backend Language** | Go | **PHP 8.3 (Laravel 12)** | 🔴 Total berbeda |
| **Backend Framework** | Gin + GORM | **Laravel 12 + Eloquent** | 🔴 Total berbeda |
| **Frontend Framework** | Vue 3 + TypeScript | **Vue 3 + Inertia.js** | 🟡 Sama Vue 3, beda approach |
| **UI Library** | PrimeVue | **Metronic Theme + Tailwind + Flowbite** | 🟡 Beda UI library |
| **State Management** | Pinia | **Inertia.js props (tidak pakai Pinia)** | 🟡 Inertia menggantikan Pinia |
| **Data Fetching** | TanStack Query | **Inertia.js (server-driven)** | 🟡 Pendekatan berbeda |
| **API Pattern** | REST API | **Inertia (no REST API)** | 🔴 Fundamental berbeda |
| **Build Tool** | Vite | **Vite** | 🟢 Sama |
| **CSS Framework** | Tailwind CSS | **Tailwind CSS + Flowbite** | 🟢 Kompatibel |
| **Auth** | JWT | **Session-based + WebAuthn (FIDO2)** | 🔴 Berbeda mekanisme |
| **Permission** | Casbin | **Spatie Laravel Permission + Custom ACL** | 🔴 Berbeda implementasi |
| **Background Jobs** | Asynq (Redis) | **Laravel Queue (database default)** | 🟡 Fungsi sama, beda tech |
| **Logging** | Zap | **Laravel Logging (stack/single)** | 🟡 Fungsi sama, beda tech |
| **Monitoring** | OpenTelemetry | ❌ **Belum ada** | 🔴 Gap |
| **Caching** | Redis | **Database cache (default) / Redis (optional)** | 🟡 Beda implementasi |

### 2.2 Detail Stack Existing (inthros-web)

#### Backend Dependencies (composer.json)

| Package | Versi | Fungsi |
|---|---|---|
| `laravel/framework` | ^12.0 | Framework utama |
| `inertiajs/inertia-laravel` | ^2.0 | Bridge frontend-backend |
| `spatie/laravel-permission` | ^6.20 | ACL & Role Management |
| `asbiin/laravel-webauthn` | 5.3 | WebAuthn/FIDO2 authentication |
| `barryvdh/laravel-dompdf` | ^3.1 | PDF generation |
| `maatwebsite/excel` | ^3.1 | Excel export/import |
| `tightenco/ziggy` | ^2.5 | Laravel route helper untuk JS |

#### Frontend Dependencies (package.json)

| Package | Versi | Fungsi |
|---|---|---|
| `@inertiajs/vue3` | ^2.0.3 | Inertia.js integration |
| `@vitejs/plugin-vue` | ^5.2.1 | Vue Vite plugin |
| `tailwindcss` | ^3.4.17 | CSS framework |
| `flowbite` | ^3.1.2 | Tailwind component library |
| `@iconify/vue` | ^5.0.0 | Icon library |
| `chart.js` | ^4.4.8 | Charts & visualizations |
| `d3-org-chart` | ^3.1.0 | Organization chart tree |
| `vue-i18n` | ^11.1.12 | Internationalization |
| `@vuepic/vue-datepicker` | ^11.0.2 | Date picker component |
| `leaflet` | ^1.9.4 | Maps |

---

## 3. Multi-Tenant Architecture

### 3.1 Perbandingan Model

| Aspek | Blueprint v3 | Existing (inthros-web) |
|---|---|---|
| **Model** | **Database per Company** (true multi-tenant) | **company_id di setiap tabel** (single database) |
| **Isolasi Data** | Separate database connections per tenant | Scoped via `HasCompany` trait & Spatie Teams |
| **Provisioning** | Tenant Provisioning Engine (auto) | **Manual setup per company** |
| **Platform DB** | companies, company_connections, modules, licenses, dll | **Hanya tabel `companies` biasa** |
| **Dynamic Connection** | Multiple DB connections at runtime | Single connection |

### 3.2 Implementasi Multi-Company di Existing

#### Trait HasCompany
**Lokasi:** `app/Traits/HasCompany.php`

```php
trait HasCompany
{
    public static function bootHasCompany()
    {
        static::creating(function ($model) {
            if (Auth::check() && empty($model->company_id)) {
                $model->company_id = Auth::user()->company_id;
            }
        });
    }

    public function company()
    {
        return $this->belongsTo(\App\Models\Settings\Company::class, 'company_id');
    }
}
```

#### Spatie Teams Configuration
**Lokasi:** `config/permission.php`

```php
'teams' => true,
'team_foreign_key' => 'team_id', // dipetakan ke company_id
```

#### Middleware SetPermissionTeam
**Lokasi:** `app/Http/Middleware/SetPermissionTeam.php`

```php
class SetPermissionTeam
{
    public function handle(Request $request, Closure $next)
    {
        $teamId = $request->user()?->company_id;
        app(PermissionRegistrar::class)->setPermissionsTeamId($teamId);
        return $next($request);
    }
}
```

#### Model yang Menggunakan HasCompany

| Model | File |
|---|---|
| Employee | `app/Models/Employees/Employee.php` |
| Organization | `app/Models/Organizations/Organization.php` |
| Employment | `app/Models/Employment.php` |
| *(dan model bisnis lainnya)* | |

### 3.3 Gap Analysis Multi-Tenant

| Gap | Detail | Dampak |
|---|---|---|
| **Single Database** | Semua tenant di DB yang sama | Risiko kebocoran data antar tenant |
| **No Platform DB** | Tidak ada tabel platform (modules, licenses, company_connections) | Tidak bisa manage tenant lifecycle |
| **No Tenant Provisioning** | Setup tenant masih manual | Tidak scalable untuk multi-company |
| **company_id everywhere** | Setiap tabel bisnis punya company_id | Overhead query & indexing |

---

## 4. Domain & Modul Coverage

### 4.1 Matriks Kesesuaian Domain

| Domain | Blueprint v3 | Existing (inthros-web) | Detail Existing | Status |
|---|---|---|---|---|
| **Organization** | ✅ Module SDK | ✅ | Organization, Zone, Summary, Level, Position, Job Family | 🟢 Lengkap |
| **Employee** | ✅ Module SDK | ✅ | Employee + 12 sub-modules (address, family, education, experience, documents, insurance, bank, BPJS, tax, payroll, salary, employment) | 🟢 Sangat lengkap |
| **Attendance** | ✅ Module SDK | ✅ | Sessions, events, shifts, exempt positions, company settings, overtime, device captures, face captures, locations | 🟢 Lengkap |
| **Leave** | ✅ Module SDK | ✅ | Types, requests, accrual policies, balances, reasons, request details | 🟢 Lengkap |
| **Payroll** | ✅ Module SDK | ✅ | Periods, runs, items, payslips, BPJS settings & rate components, PPh21 (PTKP, tax brackets, settings, calculation logs), salary components & grade components, employee adjustments, change logs | 🟢 Sangat lengkap (dengan stored procedures MySQL) |
| **Approval** | ✅ Module SDK | ✅ | Approval flows, steps, instances, tasks, actions | 🟢 Lengkap |
| **Reports** | ✅ Module SDK | ✅ | Attendance reports, employee reports (PDF + Excel), job management reports | 🟢 Ada |
| **Competency** | ❌ Tidak disebut | ✅ | Competency management, events, targets, scores & details, assignments, dashboard | 🟢 Existing lebih lengkap |
| **Job Management** | ❌ Tidak disebut | ✅ | Job titles & sub-titles, values, objectives, identifications, responsibilities, HR authorities, operational authorities, working activities, working risks, relationships & details, subordinate controls, assets, financials, competency groups, potency competencies, scores | 🟢 Existing sangat lengkap |
| **Recruitment** | ✅ Module SDK | ❌ | Tidak ada | 🔴 Gap |
| **Performance** | ✅ Module SDK | ❌ | Tidak ada (hanya dashboard card placeholder) | 🔴 Gap |
| **Training** | ✅ Module SDK | ❌ | Tidak ada | 🔴 Gap |
| **Asset** | ✅ Module SDK | ❌ | Tidak ada | 🔴 Gap |
| **Notifications** | ❌ Tidak disebut | ✅ | Account verification email notification | 🟢 Ada |

### 4.2 Detail Struktur Domain Existing

#### Controllers per Domain

| Domain | Jumlah Controller | Contoh File |
|---|---|---|
| **Auth** | 4 | AuthController, ProfileController, ForgotPasswordController, PasswordResetController |
| **Employees** | 17 | EmployeeController, EmergencyContactController, EmployeeAddressController, EmployeeInsuranceController, EmployeePayrollProfileController, EmployeeBankProfileController, EmployeeBpjsProfileController, EmployeeTaxProfileController, EmployeeSalaryComponentController, EmployeeSalaryAdjustmentController, EmployeeFamilyController, EmployeeEducationController, EmployeeExperienceController, EmployeeDocumentController, EmployeeDocumentGeneratorController, EmploymentController, EmployeeDashboardController |
| **Organizations** | 6 | OrganizationSummaryController, OrganizationController, ZoneController, JobFamilyController, OrganizationLevelController, PositionController |
| **Masters** | 12 | CompetencyController, CountryController, DistrictController, DocumentTemplateController, EducationController, EmploymentStatusController, GradingController, MaritalStatusController, ProvinceController, RegencyController, RelationshipTypeController, ReligionController, VillageController |
| **Settings** | 19 | CompanyController, UserController, RoleController, PermissionController, FeatureController, GeneralSettingController, PtkpController, TerController, SalaryComponentController, PayrollPeriodController, BpjsSettingController, SalaryGradeComponentController, BpjsRateComponentController, Pph21PtkpRateController, Pph21TaxBracketController, Pph21SettingController, CompanyHolidayController, AttendanceCompanySettingController, AttendanceCompanyShiftController, ApprovalFlowController, LeaveTypeController, LeaveAccrualPolicyController |
| **Attendance** | 5 | AttendanceEmployeeShiftController, AttendanceExemptPositionController, AttendanceReportController, AttendanceSessionController, AttendanceEventController |
| **Jobs** | 1 | JobsManagementController |
| **Competency** | 4 | CompetencyManagementController, CompetencyAssignmentController, CompetencyEventController, CompetencyDashboardController |
| **Other** | 4 | DashboardController, PlanController, ProgramPlanController, WorkUnitPlanController, RealTimeAbsensiController |

#### Service Layer per Domain

| Domain | Jumlah Service | Contoh File |
|---|---|---|
| **Employees** | 12 | EmployeeService, EmergencyContactService, EmployeeAddressService, EmployeeBankProfileService, EmployeeBpjsProfileService, EmployeeDocumentService, EmployeeEducationService, EmployeeExperienceService, EmployeeFamilyService, EmployeeInsuranceService, EmployeePayrollProfileService, EmployeeSalaryAdjustmentService, EmployeeSalaryComponentService, EmployeeTaxProfileService, EmploymentService, EmployeeAccountService |
| **Organizations** | 7 | OrganizationService, OrganizationLevelService, OrganizationSummaryService, PositionService, ZoneService, JobFamilyService |
| **Masters** | 12 | CompetencyService, CountryService, DistrictService, EducationService, EmploymentStatusService, GradingService, MaritalStatusService, ProvinceService, RegencyService, RelationshipTypeService, ReligionService, VillageService |
| **Settings** | 13 | CompanyService, UserService, RoleService, FeatureService, GeneralSettingService, PtkpService, TerService, SalaryComponentService, PayrollPeriodService, BpjsSettingService, SalaryGradeComponentService, BpjsRateComponentService, Pph21PtkpRateService, Pph21TaxBracketService, Pph21SettingService, CompanyHolidayService |
| **JobManagement** | 2 | JobManagementService, JobValueCalculator |
| **Auth** | 1 | WebauthnService |

#### Model per Domain

| Domain | Jumlah Model | Contoh File |
|---|---|---|
| **Employees** | 15 | Employee, EmergencyContact, EmployeeAddress, EmployeeBankProfile, EmployeeBpjsProfile, EmployeeDocument, EmployeeEducation, EmployeeExperience, EmployeeFamily, EmployeeInsurance, EmployeeLeaveBalance, EmployeePayrollProfile, EmployeeSalaryAdjustment, EmployeeSalaryComponent, EmployeeTaxProfile |
| **Organizations** | 5 | Organization, OrganizationLevel, OrganizationSummary, Position, Zone |
| **Masters** | 13 | Competency, CompetenceValue, CompetencyValue, Country, District, DocumentTemplate, Education, EmploymentStatus, Grading, MaritalStatus, Province, Regency, RelationshipType, Religion, Village |
| **Settings** | 24 | Company, User, Role, Permission, Feature, FeaturePermission, GeneralSetting, ApprovalFlow, ApprovalFlowStep, Ptkp, Ter, SalaryComponent, SalaryGradeComponent, PayrollPeriod, BpjsSetting, BpjsRateComponent, Pph21PtkpRate, Pph21TaxBracket, Pph21Setting, Pph21CalculationLog, PayrollPayslip, PayrollRun, PayrollRunEmployee, PayrollRunItem, SalaryChangeLog, AttendanceCompanySetting, AttendanceCompanyShift, CompanyHoliday, LeaveType, LeaveAccrualPolicy, LeaveRequest, LeaveReason, LeaveRequestDetail, EmployeeLeaveBalance |
| **Attendance** | 6 | AttendanceDeviceCapture, AttendanceEmployeeShift, AttendanceEvent, AttendanceExemptPosition, AttendanceSession, AttendanceOvertimeRequest |
| **JobManagement** | 16 | JobManagementTitle, JobManagementTitleSub, JobManagementValue, JobManagementObjective, JobManagementIdentification, JobManagementResponsibility, JobManagementEducationExperiences, JobManagementEducationExperienceField, JobManagementEducationExperienceMajor, JobManagementHrAuthority, JobManagementOperationalAuthority, JobManagementWorkingActivity, JobManagementWorkingRisk, JobManagementRelationship, JobManagementRelationshipDetail, JobManagementSubordinateControl, JobManagementAsset, JobManagementFinancial, JobManagementPotencyCompetency, JobManagementCompetencyGroup, JobManagementScore |
| **Competency** | 4 | CompetencyScore, CompetencyScoreDetail, CompetencyEvent, CompetencyEventTarget |

---

## 5. Platform Management Features

### 5.1 Matriks Fitur Platform

| Feature | Blueprint v3 | Existing (inthros-web) | Detail Existing | Status |
|---|---|---|---|---|
| **Company Management** | ✅ | ✅ | CRUD companies dengan soft delete | 🟢 Ada |
| **Module Management** | ✅ (Module SDK + manifest.yaml) | ❌ | **Tidak ada konsep module activation** | 🔴 Gap |
| **License Management** | ✅ | ❌ | **Tidak ada lisensi** | 🔴 Gap |
| **Feature Flag** | ✅ | ✅ Partial | Tabel `features` + `feature_permission` pivot | 🟡 Statis (bukan runtime flag) |
| **System Config** | ✅ | ✅ | Tabel `general_settings` | 🟢 Ada |
| **Tenant Monitoring** | ✅ | ❌ | **Tidak ada monitoring tenant** | 🔴 Gap |
| **Tenant Backup/Restore** | ✅ | ❌ | **Tidak ada backup per tenant** | 🔴 Gap |
| **User Management** | ✅ | ✅ | Dengan roles & permissions | 🟢 Lengkap |
| **Audit Log** | ✅ | Partial | Hanya `created_by` / `updated_by` traits | 🟡 Belum audit log detail |
| **MFA / Passkey** | ✅ (MFA Ready) | ✅ | WebAuthn (FIDO2) sudah terintegrasi via `asbiin/laravel-webauthn` | 🟢 Existing sudah punya! |

### 5.2 Detail Implementasi Feature Flag Existing

**Lokasi:** `app/Support/AclFeatureRegistry.php`

```php
class AclFeatureRegistry
{
    public static function items(): array
    {
        return [
            ['name' => 'Dashboard', 'slug' => 'dashboard', 'group' => 'Dashboard'],
            ['name' => 'Manajemen Pengguna', 'slug' => 'user', 'group' => 'Settings'],
            ['name' => 'Role', 'slug' => 'role', 'group' => 'Settings'],
            // ... 40+ features terdaftar
        ];
    }
}
```

**Tabel Database:**
- `features` — menyimpan daftar fitur/module
- `feature_permission` — pivot table feature ↔ permission
- `permissions` — menyimpan permission (view/create/update/delete/approve)

**Command Sync:** `app/Console/Commands/SyncAclPermissions.php` — `php artisan acl:sync-permissions`

### 5.3 Detail ACL & Permission Existing

#### AclPermissionMap
**Lokasi:** `app/Support/AclPermissionMap.php`

Custom mapping yang meng-handle:
1. **Action aliases** — `read` → `view`, `edit`/`modify` → `update`
2. **Legacy compatibility** — Permission lama tetap bisa dipakai
3. **Feature actions** — `view`, `create`, `update`, `delete`, `approve`
4. **Gate integration** — `Gate::before()` di `AppServiceProvider`

#### AuthorizeRoutePermission Middleware
**Lokasi:** `app/Http/Middleware/AuthorizeRoutePermission.php`

Memetakan route name ke resource slug, lalu mengecek permission user.
Contoh: `employees.index` → `view employee`, `employees.store` → `create employee`

---

## 6. Arsitektur & Alur Data

### 6.1 Perbandingan Arsitektur

| Aspek | Blueprint v3 (Target) | Existing (inthros-web) |
|---|---|---|
| **Request Flow** | Client → REST API → Controller → Service → Repository → DB | Browser → Inertia → Laravel Controller → Service → Model (Eloquent) → DB |
| **Response** | JSON | Inertia page render / redirect back with flash |
| **Routing** | Gin router | Laravel routes (`routes/web.php`) |
| **ORM** | GORM | Eloquent ORM |
| **Validation** | DTO validation | Laravel Form Request |
| **Business Logic** | Services | Services (sesuai AI_CODING_RULES.md) |
| **Database** | PostgreSQL / MySQL | MySQL |
| **Stored Procedures** | Tidak disebut | ✅ Ada untuk payroll (generate, lock, validate, publish) |

### 6.2 Alur Data Existing

#### Alur Menampilkan Data (Index)

```text
User buka halaman
    ↓
routes/web.php
    ↓
Controller@index
    ↓
Service@table / Service@index
    ↓
Model query + filter + pagination
    ↓
Inertia::render()
    ↓
Vue Page menerima props
    ↓
Table tampil di browser
```

#### Alur Menambah Data (Store)

```text
User submit form Vue
    ↓
form.post(route(...))
    ↓
routes/web.php
    ↓
Form Request validasi input
    ↓
Controller@store
    ↓
Service@store
    ↓
Model::create()
    ↓
redirect()->back()->with('success', ...)
    ↓
Inertia reload props
    ↓
Toast / flash message tampil
```

#### Alur Mengubah Data (Update)

```text
User edit data di modal / form
    ↓
form.put(route(...))
    ↓
Route Model Binding
    ↓
Form Request validasi input
    ↓
Controller@update
    ↓
Service@update
    ↓
$model->update()
    ↓
redirect()->back()->with('success', ...)
    ↓
Vue re-render
```

### ### 6.3 Stored Procedures Payroll

Database `database/procedures/` berisi stored procedures MySQL untuk menjalankan logika perhitungan penggajian di sisi database — memastikan konsistensi dan performa perhitungan payroll bulanan.

| File | Konteks Penggunaan | Fungsi |
|---|---|---|
| `sp_generate_bpjs_items.sql` | **Saat run payroll** | Generate BPJS items untuk payroll run |
| `sp_generate_employee_salary_components.sql` | **Saat aktivasi employee** | Generate komponen gaji per employee |
| `sp_generate_payroll_run.sql` | **Saat run payroll** | Generate payroll run utama (memanggil SP lain) |
| `sp_generate_payslips.sql` | **Saat run payroll** | Generate payslip per employee |
| `sp_generate_pph21_items.sql` | **Saat run payroll** | Generate perhitungan PPh21 |
| `sp_lock_payroll_run.sql` | **Setelah validasi** | Lock payroll run agar tidak bisa diubah |
| `sp_precheck_payroll_profile.sql` | **Sebelum run payroll** | Validasi payroll profile sebelum run |
| `sp_publish_payslips.sql` | **Setelah approval** | Publish payslip ke employee |
| `sp_recalculate_payroll_run_totals.sql` | **Setelah koreksi** | Rekalkulasi total payroll |
| `sp_validate_bpjs_payroll_run.sql` | **Saat validasi** | Validasi BPJS payroll |
| `sp_validate_payroll_run.sql` | **Saat validasi** | Validasi payroll run |
| `sp_validate_pph21_payroll_run.sql` | **Saat validasi** | Validasi PPh21 payroll |

---

## 7. Struktur Kode & Pola

### 7.1 Perbandingan Struktur Folder

| Aspek | Blueprint v3 (Target) | Existing (inthros-web) |
|---|---|---|
| **Repository** | Monorepo: `backend/`, `frontend/`, `shared/`, `installer/`, `docs/`, `docker/` | **Single Laravel app** |
| **Backend Pattern** | Controller → Service → Repository | Controller → **Form Request** → **Service** → Model (Eloquent) |
| **Module Structure** | Setiap modul: `manifest.yaml`, routes, handlers, services, repositories, entities, dto, permissions, menus, migrations, seeders | **Domain-based folders** (Controllers, Models, Services, Requests per domain) |
| **Frontend** | `frontend/platform-admin/`, `frontend/tenant/`, `frontend/shared/` | `resources/js/Pages/` (single frontend) |

### 7.2 Struktur Folder Existing

```text
inthros-web/
├── app/
│   ├── Console/Commands/          # Artisan commands
│   ├── Exports/                   # Excel exports
│   ├── Http/
│   │   ├── Controllers/           # Controllers per domain
│   │   │   ├── Attendance/
│   │   │   ├── Auth/
│   │   │   ├── Employees/
│   │   │   ├── Jobs/
│   │   │   ├── Masters/
│   │   │   ├── Organizations/
│   │   │   └── Settings/
│   │   ├── Middleware/            # Custom middleware
│   │   ├── Requests/              # Form Request validation
│   │   │   ├── Attendance/
│   │   │   ├── Employees/
│   │   │   ├── Masters/
│   │   │   ├── Organizations/
│   │   │   ├── Plans/
│   │   │   └── Settings/
│   ├── Models/                    # Eloquent models
│   │   ├── Attendance/
│   │   ├── Employees/
│   │   ├── JobManagement/
│   │   ├── Masters/
│   │   ├── Organizations/
│   │   └── Settings/
│   ├── Notifications/             # Email notifications
│   ├── Policies/                  # Authorization policies
│   ├── Providers/                 # Service providers
│   ├── Rules/                     # Custom validation rules
│   ├── Services/                  # Business logic layer
│   │   ├── Attendance/
│   │   ├── Auth/
│   │   ├── Employees/
│   │   ├── JobManagements/
│   │   ├── Masters/
│   │   ├── Organizations/
│   │   └── Settings/
│   ├── Support/                   # ACL mapping & feature registry
│   └── Traits/                    # Reusable traits
├── bootstrap/
├── config/                        # Laravel config files
├── database/
│   ├── factories/
│   ├── migrations/                # ~80 migration files
│   ├── procedures/                # 11 stored procedures
│   └── seeders/                   # ~20 seeder files
├── resources/
│   ├── css/
│   ├── js/
│   │   ├── Components/           # Vue shared components
│   │   │   ├── Forms/
│   │   │   ├── Kanbans/
│   │   │   ├── Layouts/
│   │   │   └── Tables/
│   │   ├── Layouts/              # Vue layouts
│   │   ├── Pages/                # Inertia pages per domain
│   │   │   ├── Attendance/
│   │   │   ├── Auth/
│   │   │   ├── Competencies/
│   │   │   ├── Employees/
│   │   │   ├── Jobs/
│   │   │   └── Overtime/
│   │   └── Plugins/              # Vue plugins (i18n)
│   ├── metronic/                 # Metronic theme assets
│   └── views/                    # Blade templates
├── routes/
│   └── web.php                   # ~800+ lines routes
├── public/
├── storage/
└── tests/                        # Unit & feature tests
```

### 7.3 Traits Reusable Existing

| Trait | File | Fungsi |
|---|---|---|
| `HasCompany` | `app/Traits/HasCompany.php` | Auto-fill company_id pada model scoped |
| `HasCreatedBy` | `app/Traits/HasCreatedBy.php` | Auto-fill created_by user |
| `HasUpdatedBy` | `app/Traits/HasUpdatedBy.php` | Auto-fill updated_by user |
| `HasUuid` | `app/Traits/HasUuid.php` | Auto-generate UUID & route key |
| `HasFeaturePermission` | `app/Traits/HasFeaturePermission.php` | Feature-permission relationship |
| `HasPermission` | `app/Traits/HasPermission.php` | Custom permission helpers |
| `HasProgramPlanReport` | `app/Traits/HasProgramPlanReport.php` | Program plan report relationship |

### 7.4 Coding Standards (AI_CODING_RULES.md)

Existing memiliki panduan coding yang sangat detail di `AI_CODING_RULES.md` yang mencakup:

1. **Alur data baku**: Controller → Form Request → Service → Model → Inertia Page
2. **Pola CRUD**: Migration → Model → Service → Form Request → Controller → Route → Vue Page → Vue Form
3. **Larangan**: Business logic di controller, query kompleks di Vue, validasi manual, mengabaikan traits
4. **Standar response**: `redirect()->back()->with('success', ...)`
5. **Checklist CRUD baru**: 11 komponen yang wajib dibuat

---

## 8. Fitur Existing yang Tidak Ada di Blueprint

Existing memiliki beberapa fitur yang tidak tercantum di Blueprint v3:

### 8.1 Job Management (Analisis Jabatan)

Module komprehensif untuk analisis jabatan dengan 16+ model:

| Model | Deskripsi |
|---|---|
| JobManagementTitle | Daftar jabatan/posisi |
| JobManagementTitleSub | Sub-jabatan |
| JobManagementValue | Nilai jabatan (dengan scoring) |
| JobManagementObjective | Tujuan jabatan |
| JobManagementIdentification | Identifikasi jabatan |
| JobManagementResponsibility | Tanggung jawab jabatan |
| JobManagementEducationExperiences | Persyaratan pendidikan & pengalaman |
| JobManagementHrAuthority | Wewenang HR |
| JobManagementOperationalAuthority | Wewenang operasional |
| JobManagementWorkingActivity | Aktivitas kerja |
| JobManagementWorkingRisk | Risiko kerja |
| JobManagementRelationship | Relasi jabatan |
| JobManagementSubordinateControl | Kendali bawahan |
| JobManagementAsset | Aset yang dikelola |
| JobManagementFinancial | Kewenangan finansial |
| JobManagementPotencyCompetency | Kompetensi potensi |
| JobManagementCompetencyGroup | Group kompetensi |
| JobManagementScore | Skoring jabatan |

### 8.2 Competency Management

Module penuh untuk manajemen kompetensi:

| Fitur | Deskripsi |
|---|---|
| Competency CRUD | Master data kompetensi |
| Competency Values | Nilai-nilai kompetensi |
| Competency Events | Event penilaian kompetensi (auto/manual) |
| Competency Event Targets | Target organisasi per event |
| Competency Scores | Scoring kompetensi per organisasi |
| Competency Score Details | Detail scoring per kompetensi |
| Competency Dashboard | Dashboard hasil kompetensi |
| Competency Assignments | Assignment kompetensi ke organisasi |
| Job Family Competencies | Kompetensi per job family |
| Job Management Potency Competencies | Kompetensi potensi per jabatan |

### 8.3 Pajak & BPJS Indonesia (Payroll Spesifik Lokal)

| Fitur | Deskripsi |
|---|---|
| PTKP | Penghasilan Tidak Kena Pajak |
| TER | Tarif Efektif Rata-rata |
| PPh21 Settings | Konfigurasi PPh21 |
| PPh21 Tax Brackets | Lapisan tarif pajak |
| PPh21 PTKP Rates | Rate PTKP per status |
| PPh21 Calculation Log | Log perhitungan PPh21 |
| BPJS Settings | Konfigurasi BPJS Kesehatan & Ketenagakerjaan |
| BPJS Rate Components | Komponen rate BPJS |
| Salary Grade Components | Komponen gaji per grade |
| Salary Change Logs | Log perubahan gaji |
| Payroll Profile Change Logs | Log perubahan payroll profile |

### 8.4 Fitur Tambahan Lainnya

| Fitur | Deskripsi |
|---|---|
| Organization Chart (d3-org-chart) | Tree chart struktur organisasi |
| Kanban Board | Komponen Kanban untuk workflow |
| Tree View & Tree Table | View hierarki |
| Real-Time Attendance Dashboard | Dashboard absensi real-time |
| WebAuthn/Passkey Auth | Autentikasi FIDO2 tanpa password |
| Document Template | Template dokumen untuk generate |
| Document Generator | Generate dokumen employee (PDF) |
| Excel Reports | Export Excel untuk employee & jobs management |
| PDF Reports | Export PDF untuk employee & jobs management |

---

## 9. CI/CD & Deployment

### 9.1 Perbandingan

| Aspek | Blueprint v3 | Existing (inthros-web) |
|---|---|---|
| **CI/CD Platform** | GitHub Actions | ✅ GitHub Actions |
| **Containerization** | Docker | ❌ Deploy langsung ke VPS (SSH + SCP) |
| **Testing** | Unit, Integration, API, E2E, Performance | ✅ Unit test (PHPUnit) — **masih minimal** |
| **Auto Build** | ✅ | ✅ (Composer install + npm build) |
| **Auto Deploy** | ✅ | ✅ (SSH ke VPS + artisan commands) |
| **Health Check** | ✅ Tenant health check | ✅ curl health check after deploy |
| **E2E Testing** | Playwright | ❌ Belum ada |

### 9.2 Detail CI/CD Existing

**File:** `.github/workflows/deploy-laravel.yml`

**Trigger:** Push ke branch `main`

**Steps:**
1. Checkout source
2. Setup PHP 8.3 + extensions
3. Composer install
4. Environment setup
5. Migration test
6. PHPUnit test
7. Setup Node 20
8. Npm install & build
9. Check required secrets
10. SSH to VPS → git pull + composer install
11. SCP upload frontend assets
12. SSH finalize: optimize, migrate, cache, storage link, queue restart
13. Health check curl

---

## 10. Resource Count Summary

Ringkasan jumlah resource berdasarkan hasil eksplorasi direktori:

### Controllers per Domain

| Domain | Count |
|---|---|
| Auth | 4 |
| Employees | 17 |
| Organizations | 6 |
| Masters | 13 |
| Settings | 19 |
| Attendance | 5 |
| Jobs | 1 |
| Competency | 4 |
| Other (Dashboard, Plans, etc.) | 5 |
| **Total** | **~74** |

### Models per Domain

| Domain | Count |
|---|---|
| Employees | 16 (incl. PayrollProfileChangeLog, EmployeeLeaveBalance) |
| Organizations | 5 |
| Masters | 15 |
| Settings | 29 (incl. Approval, Payroll, Leave, Attendance settings) |
| Attendance | 6 |
| JobManagement | 21 |
| Competency | 4 |
| Other (Plan, ProgramPlan, etc.) | 5 |
| **Total** | **~101** |

### Services per Domain

| Domain | Count |
|---|---|
| Employees | 15 |
| Organizations | 7 |
| Masters | 12 |
| Settings | 16 |
| JobManagement | 2 |
| Auth | 1 |
| Attendance | 2 |
| Competency | 1 |
| **Total** | **~56** |

### Vue Pages & Components

| Kategori | Count |
|---|---|
| Pages (Inertia) | ~30+ |
| Shared Components | ~60+ |
| Layouts | ~7 |

---

## 11. Gap Analysis & Prioritas

### 11.1 🔴 Critical Gaps (Perlu Rewrite Total)

| # | Gap | Detail | Dampak | Prioritas |
|---|---|---|---|---|
| 1 | **Stack backend berbeda** | Laravel (PHP) vs Go | Rewrite 100% backend | High |
| 2 | **Arsitektur API** | Inertia (server-driven) vs REST API | Ubah total frontend-backend communication | High |
| 3 | **Multi-tenant database** | Single DB vs DB per company | Ubah fundamental arsitektur data | High |
| 4 | **Auth mechanism** | Session + WebAuthn vs JWT | Ubah authentication flow | High |

### 11.2 🟡 Moderate Gaps

| # | Gap | Detail | Dampak | Prioritas |
|---|---|---|---|---|
| 5 | **Module Management** | Belum ada konsep module + manifest | Tidak bisa modular activation | Medium |
| 6 | **License Management** | Belum ada sistem lisensi | Tidak bisa manage subscription | Medium |
| 7 | **Tenant Provisioning** | Setup tenant masih manual | Tidak scalable | Medium |
| 8 | **Monitoring** | Belum ada OpenTelemetry | Tidak ada observability | Medium |
| 9 | **Audit Log** | Hanya created_by/updated_by | Kurang detail untuk compliance | Medium |
| 10 | **Unit Test** | Masih minimal | Risiko regresi tinggi | Medium |
| 11 | **E2E Test** | Belum ada Playwright | Tidak ada automated UI test | Medium |
| 12 | **Docker** | Belum dipakai | Environment inconsistency | Medium |

### 11.3 🟢 Non-Gaps (Sudah Sesuai / Lebih Baik)

| # | Aspek | Keterangan |
|---|---|---|
| 1 | **Vue 3 frontend** | Sama-sama Vue 3 |
| 2 | **Service layer pattern** | Existing sudah memisahkan business logic |
| 3 | **Form Request validation** | Best practice Laravel |
| 4 | **UUID for all entities** | Siap untuk distributed system |
| 5 | **Multi-company scoping** | Meski single DB, scoping sudah benar via HasCompany trait |
| 6 | **ACL system** | Spatie + custom mapping sudah mature |
| 7 | **WebAuthn/Passkey** | Security modern sudah terintegrasi |
| 8 | **Domain coverage** | Existing lebih luas (job management, competency, payroll lokal) |
| 9 | **Coding standards** | AI_CODING_RULES.md sangat comprehensive |
| 10 | **CI/CD** | GitHub Actions sudah jalan |

---

## 12. Strengths Existing

Kelebihan aplikasi existing dibandingkan target blueprint:

### 11.1 Maturity Payroll System
- ✅ Perhitungan BPJS Kesehatan & Ketenagakerjaan
- ✅ Perhitungan PPh21 dengan PTKP, TER, dan tax brackets
- ✅ Stored procedures MySQL untuk performa perhitungan
- ✅ Payroll run dengan lock/unlock mechanism
- ✅ Payslip generation & publishing
- ✅ Change logs untuk audit

### 11.2 Comprehensive Job Analysis
- ✅ 16+ model untuk analisis jabatan
- ✅ Scoring system untuk nilai jabatan
- ✅ Competency mapping per jabatan
- ✅ Responsibility & authority matrix
- ✅ Working risk & activity analysis

### 11.3 Clean Architecture
- ✅ Separation of concerns (Controller → Service → Model)
- ✅ Trait reuse (HasCompany, HasCreatedBy, HasUpdatedBy, HasUuid)
- ✅ Form Request untuk validasi
- ✅ Service layer untuk business logic
- ✅ Rich set of Vue shared components

### 11.4 Production-Ready Features
- ✅ Soft delete di major entities
- ✅ UUID sebagai primary key
- ✅ Spatie permission system dengan custom ACL mapping
- ✅ i18n support (vue-i18n)
- ✅ WebAuthn untuk passkey authentication
- ✅ Export PDF & Excel

---

## 13. Rekomendasi Strategis

### Opsi A: Evolusi (Lanjutkan Laravel)

Pertahankan Laravel stack yang sudah mature dan fokus pada gap yang tidak perlu stack change.

**Rekomendasi Implementasi:**
1. ✅ Tambah **Module Management** — buat sistem plugin/module dengan manifest
2. ✅ Tambah **License Management** — fitur licensing per company
3. ✅ Upgrade **Audit Log** — gunakan `spatie/laravel-activitylog`
4. ✅ Implementasi **unit test + feature test** untuk coverage
5. ✅ Implementasi **E2E test** dengan Playwright
6. ✅ **Dockerize** aplikasi untuk environment consistency
7. ✅ Upgrade monitoring — Laravel Telescope atau OpenTelemetry

**Estimasi effort:** Medium (3-6 bulan)

### Opsi B: Blueprint Baru (Go Rewrite)

Rewrite total mengikuti blueprint, hanya disarankan jika performa/scalability menjadi bottleneck.

**Rekomendasi Implementasi:**
1. Bangun Go API untuk module baru terlebih dahulu
2. Pertahankan Laravel untuk module existing yang already mature
3. Integrasi via API gateway
4. Migrasi bertahap module per module

**Estimasi effort:** Very High (12-24 bulan)

### Opsi C: Hybrid (Direkomendasikan ✅)

Backend Laravel tetap untuk existing modules. Module baru dan Platform Management dibangun di Go.

**Langkah Implementasi:**

| Fase | Module | Tech Stack | Timeline |
|---|---|---|---|
| **Fase 1** | Platform Management (Module, License, Tenant) | **Go + Gin** | 3-4 bulan |
| **Fase 2** | Recruitment, Performance, Training, Asset | **Go + Gin** | 4-6 bulan |
| **Fase 3** | Upgrade existing (Audit Log, Testing, Docker) | **Laravel** | 2-3 bulan |
| **Fase 4** | Migration bertahap module existing ke Go | **Go** | 6-12 bulan |

**Keuntungan:**
- Tidak perlu rewrite semua module sekaligus
- Module payroll & job management tetap bisa berjalan di Laravel
- Module baru bisa langsung pakai arsitektur baru
- Platform Management (multi-tenant) bisa pakai Go untuk performance

---

## 14. Lampiran: Struktur File Key

### File Konfigurasi

| File | Deskripsi |
|---|---|
| `composer.json` | PHP dependencies |
| `package.json` | Node.js dependencies |
| `.env.example` | Environment template |
| `vite.config.js` | Vite build configuration |
| `tailwind.config.js` | Tailwind CSS configuration |
| `phpunit.xml` | PHPUnit configuration |
| `config/app.php` | App configuration |
| `config/database.php` | Database connections (MySQL, SQLite, Redis) |
| `config/auth.php` | Authentication (session-based) |
| `config/permission.php` | Spatie permission (teams enabled) |
| `config/session.php` | Session (database driver) |
| `config/queue.php` | Queue (database driver) |
| `config/cache.php` | Cache (database driver) |
| `config/mail.php` | Mail configuration |
| `config/filesystems.php` | File storage (local, S3) |

### File Backend Utama

| File | Deskripsi |
|---|---|
| `app/Providers/AppServiceProvider.php` | Gate registration, ACL integration |
| `app/Http/Middleware/HandleInertiaRequests.php` | Shared Inertia props (auth, flash) |
| `app/Http/Middleware/SetPermissionTeam.php` | Multi-company permission scoping |
| `app/Http/Middleware/AuthorizeRoutePermission.php` | Route-based permission check |
| `app/Support/AclPermissionMap.php` | Permission alias mapping |
| `app/Support/AclFeatureRegistry.php` | Feature registry for ACL |
| `app/Console/Commands/SyncAclPermissions.php` | Sync ACL seeder command |
| `app/Traits/HasCompany.php` | Multi-company model trait |
| `routes/web.php` | ~800+ lines route definitions |

### File Database

| File | Jumlah |
|---|---|
| `database/migrations/` | ~80 migration files |
| `database/seeders/` | ~20 seeder files |
| `database/procedures/` | 11 stored procedures |
| `database/factories/` | 1 factory file |

### File Frontend

| File | Deskripsi |
|---|---|
| `resources/js/app.js` | Vue app bootstrap with Inertia |
| `resources/js/Layouts/AdminLayout.vue` | Main admin layout |
| `resources/js/Pages/Dashboard.vue` | Dashboard page |
| `resources/js/Pages/Auth/*.vue` | Auth pages (login, forgot password, etc.) |
| `resources/js/Pages/Employees/Form.vue` | Employee multi-tab form |
| `resources/js/Pages/Attendance/Report.vue` | Attendance report |
| `resources/js/Components/Tables/Table.vue` | Reusable table component |
| `resources/js/Components/Layouts/` | Layout components (Navbar, Sidebar, Aside, etc.) |
| `resources/js/Components/Forms/` | Form components (~30+ components) |

---

**Dokumen ini disusun berdasarkan analisis:**
- `HRIS_Enterprise_Architecture_Blueprint_v3.md` — Dokumen target arsitektur
- `inthros-web/` — Kode sumber aplikasi existing
- `AI_CODING_RULES.md` — Standar coding project

Untuk pertanyaan atau diskusi lebih lanjut, silakan hubungi tim arsitektur.
