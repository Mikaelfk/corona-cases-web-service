package utils

import (
	"assignment-2/structs"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

// GetBody returns the response body of a get request
func GetBody(request string, w http.ResponseWriter) ([]byte, error) {
	resp, err := http.Get(request)

	if err != nil {
		// Handles retrieval errors
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return nil, err
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

// SplitDate splits one string with two dates into two strings
func SplitDate(date string) (string, string, error) {
	splitdate := strings.Split(date, "-")
	if len(splitdate) < 6 {
		// Handles string error
		err := errors.New("date is not valid")
		return "", "", err
	}
	// Splits the dates into two different strings
	beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
	endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
	// Check if the dates are valid
	if !ValidDate(beginDate) || !ValidDate(endDate) {
		err := errors.New("date is not valid")
		return "", "", err
	}
	return beginDate, endDate, nil
}

func Remove(slice []structs.WebhookRegistration, s int) []structs.WebhookRegistration {
	return append(slice[:s], slice[s+1:]...)
}
