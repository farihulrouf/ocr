       [ User / Frontend ]
               |
               v
         [ Backend API ]
               |  Create Receipt + OCRJob
               v
         [ PostgreSQL ]
               |
               | Push receiptID
               v
          [ Redis Queue ]
               |
               | Worker ambil job
               v
           [ OCR Worker ]
               |
               | Download image
               v
           [ MinIO / S3 ]
               |
               | OCR & AI Agent
               v
           [ AI Agent ]
               |
               | Update DB
               v
         [ PostgreSQL ]