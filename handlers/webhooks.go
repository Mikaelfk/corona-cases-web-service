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

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var SignatureKey = "X-SIGNATURE"
var ClientSignatureKey = "X-SIGNATURE"

//var Mac hash.Hash
var Secret []byte

// Webhook DB
var Webhooks []structs.WebhookRegistration

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
		}
		out := uuid.NewV4()
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		}
		idString := out.String()
		webhook.ID = idString
		Webhooks = append(Webhooks, webhook)

		fmt.Println("Webhook " + webhook.Url + " has been registered.")
		http.Error(w, webhook.ID, http.StatusCreated)
	case http.MethodGet:
		// Returns all webhooks in JSON format
		err := json.NewEncoder(w).Encode(Webhooks)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
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
		for _, v := range Webhooks {
			if v.ID == id {
				err := json.NewEncoder(w).Encode(v)
				if err != nil {
					http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
				}
			}
		}
	case http.MethodDelete:
		i := 0
		for _, v := range Webhooks {
			if v.ID == id {
				fmt.Printf("Deleted Webhook with ID: %s", v.ID)
				Webhooks = utils.Remove(Webhooks, i)
				i--
			}
			i++
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

// Invokes the web service to trigger event. Currently only responds to POST requests.
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Received POST request...")
		for _, v := range Webhooks {
			go CallUrl(v.Url, "Trigger event")
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
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
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write(content)
	if err != nil {
		http.Error(w, "Error during message decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Content: " + hex.EncodeToString(mac.Sum(nil)))

	// Compare HMAC with received request
	if hmac.Equal(signatureByte, mac.Sum(nil)) {
		fmt.Println("Valid invocation (with validated content) on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
		if err != nil {
			fmt.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		fmt.Println("Invalid invocation (tampered content?) on " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
	}
}
