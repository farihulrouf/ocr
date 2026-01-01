                         ┌────────────────────┐
                         │     USER LOGIN      │
                         └─────────┬──────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │  AUTHENTICATION (JWT + RT)   │
                    └──────┬─────────┬────────────┘
                           │         │
                           │         ▼
                           │   Refresh Token
                           │
                           ▼
         ┌────────────────────────────────────────────┐
         │          LOAD TENANT ENVIRONMENT            │
         └───────────────┬──────────────┬─────────────┘
                         │              │
                         │              ▼
                         │       Get Settings
                         │
                         ▼
                Get Tenant Info
                         │
                         ▼
               Get Subscription Plan
                         │
                         ▼
                  Get Usage Stats
                         │
                         ▼
           ┌────────────────────────────────┐
           │  ORGANIZATION MANAGEMENT FLOW   │
           └──────────┬──────────┬──────────┘
                      │          │
                      ▼          ▼
            Manage Departments   Manage Users
                      │          │
                      ▼          ▼
               Set Approvers     Set Roles
                      │          │
                      └──────┬───┘
                             ▼
             ┌─────────────────────────────────┐
             │       RECEIPT AI-OCR FLOW        │
             └──────────────┬───────────────────┘
                            │
                            ▼
                 STEP 1: Upload Receipt
                            │
                            ▼
                Backend → OCR Provider
                            │
                            ▼
                STEP 2: Webhook OCR Result
                            │
                            ▼
                Update Receipt in Database
                            │
                            ▼
                 STEP 3: User Review/Confirm
                            │
      ┌──────────────┬──────────────┬──────────────┬───────────────┐
      ▼              ▼              ▼               ▼
   Edit Data     Add Items    Change Category    Delete Receipt
      └──────────────┬──────────────┬──────────────┬───────────────┘
                     ▼
           ┌───────────────────────────────────┐
           │     EXPENSE REPORT CREATION        │
           └───────────────────┬────────────────┘
                               │
                               ▼
                     Create Expense Report
                               │
                               ▼
                    Add Receipts to Report
                               │
                               ▼
                            Submit
                               │
                               ▼
           ┌──────────────────────────────────────────┐
           │        APPROVAL WORKFLOW SYSTEM          │
           └─────────────────┬──────────┬─────────────┘
                             │          │
                             ▼          ▼
                       Approver Checks  Approver Reviews
                             │          │
                             ▼          ▼
                      Approve / Reject / Remand
                             │
                             ▼
                     Update Report Status
                             │
                             ▼
                 If Approved → Lock & Archive
                             │
                             ▼
         ┌────────────────────────────────────────┐
         │           AUDIT & COMPLIANCE            │
         └──────────────────┬──────────┬───────────┘
                            │          │
                            ▼          ▼
                    Audit Trails     Export PDF/CSV
                            │
                            ▼
                       Tax Summary
                            │
                            ▼
           ┌───────────────────────────────────────────┐
           │                SYSTEM ADMIN                │
           └─────────────┬───────────────┬────────────┘
                         │               │
                         ▼               ▼
                 Manage Plans       Manage Tenants
                         │               │
                         ▼               ▼
                 Maintenance Mode   OCR Monitoring
