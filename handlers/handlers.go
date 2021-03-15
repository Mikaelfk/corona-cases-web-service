package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"assignment-2/structs"
	"assignment-2/utils"

	"github.com/gorilla/mux"
)

const covidCasesAPI = "https://covid-api.mmediagroup.fr/v1"
const dataAPI = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2"

// CasesPerCountry gets amount of cases for a country
func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Reached country endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	date, ok := vars["begin_date-end_date"]
	if ok {
		splitdate := strings.Split(date, "-")
		if len(splitdate) < 6 {
			// Handles string error
			err := errors.New("Error in date")
			log.Printf("Error, %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Splits the dates into two differen strings
		beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
		endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
		if !utils.ValidDate(beginDate) || !utils.ValidDate(endDate) {
			log.Println("Error in date")
			http.Error(w, "Error in date", http.StatusBadRequest)
			return
		}
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
		var returnJSON structs.ReturnConfirmedCases
		returnJSON.Confirmed = information.All.Confirmed
		returnJSON.Recovered = information.All.Recovered
		returnJSON.Continent = information.All.Continent
		returnJSON.Country = information.All.Country
		returnJSON.Scope = "total"
		returnJSON.PopulationPercentage = float32(int((float32(returnJSON.Confirmed)/float32(information.All.Population)*100)*100)) / 100
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
