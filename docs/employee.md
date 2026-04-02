# SEIDO OCR Enterprise - Employee API Specification

Dokumentasi ini mendefinisikan API untuk role **Employee** yang mencakup pengelolaan Receipt (OCR) dan Expense Reports (Laporan Pengeluaran).

---

## 1. Summary Table

| ID | Endpoint | Method | Description |
| :--- | :--- | :--- | :--- |
| **RECEIPT** | | | |
| E-API-01 | `/v0/api/emp/receipt` | GET | Ambil daftar receipt milik sendiri (Paging/Filter). |
| E-API-02 | `/v0/api/emp/receipt/:id` | GET | Detail receipt (OCR result & items). |
| E-API-03 | `/v0/api/emp/receipt/upload` | POST | Upload gambar receipt untuk proses AI/OCR. |
| E-API-04 | `/v0/api/emp/receipt/:id` | PUT | Update data dasar receipt (Store, Date, Total). |
| E-API-05 | `/v0/api/emp/receipt/:id` | DELETE | Hapus receipt milik sendiri. |
| E-API-06 | `/v0/api/emp/receipt/:id/status` | GET | Cek status proses OCR (Polling). |
| E-API-07 | `/v0/api/emp/receipt/:id/items` | POST | Tambah item baru ke dalam receipt. |
| E-API-08 | `/v0/api/emp/receipt/items/:itemId` | PUT | Update nama/harga item spesifik. |
| E-API-09 | `/v0/api/emp/receipt/items/:itemId` | DELETE | Hapus item dari receipt. |
| **REPORT** | | | |
| E-API-10 | `/v0/api/emp/reports/` | GET | Daftar laporan pengeluaran milik sendiri. |
| E-API-11 | `/v0/api/emp/reports/` | POST | Buat laporan pengeluaran baru (Draft). |
| E-API-12 | `/v0/api/emp/reports/:id` | GET | Detail laporan & daftar receipt di dalamnya. |
| E-API-13 | `/v0/api/emp/reports/:id` | PUT | Update informasi dasar laporan. |
| E-API-14 | `/v0/api/emp/reports/:id/submit` | POST | Submit laporan (Kunci data untuk verifikasi). |
| E-API-15 | `/v0/api/emp/reports/:id/receipts` | POST | Tambah banyak receipt sekaligus ke laporan. |
| E-API-16 | `/v0/api/emp/reports/:id/receipts/:receiptId` | DELETE | Lepas receipt dari laporan tertentu. |
| **DASHBOARD**| | | |
| E-API-17 | `/v0/api/emp/dashboard` | GET | Ambil ringkasan statistik user (Total expense, dll). |

---

## 2. Detail Endpoints - Receipt Section

### E-API-01: List My Receipts
- **Method**: `GET`
- **Path**: `/v0/api/emp/receipt`
- **Query Params**:
  - `page`: number (default 1)
  - `page_size`: number (default 10)
  - `q`: string (Search by store name)
  - `status`: string (`PENDING` | `SUCCESS` | `REPORTED`)

### E-API-03: Upload Receipt (OCR)
- **Method**: `POST`
- **Path**: `/v0/api/emp/receipt/upload`
- **Content-Type**: `multipart/form-data`
- **Body**: 
  - `file`: Image/PDF binary.

### E-API-06: Check OCR Status
- **Method**: `GET`
- **Path**: `/v0/api/emp/receipt/:id/status`
- **Response**:
```json
{
  "id": "uuid-string",
  "status": "SUCCESS"
}