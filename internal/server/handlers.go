package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// API is used to store the database pointer
type API struct {
	DB        *sql.DB
	DataDir   string
	RelayNode string
	QaasHost  string
	QaasPort  string
}

// WebhookPath is the first part of the webhook payload URL
const WebhookPath = "/webhook/"

// ConfigurationPath is the URL path to add a new webhook
const ConfigurationPath = "/configuration"

// RunsWithinContainer checks if the program runs in a Docker container or not
func RunsWithinContainer() bool {
	file, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(file), "docker")
}

// SetDataDir sets the filepath of the qaas data files
func (a *API) SetDataDir(elem ...string) {
	a.DataDir = filepath.Join(elem...)
}

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

	// Create the key dir
	keyDir := path.Join(a.DataDir, "keys", configuration.Username)
	err = os.MkdirAll(keyDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Check if key pair exists for user
	privateKeyFilename := path.Join(keyDir, "id_rsa")
	publicKeyFilename := path.Join(keyDir, "id_rsa.pub")
	hasPrivateKeyFilename, _ := checkFile(privateKeyFilename)
	hasPublicKeyFilename, _ := checkFile(publicKeyFilename)
	if !hasPrivateKeyFilename || !hasPublicKeyFilename {
		err = generateKeyPair(privateKeyFilename, publicKeyFilename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println(err)
			fmt.Fprint(w, "Error 404 - Not found: ", err)
			return
		}
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

	// Create the payload dir
	payloadDir := path.Join(a.DataDir, "payloads", configuration.Username)
	err = os.MkdirAll(payloadDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Write the payload to file
	err = writeWebhookPayloadToFile(payloadDir, payload, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Execute the script
	fmt.Printf("Webhook: %+v\n", webhook)
	if err := ExecuteScript(a.RelayNode, a.DataDir, webhookID, payload, username); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	webhookPayloadURL := fmt.Sprintf("https://%s:%s/webhook/%s", a.QaasHost, a.QaasPort, webhookID)
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
