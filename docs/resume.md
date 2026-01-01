Berdasarkan struktur model yang kamu buat, pengguna (user) dalam sistem ini terbagi menjadi beberapa tingkatan peran (Role) dengan tanggung jawab yang berbeda.
Berikut adalah daftar pengguna yang akan menggunakan aplikasi ini:
1. Super Admin (SaaS Owner)
Ini adalah tim internal kamu (pemilik aplikasi). Mereka tidak terikat pada satu Tenant saja.
Tugas:
Mengelola SubscriptionPlan (membuat paket harga).
Memantau Tenant (perusahaan mana saja yang daftar).
Melihat TenantUsage (siapa yang sudah habis kuota OCR-nya).
Melakukan maintenance sistem.
2. Tenant Admin (Company Admin)
Biasanya orang HR atau IT di perusahaan yang menggunakan jasa aplikasi kamu.
Tugas:
Mengatur CompanySetting (format tanggal, mata uang).
Mengelola struktur Department.
Menambahkan/mengundang User (Karyawan) ke dalam sistem.
Menentukan siapa yang menjadi atasan siapa di tabel UserApprover.
Mengatur ApprovalWorkflow (alur persetujuan di perusahaan tersebut).
3. Employee (Karyawan Biasa)
Pengguna paling banyak dalam sistem.
Tugas:
Mengunggah foto struk belanja ke tabel Receipt.
Melakukan verifikasi hasil OCR (apakah nama toko dan jumlah uang sudah benar).
Mengelompokkan struk-struk tersebut ke dalam satu Expense Report.
Mengajukan laporan tersebut untuk di-approve.
4. Approver (Manager/Kepala Divisi)
User yang memiliki bawahan. Di model kamu, ini diatur dalam tabel UserApprover.
Tugas:
Menerima notifikasi pengajuan pengeluaran dari bawahan.
Mengecek detail Receipt dan ReceiptItem.
Memberikan keputusan: APPROVE atau REJECT.
Membubuhkan Digital Hanko (stempel digital) yang tercatat di ApprovalLog.
5. Finance / Accountant (Bagian Keuangan)
User khusus yang memiliki akses ke data finansial seluruh perusahaan.
Tugas:
Mengelola AccountCategory (COA - Chart of Accounts) dan TaxRate.
Mengelola VendorMaster (daftar toko/supplier).
Melakukan pengecekan akhir apakah pajak (TaxRegistrationID) sudah sesuai aturan pemerintah (Qualified Invoice).
Melakukan Export data pengeluaran (melalui ExportLog) untuk dimasukkan ke software akuntansi perusahaan (seperti Freee, Money Forward, atau SAP).
Ringkasan dalam Tabel Peran (Role):
Role	Akses Utama	Tujuan Utama
Super Admin	Semua Tenant	Kelola Billing & Platform
Admin Tenant	User, Dept, Workflow	Konfigurasi Perusahaan
Employee	Receipt, Report	Klaim Pengeluaran (Reimbursement)
Approver	ApprovalLog, UserApprover	Validasi & Kontrol Budget
Finance	Tax, Vendor, Export	Audit & Pembukuan
Tips Tambahan:
Karena di model User kamu sudah ada field Role string, kamu bisa menentukan nilai isinya seperti: "ADMIN", "FINANCE", "APPROVER", atau "EMPLOYEE".
Apakah kamu butuh bantuan untuk membuat logic pengecekan role (RBAC) di bagian backend-nya