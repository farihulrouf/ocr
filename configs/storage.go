package configs

import (
    "log"
    "os"
    "strconv"
)

type MinioConfigStruct struct {
    Endpoint  string
    AccessKey string
    SecretKey string
    Bucket    string
    UseSSL    bool
}

var MinioConfig *MinioConfigStruct

func InitMinioConfig() {
    useSSL, err := strconv.ParseBool(os.Getenv("STORAGE_USE_SSL"))
    if err != nil {
        useSSL = false
    }

    MinioConfig = &MinioConfigStruct{
        Endpoint:  os.Getenv("STORAGE_ENDPOINT"),
        AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
        SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
        Bucket:    os.Getenv("STORAGE_BUCKET"),
        UseSSL:    useSSL,
    }

    if MinioConfig.Endpoint == "" || MinioConfig.Bucket == "" {
        log.Fatalln("MinIO config is incomplete")
    }
}