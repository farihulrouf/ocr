# Zip file hasil build go (kamu harus build manual dulu atau pakai resource)
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "../cmd/export-worker/main" # Pastikan ini path binary hasil build
  output_path = "export_worker.zip"
}

resource "aws_lambda_function" "export_worker" {
  filename      = data.archive_file.lambda_zip.output_path
  function_name = "export-expense-lambda"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "main" # Karena Go menggunakan binary bernama main
  runtime       = "go1.x" # Atau "provided.al2" jika menggunakan AL2

  # Teruskan environment variable agar Lambda bisa konek ke Postgres
  environment {
    variables = {
      DB_HOST = "postgres_db" # Nama container di docker-compose
      DB_PORT = "5432"
      S3_ENDPOINT = "http://localstack:4566"
    }
  }

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
}

# Dummy Role untuk Localstack
resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}