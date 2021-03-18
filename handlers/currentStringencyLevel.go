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

	"github.com/gorilla/mux"
)

// CurrentStringencyLevel shows the stringency level
func CurrentStringencyLevel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Reached country endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	date, ok := vars["begin_date-end_date"]
	body, err := utils.GetBody(utils.CountryAPI+"/name/"+countryName, w)
	if err != nil {
		// Errors are reported to the user in the GetBody function, and logged here
		log.Printf("Error: %v", err)
		return
	}
	var countryInformation []structs.CountryResponse
	if err = json.Unmarshal([]byte(string(body)), &countryInformation); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
		return
	}
	countryCode := countryInformation[0].Alpha3Code

	if ok {
		fmt.Println(date)
	} else {
		y, m, d := time.Now().AddDate(0, 0, -10).Date()
		dateString := strconv.Itoa(y) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(d)
		body, err := utils.GetBody(utils.DataAPI+"/stringency/actions/"+countryCode+"/"+dateString, w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var information structs.StringencyResponse
		if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
			return
		}
		var returnStringency structs.ReturnStringency
		returnStringency.Country = countryInformation[0].Name
		returnStringency.Scope = "total"
		returnStringency.Stringency = information.StringencyData.StringencyActual
		returnStringency.Trend = 0
		b, err := json.Marshal(returnStringency)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
		}
		fmt.Fprint(w, string(b))
	}
}
