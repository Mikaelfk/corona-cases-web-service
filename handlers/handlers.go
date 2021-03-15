package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"assignment-2/structs"
	"assignment-2/utils"

	"github.com/gorilla/mux"
)

const covidCasesApi = "https://covid-api.mmediagroup.fr/v1"
const dataAPI = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2"

// CasesPerCountry gets amount of cases for a country
func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Reached country endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	val, ok := vars["begin_date-end_date"]
	if ok {
		fmt.Fprintln(w, val)
	} else {
		body, err := utils.GetBody(covidCasesApi+"/cases?country="+countryName, w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var information structs.CovidApiResponse
		if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var returnJSON structs.ReturnConfirmedCases
		returnJSON.Confirmed = information.All.Confirmed
		returnJSON.Recovered = information.All.Recovered
		returnJSON.Continent = information.All.Continent
		returnJSON.Country = information.All.Country
		returnJSON.Scope = "total"
		returnJSON.PopulationPercentage = float32(returnJSON.Confirmed) / float32(information.All.Population) * 100
		b, err := json.Marshal(returnJSON)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(b))
	}

}

// CurrentStringencyLevel shows the stringency level
func CurrentStringencyLevel(w http.ResponseWriter, r *http.Request) {

}

// Diag shows a diagnostics interface
func Diag(w http.ResponseWriter, r *http.Request) {

}
