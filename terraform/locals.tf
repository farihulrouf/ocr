locals {
  env_file = file("../.env")
  env_lines = [
    for line in split("\n", local.env_file) : 
    line if length(trimspace(line)) > 0 && !startswith(line, "#")
  ]
  
  s3_raw      = [for l in local.env_lines : l if startswith(l, "S3_ENDPOINT=")][0]
  s3_val      = split("=", local.s3_raw)[1]
  s3_endpoint = replace(local.s3_val, "localstack", "localhost")
}