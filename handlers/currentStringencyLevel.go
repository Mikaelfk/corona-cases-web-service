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
	fmt.Println("Reached policy stringency endpoint")
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
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	countryCode := countryInformation[0].Alpha3Code

	country := countryInformation[0].Name
	scope := "total"
	stringency := float32(0)
	trend := float32(0)
	if ok {
		beginDate, endDate, err := utils.SplitDate(date)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		scope = beginDate + "-" + endDate
		body, err := utils.GetBody(utils.DataAPI+"/stringency/actions/"+countryCode+"/"+endDate, w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}
		var informationEndDate structs.StringencyResponse
		if err = json.Unmarshal([]byte(string(body)), &informationEndDate); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		body, err = utils.GetBody(utils.DataAPI+"/stringency/actions/"+countryCode+"/"+beginDate, w)
		if err != nil {
			// Errors are reported to the user in the GetBody function, and logged here
			log.Printf("Error: %v", err)
			return
		}

		var informationBeginDate structs.StringencyResponse
		if err = json.Unmarshal([]byte(string(body)), &informationBeginDate); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		stringencyEndDate := informationEndDate.StringencyData.StringencyActual
		stringencyBeginDate := informationBeginDate.StringencyData.StringencyActual
		if stringencyEndDate == 0 {
			stringencyEndDate = informationEndDate.StringencyData.Stringency
			if informationEndDate.StringencyData.Msg != "" {
				stringencyEndDate = -1
			}
		}
		if stringencyBeginDate == 0 {
			stringencyBeginDate = informationBeginDate.StringencyData.Stringency
			if informationBeginDate.StringencyData.Msg != "" {
				stringencyEndDate = -1
			}
		}
		trend = stringencyEndDate - stringencyBeginDate
		if stringencyBeginDate == 0 || stringencyEndDate == 0 {
			trend = 0
		}
		stringency = stringencyEndDate
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
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if information.StringencyData.Msg != "" {
			stringency = -1
		} else {
			stringency = information.StringencyData.StringencyActual
			if stringency == 0 {
				stringency = information.StringencyData.Stringency
			}
		}
	}
	var returnStringency structs.ReturnStringency
	returnStringency.Country = country
	returnStringency.Scope = scope
	returnStringency.Stringency = stringency
	returnStringency.Trend = trend
	b, err := json.Marshal(returnStringency)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(b))
}
