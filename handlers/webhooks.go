package handlers

import (
	"assignment-2/structs"
	"assignment-2/utils"
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
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var SignatureKey = "X-SIGNATURE"

//var Mac hash.Hash
var Secret []byte

// Webhook DB
var Webhooks map[string]structs.WebhookRegistration
var WebhookPreviousInfo map[string]float32

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
		if !utils.ValidateWebhook(webhook, w) {
			fmt.Println("Not a valid webhook body")
			http.Error(w, "Not a valid webhook body", http.StatusBadRequest)
			return
		}
		// Gets a random uuid, there is no checking if this id already exists.
		// Considering there are so many different uuid's that could be generated
		// the chance that a duplicate occurs is basically zero
		out := uuid.NewV4()
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		idString := out.String()
		webhook.ID = idString
		Webhooks[idString] = webhook
		fmt.Println("Webhook " + webhook.Url + " has been registered.")
		go CallWebhookAfterSetTime(webhook.Timeout, Webhooks[idString])
		http.Error(w, idString, http.StatusCreated)
	case http.MethodGet:
		// Returns all webhooks in JSON format
		err := json.NewEncoder(w).Encode(Webhooks)
		if err != nil {
			http.Error(w, "Something went wrong when parsing to JSON: "+err.Error(), http.StatusInternalServerError)
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
		// Checks if the webhook exists
		val, ok := Webhooks[id]
		if ok {
			webhook := structs.ReturnWebhook{}
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
			delete(WebhookPreviousInfo, id)
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
		for key, v := range Webhooks {
			CallWebhook(&v)
			Webhooks[key] = v
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
		return
	}
}

// CallUrl is used to invoke the url in the webhook
func CallUrl(url string, content string) {
	fmt.Println("Attempting invocation of url " + url + " ...")
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

// CallWebhook makes a get request to the correct API endpoint and sends this to the client with
// CallUrl
func CallWebhook(webhook *structs.WebhookRegistration) {
	if strings.ToLower(webhook.Field) == "stringency" {
		// Makes a get request to the policy endpoint
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
		var informationStringency structs.ReturnStringency
		if err = json.Unmarshal([]byte(string(body)), &informationStringency); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			return
		}

		policyResponse := string(body)
		// These if statements check if the trigger is ON_CHANGE, everything that is
		// not ON_CHANGE is considered ON_TIMEOUT
		if strings.ToLower(webhook.Trigger) == "on_change" {
			if informationStringency.Stringency == WebhookPreviousInfo[webhook.ID] {
				fmt.Println("Information is the same as previously")
			} else {
				go CallUrl(webhook.Url, policyResponse)
			}
		} else {
			go CallUrl(webhook.Url, policyResponse)
		}
		// Set the previous info of to the information recieved
		WebhookPreviousInfo[webhook.ID] = informationStringency.Stringency
	} else if strings.ToLower(webhook.Field) == "confirmed" {
		// Makes a get request to the country endpoint
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
		var informationConfirmedCases structs.ReturnConfirmedCases
		if err = json.Unmarshal([]byte(string(body)), &informationConfirmedCases); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			return
		}
		casesResponse := string(body)
		// These if statements check if the trigger is ON_CHANGE, everything that is
		// not ON_CHANGE is considered ON_TIMEOUT
		if strings.ToLower(webhook.Trigger) == "on_change" {
			if float32(informationConfirmedCases.Confirmed) == WebhookPreviousInfo[webhook.ID] {
				fmt.Println("Information is the same as previously")
			} else {
				go CallUrl(webhook.Url, casesResponse)
			}
		} else {
			go CallUrl(webhook.Url, casesResponse)
		}
		// Set the previous info of to the information recieved
		WebhookPreviousInfo[webhook.ID] = float32(informationConfirmedCases.Confirmed)
	}
}

// CallWebhookAfterSetTime is sent a timeout in seconds and a webhook,
// it calls the webhook when the timeout is reached.
func CallWebhookAfterSetTime(timeout int, webhook structs.WebhookRegistration) {
	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				webhook = Webhooks[webhook.ID]
				go CallWebhook(&webhook)
			}
		}
	}()
	go func() {
		for {
			if (structs.WebhookRegistration{}) == webhook {
				ticker.Stop()
				quit <- true
			}
		}
	}()
}
