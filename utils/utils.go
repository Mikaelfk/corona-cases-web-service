package utils

import (
	"assignment-2/handlers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const version = "/corona/v1/"

// A handler function for handling what functions
// are called when the different urls are visited
func Handler() {
	r := mux.NewRouter()
	r.HandleFunc(version+"{country_name}", handlers.CasesPerCountry).Queries("begin_date-end_date", "{begin_date-end_date}")
	r.HandleFunc(version+"{country_name}", handlers.CasesPerCountry)
	r.HandleFunc(version+"{country_name}", handlers.CurrentStringencyLevel).Queries("begin_date-end_date", "{begin_date-end_date}")
	r.HandleFunc(version+"{country_name}", handlers.CurrentStringencyLevel)
	r.HandleFunc(version+"diag/", handlers.Diag)
	r.HandleFunc(version+"diag", handlers.Diag)
	http.Handle("/", r)
}

// Gets a port from the environmental variable PORT, if this is empty, set port to 8181
func GetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8181"
	}
	return port
}
