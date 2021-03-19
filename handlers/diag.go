package handlers

import (
	"assignment-2/structs"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// StartTime is the time the application is started
var StartTime time.Time

// RegisteredWebhooks is how many webooks are registered
var RegisteredWebhooks int

// URLBase is the start of the url
const URLBase = "/corona/"

// Version is the application version
const Version = "v1"

// Diag shows a diagnostics interface
func Diag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Diag Endpoint")
	var mMediaGroupStatusCode int
	var covidTrackerStatusCode int
	var countriesStatusCode int
	// Does a basic request to the M Media Group API.
	respMMedia, err := http.Get(utils.CovidCasesAPI + "/cases?country=Norway")
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
	respCovidTracker, err := http.Get(utils.DataAPI + "/stringency/actions/NOR/2020-12-31")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the covid tracker api, %v", err)
		covidTrackerStatusCode = 500
	} else {
		covidTrackerStatusCode = respCovidTracker.StatusCode
		defer respCovidTracker.Body.Close()
	}
	// Does a basic request to the countries information api.
	respCountries, err := http.Get("https://restcountries.eu")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the countries api, %v", err)
		countriesStatusCode = 500
	} else {
		countriesStatusCode = respCountries.StatusCode
		defer respCountries.Body.Close()
	}

	var returnJSON structs.ReturnDiag
	returnJSON.MMediaGroupApi = strconv.Itoa(mMediaGroupStatusCode)
	returnJSON.CovidTrackerAPI = strconv.Itoa(covidTrackerStatusCode)
	returnJSON.CountryAPI = strconv.Itoa(countriesStatusCode)
	returnJSON.Registered = RegisteredWebhooks
	returnJSON.Version = "v1"
	returnJSON.Uptime = int(time.Since(StartTime) / time.Second)
	b, err := json.Marshal(returnJSON)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(b))
}
