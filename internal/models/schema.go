package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model dengan UUID dan JSON Tags snake_case
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Hook GORM untuk generate UUID otomatis sebelum create
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

// --- GROUP 1: TENANT & AUTH ---

type SubscriptionPlan struct {
	Base
	Name        string `json:"name"`
	MaxReceipts int    `json:"max_receipts"`
	Price       int64  `json:"price"`
}

type Tenant struct {
	Base
	Name               string           `gorm:"not null" json:"name"`
	Subdomain          string           `gorm:"uniqueIndex;not null" json:"subdomain"`
	SubscriptionPlanID uuid.UUID        `gorm:"type:uuid" json:"subscription_plan_id"`
	SubscriptionPlan   SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscription_plan"`
	BusinessNumber     string           `json:"business_number"`
	Status             string           `gorm:"default:'ACTIVE'" json:"status"`
}

type CompanySetting struct {
	Base
	TenantID   uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID" json:"-"`
	DateFormat string    `gorm:"default:'YYYY/MM/DD'" json:"date_format"`
	Currency   string    `gorm:"default:'JPY'" json:"currency"`
	AutoOCR    bool      `gorm:"default:true" json:"auto_ocr"`
}

type Department struct {
	Base
	TenantID uuid.UUID  `gorm:"type:uuid" json:"tenant_id"`
	ParentID *uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	Name     string     `json:"name"`
	Code     string     `json:"code"`
}

type User struct {
	Base
	TenantID     uuid.UUID   `gorm:"type:uuid" json:"tenant_id"`
	Tenant       Tenant      `gorm:"foreignKey:TenantID" json:"tenant"` // RELASI KE TENANT
	DepartmentID *uuid.UUID  `gorm:"type:uuid" json:"department_id"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department"` // RELASI KE DEPT
	Email        string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string      `json:"-"` // Hidden
	Name         string      `json:"name"`
	Avatar       string      // <--- pastikan ada ini
	Role         string      `gorm:"default:'EMPLOYEE'" json:"role"`
}

type UserApprover struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uuid.UUID `gorm:"type:uuid;index" json:"employee_id"`
	ApproverID uuid.UUID `gorm:"type:uuid;index" json:"approver_id"`
}

// --- GROUP 2: FINANCE SETUP ---

type AccountCategory struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
}

type TaxRate struct {
	Base
	TenantID   uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Name       string    `json:"name"`
	Percentage int       `json:"percentage"`
}

type PaymentMethod struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Name     string    `json:"name"`
}

type VendorMaster struct {
	Base
	TenantID  uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Name      string    `json:"name"`
	TaxNumber string    `json:"tax_number"`
}

// --- GROUP 3: DOCUMENT & OCR ---

type Receipt struct {
	Base
	TenantID          uuid.UUID        `gorm:"type:uuid" json:"tenant_id"`
	UserID            uuid.UUID        `gorm:"type:uuid" json:"user_id"`
	User              User             `gorm:"foreignKey:UserID" json:"user"` // RELASI KE USER
	ReportID          *uuid.UUID       `gorm:"type:uuid" json:"report_id"`
	AccountCategoryID *uuid.UUID       `gorm:"type:uuid" json:"account_category_id"`
	AccountCategory   *AccountCategory `gorm:"foreignKey:AccountCategoryID" json:"account_category"`
	ImageURL          string           `json:"image_url"`
	Status            string           `gorm:"default:'PENDING'" json:"status"`
	StoreName         string           `json:"store_name"`
	TransactionDate   *time.Time       `json:"transaction_date"`
	TotalAmount       int64            `json:"total_amount"`
	TaxRegistrationID string           `json:"tax_id"`
	IsQualified       bool             `gorm:"default:false" json:"is_qualified"`
	LineItems         []ReceiptItem    `gorm:"foreignKey:ReceiptID" json:"line_items"` // RELASI KE ITEMS
}

type ReceiptItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReceiptID   uuid.UUID `gorm:"type:uuid" json:"receipt_id"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	TaxAmount   int64     `json:"tax_amount"`
	TaxRate     int       `json:"tax_rate"`
}

// --- GROUP 4: WORKFLOW ---

type ExpenseReport struct {
	Base
	TenantID    uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Title       string    `json:"title"`
	TotalAmount int64     `json:"total_amount"`
	Status      string    `gorm:"default:'PENDING'" json:"status"`
	Receipts    []Receipt `gorm:"foreignKey:ReportID" json:"receipts"`
}

type ApprovalWorkflow struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}

type ApprovalStep struct {
	Base
	WorkflowID uuid.UUID `gorm:"type:uuid" json:"workflow_id"`
	StepOrder  int       `json:"step_order"`
	ApproverID uuid.UUID `gorm:"type:uuid" json:"approver_id"`
}

// --- GROUP 5: LOGS & AUDIT ---

type ApprovalLog struct {
	Base
	ReceiptID       *uuid.UUID `gorm:"type:uuid" json:"receipt_id"`
	ExpenseReportID *uuid.UUID `gorm:"type:uuid" json:"expense_report_id"`
	UserID          uuid.UUID  `gorm:"type:uuid" json:"user_id"` // Approver
	User            User       `gorm:"foreignKey:UserID" json:"approver"`
	Action          string     `json:"action"` // APPROVE, REJECT
	Comment         string     `json:"comment"`
	DigitalHankoURL string     `json:"digital_hanko_url"`
}

type AuditTrail struct {
	Base
	TenantID  uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Action    string    `json:"action"`
	TableName string    `json:"table_name"`
	RecordID  string    `json:"record_id"`
	OldData   string    `gorm:"type:jsonb" json:"old_data"`
	NewData   string    `gorm:"type:jsonb" json:"new_data"`
}

type ExportLog struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	UserID   uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Format   string    `json:"format"`
	FileURL  string    `json:"file_url"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time
}

// TenantUsage untuk menyimpan limit & pemakaian OCR per tenant
type TenantUsage struct {
	TenantID uuid.UUID `gorm:"type:uuid;primaryKey" json:"tenant_id"`
	OCRLimit int64     `json:"ocr_limit"`
	OCRUsed  int64     `json:"ocr_used"`

	// Optional: relasi ke Tenant
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
}
