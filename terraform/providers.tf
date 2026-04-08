terraform {
  required_providers {
    docker = { source = "kreuzwerker/docker", version = "~> 3.0.1" }
    aws    = { source = "hashicorp/aws", version = "~> 5.0" }
  }
}

provider "aws" {
  access_key = "test"
  secret_key = "test"
  region     = "us-east-1"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    s3     = local.s3_endpoint
    lambda = "http://localhost:4566" # Tambahkan ini!
    iam    = "http://localhost:4566" # Tambahkan ini!
  }
}

provider "docker" {}


