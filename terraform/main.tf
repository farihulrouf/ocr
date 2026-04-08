# Resource untuk Postgres
resource "docker_container" "postgres_db" {
  name  = "postgres_db"
  image = "postgres:15"
  env   = [
    "POSTGRES_USER=postgres",
    "POSTGRES_PASSWORD=postgres",
    "POSTGRES_DB=ocr"
  ]
  networks_advanced {
    name = data.docker_network.ocr_network.name
  }
  # Port mapping supaya kamu bisa akses dari psql laptop
  ports {
    internal = 5432
    external = 5432
  }
}

# Resource untuk Redis
resource "docker_container" "redis_db" {
  name  = "redis_db"
  image = "redis:7"
  networks_advanced {
    name = data.docker_network.ocr_network.name
  }
  ports {
    internal = 6379
    external = 6379
  }
}