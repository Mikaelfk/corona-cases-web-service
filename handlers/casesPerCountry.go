package handlers

import (
	"assignment-2/structs"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
		// Splits the date variable into two date variables
		beginDate, endDate, err := utils.SplitDate(date)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		bodyCases, err := utils.GetBody(utils.CovidCasesAPI+"/history?country="+countryName+"&status=Confirmed", w)
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
		bodyRecovered, err := utils.GetBody(utils.CovidCasesAPI+"/history?country="+countryName+"&status=Recovered", w)
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
		populationPercentage = float32(int((float32(casesSince)/float32(informationRecovered.All.Population)*100)*100)) / 100
	} else {
		body, err := utils.GetBody(utils.CovidCasesAPI+"/cases?country="+countryName, w)
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
		populationPercentage = float32(int((float32(casesSince)/float32(information.All.Population)*100)*100)) / 100
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
	fmt.Fprint(w, string(b))
}
