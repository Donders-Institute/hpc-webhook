package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// API is used to store the database pointer
type API struct {
	DB *sql.DB
}

// WebhookPath is the first part of the webhook payload URL
const WebhookPath = "/webhook/"

// ConfigurationPath is the URL path to add a new webhook
const ConfigurationPath = "/configuration"

// QaasHost stores the hostname of the qaas server
const QaasHost = "qaas.dccn.nl"

// QaasPort stores the port of the qaas server
const QaasPort = "5111"

// ConfigurationHandler handles a webhook registration HTTP PUT request
// with the hash and username in its body
func (a *API) ConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Add a row in the database
	err = addRow(a.DB, configuration.Hash, configuration.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Webhook added successfully")
	return
}

// WebhookHandler handles a HTTP POST request containing the webhook payload in its body
func (a *API) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	webhook, webhookID, err := parseWebhookRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Check if webhookID exists
	username, err := checkWebhookID(a.DB, webhookID)
	if err != nil {
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

	// Execute the script
	fmt.Printf("Webhook: %+v\n", webhook)
	if err := ExecuteScript(webhookID, payload, username); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	webhookPayloadURL := fmt.Sprintf("https://%s:%s/webhook/%s", QaasHost, QaasPort, webhookID)
	configurationResponse := ConfigurationResponse{
		Webhook: webhookPayloadURL,
	}
	js, err := json.Marshal(configurationResponse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
