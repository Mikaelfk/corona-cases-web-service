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

func init() {
	handlers.StartTime = time.Now()
	handlers.RegisteredWebhooks = 0
}

/*
 * Main function that initialized the application
 */
func main() {
	port := utils.GetPort()
	handler()
	fmt.Printf("Listening on port :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// handler() handles what functions
// are called when the different urls are visited
func handler() {
	r := mux.NewRouter()
	r.HandleFunc(handlers.URLBase+handlers.Version+"/country/{country_name}", handlers.CasesPerCountry).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(handlers.URLBase+handlers.Version+"/country/{country_name}", handlers.CasesPerCountry)
	r.HandleFunc(handlers.URLBase+handlers.Version+"/policy/{country_name}", handlers.CurrentStringencyLevel).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(handlers.URLBase+handlers.Version+"/policy/{country_name}", handlers.CurrentStringencyLevel)
	r.HandleFunc(handlers.URLBase+handlers.Version+"/diag/", handlers.Diag)
	r.HandleFunc(handlers.URLBase+handlers.Version+"/diag", handlers.Diag)
	http.Handle("/", r)
}
