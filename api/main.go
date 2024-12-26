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
)

const (
    IMAGE_DIR = "content/images/posts"
    MAX_FILE_SIZE = 10 << 20 // 10MB
    MAX_FILES = 10
)

// Allowed MIME types and their extensions
var allowedTypes = map[string][]string{
    "image/jpeg": {".jpg", ".jpeg"},
    "image/png":  {".png"},
    "image/gif":  {".gif"},
    "image/webp": {".webp"},
}

type UploadResponse struct {
    Success bool     `json:"success"`
    Images  []string `json:"images"`
    Error   string  `json:"error,omitempty"`
}

func main() {
    // Ensure upload key is set
    uploadKey := os.Getenv("UPLOAD_KEY")
    if uploadKey == "" {
        panic("UPLOAD_KEY environment variable not set")
    }

    // Ensure image directory exists
    if err := os.MkdirAll(IMAGE_DIR, 0755); err != nil {
        panic(fmt.Sprintf("Failed to create image directory: %v", err))
    }

    http.HandleFunc("/api/upload", handleUpload)
    
    port := os.Getenv("UPLOAD_PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Printf("Upload service running on :%s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        panic(err)
    }
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Check method
    if r.Method != http.MethodPost {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Method not allowed",
        })
        return
    }

    // Check API key
    if r.Header.Get("X-API-Key") != os.Getenv("UPLOAD_KEY") {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Unauthorized",
        })
        return
    }

    // Parse multipart form
    if err := r.ParseMultipartForm(MAX_FILE_SIZE * MAX_FILES); err != nil {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Failed to parse form",
        })
        return
    }

    var savedImages []string

    // Process each file
    for _, fileHeaders := range r.MultipartForm.File {
        for _, fileHeader := range fileHeaders {
            // Open uploaded file
            file, err := fileHeader.Open()
            if err != nil {
                continue
            }
            defer file.Close()

            // Generate unique filename
            timestamp := time.Now().Format("2006-01-02_150405.000")
            ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
            filename := fmt.Sprintf("IMG_%s%s", timestamp, ext)

            // Create destination file
            imgPath := filepath.Join(IMAGE_DIR, filename)
            dst, err := os.Create(imgPath)
            if err != nil {
                continue
            }
            defer dst.Close()

            // Copy file content
            if _, err := io.Copy(dst, file); err != nil {
                continue
            }

            savedImages = append(savedImages, fmt.Sprintf("/images/posts/%s", filename))
        }
    }

    // Return response
    json.NewEncoder(w).Encode(UploadResponse{
        Success: len(savedImages) > 0,
        Images:  savedImages,
    })
}