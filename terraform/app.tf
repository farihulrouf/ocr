# =====================================================
# 1. RESOURCE: BUILD IMAGE DOCKER
# =====================================================
resource "docker_image" "seido_app_image" {
  name = "seido-app:latest"
  build {
    context    = ".."
    dockerfile = "Dockerfile"
    remove     = true
  }
}

# =====================================================
# 2. RESOURCE: CONTAINER SEIDO API
# =====================================================
resource "docker_container" "seido_api" {
  name  = "seido-api-instance"
  image = docker_image.seido_app_image.image_id
  
  ports {
    internal = 8080
    external = 8080
  }

  networks_advanced {
    name = data.docker_network.ocr_network.name
  }

  # KUNCI PERBAIKAN:
  # Kita paksa API menggunakan koordinat LocalStack dan nama fungsi yang BENAR
  env = concat(local.env_lines, [
    "S3_ENDPOINT=http://localstack:4566",
    "AWS_ENDPOINT=http://localstack:4566",
    "AWS_LAMBDA_ENDPOINT=http://localstack:4566",
    "AWS_REGION=us-east-1",
    # Nama ini HARUS sama dengan function_name di lambda.tf
    "AWS_LAMBDA_FUNCTION_NAME=seido-export-service" 
  ])

  restart = "always"

  # API baru jalan setelah Bucket S3 dan Lambda siap
  depends_on = [aws_s3_bucket.ocr_bucket, aws_lambda_function.export_worker]
}

# =====================================================
# 3. RESOURCE: CONTAINER SEIDO WORKER (SCALING)
# =====================================================
resource "docker_container" "seido_worker" {
  count = 2  

  name  = "seido-worker-${count.index}"
  image = docker_image.seido_app_image.image_id
  
  # Override CMD untuk menjalankan binary worker
  command = ["./seido-worker"]

  networks_advanced {
    name = data.docker_network.ocr_network.name
  }

  # Worker juga butuh koneksi yang sama ke database dan localstack
  env = concat(local.env_lines, [
    "S3_ENDPOINT=http://localstack:4566",
    "AWS_ENDPOINT=http://localstack:4566",
    "AWS_REGION=us-east-1"
  ])

  restart = "always"

  depends_on = [aws_s3_bucket.ocr_bucket]
}