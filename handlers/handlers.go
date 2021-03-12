package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CasesPerCountry gets amount of cases for a country
func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Reached country endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	fmt.Fprintln(w, countryName)

}

// CurrentStringencyLevel shows the stringency level
func CurrentStringencyLevel(w http.ResponseWriter, r *http.Request) {

}

// Diag shows a diagnostics interface
func Diag(w http.ResponseWriter, r *http.Request) {

}
