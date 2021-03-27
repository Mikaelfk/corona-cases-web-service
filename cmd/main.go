package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"assignment-2/handlers"
	"assignment-2/structs"
	"assignment-2/utils"

	"github.com/gorilla/mux"
)

func init() {
	handlers.StartTime = time.Now()
	handlers.Secret = []byte{143, 32, 33, 46, 59}
	handlers.Webhooks = make(map[string]structs.WebhookRegistration)
	handlers.WebhookPreviousInfo = make(map[string]float32)
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
	// for registration
	webhookEndpoint := "/notifications/"
	// for invocation
	serviceEndpoint := "/service"

	r := mux.NewRouter()
	r.HandleFunc(utils.URLBase+utils.Version+"/country/{country_name}", handlers.CasesPerCountry).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(utils.URLBase+utils.Version+"/country/{country_name}", handlers.CasesPerCountry)
	r.HandleFunc(utils.URLBase+utils.Version+"/policy/{country_name}", handlers.CurrentStringencyLevel).Queries("scope", "{begin_date-end_date}")
	r.HandleFunc(utils.URLBase+utils.Version+"/policy/{country_name}", handlers.CurrentStringencyLevel)
	r.HandleFunc(utils.URLBase+utils.Version+"/diag/", handlers.Diag)
	r.HandleFunc(utils.URLBase+utils.Version+"/diag", handlers.Diag)
	r.HandleFunc(utils.URLBase+utils.Version+webhookEndpoint, handlers.WebhookRegistrationHandler)
	r.HandleFunc(utils.URLBase+utils.Version+webhookEndpoint+"{id}", handlers.WebhookIDHandler)
	r.HandleFunc(serviceEndpoint, handlers.ServiceHandler)

	http.Handle("/", r)
}
