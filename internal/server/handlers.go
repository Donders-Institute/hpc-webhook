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
	DB                        *sql.DB
	Connector                 Connector
	DataDir                   string
	VaultDir                  string
	HomeDir                   string
	RelayNode                 string
	RelayNodeTestUser         string
	RelayNodeTestUserPassword string
	QaasHost                  string
	QaasPort                  string
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

// SetVaultDir sets the filepath of the private key files with root read and write file permissions
func (a *API) SetVaultDir(elem ...string) {
	a.VaultDir = filepath.Join(elem...)
}

// SetHomeDir sets the filepath to the mounted /home dir
func (a *API) SetHomeDir(elem ...string) {
	a.HomeDir = filepath.Join(elem...)
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

		err = addAuthorizedPublicKey(a.Connector, privateKeyFilename, publicKeyFilename, a.RelayNodeTestUser, a.RelayNodeTestUserPassword, a.RelayNode)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println(err)
			fmt.Fprint(w, "Error 404 - Not found: ", err)
			return
		}
	}

	// Add a row in the database
	err = addRow(a.DB, configuration.Hash, configuration.Groupname, configuration.Username)
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
	_, webhookID, err := parseWebhookRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Check if webhookID exists
	groupname, username, err := checkWebhookID(a.DB, webhookID)
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
	payloadDir := path.Join(a.DataDir, "payloads", username)
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

	// Prepare the execution of the script
	privateKeyFilename := path.Join(a.DataDir, "keys", username, "id_rsa")
	payloadFilename := path.Join(a.DataDir, "payloads", username, "payload")
	targetPayloadDir := path.Join(a.HomeDir, groupname, username, ".qaas", webhookID)
	targetPayloadFilename := path.Join(targetPayloadDir, "payload")
	userScriptPathFilename := path.Join(a.HomeDir, groupname, username, ".qaas", webhookID, "script.sh")
	tempPrivateKeyDir := path.Join(a.VaultDir, username)
	tempPrivateKeyFilename := path.Join(tempPrivateKeyDir, "id_rsa")

	executeConfig := executeConfiguration{
		privateKeyFilename:     privateKeyFilename,
		tempPrivateKeyDir:      tempPrivateKeyDir,
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		targetPayloadDir:       targetPayloadDir,
		targetPayloadFilename:  targetPayloadFilename,
		userScriptPathFilename: userScriptPathFilename,
		username:               username,
		groupname:              groupname,
		password:               a.RelayNodeTestUserPassword,
		relayNodeName:          a.RelayNode,
		dataDir:                a.DataDir,
		vaultDir:               a.VaultDir,
		homeDir:                a.HomeDir,
		webhookID:              webhookID,
		payload:                payload,
	}

	// Execute the script
	if err := ExecuteScript(a.Connector, executeConfig); err != nil {
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
