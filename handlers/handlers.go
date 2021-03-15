package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	// Makes sure that the first letter of the country name is capitalized
	countryName = strings.Title(strings.ToLower(countryName))
	fmt.Fprintln(w, countryName)
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
		fmt.Fprintln(w, string(body))
	}

}

// CurrentStringencyLevel shows the stringency level
func CurrentStringencyLevel(w http.ResponseWriter, r *http.Request) {

}

// Diag shows a diagnostics interface
func Diag(w http.ResponseWriter, r *http.Request) {

}
