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
    "crypto/subtle"
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

type Config struct {
    Params struct {
        API struct {
            UploadKey string `yaml:"uploadKey"`
        } `yaml:"api"`
    } `yaml:"params"`
}

type UploadResponse struct {
    Success bool     `json:"success"`
    Images  []string `json:"images"`
    Error   string  `json:"error,omitempty"`
}

var config Config

func init() {
    // Read config file using safe path
    configPath := filepath.Clean("config.toml")
    if !strings.HasSuffix(configPath, "config.toml") {
        panic("Invalid config path")
    }

    configData, err := os.ReadFile(configPath)
    if err != nil {
        panic(fmt.Sprintf("Failed to read config file: %v", err))
    }

    // Parse TOML config
    if err := yaml.Unmarshal(configData, &config); err != nil {
        panic(fmt.Sprintf("Failed to parse config: %v", err))
    }

    // Verify API key exists and has minimum length
    if len(config.Params.API.UploadKey) < 32 {
        panic("API upload key must be at least 32 characters")
    }

    // Create image directory with secure permissions
    if err := os.MkdirAll(IMAGE_DIR, 0755); err != nil {
        panic(fmt.Sprintf("Failed to create image directory: %v", err))
    }
}

func main() {
    // Use a mux for better route handling
    mux := http.NewServeMux()
    mux.HandleFunc("/api/upload", handleMultiUpload)

    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 20 * time.Second,
        ReadTimeout:       5 * time.Minute,
        WriteTimeout:      5 * time.Minute,
        MaxHeaderBytes:    1 << 20, // 1MB
    }

    fmt.Println("Server running on :8080")
    server.ListenAndServe()
}

func verifyAPIKey(key string) bool {
    // Constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare([]byte(key), []byte(config.Params.API.UploadKey)) == 1
}

func isAllowedFile(filename string, header *mime.FileHeader) bool {
    // Check file extension
    ext := strings.ToLower(filepath.Ext(filename))
    
    // Get content type
    contentType := header.Header.Get("Content-Type")
    
    // Verify content type and extension match
    allowedExts, ok := allowedTypes[contentType]
    if !ok {
        return false
    }
    
    for _, allowedExt := range allowedExts {
        if ext == allowedExt {
            return true
        }
    }
    return false
}

func sanitizeFilename(filename string) string {
    // Keep only alphanumeric characters and some safe symbols
    safe := strings.Map(func(r rune) rune {
        switch {
        case r >= 'a' && r <= 'z',
             r >= 'A' && r <= 'Z',
             r >= '0' && r <= '9',
             r == '-' || r == '_' || r == '.':
            return r
        default:
            return '_'
        }
    }, filename)
    
    // Ensure no directory traversal
    return filepath.Base(safe)
}

func handleMultiUpload(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Check API key
    if !verifyAPIKey(r.Header.Get("X-API-Key")) {
        time.Sleep(time.Second) // Prevent timing attacks
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Unauthorized",
        })
        return
    }

    if r.Method != "POST" {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Method not allowed",
        })
        return
    }

    // Limit request size
    r.Body = http.MaxBytesReader(w, r.Body, MAX_FILE_SIZE*MAX_FILES)

    // Parse multipart form
    if err := r.ParseMultipartForm(MAX_FILE_SIZE * 2); err != nil {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   "Request too large or invalid",
        })
        return
    }
    defer r.MultipartForm.RemoveAll()

    files := r.MultipartForm.File["images"]
    if len(files) == 0 || len(files) > MAX_FILES {
        json.NewEncoder(w).Encode(UploadResponse{
            Success: false,
            Error:   fmt.Sprintf("Must provide 1-%d images", MAX_FILES),
        })
        return
    }

    savedImages := []string{}

    for _, fileHeader := range files {
        // Check file size
        if fileHeader.Size > MAX_FILE_SIZE {
            continue
        }

        // Validate file type
        if !isAllowedFile(fileHeader.Filename, fileHeader) {
            continue
        }

        // Open uploaded file
        file, err := fileHeader.Open()
        if err != nil {
            continue
        }
        defer file.Close()

        // Generate unique filename with timestamp and sanitize
        timestamp := time.Now().Format("2006-01-02_150405.000")
        ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
        filename := sanitizeFilename(fmt.Sprintf("IMG_%s%s", timestamp, ext))

        // Create destination path
        imgPath := filepath.Join(IMAGE_DIR, filename)
        
        // Verify the resolved path is within IMAGE_DIR
        if !strings.HasPrefix(filepath.Clean(imgPath), filepath.Clean(IMAGE_DIR)) {
            continue
        }

        // Create destination file with restricted permissions
        dst, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
        if err != nil {
            continue
        }
        defer dst.Close()

        // Copy file content with limited reader
        if _, err := io.Copy(dst, io.LimitReader(file, MAX_FILE_SIZE)); err != nil {
            os.Remove(imgPath) // Clean up on error
            continue
        }

        savedImages = append(savedImages, fmt.Sprintf("/images/posts/%s", filename))
    }

    // Return response
    json.NewEncoder(w).Encode(UploadResponse{
        Success: len(savedImages) > 0,
        Images:  savedImages,
    })
}