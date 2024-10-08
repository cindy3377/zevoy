package tests

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"zevoy-project/handlers"
)

func TestUploadHandler(t *testing.T) {
	tests := []struct {
		name           string
		imagePath      string
		userToken      string
		expectedStatus int
	}{
		{
			name:           "Invalid User Token",
			imagePath:      "files/receipt.jpg",
			userToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Image Content Type",
			imagePath:      "files/receipt.jpg",
			userToken:      "Bearer valid_token",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.imagePath)
			if err != nil {
				t.Fatalf("Failed to open image file: %v", err)
			}
			defer file.Close()

			fileContent, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatalf("Failed to read image file: %v", err)
			}

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", filepath.Base(tt.imagePath))
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			part.Write(fileContent)
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Authorization", tt.userToken)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.UploadHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
