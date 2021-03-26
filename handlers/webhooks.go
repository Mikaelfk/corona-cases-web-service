package handlers

import (
	"assignment-2/structs"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var SignatureKey = "X-SIGNATURE"

//var Mac hash.Hash
var Secret []byte

// Webhook DB
var Webhooks map[string]structs.WebhookRegistration

// WebhookRegistrationHandler handles webhook registration (POST) and lookup (GET) requests.
// Expects WebhookRegistration struct body in request.
func WebhookRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		// Expects incoming body in terms of WebhookRegistration struct
		webhook := structs.WebhookRegistration{}
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
			return
		}
		out := uuid.NewV4()
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		idString := out.String()
		Webhooks[idString] = webhook

		fmt.Println("Webhook " + webhook.Url + " has been registered.")
		http.Error(w, idString, http.StatusCreated)
	case http.MethodGet:
		// Returns all webhooks in JSON format
		err := json.NewEncoder(w).Encode(Webhooks)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

// WebhookIDHandler handles webhook lookup (GET) and deletion (DELETE) requests by ID
func WebhookIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		val, ok := Webhooks[id]
		if ok {
			webhook := structs.WebhookResponse{}
			webhook.Country = val.Country
			webhook.Field = val.Field
			webhook.ID = id
			webhook.Timeout = val.Timeout
			webhook.Trigger = val.Trigger
			webhook.Url = val.Url
			err := json.NewEncoder(w).Encode(webhook)
			if err != nil {
				http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "No webhook with this ID", http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		_, ok := Webhooks[id]
		if ok {
			fmt.Printf("Deleted webhook with ID: %s", id)
			fmt.Fprintf(w, "Deleted webhook with ID: %s", id)
			delete(Webhooks, id)
		} else {
			http.Error(w, "No webhook with this ID", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
		return
	}
}

// Invokes the web service to trigger event. Currently only responds to POST requests.
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Received POST request...")
		for _, v := range Webhooks {
			go CallWebhook(v)
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
		return
	}
}

func CallUrl(url string, content string) {
	fmt.Println("Attempting invocation of url " + url + " ...")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(content)))
	if err != nil {
		fmt.Errorf("%v", "Error during request creation.")
		return
	}
	// Hash content
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		fmt.Errorf("%v", "Error during content hashing.")
		return
	}
	// Convert to string & add to header
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in HTTP request: " + err.Error())
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Something is wrong with invocation response: " + err.Error())
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}

func CallWebhook(webhook structs.WebhookRegistration) {
	if strings.ToLower(webhook.Field) == "stringency" {
		resp, err := http.Get("http://localhost:8080/corona/v1/policy/" + webhook.Country)
		if err != nil {
			log.Println("Error when making get request")
			go CallUrl(webhook.Url, "Error in request")
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error when reading body")
			go CallUrl(webhook.Url, "Error when reading body")
			return
		}
		policyResponse := string(body)
		go CallUrl(webhook.Url, policyResponse)
	} else if strings.ToLower(webhook.Field) == "confirmed" {
		resp, err := http.Get("http://localhost:8080/corona/v1/country/" + webhook.Country)
		if err != nil {
			log.Println("Error when making get request")
			go CallUrl(webhook.Url, "Error in request")
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Println("Error when reading body")
			go CallUrl(webhook.Url, "Error when reading body")
			return
		}
		casesResponse := string(body)
		go CallUrl(webhook.Url, casesResponse)
	}
}

func checkIfTimeoutReached() {

}
