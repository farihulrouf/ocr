Tentu, ini adalah dokumentasi lengkap untuk **96 API** yang Anda miliki. Karena jumlahnya sangat banyak, saya akan menyajikannya secara sistematis per grup dengan detail **Method, Endpoint, Request Body (jika ada), dan Response Body**.

---

### **1. AUTHENTICATION & SECURITY (10 API)**
Fokus pada akses pengguna dan keamanan akun.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | Login | POST | `/auth/login` | `{"email", "password"}` | `{"token", "refresh_token", "user_obj"}` |
| 2 | Refresh Token | POST | `/auth/refresh-token` | `{"refresh_token"}` | `{"token"}` |
| 3 | Get Profile | GET | `/auth/me` | - | `{"id", "name", "email", "role", "dept"}` |
| 4 | Update Profile | PUT | `/auth/profile` | `{"name", "avatar"}` | `{"message": "Profile updated"}` |
| 5 | Update Pass | PUT | `/auth/password` | `{"old_pass", "new_pass"}` | `{"message": "Password updated"}` |
| 6 | Forgot Pass | POST | `/auth/forgot-password`| `{"email"}` | `{"message": "Link sent"}` |
| 7 | Reset Pass | POST | `/auth/reset-password` | `{"token", "new_pass"}`| `{"message": "Password reset"}` |
| 8 | Verify Email | POST | `/auth/verify-email` | `{"token"}` | `{"message": "Email verified"}` |
| 9 | Setup 2FA | POST | `/auth/2fa/setup` | - | `{"qr_code", "secret_key"}` |
| 10| Logout | POST | `/auth/logout` | - | `{"message": "Session cleared"}` |

---

### **2. TENANT & COMPANY (8 API)**
Manajemen entitas perusahaan/SaaS.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | Get Info | GET | `/tenant/info` | - | `{"company_name", "tax_id", "address"}` |
| 2 | Update Info | PUT | `/tenant/info` | `{"company_name", ...}` | `{"message": "Info updated"}` |
| 3 | Get Settings | GET | `/tenant/settings` | - | `{"currency", "timezone", "ocr_auto"}` |
| 4 | Update Settings| PUT | `/tenant/settings` | `{"settings_obj"}` | `{"message": "Settings updated"}` |
| 5 | Get Plan | GET | `/tenant/subscription`| - | `{"plan_id", "status", "exp_date"}` |
| 6 | Upgrade Plan | POST | `/tenant/subscription/upgrade`| `{"plan_id"}` | `{"checkout_url"}` |
| 7 | Usage Stats | GET | `/tenant/usage` | - | `{"ocr_limit", "ocr_used", "users"}` |
| 8 | Terminate | DELETE | `/tenant/terminate` | - | `{"message": "Deletion requested"}` |

---

### **3. ORGANIZATION & USERS (15 API)**
Struktur organisasi dan manajemen karyawan.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | List Depts | GET | `/org/departments` | - | `[{"id", "name", "manager_id"}]` |
| 2 | Create Dept | POST | `/org/departments` | `{"name"}` | `{"id", "name"}` |
| 3 | Detail Dept | GET | `/org/departments/:id` | - | `{"id", "name", "users": []}` |
| 4 | Update Dept | PUT | `/org/departments/:id` | `{"name"}` | `{"message": "Updated"}` |
| 5 | Delete Dept | DELETE | `/org/departments/:id` | - | `{"message": "Deleted"}` |
| 6 | List Users | GET | `/org/users` | - | `[{"id", "name", "role", "dept"}]` |
| 7 | Invite User | POST | `/org/users` | `{"email", "role"}` | `{"message": "Invite sent"}` |
| 8 | Detail User | GET | `/org/users/:id` | - | `{"id", "name", "history": []}` |
| 9 | Update User | PUT | `/org/users/:id` | `{"role", "dept_id"}` | `{"message": "User updated"}` |
| 10| Update Status | PUT | `/org/users/:id/status`| `{"is_active"}` | `{"status": "active/inactive"}` |
| 11| Delete User | DELETE | `/org/users/:id` | - | `{"message": "Soft deleted"}` |
| 12| Org Structure | GET | `/org/hierarchy` | - | `{"tree_data": {...}}` |
| 13| Get Approvers | GET | `/org/users/:id/approvers`| - | `[{"id", "name", "level"}]` |
| 14| Set Approver | POST | `/org/users/:id/approvers`| `{"approver_id"}` | `{"message": "Assigned"}` |
| 15| Rem Approver | DELETE | `/org/users/:id/approvers/:aId`| - | `{"message": "Removed"}` |

