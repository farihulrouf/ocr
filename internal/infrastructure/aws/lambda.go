package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// ExportPayload adalah struktur data yang dikirim ke Lambda.
// Pastikan tag JSON ini sinkron dengan struct di main.go export-worker.
type ExportPayload struct {
	ExportLogID string `json:"export_log_id"`
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	Status      string `json:"status"`
}

// InvokeExportLambda memicu fungsi Lambda secara asynchronous (Fire and Forget)
func InvokeExportLambda(ctx context.Context, payload ExportPayload) error {
	// 1. Ambil Endpoint dari environment variable .env
	awsEndpoint := os.Getenv("S3_ENDPOINT")
	if awsEndpoint == "" {
		// Fallback ke localhost jika tidak didefinisikan
		awsEndpoint = "http://localhost:4566"
	}

	// 2. Load Konfigurasi SDK dengan Kredensial Statis untuk Localstack
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		// Pakai kredensial dummy "test" agar SDK tidak mencoba mencari IAM Role asli
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		return fmt.Errorf("gagal load SDK config: %v", err)
	}

	// 3. Inisialisasi Lambda Client dengan BaseEndpoint ke Localstack
	client := lambda.NewFromConfig(cfg, func(o *lambda.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	// 4. Serialisasi Payload ke JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload ke JSON: %v", err)
	}

	// 5. Invoke Lambda secara Async (InvocationTypeEvent)
	log.Printf("[AWS-SDK] Memanggil Lambda 'seido-export-service' di %s", awsEndpoint)
	log.Printf("[AWS-SDK] Payload: %s", string(payloadBytes))

	input := &lambda.InvokeInput{
		FunctionName:   aws.String("seido-export-service"),
		InvocationType: types.InvocationTypeEvent, // Async: API tidak menunggu Lambda selesai
		Payload:        payloadBytes,
	}

	output, err := client.Invoke(ctx, input)
	if err != nil {
		return fmt.Errorf("gagal memicu lambda: %v", err)
	}

	// Status 202 berarti Event berhasil masuk ke antrian Lambda
	log.Printf("[AWS-SDK] Lambda dipicu! Response Status: %d", output.StatusCode)

	return nil
}
