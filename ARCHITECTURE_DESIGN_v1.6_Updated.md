# Architecture Design Document (ADD)
**Blueprint Roadmap v1.6 & System Specification**

---

## 1. Executive Summary & Core Architecture

Sistem ini dirancang menggunakan pendekatan **Modular Monolith** berbasis bahasa pemrograman **Go (Golang)** dengan pola isolasi data **Database per Tenant**. Arsitektur ini menggabungkan fleksibilitas pemeliharaan kode berbasis modul dengan keamanan isolasi data tingkat tinggi antar tenant.

### Tech Stack Utama
* **Backend:** Go (Golang) dengan Web Framework (Gin/Echo) & GORM ORM
* **Frontend:** Next.js / Vue.js
* **Database:** PostgreSQL / MySQL (Multi-tenant isolation)
* **Caching & Message Bus:** Redis & Asynq (Distributed Task Queue)
* **Connection Pooler:** PgBouncer (Rekomendasi Production)

---

## 2. Analisis Arsitektur & Rekomendasi Teknis (Gap Analysis)

Berdasarkan tinjauan arsitektur teknis, beberapa area kritis berikut harus diterapkan untuk menjamin keandalan sistem pada skala produksi (*production-ready*):

### 2.1. Tenant Lifecycle & Resource Cleanup
* **Masalah:** Penanganan status tenant (*Suspend*, *Soft Delete*, *Terminate*) berisiko menyisakan koneksi TCP/GORM yang menggantung di memori aplikasi (`Manager.tenants map`).
* **Solusi:** Implementasikan fungsi `CloseTenantDB(companyID string) error` secara eksplisit pada `database.Manager` untuk memanggil `sqlDB.Close()` dan menghapus entri *map* guna mengantisipasi kebocoran memori (*memory leak*).

### 2.2. Keamanan Kredensial Database Tenant
* **Masalah:** Kredensial koneksi database tenant pada skema Platform/Master tidak boleh disimpan dalam bentuk *plain text*.
* **Solusi:** Terapkan enkripsi simetris **AES-256-GCM** pada layer *repository* sebelum kolom `TenantConnection.Password` disimpan atau dibaca dari database master. Kunci enkripsi dikelola melalui *environment variable* `HRIS_ENCRYPTION_KEY`.

### 2.3. Optimasi Connection Pooling & Database Limit
* **Masalah:** Alokasi koneksi maksimum yang terlalu besar per tenant (misal `SetMaxOpenConns(100)`) berisiko menghabiskan kuota `max_connections` pada host database jika jumlah tenant meningkat.
* **Solusi:**
  * Turunkan batas standar *pool* per tenant (misal: 10–20 koneksi).
  * Terapkan kebijakan pengosongan koneksi (*idle connection eviction policy*).
  * Manfaatkan connection pooler eksternal seperti **PgBouncer** pada lingkungan produksi.

### 2.4. Penanganan Dialek SQL Migrasi
* **Masalah:** Penggunaan eksekusi file `.sql` mentah memerlukan penanganan khusus jika sistem mendukung dual-driver (PostgreSQL & MySQL) karena perbedaan sintaksis DDL (seperti `UUID` vs `VARCHAR/CHAR(36)`).
* **Solusi:** Pisahkan struktur migrasi berdasarkan jenis driver (`migrations/tenant/postgres/` dan `migrations/tenant/mysql/`) atau manfaatkan perkakas migrasi agnostik (seperti Goose atau Atlas).

### 2.5. Sinkronisasi Cache Terdistribusi
* **Masalah:** Pembaruan *Feature Flags* atau *Permissions* di Redis pada satu instance Go server perlu dikonsumsi secara konsisten oleh instance lainnya pada lingkungan *multi-node deployment*.
* **Solusi:** Manfaatkan **Redis Pub/Sub** untuk menyiarkan *event* pembaruan data agar cache lokal pada seluruh replika server langsung ter-involidasi secara *real-time*.

---

## 3. Status & Roadmap Kelengkapan Modul HRIS

### 3.1. Modul Inti (Completed / Core Engine)
Modul-modul berikut telah selesai didesain dan diimplementasikan sebagai fondasi utama sistem:

* [x] **Platform & Tenant Management:** Provisioning DB multi-tenant, isolasi database, switching context tenant.
* [x] **Organization Management:** Organization Summary, Positions/Jabatan, Zone/Wilayah.
* [x] **Employee Management:** Data personal, kontak, alamat, keluarga, pendidikan, dokumen, riwayat kerja, rekening/pajak, serta pengaturan akun.
* [x] **Job Management:** Deskripsi jabatan lengkap (*Responsibilities, Working Activities, Operational/HR Authorities, Working Risks, Title Subs, Assets, Financials*).
* [x] **Competency Management:** *Competency Groups*, Kamus Kompetensi, dan *Potency Competencies*.

### 3.2. Modul Operasional & Siklus Karier (Planned / Phase 2 Roadmap)
Untuk melengkapi cakupan *Full-Suite HRIS*, modul-modul operasional berikut masuk dalam skala prioritas pengembangan tahap berikutnya:

* [ ] **Employee Movement & Career Management:**
  * **Promosi & Demosi:** Perubahan jenjang jabatan, posisi, dan penyesuaian kelas gaji.
  * **Perpanjangan Kontrak (PKWT):** Pengelolaan masa berlaku kerja, peringatan jatuh tempo, dan adendum kontrak.
  * **Pensiun & Offboarding/PHK:** Pengelolaan masa purna bakti, pemutusan hubungan kerja, dan integrasi perhitungan pesangon/uang jasa pada modul Payroll.
* [ ] **Time & Attendance:** Perekaman presensi, penjadwalan *shift*, lembur (*overtime*), dan kalkulasi keterlambatan.
* [ ] **Leave & Time Off:** Pengajuan cuti, sakit, izin, manajemen kuota cuti tahunan, dan *multi-level approval*.
* [ ] **Payroll & Compensation Engine:** Kalkulasi gaji bersih/kotor, tunjangan/potongan, PPh 21, BPJS Ketenagakerjaan/Kesehatan, dan slip gaji digital.
* [ ] **Performance Management:** Penilaian kinerja (KPI, OKR, review 360) yang terintegrasi langsung dengan modul *Job Management* dan *Competency*.
* [ ] **Reimbursement & Claim:** Pengajuan dan verifikasi klaim kesehatan maupun operasional dinas.
* [ ] **Recruitment & Onboarding (ATS):** Manajemen kandidat, alur seleksi, hingga otomatisasi pendaftaran karyawan baru (*onboarding*).

---

## 4. Matriks Prioritas Eksekusi

| Area | Komponen | Prioritas | Action Item Utama |
| :--- | :--- | :---: | :--- |
| **Security** | Tenant Credentials | 🔴 High | Enkripsi kolom `password` pada tabel `tenant_connections` menggunakan AES-256-GCM. |
| **Database** | SQL Dialect | 🔴 High | Pemisahan skrip DDL/Migration untuk PostgreSQL dan MySQL. |
| **Resource** | Lifecycle Tenant | 🟡 Medium | Implementasi penutupan koneksi DB aktif (`CloseTenantDB`) saat status tenant non-aktif. |
| **Performance**| Connection Pool | 🟡 Medium | Penyesuaian `SetMaxOpenConns` per tenant & konfigurasi PgBouncer. |
| **Architecture**| Modul Operasional & Career | 🟢 Low | Inisiasi desain skema database modul *Employee Movement*, *Time & Attendance*, dan *Payroll*. |

---
*Document Version: 1.6-Updated (v2)*  
*Status: Approved for Architecture Enhancement*
