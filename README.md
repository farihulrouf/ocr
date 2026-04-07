# 🚀 SEIDO OCR Platform
### AI-Powered Multi-tenant Expense Management SaaS
*従業員向け領収書OCR & 経費精算プラットフォーム*

SEIDO OCR (精励) adalah platform **Multi-tenant SaaS** yang mendigitalisasi pengelolaan struk belanja dan laporan pengeluaran. Dengan teknologi **AI OCR (Mistral AI)**, platform ini mengotomatisasi ekstraksi data dari struk untuk audit keuangan perusahaan secara real-time.

---

## 🛠️ Tech Stack & Infrastructure
* **Backend:** Go 1.24 (Fiber Framework)
* **Database:** PostgreSQL (GORM)
* **Cache & Queue:** Redis
* **Storage:** AWS S3 (Simulated via LocalStack)
* **Infrastructure:** Terraform & Docker (Multi-stage Build)
* **AI Engine:** Mistral AI (OCR Extraction)

---

## ✨ Fitur Utama (Key Features)

* 📸 **AI OCR Extraction:** Ekstraksi otomatis (Merchant, Tax, Total, Items) dari JPG/PNG via Mistral AI.
* 🏢 **Multi-tenancy:** Isolasi data antar perusahaan menggunakan middleware tenant yang ketat.
* ⚖️ **Approval Workflow:** Manager dapat melakukan *Bulk Approve/Reject/Restore* pada laporan karyawan.
* 🔐 **RBAC Security:** Akses kontrol berbasis peran (SuperAdmin, TenantAdmin, Employee).
* ⚙️ **Automated Infrastructure:** Setup environment (DB, Redis, S3) secara otomatis menggunakan Terraform.

---

## 📂 Struktur Proyek
```text
ocr-saas-backend/
├── cmd/                # Entry points (API & OCR Worker)
├── configs/            # Koneksi Database, S3, Redis & Seeder
├── internal/
│   ├── handler/        # Controller API (auth, receipts, reports)
│   ├── models/         # GORM Entities & Database Schema
│   ├── repository/     # Data Access Layer (SQL Queries)
│   └── service/        # Business Logic & AI OCR Integration
├── terraform/          # Infrastructure as Code (Docker & S3 Setup)
├── middleware/         # JWT, Tenant, & Role-based Auth
├── pkg/                # Utility (JWT, Response Helpers)
└── Dockerfile          # Multi-stage build (Alpine-based)

## 🚀 Quick Start (Automated Setup)

1. Konfigurasi Environment
Pastikan file .env di root folder sudah sesuai dengan mode Docker:

# Database (Internal Docker Name)
DB_HOST=postgres_db
DB_PORT=5432

# App
APP_PORT=8080
MISTRAL_API_KEY=your_mistral_api_key

# Services
REDIS_ADDR=redis_db:6379
S3_ENDPOINT=http://localstack:4566
S3_BUCKET=ocr-bucket

---
