terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# --- LOGIC: Membaca .env secara otomatis ---
locals {
  env_file = file("../.env")
  env_lines = [
    for line in split("\n", local.env_file) : 
    line if length(trimspace(line)) > 0 && !startswith(line, "#")
  ]
  
  # Otomatis ambil S3_ENDPOINT dari .env untuk provider AWS
  # Jika di .env: http://localstack:4566, terraform ganti ke localhost agar bisa akses dari luar
  s3_raw      = [for l in local.env_lines : l if startswith(l, "S3_ENDPOINT=")][0]
  s3_val      = split("=", local.s3_raw)[1]
  s3_endpoint = replace(local.s3_val, "localstack", "localhost")
}

# 1. Provider AWS (Untuk Setup S3 di LocalStack)
provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = "us-east-1"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3 = local.s3_endpoint 
  }
}

provider "docker" {}

# 2. Resource: Membuat Bucket S3 secara Otomatis
resource "aws_s3_bucket" "ocr_bucket" {
  bucket = "ocr-bucket"
}

# 3. Data: Mengambil Network Docker yang sudah ada
data "docker_network" "ocr_network" {
  name = "ocr_seido-network"
}

# 4. Resource: Build Image Docker (Multi-Stage)
# Image ini berisi seido-api dan seido-worker
resource "docker_image" "seido_app_image" {
  name = "seido-app:latest"
  build {
    context    = ".."
    dockerfile = "Dockerfile"
    remove     = true
  }
}

# 5. Resource: Container SEIDO API
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

  env = local.env_lines

  restart = "always"

  # API baru jalan setelah Bucket S3 siap
  depends_on = [aws_s3_bucket.ocr_bucket]
}

# 6. Resource: Container SEIDO WORKER (Untuk OCR Processing)
resource "docker_container" "seido_worker" {
  name  = "seido-worker-instance"
  image = docker_image.seido_app_image.image_id
  
  # Jalankan binary worker (override default CMD)
  command = ["./seido-worker"]

  networks_advanced {
    name = data.docker_network.ocr_network.name
  }

  env = local.env_lines

  restart = "always"

  # Worker juga butuh S3 Bucket siap
  depends_on = [aws_s3_bucket.ocr_bucket]
}