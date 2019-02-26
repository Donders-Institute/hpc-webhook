package server

import (
	"fmt"
	"net/http"
)

// WebhookPath is the first part of the webhook payload URL
const WebhookPath = "/webhook/"

// ConfigurationPath is the URL path to add a new webhook
const ConfigurationPath = "/configuration"

// ConfigurationHandler handles a webhook registration HTTP PUT request
// with the hash and username in its body
func ConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// TODO: Add a row to database
	fmt.Printf("%+v\n", configuration)

	// Succes
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Webhook added successfully")
	return
}

// WebhookHandler handles a HTTP POST request containing the webhook payload in its body
func WebhookHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	webhook, webhookID, err := parseWebhookRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Check if webhookID exists
	if err := checkWebhookID(webhookID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Parse the webhook payload
	var payload []byte
	payload, err = parseWebhookPayload(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	fmt.Printf("Payload: %+v\n", payload)

	// Execute the script
	if err := ExecuteScript(webhook); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Webhook payload delivered successfully")
	return
}
