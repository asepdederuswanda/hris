# Panduan & Standar UI/UX HRIS Enterprise Grade

> **Tech Stack:** Vue 3 (`<script setup>`) + PrimeVue v4+ + Tailwind CSS  
> **Konsep:** High-Density, Compact, Modal-First, & Single-Page Dashboard (Bebas Pindah Halaman)

---

## 1. Konsep Utama & Prinsip Desain

Dalam sistem HRIS skala Enterprise, produktivitas pengguna sangat bergantung pada kecepatan eksekusi data. **Prinsip Utama:** Pengguna tidak boleh sering berpindah halaman penuh (*page navigation/routing*) hanya untuk melihat detail, mengedit data, atau melakukan persetujuan (*approval*).

| Prinsip UX | Implementasi Teknikal | Manfaat Utama |
| :--- | :--- | :--- |
| **Modal & Dialog First** | `Dialog` (Modal) & `Drawer` (Slide-over) dari PrimeVue. | Mencegah perpindahan halaman & menjaga konteks pengguna di tabel utama. |
| **Single Viewport (Zero Scroll)** | Master-Detail (Split View) & Tabbed Interface. | Semua informasi penting terlihat tanpa perlu *scroll* panjang. |
| **High Density Layout** | Padding ketat (`py-1.5`, `px-3`), font size `12px`–`14px` (`text-xs`–`text-sm`). | Memaksimalkan data yang tampil dalam 1 layar. |

---

## 2. Hirarki Penggunaan Modal & Overlay (Kapan Menggunakan Apa?)

Agar tidak terjadi *modal-cluttering* (terlalu banyak popup menumpuk), ikuti aturan penggunaan overlay berikut:

### A. Compact `Dialog` (Modal Standar)
* **Penggunaan:** Form ringkas (1-2 kolom), konfirmasi aksi (*Approve/Reject*), ubah status, atau input catatan cepat.
* **Properti PrimeVue:** `<Dialog modal :style="{ width: '400px' }" header="...">`
* **Keunggulan:** Fokus instan ke aksi spesifik tanpa mengganggu layar latar belakang.

### B. Slide-over `Drawer` (Modal Samping)
* **Penggunaan:** Form data kompleks (misal: *Tambah Karyawan Baru*, *Detail Gaji & Komponen*).
* **Properti PrimeVue:** `<Drawer v-model:visible="open" position="right" class="!w-[500px]">`
* **Keunggulan:** Memberikan area vertikal lebih luas untuk form ber-tab tanpa menutupi seluruh layar utama.

### C. `Popover` / Context Menu
* **Penggunaan:** Menu tindakan cepat pada baris tabel (Edit, Hapus, Kirim Slip Gaji, Cetak ID).
* **Keunggulan:** Menghemat kolom pada `DataTable` sehingga tabel tidak terasa sesak.

---

## 3. Master Prompt untuk AI Assistant (Diperbarui)

Gunakan prompt acuan di bawah ini untuk disalin ke AI (ChatGPT / Claude / Gemini) agar seluruh kode yang dihasilkan selalu konsisten dengan standar HRIS Anda:

```text
[CONTEXT & ROLE]
Anda adalah Senior Full-Stack Frontend Engineer & UX Designer spesialis sistem Enterprise ERP/HRIS. Tugas Anda adalah membantu saya membangun aplikasi Human Resource Information System (HRIS) Enterprise Grade yang modern, minimalis, efisien ruang (high-density), user-friendly, bebas dari scroll berlebih, dan SANGAT MINIM PINDAH HALAMAN.

[TECH STACK]
- Framework: Vue 3 (Composition API / <script setup>)
- UI Library: PrimeVue (Gunakan skema komponen modern v4+, Pass-Through, atau Design Tokens)
- Utility Styling: Tailwind CSS
- Icons: PrimeIcons / Lucide Icons

[DESIGN & UX PRINCIPLES - HARUS DITURUTI]
1. Minimum Page Navigation & Modal-First Architecture:
   - UTAMAKAN penggunaan Dialog/Modal, Drawer (Slide-over), dan Popover alih-alih berpindah halaman/route baru.
   - Semua aksi CRUD (Create, Read, Update, Delete) dan Approval HARUS dilakukan di atas halaman utama menggunakan Modal/Drawer.
   - Pengguna harus tetap berada di tabel/dashboard utama saat melakukan operasi data.

2. Compact & High-Density First:
   - Hindari ruang kosong (whitespace) yang tidak perlu.
   - Gunakan padding dan margin yang ketat (misal: py-1.5, py-2, px-3, p-3).
   - Ukuran teks dominan adalah text-xs (12px) sampai text-sm (14px).
   - Gunakan border-radius yang kecil/presisi (rounded-md atau rounded-sm).

3. Zero / Minimal Scroll Architecture:
   - Manfaatkan layout 'h-screen' atau 'h-[calc(100vh-...)]' dengan 'overflow-hidden' pada container utama.
   - Gunakan pola Split-Screen / Master-Detail (Panel Kiri: DataTable, Panel Kanan: Detail View/Tabs).
   - Gunakan komponen 'Tabs' horizontal di dalam Modal/Drawer untuk memecah form berukuran besar.

4. Standar Komponen PrimeVue & Tailwind:
   - DataTable: Gunakan opsi compact/small (`p-datatable-sm`), scrollable dengan `scrollHeight="flex"`, virtual scroll jika perlu, dan pinned/frozen columns untuk aksi.
   - Dialog & Drawer: Selalu atur ukuran lebar eksplisit (misal: !w-[450px] untuk Dialog, !w-[550px] untuk Drawer).
   - Filters & Actions: Satukan Search Bar, Filter Select, dan Quick Action Buttons dalam satu baris horizontal menggunakan `Toolbar` atau Flexbox ringkas.
   - Warna Status: Gunakan Soft Badge dengan Tailwind (misal: `bg-emerald-50 text-emerald-700` untuk aktif/disetujui, `bg-amber-50 text-amber-700` untuk pending).

[OUTPUT FORMAT REQUIREMENT]
Setiap kali memberikan kodingan atau solusi UI/UX:
1. Tuliskan kode Vue 3 lengkap menggunakan Single File Component (SFC) <template> dan <script setup>.
2. Pastikan gabungan kelas Tailwind CSS dan atribut PrimeVue mematuhi aturan high-density & modal-first di atas.
3. Berikan penjelasan singkat mengenai struktur layout jika ada pola UX khusus yang diterapkan.
```

---

## 4. Standar Warna Badge Status

Gunakan palet warna status yang lembut (*soft badges*) agar tampilan tetap profesional dan tidak mencolok secara berlebihan:

* **Active / Approved:** `bg-emerald-50 text-emerald-700 border border-emerald-200`
* **Pending / Review:** `bg-amber-50 text-amber-700 border border-amber-200`
* **Draft / Info:** `bg-indigo-50 text-indigo-700 border border-indigo-200`
* **Inactive / Rejected:** `bg-rose-50 text-rose-700 border border-rose-200`
