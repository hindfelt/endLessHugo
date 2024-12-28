package main

import (
	"api/auth"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	MAX_FILE_SIZE = 500 << 20 // 500MB to accommodate videos
	API_KEY       = "agnes"
)

// Map of allowed file extensions and their media types
var allowedTypes = map[string]string{
	".jpg":  "image",
	".jpeg": "image",
	".png":  "image",
	".gif":  "image",
	".webp": "image",
	".mov":  "video",
	".mp4":  "video",
}

type UploadResponse struct {
	Success  bool     `json:"success"`
	Files    []string `json:"files"`
	Markdown string   `json:"markdown"`
	Error    string   `json:"error,omitempty"`
}

func getMarkdownForFile(filename, title, mediaType string) string {
	switch mediaType {
	case "video":
		return fmt.Sprintf(`<video controls>
  <source src="%s" type="video/%s">
  Your browser does not support the video tag.
</video>`, filename, strings.TrimPrefix(filepath.Ext(filename), "."))
	default: // "image"
		return fmt.Sprintf("![%s](%s)", title, filename)
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize auth
	auth.Init()

	// Create a protected mux for Hugo content
	hugoMux := http.NewServeMux()
	hugoMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hugoURL := "http://localhost:1313" + r.URL.Path
		resp, err := http.Get(hugoURL)
		if err != nil {
			http.Error(w, "Failed to reach Hugo", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	})

	// Auth routes (unprotected)
	http.HandleFunc("/auth/google/callback", auth.HandleCallback)
	http.HandleFunc("/login", auth.HandleLogin)
	http.HandleFunc("/logout", auth.HandleLogout)

	// Everything else goes through auth
	http.Handle("/", auth.RequireAuth(hugoMux))

	fmt.Printf("Server running on :8000\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check API key
	if r.Header.Get("X-API-Key") != API_KEY {
		json.NewEncoder(w).Encode(UploadResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get text fields
	title := r.FormValue("title")
	description := r.FormValue("description")
	tags := r.FormValue("tags")

	fmt.Printf("Received: title=%s, description=%s, tags=%s\n", title, description, tags)

	// Parse multipart form with increased size limit
	if err := r.ParseMultipartForm(MAX_FILE_SIZE * 10); err != nil {
		json.NewEncoder(w).Encode(UploadResponse{
			Success: false,
			Error:   "Failed to parse form",
		})
		return
	}

	var savedFiles []string
	var markdownLinks []string

	// Process files
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

			// Check if file type is allowed
			mediaType, allowed := allowedTypes[ext]
			if !allowed {
				fmt.Printf("Skipping file with unsupported type: %s\n", ext)
				continue
			}

			// Open uploaded file
			file, err := fileHeader.Open()
			if err != nil {
				continue
			}
			defer file.Close()

			// Generate unique filename
			timestamp := time.Now().Format("2006-01-02_150405")
			var filename string

			switch mediaType {
			case "video":
				filename = fmt.Sprintf("VID_%s%s", timestamp, ext)
			default:
				filename = fmt.Sprintf("IMG_%s%s", timestamp, ext)
			}

			// Create destination path
			mediaDir := "/Users/mathin/Code/blogg/content/images/posts"
			filePath := filepath.Join(mediaDir, filename)

			// Create directory if needed
			if err := os.MkdirAll(mediaDir, 0755); err != nil {
				fmt.Printf("Failed to create directory: %v\n", err)
				continue
			}

			// Create destination file
			dst, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Failed to create file: %v\n", err)
				continue
			}
			defer dst.Close()

			// Copy file content
			if _, err := io.Copy(dst, file); err != nil {
				fmt.Printf("Failed to save file: %v\n", err)
				continue
			}

			filePath = fmt.Sprintf("/images/posts/%s", filename)
			savedFiles = append(savedFiles, filePath)
			markdownLinks = append(markdownLinks, getMarkdownForFile(filePath, title, mediaType))

			fmt.Printf("Saved %s: %s\n", mediaType, filePath)
		}
	}

	// Create post file if title is provided
	if title != "" {
		postsDir := "/Users/mathin/Code/blogg/content/posts"
		now := time.Now()
		postFilename := fmt.Sprintf("%s-%s.md",
			now.Format("2006-01-02"),
			strings.ToLower(strings.ReplaceAll(title, " ", "-")))

		// Create YAML frontmatter with exact spacing and alignment
		frontmatter := fmt.Sprintf(`---
date: %s
title: %s
type: post
---`, // No trailing newline here
			now.Format("2006-01-02 15:04:05-07:00"),
			title)

		// Add double newline after frontmatter
		frontmatter += "\n\n"

		if description != "" {
			frontmatter += description + "\n\n"
		}

		postContent := frontmatter + strings.Join(markdownLinks, "\n\n")

		if err := os.MkdirAll(postsDir, 0755); err != nil {
			fmt.Printf("Failed to create posts directory: %v\n", err)
		} else {
			postPath := filepath.Join(postsDir, postFilename)
			if err := os.WriteFile(postPath, []byte(postContent), 0644); err != nil {
				fmt.Printf("Failed to create post file: %v\n", err)
			} else {
				fmt.Printf("Created post file: %s\n", postPath)
			}
		}
	}

	// Return response
	json.NewEncoder(w).Encode(UploadResponse{
		Success:  len(savedFiles) > 0,
		Files:    savedFiles,
		Markdown: strings.Join(markdownLinks, "\n\n"),
	})
}
