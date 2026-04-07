package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// ExportPayload adalah struktur data yang dikirim ke Lambda
// Harus sinkron dengan struct ExportEvent di cmd/export-worker/main.go
type ExportPayload struct {
	ExportLogID string `json:"export_log_id"`
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	Status      string `json:"status"`
}

// InvokeExportLambda memicu fungsi Lambda secara asynchronous di LocalStack/AWS
func InvokeExportLambda(ctx context.Context, payload ExportPayload) error {
	// 1. Load Konfigurasi SDK
	// Kita set region default us-east-1 untuk LocalStack
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return fmt.Errorf("gagal load SDK config: %v", err)
	}

	// 2. Inisialisasi Lambda Client
	// Point penting: BaseEndpoint diarahkan ke LocalStack (port 4566)
	client := lambda.NewFromConfig(cfg, func(o *lambda.Options) {
		// Ganti localhost dengan host.docker.internal jika API jalan di dalam Docker
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})

	// 3. Serialisasi Payload ke JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload ke JSON: %v", err)
	}

	// 4. Invoke Lambda (Async menggunakan InvocationTypeEvent)
	log.Printf("[AWS-SDK] Memicu Lambda 'seido-export-service' untuk LogID: %s", payload.ExportLogID)

	input := &lambda.InvokeInput{
		FunctionName:   aws.String("seido-export-service"), // Harus sama dengan nama saat 'awslocal lambda create'
		InvocationType: types.InvocationTypeEvent,          // Event = Tidak menunggu response (Fire and Forget)
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
