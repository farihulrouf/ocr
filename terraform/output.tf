# ---- OUTPUTS: Ringkasan Informasi Infrastruktur ----

output "api_endpoint" {
  description = "URL utama untuk mengakses API"
  value       = "http://localhost:${docker_container.seido_api.ports[0].external}"
}

output "s3_bucket_name" {
  description = "Nama bucket S3 yang berhasil dibuat"
  value       = aws_s3_bucket.ocr_bucket.id
}

output "active_workers" {
  description = "Jumlah worker yang sedang berjalan"
  value       = length(docker_container.seido_worker)
}

output "docker_network" {
  description = "Network yang digunakan oleh kontainer"
  value       = data.docker_network.ocr_network.name
}