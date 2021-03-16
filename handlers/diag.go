package handlers

import (
	"assignment-2/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var StartTime time.Time

var RegisteredWebhooks int

const UrlBase = "/corona/"
const Version = "v1"

// Diag shows a diagnostics interface
func Diag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Diag Endpoint")
	var mMediaGroupStatusCode int
	var covidTrackerStatusCode int
	// Does a basic request to the M Media Group API.
	respMMedia, err := http.Get(covidCasesAPI + "/cases?country=Norway")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the M Media Group API, %v", err)
		mMediaGroupStatusCode = 500
	} else {
		mMediaGroupStatusCode = respMMedia.StatusCode
		defer respMMedia.Body.Close()
	}
	// Does a basic request to the covid tracker API.
	respCovidTracker, err := http.Get(dataAPI + "/stringency/actions/NOR/2020-12-31")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the covid tracker api, %v", err)
		covidTrackerStatusCode = 500
	} else {
		covidTrackerStatusCode = respCovidTracker.StatusCode
		defer respCovidTracker.Body.Close()
	}
	var returnJSON structs.ReturnDiag
	returnJSON.MMediaGroupApi = strconv.Itoa(mMediaGroupStatusCode)
	returnJSON.CovidTrackerAPI = strconv.Itoa(covidTrackerStatusCode)
	returnJSON.Registered = RegisteredWebhooks
	returnJSON.Version = "v1"
	returnJSON.Uptime = int(time.Since(StartTime) / time.Second)
	b, err := json.Marshal(returnJSON)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(b))
}
