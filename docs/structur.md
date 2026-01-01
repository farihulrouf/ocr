# ocr
.
├── cmd/
│   └── api/
│       └── main.go          # Entry point aplikasi
├── configs/
│   └── config.go            # Load .env dan konfigurasi sistem
├── internal/
│   ├── models/              # Struct GORM (Sudah Anda buat)
│   ├── middleware/          # JWT, RBAC, Multi-tenancy check
│   ├── repository/          # Layer Database (Query GORM)
│   │   ├── auth_repo.go
│   │   ├── receipt_repo.go
│   │   └── ...
│   ├── service/             # Layer Bisnis Logic (Validasi berat, hitung pajak, call AI)
│   │   ├── auth_service.go
│   │   ├── receipt_service.go
│   │   └── ...
│   ├── handler/             # Layer HTTP (Parsing request & kirim response JSON)
│   │   ├── auth_handler.go
│   │   ├── tenant_handler.go
│   │   ├── org_handler.go
│   │   ├── finance_handler.go
│   │   ├── receipt_handler.go
│   │   ├── report_handler.go
│   │   ├── approval_handler.go
│   │   ├── audit_handler.go
│   │   └── system_handler.go
│   └── routes/
│       ├── routes.go        # Master route (v0 group)
│       └── ...              # Opsional: split route per file jika terlalu panjang
├── pkg/                     # Helper universal (non-bisnis logic)
│   ├── utils/               # JWT, Bcrypt, UUID generator
│   └── ocr/                 # Integrasi OpenAI / Vision API
├── docs/                    # Swagger / API Documentation
├── uploads/                 # File sementara (struk)
├── .env
├── go.mod
└── go.sum