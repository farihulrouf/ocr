package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os" // Tambahkan ini untuk membaca environment variables

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// ExportPayload adalah struktur data yang dikirim ke Lambda
type ExportPayload struct {
	ExportLogID string `json:"export_log_id"`
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	Status      string `json:"status"`
}

// InvokeExportLambda memicu fungsi Lambda secara asynchronous
func InvokeExportLambda(ctx context.Context, payload ExportPayload) error {
	// 1. Ambil Endpoint dari Environment Variable (Disuntikkan oleh Terraform)
	// Jika kosong (saat dev lokal), gunakan localhost sebagai fallback.
	awsEndpoint := os.Getenv("S3_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localhost:4566"
	}

	// 2. Load Konfigurasi SDK
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return fmt.Errorf("gagal load SDK config: %v", err)
	}

	// 3. Inisialisasi Lambda Client dengan Endpoint Dinamis
	client := lambda.NewFromConfig(cfg, func(o *lambda.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	// 4. Serialisasi Payload ke JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload ke JSON: %v", err)
	}

	// 5. Invoke Lambda (Async menggunakan InvocationTypeEvent)
	log.Printf("[AWS-SDK] Memicu Lambda 'seido-export-service' melalui %s untuk LogID: %s", awsEndpoint, payload.ExportLogID)

	input := &lambda.InvokeInput{
		FunctionName:   aws.String("seido-export-service"),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payloadBytes,
	}

	output, err := client.Invoke(ctx, input)
	if err != nil {
		return fmt.Errorf("gagal memicu lambda: %v", err)
	}

	// HTTP 202 Accepted berarti Lambda sukses menerima event tersebut
	log.Printf("[AWS-SDK] Lambda berhasil dipicu. Status Code: %d", output.StatusCode)

	return nil
}
