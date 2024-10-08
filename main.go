package main

import (
	"fmt" // Update the import path to the correct location
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"zevoy-project/handlers"
)

func main() {
    // Create uploads directory if it doesn't exist
    if _, err := os.Stat("uploads"); os.IsNotExist(err) {
        os.Mkdir("uploads", 0755)
    }

    // Register the upload handler
    http.HandleFunc("/upload", handlers.UploadHandler)

   // Serve uploaded receipt files using RetrieveHandler
    http.HandleFunc("/receipts/", handlers.RetrieveHandler)

    // Start the server
    fmt.Println("Server is starting at :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

