package dynamodb

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func SaveToLocalDynamo(result interface{}) error {
	// 1. Inisialisasi AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000"), // DynamoDB Local
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := dynamodb.New(sess)

	// 2. Buat item
	item := map[string]interface{}{
		"receipt_id": fmt.Sprintf("r-%d", time.Now().UnixNano()),
		"ocr_data":   result,
		"created_at": time.Now().Format(time.RFC3339),
	}

	// 3. Marshal ke AttributeValue
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %v", err)
	}

	// 4. Put item ke DynamoDB Local
	input := &dynamodb.PutItemInput{
		TableName: aws.String("ocr_receipts"), // pastikan table sudah ada
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	fmt.Println("[INFO] OCR result saved to local DynamoDB successfully")
	return nil
}

// ========================
// Fungsi ScanDynamoDB untuk verifikasi
// ========================
func ScanDynamoDB() error {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("ocr_receipts"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return fmt.Errorf("failed to scan table: %v", err)
	}

	fmt.Printf("[INFO] Total items in table: %d\n", *result.Count)
	for i, item := range result.Items {
		fmt.Printf("Item %d: %v\n", i+1, item)
	}

	return nil
}
