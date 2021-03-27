package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Invoked Hash to be accepted
var secret = []byte{143, 32, 33, 46, 59}

// used for content-based validation
var ClientSignatureKey = "X-SIGNATURE"

// ContentValidatingHandler validates the incoming content and if it succeeds prints the content to the console
func ContentValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Simply print body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: "+err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Extract signature from header based on known key
	signature := r.Header.Get(ClientSignatureKey)

	// Convert string to []byte
	signatureByte, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, "Error during Signature decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Signature: " + signature)
	// Hash content of body
	mac := hmac.New(sha256.New, secret)
	_, err = mac.Write(content)
	if err != nil {
		http.Error(w, "Error during message decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Content: " + hex.EncodeToString(mac.Sum(nil)))

	// Compare HMAC with received request
	if hmac.Equal(signatureByte, mac.Sum(nil)) {
		fmt.Println("Valid invocation (with validated content) on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked web service.")
		if err != nil {
			fmt.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		fmt.Println("Invalid invocation (tampered content?) on " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
	}
}

func main() {

	port := "8081"

	// Environment variable constant for Heroku support
	PORT := "PORT"

	if os.Getenv(PORT) != "" {
		//port = os.Getenv(PORT)
	}

	endpoint := "/invoked/"

	r := mux.NewRouter()
	fmt.Println("Service listening on port " + port)
	fmt.Println("Service URL (content-validating): http://localhost:" + port + endpoint)
	r.HandleFunc(endpoint, ContentValidatingHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
