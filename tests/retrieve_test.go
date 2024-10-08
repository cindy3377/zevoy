package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"zevoy-project/handlers"
)

func TestRetrieveHandler(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a sample image file in the temporary directory
	userToken := "test-token"
	fileName := "receipt.jpg"
	filePath := filepath.Join(tempDir, userToken, fileName)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	file.Close()

	tests := []struct {
		name           string
		url            string
		userToken      string
		expectedStatus int
	}{
		{
			name:           "Missing user token",
			url:            "/uploads/receipt.jpg",
			userToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid file name",
			url:            "/uploads/../test-image.png",
			userToken:      "test-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "File not found",
			url:            "/uploads/nonexistent.png",
			userToken:      "test-token",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.userToken != "" {
				req.Header.Set("Authorization", tt.userToken)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.RetrieveHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
