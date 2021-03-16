package handlers

import (
	"assignment-2/structs"
	"assignment-2/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const covidCasesAPI = "https://covid-api.mmediagroup.fr/v1"
const dataAPI = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2"

// CasesPerCountry gets amount of cases for a country
func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Reached country endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	date, ok := vars["begin_date-end_date"]
	casesSince := 0
	recoveredSince := 0
	continent := ""
	country := ""
	scope := "total"
	populationPercentage := float32(0.0)
	if ok {
		splitdate := strings.Split(date, "-")
		if len(splitdate) < 6 {
			// Handles string error
			err := errors.New("Error in date")
			log.Printf("Error, %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Splits the dates into two different strings
		beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
		endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
		// Check if the dates are valid
		if !utils.ValidDate(beginDate) || !utils.ValidDate(endDate) {
			log.Println("Error in date")
			http.Error(w, "Error in date", http.StatusBadRequest)
			return
		}
		bodyCases, err := utils.GetBody(covidCasesAPI+"/history?country="+countryName+"&status=Confirmed", w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var informationCases structs.CovidAPIResponse
		if err = json.Unmarshal([]byte(string(bodyCases)), &informationCases); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Check if the country requested exists
		if informationCases.All.Country == "" {
			http.Error(w, "Country does not exists, remember to capitalize the first letter of the country name", http.StatusBadRequest)
			return
		}
		bodyRecovered, err := utils.GetBody(covidCasesAPI+"/history?country="+countryName+"&status=Recovered", w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var informationRecovered structs.CovidAPIResponse
		if err = json.Unmarshal([]byte(string(bodyRecovered)), &informationRecovered); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		casesSince = informationCases.All.Dates[endDate] - informationCases.All.Dates[beginDate]
		recoveredSince = informationRecovered.All.Dates[endDate] - informationRecovered.All.Dates[beginDate]
		continent = informationRecovered.All.Continent
		country = informationRecovered.All.Country
		scope = beginDate + "-" + endDate
		// The expression is complicated to make sure there are only two decimals.
		populationPercentage = float32(int((float32(casesSince)/float32(informationRecovered.All.Population)*100)*100)) / 10
	} else {
		body, err := utils.GetBody(covidCasesAPI+"/cases?country="+countryName, w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var information structs.CovidAPIResponse
		if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if information.All.Country == "" {
			http.Error(w, "Country does not exists, remember to capitalize the first letter of the country name", http.StatusBadRequest)
			return
		}
		casesSince = information.All.Confirmed
		recoveredSince = information.All.Recovered
		continent = information.All.Continent
		country = information.All.Country
		populationPercentage = float32(int((float32(casesSince)/float32(information.All.Population)*100)*100)) / 1
	}
	var returnJSON structs.ReturnConfirmedCases
	returnJSON.Confirmed = casesSince
	returnJSON.Recovered = recoveredSince
	returnJSON.Continent = continent
	returnJSON.Country = country
	returnJSON.Scope = scope
	returnJSON.PopulationPercentage = float32(populationPercentage)
	b, err := json.Marshal(returnJSON)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(b))

}
