package main

import (
	"log"
	"net/http"
	"time"

	"assignment-2/utils"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

/*
 * Main function that initialized the application
 */
func main() {
	port := utils.GetPort()
	utils.Handler()
	log.Printf("Listening on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
