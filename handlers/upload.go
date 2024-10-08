package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadHandler handles image uploads and adds a UserToken-based permissions check
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Get("Authorization") // Simple token for user identification

	// If no token is provided, reject the request
	if userToken == "" {
		http.Error(w, "Unauthorized: Missing user token", http.StatusUnauthorized)
		return
	}

	// Set the maximum size for uploaded files
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Limit to 10MB

	// Ensure user-specific directory exists
	userDir := filepath.Join("uploads", userToken)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		err := os.MkdirAll(userDir, 0755) // Create user-specific directory
		if err != nil {
			http.Error(w, "Failed to create user directory", http.StatusInternalServerError)
			return
		}
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limit to 10MB
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	// Save the file to the user's directory
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Content-Type validation
	contentType := r.Header.Get("Content-Type")
	if contentType != "multipart/form-data" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	// Reset the file pointer to the start, so we can save it
	file.Seek(0, 0)

	// Sanitize the file name to prevent directory traversal attacks
	fileName := filepath.Base(header.Filename)

	// Ensure the file name does not contain ".." (preventing traversal)
	if strings.Contains(fileName, "..") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	// Create the full file path (user-specific directory + sanitized file name)
	filePath := filepath.Join(userDir, fileName)

	// Save the image to the local filesystem
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	// Respond with the file info, including the user token in the URL path
	response := map[string]string{
		"filename": header.Filename,
		"url":      fmt.Sprintf("/receipts/%s/%s", userToken, fileName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
