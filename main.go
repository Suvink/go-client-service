package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	// Get the URL from environment variable, default to localhost:8080/hello if not set
	url := os.Getenv("CHOREO_HELLO_SERVICE_MAIN_CONNECTION_SERVICEURL")
	if url == "" {
		url = "http://localhost:8080/hello"
	}

	// Make GET request to the configured URL
	resp, err := http.Get(url + "/hello")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
		return
	}

	// Create response with the template format
	response := Response{
		Type:    "main",
		Message: string(body),
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/greet", greetHandler)

	port := ":8081"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
