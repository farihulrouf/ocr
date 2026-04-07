resource "aws_s3_bucket" "ocr_bucket" {
  bucket = "ocr-bucket"
}

data "docker_network" "ocr_network" {
  name = "ocr_seido-network"
}