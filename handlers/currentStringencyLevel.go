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
		t, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
			return
		}
		y, m, d := t.AddDate(0, 0, -10).Date()
		endDate = strconv.Itoa(y) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(d)
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
			http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
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
			http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
			return
		}
		stringencyEndDate := informationEndDate.StringencyData.StringencyActual
		stringencyBeginDate := informationBeginDate.StringencyData.StringencyActual
		trend = stringencyEndDate - stringencyBeginDate
		if informationEndDate.StringencyData.Msg != "" {
			stringencyEndDate = -1
			trend = 0
		}
		if informationBeginDate.StringencyData.Msg != "" {
			stringencyBeginDate = -1
			trend = 0
		}
		fmt.Println(stringency)
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
		if information.StringencyData.Msg != "" {
			stringency = -1
		} else {
			stringency = information.StringencyData.StringencyActual
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
		http.Error(w, "Error: "+err.Error(), utils.InternalServerError)
	}
	fmt.Fprint(w, string(b))
}
