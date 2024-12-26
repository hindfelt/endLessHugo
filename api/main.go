package main

import (
    "encoding/json"
    "fmt"
    "io"
    "mime"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
    "gopkg.in/yaml.v2"
)

func init() {
    // Get upload key from environment
    uploadKey := os.Getenv("UPLOAD_KEY")
    if uploadKey == "" {
        panic("UPLOAD_KEY environment variable not set")
    }

    // Optional S3 configuration
    if awsBucket := os.Getenv("AWS_BUCKET_NAME"); awsBucket != "" {
        // Initialize S3 client
        sess := session.Must(session.NewSession(&aws.Config{
            Region: aws.String(os.Getenv("AWS_REGION")),
        }))
        s3Client = s3.New(sess)
    }
}