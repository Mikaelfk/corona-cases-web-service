package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// GetPort gets a port from the environmental variable PORT, if this is empty, set port to 8181
func GetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	return port
}

// Returns the response body of a get request
func GetBody(request string, w http.ResponseWriter) ([]byte, error) {
	resp, err := http.Get(request)

	if err != nil {
		// Handles retrieval errors
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Handles user input error
		http.Error(w, "Error: Status code not 2xx", resp.StatusCode)
		return nil, errors.New("Status code not 2xx")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return body, nil
}

// ValidDate checks if a date is valid
func ValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}