---

### **4. FINANCE MASTER DATA (16 API)**
Data master untuk kebutuhan akuntansi (Invois Seido Ready).

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1-4| Categories | CRUD | `/finance/categories` | `{"name", "code"}` | List/Detail/Msg |
| 5 | Import Cat | POST | `/finance/categories/import`| `csv_file` | `{"imported_count": 100}` |
| 6-8| Tax Rates | CRUD | `/finance/tax-rates` | `{"name", "percent"}` | List/Detail/Msg |
| 9-10| Pay Methods | G/P | `/finance/payments` | `{"name"}` | List/Detail/Msg |
| 11-15| Vendors | CRUD | `/finance/vendors` | `{"name", "t_number"}` | List/Detail/Msg |
| 16 | Verify T-No | GET | `/finance/vendors/verify/:tn`| - | `{"is_valid", "official_name"}` |

---

### **5. RECEIPTS & AI-OCR (15 API)**
Inti dari aplikasi: Pengolahan struk dengan AI.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | My Receipts | GET | `/receipts` | - | `[{"id", "total", "status"}]` |
| 2 | All (Admin) | GET | `/receipts/all` | - | `[{"user", "receipt_obj"}]` |
| 3 | Upload OCR | POST | `/receipts/upload` | `file_binary` | `{"id", "status": "processing"}` |
| 4 | Webhook OCR | POST | `/receipts/webhook/ocr` | `json_from_ai` | `{"message": "Data received"}` |
| 5 | Detail Struk | GET | `/receipts/:id` | - | `{"id", "items": [], "ocr_raw"}` |
| 6 | Confirm Data | PUT | `/receipts/:id` | `{"total", "date"}` | `{"status": "confirmed"}` |
| 7 | Delete Struk | DELETE | `/receipts/:id` | - | `{"message": "Deleted"}` |
| 8 | Bulk Delete | POST | `/receipts/bulk/delete` | `{"ids": []}` | `{"deleted": 5}` |
| 9 | Bulk Update | POST | `/receipts/bulk/update-category`| `{"ids", "cat_id"}` | `{"updated": 5}` |
| 10| Secure Image | GET | `/receipts/:id/image` | - | `{"url": "signed-s3-url"}` |
| 11| Re-trigger | POST | `/receipts/:id/re-ocr` | - | `{"status": "re-processing"}` |
| 12| Adv Search | GET | `/receipts/search/advanced`| `?min=100&cat=1` | `[results]` |
| 13| Add Item | POST | `/receipts/:id/items` | `{"name", "price"}` | `{"item_id": 1}` |
| 14| Update Item | PUT | `/receipts/items/:itemId`| `{"price"}` | `{"message": "Item updated"}` |
| 15| Delete Item | DELETE | `/receipts/items/:itemId`| - | `{"message": "Item deleted"}` |

---

### **6. EXPENSE REPORTS (10 API)**
Bundling struk menjadi laporan klaim.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | List Reports | GET | `/reports` | - | `[{"id", "title", "total"}]` |
| 2 | Create Report | POST | `/reports` | `{"title", "r_ids":[]}`| `{"report_id": 101}` |
| 3-5| Detail/U/D | CRUD | `/reports/:id` | `{"title"}` | Detail/Msg |
| 6 | Submit | POST | `/reports/:id/submit` | - | `{"status": "waiting_approval"}`|
| 7 | Cancel | POST | `/reports/:id/cancel` | - | `{"status": "draft"}` |
| 8 | Monthly Stat | GET | `/reports/stats/monthly`| - | `{"month": "Jan", "total": 500}`|
| 9 | Preview PDF | GET | `/reports/export/preview`| - | `{"pdf_url": "..."}` |
| 10| Bulk Submit | POST | `/reports/bulk/submit` | `{"ids": []}` | `{"message": "3 submitted"}` |

---

