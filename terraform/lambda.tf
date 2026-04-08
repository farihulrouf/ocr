# =====================================================
# 1. OTOMATISASI BUILD & ZIP (GO)
# =====================================================
resource "null_resource" "prepare_lambda" {
  provisioner "local-exec" {
    command = <<EOT
      cd ${path.module}/../cmd/export-worker && \
      GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go && \
      zip -j ../../terraform/export_worker.zip main
    EOT
  }

  triggers = {
    main_go_hash = filebase64sha256("${path.module}/../cmd/export-worker/main.go")
  }
}

# =====================================================
# 2. DEFINISI FUNGSI LAMBDA
# =====================================================
resource "aws_lambda_function" "export_worker" {
  filename      = "${path.module}/export_worker.zip"
  function_name = "seido-export-service"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "main"
  runtime       = "go1.x" 

  environment {
    variables = {
      # GANTI 'postgres_db' jadi IP Gateway Docker kamu
      DB_HOST      = "172.17.0.1" 
      DB_PORT      = "5432"
      DB_USER      = "postgres"
      DB_PASSWORD  = "postgres"
      DB_NAME      = "ocr"
      # GANTI 'localstack' jadi IP Gateway Docker kamu
      S3_ENDPOINT  = "http://172.17.0.1:4566"
      S3_BUCKET    = aws_s3_bucket.ocr_bucket.id
    }
  }

  source_code_hash = null_resource.prepare_lambda.id
  depends_on       = [null_resource.prepare_lambda]
}

# =====================================================
# 3. IAM ROLE & LOGGING
# =====================================================
resource "aws_iam_role" "lambda_exec" {
  name = "seido_lambda_exec_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}