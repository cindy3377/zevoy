package handlers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// RetrieveHandler serves images and supports resizing based on query parameters (width and height)
func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Get("Authorization")

	// Ensure the token is provided
	if userToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the requested file name
	fileName := filepath.Base(r.URL.Path)
	if strings.Contains(fileName, "..") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	// Build the file path for the specific user
	filePath := filepath.Join("uploads", userToken, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Invalid image format", http.StatusInternalServerError)
		return
	}

	// Log the format of the image
	fmt.Println("Image format:", format)

	// Read query parameters for width and height
	widthStr := r.URL.Query().Get("width")
	heightStr := r.URL.Query().Get("height")

	if widthStr != "" && heightStr != "" {
		// Parse width and height from the query parameters
		width, err := strconv.Atoi(widthStr)
		if err != nil {
			http.Error(w, "Invalid width parameter", http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(heightStr)
		if err != nil {
			http.Error(w, "Invalid height parameter", http.StatusBadRequest)
			return
		}

		// Resize the image while maintaining aspect ratio
		resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

		// Serve the resized image
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", format))

		// Encode the resized image and send it as a response
		switch format {
		case "jpeg":
			jpeg.Encode(w, resizedImage, nil)
		case "png":
			png.Encode(w, resizedImage)
		}
	} else {
		// If no width/height specified, serve the original image
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", format))

		// Reset the file pointer and serve the original image
		file.Seek(0, 0)
		switch format {
		case "jpeg":
			jpeg.Encode(w, img, nil)
		case "png":
			png.Encode(w, img)
		}
	}
}