### **7. WORKFLOW & APPROVALS (12 API)**
Manajemen persetujuan laporan.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | Pending Tasks | GET | `/approvals/pending` | - | `[{"report_id", "from_user"}]`|
| 2 | Action | POST | `/approvals/action` | `{"id", "status"}` | `{"message": "Approved"}` |
| 3 | Remand | POST | `/approvals/remand` | `{"id", "note"}` | `{"message": "Returned"}` |
| 4 | History | GET | `/approvals/history` | - | `[{"report_id", "decision"}]`|
| 5-7| Workflows | CRUD | `/approvals/workflows` | `{"name", "steps"}` | List/Detail/Msg |
| 8 | Detail Steps | GET | `/approvals/workflows/:id/steps`| - | `[{"order", "approver"}]` |
| 9 | Add Step | POST | `/approvals/workflows/:id/steps`| `{"approver_id"}` | `{"step_id": 1}` |
| 10| Update Step | PUT | `/approvals/steps/:stepId`| `{"order": 2}` | `{"message": "Updated"}` |
| 11| Delete Step | DELETE | `/approvals/steps/:stepId`| - | `{"message": "Removed"}` |
| 12| Delegations | GET | `/approvals/delegations`| - | `{"proxy_user_id": 99}` |

---

### **8. COMPLIANCE & AUDIT (10 API)**
Log dan ekspor data untuk audit pajak.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1 | Dashboard Sum| GET | `/audit/summary` | - | `{"compliance_rate": "99%"}` |
| 2 | Audit Trails | GET | `/audit/trails` | - | `[{"action", "user", "time"}]`|
| 3 | Receipt Hist | GET | `/audit/trails/receipt/:id`| - | `[{"status_change"}]` |
| 4 | User Hist | GET | `/audit/trails/user/:id` | - | `[{"last_login", "actions"}]`|
| 5 | Export CSV | POST | `/audit/export/csv` | `{"date_range"}` | `{"file_url": "..."}` |
| 6 | Export PDF | POST | `/audit/export/pdf-bundle`| `{"ids"}` | `{"file_url": "..."}` |
| 7 | Export Logs | GET | `/audit/export/logs` | - | `[{"filename", "status"}]` |
| 8 | Download | GET | `/audit/export/download/:id`| - | `binary_file` |
| 9 | Tax Summary | GET | `/audit/tax-summary` | - | `{"tax_8": 100, "tax_10": 200}`|
| 10| Integrity Chk | POST | `/audit/data-integrity/check`| - | `{"is_tampered": false}` |

---

### **9. SYSTEM ADMIN (10 API)**
Khusus Super Admin untuk manajemen platform.

| No | Nama API | Method | Endpoint | Request Body | Response (Data Only) |
|:---|:---|:---:|:---|:---|:---|
| 1-3| Plans | CRUD | `/system/plans` | `{"name", "price"}` | List/Detail/Msg |
| 4 | All Tenants | GET | `/system/tenants` | - | `[{"id", "name", "plan"}]` |
| 5 | Suspend Tenan| PUT | `/system/tenants/:id/status`| `{"status"}` | `{"message": "Suspended"}` |
| 6 | OCR Monitor | GET | `/system/monitoring/ocr` | - | `{"ai_health": "good"}` |
| 7 | Error Logs | GET | `/system/monitoring/errors`| - | `[{"error", "stacktrace"}]` |
| 8 | Maint. Mode | POST | `/system/maintenance/mode` | `{"is_on": true}` | `{"status": "maintenance"}` |
| 9 | Backup Config| GET | `/system/backup/configs` | - | `{"last_backup": "2024..."}` |
| 10| Version | GET | `/system/version` | - | `{"version": "v1.0.0"}` |

---

### **Catatan Implementasi:**
1.  **Status Code:** Gunakan `200` untuk sukses, `201` untuk Create, `204` untuk Delete, `401` untuk Auth Error, dan `403` untuk Role Error.
2.  **Wrappers:** Sangat disarankan membungkus semua respon di atas dalam format:
    `{ "status": "success", "data": { ... } }`
3.  **Pagination:** Untuk semua API `GET` (List), tambahkan query param `?page=1&limit=10`.

Apakah ada bagian tertentu (misal: **Receipts**) yang ingin Anda buatkan kode Handler Go-nya sekarang?