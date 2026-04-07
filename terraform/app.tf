# 1. Resource: Build Image Docker
# Terraform akan build image dari Dockerfile di folder utama (..)
resource "docker_image" "seido_app_image" {
  name = "seido-app:latest"
  build {
    context    = ".."
    dockerfile = "Dockerfile"
    remove     = true
  }
}

# 2. Resource: Container SEIDO API (Hanya 1 Instance)
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

  # API baru jalan setelah Bucket S3 di Localstack siap
  depends_on = [aws_s3_bucket.ocr_bucket]
}

# 3. Resource: Container SEIDO WORKER (Scaling Horizontal)
resource "docker_container" "seido_worker" {
  # Kita set 2 worker saja dulu biar aman
  count = 2  

  # Nama akan otomatis jadi seido-worker-0 dan seido-worker-1
  name  = "seido-worker-${count.index}"
  
  image = docker_image.seido_app_image.image_id
  
  # Override CMD untuk menjalankan binary worker
  command = ["./seido-worker"]

  networks_advanced {
    name = data.docker_network.ocr_network.name
  }

  env = local.env_lines

  restart = "always"

  # Worker juga butuh S3 Bucket siap
  depends_on = [aws_s3_bucket.ocr_bucket]
}