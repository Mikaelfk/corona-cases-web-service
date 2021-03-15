package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"assignment-2/handlers"
	"assignment-2/utils"

	"github.com/gorilla/mux"
)

var startTime time.Time

const version = "/corona/v1"

func init() {
	startTime = time.Now()
}

/*
 * Main function that initialized the application
 */
func main() {
	port := utils.GetPort()
	handler()
	fmt.Printf("Listening on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// handler() handles what functions
// are called when the different urls are visited
func handler() {
	r := mux.NewRouter()
	r.HandleFunc(version+"/country/{country_name}", handlers.CasesPerCountry).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(version+"/country/{country_name}", handlers.CasesPerCountry)
	r.HandleFunc(version+"/policy/{country_name}", handlers.CurrentStringencyLevel).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(version+"/policy/{country_name}", handlers.CurrentStringencyLevel)
	r.HandleFunc(version+"/diag/", handlers.Diag)
	r.HandleFunc(version+"/diag", handlers.Diag)
	http.Handle("/", r)
}
