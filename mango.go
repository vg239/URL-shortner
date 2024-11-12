package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"time"
	"encoding/json"
)

// Create a DS for mapping between shortened version and normal URL
var urlMap = make(map[string]string)

// Generate a random short code
func generateShortCode() string {
    length := 7
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length+2)
    rand.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}

// Shorten URL handler
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
	// Parse original URL from request body
	var requestData struct{
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
	// Generate unique short code
	shortCode := generateShortCode()

	// Store the mapping
    urlMap[shortCode] = requestData.URL

    responseData := map[string]string{
        "short_url": fmt.Sprintf("http://localhost:8080/r/%s", shortCode),
    }
	
	

	// Send back the short URL as response
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)

}

// Redirect handler
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	// read the request, access your DS to find the full version
	// of URL, redirect to that URL
	shortCode := r.URL.Path[len("/r/"):]
    originalURL, exists := urlMap[shortCode]
    if !exists {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }
    http.Redirect(w, r, originalURL, http.StatusFound)
}

// Serve frontend index.html file
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"index.html")
}

func main() {

	// Route for serving the frontend page
	http.HandleFunc("/", indexHandler)

	// Route for the API to shorten URLs
	http.HandleFunc("/shorten", shortenURLHandler)

	// Route for handling redirects
	http.HandleFunc("/r/", redirectHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
