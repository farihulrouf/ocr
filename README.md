# SEIDO OCR Platform - AI-Powered Expense Management SaaS
#### 従業員向け領収書OCR & 経費精算プラットフォーム

SEIDO OCR Platform (精励) adalah platform Multi-tenant SaaS yang memungkinkan karyawan, manajer, dan administrator sistem untuk mengelola struk belanja dan laporan pengeluaran secara otomatis menggunakan teknologi AI OCR. Dibangun dengan Go (Fiber) sebagai backend utama, aplikasi ini menangani proses upload, ekstraksi teks, hingga alur persetujuan laporan keuangan secara terintegrasi.

## 📌 プロジェクト概要 (Project Overview)

* **対象ステークホルダー (Stakeholders):** 従業員 (Employee), テナント管理者 (Manager), システム管理者 (Super Admin)
* **提供価値 (Value Proposition):** Struk OCR otomatis, pelacakan pengeluaran real-time, manajemen tenant, dan approval workflow laporan keuangan.
* **データ運用 (Data Operation):** Data transaksi disimpan di PostgreSQL (Aiven Cloud), gambar struk disimpan di S3 (LocalStack), dan antrean ekstraksi dikelola oleh Redis.
* **API ドキュメント:** * `docs/api/auth.md` (認証)
    * `docs/api/employee.md` (従業員向け)
    * `docs/api/manager.md` (マネージャー向け)
    * `docs/api/system.md` (システム管理者向け)
* **運用方針 (Operation Policy):** JWT + Role-based middleware (SuperAdmin, TenantAdmin, Employee) untuk akses kontrol API via prefix `/v0/api`.

## ✨ 主な特徴 (Key Features)

* 📊 **AI OCR Upload:** Upload struk belanja (JPG/PNG) dan ekstraksi data otomatis (Merchant, Tax, Total, Items) via Mistral AI.
* 🧑‍💻 **ロール別ダッシュボード (Role Dashboard):** Dashboard khusus untuk Employee (histori pribadi) dan Manager (statistik departemen/tenant).
* 🏢 **マルチテナント (Multi-tenant):** Isolasi data antar perusahaan melalui middleware tenant dan `tenant_id` yang ketat.
* 📂 **レポート管理 (Expense Reports):** Pengelompokan struk ke dalam laporan mingguan/bulanan untuk diajukan ke manajer.
* ✅ **一括操作 (Bulk Operations):** Manajer dapat menyetujui, menolak, atau memulihkan struk secara massal (Bulk Approve/Reject/Restore).

## 🚀 クイックスタート (Quick Start)

### 前提条件 (Prerequisites)

* **Go 1.21+**
* **Docker & Docker Compose** (untuk Redis & LocalStack)
* **Python 3** (untuk `localstack-cli` dalam venv)

### 最速セットアップ (Setup Steps)

```bash
# 1. リポジトリの取得
git clone <repository-url>
cd ocr-saas-backend

# 2. 依存関係のインストール
go mod tidy

# 3. インフラの起動 (Redis & LocalStack)
# Pastikan Token LocalStack sudah di-export di terminal
source venv/bin/activate
localstack start -d

# 4. S3バケットの作成
awslocal s3 mb s3://ocr-bucket

# 5. アプリ起動
go run main.go

ocr-saas-backend/
├── cmd/                # Entry point aplikasi (main.go)
├── configs/            # Koneksi Database, S3, Redis & Database Seeder
├── docs/               # Dokumentasi API & Requirements.md
├── internal/
│   ├── handler/        # Controller API (auth, receipts, ocr, dsb)
│   ├── models/         # Entity GORM & Schema Database
│   ├── repository/     # Data Access Layer (SQL Queries)
│   ├── service/        # Business Logic (AI OCR Integration, Logic)
│   └── routes/         # Routing definition (SetupRoutes v0/api)
├── middleware/         # Protected(), EmployeeOnly(), SuperAdminOnly()
├── pkg/                # Utility (JWT, Response Helpers)
└── README.md




# Database (Aiven Cloud / Local Docker)
DB_HOST=pg-2a32409f-farihulrouf-0d16.a.aivencloud.com
DB_PORT=21092
DB_USER=avnadmin
DB_PASSWORD=********
DB_NAME=ocr
DB_SSLMODE=require

# AWS S3 (LocalStack)
S3_ENDPOINT=http://localhost:4566
S3_BUCKET=ocr-bucket
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
STORAGE_USE_SSL=false

# Redis (Queue & Cache)
REDIS_ADDR=localhost:6379

# AI Engine
MISTRAL_API_KEY=********