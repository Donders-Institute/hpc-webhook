package server

import (
	"database/sql"
	"io/ioutil"
	"strings"
)

// Setup of user's workspace directories and files
const (
	WebhooksWorkDir = ".webhook" // WebhooksWorkDir denotes the user's work directory
	PayLoadName     = "payload"  // PayLoadName is the name of the payload file in user's work directory
	ScriptName      = "script"   // ScriptName is the name of the script in the user's work directory
)

// API is used to store the database pointer
type API struct {
	DB                        *sql.DB
	Connector                 Connector
	DataDir                   string
	HomeDir                   string
	RelayNode                 string
	RelayNodeTestUser         string
	RelayNodeTestUserPassword string
	ConnectionTimeoutSeconds  int
	HPCWebhookHost            string
	HPCWebhookInternalPort    string // Port for internal use
	HPCWebhookExternalPort    string // Port for the outside world
	PrivateKeyFilename        string
	PublicKeyFilename         string
}

// WebhookPath is the basic part of the webhook payload URL
const WebhookPath = "/webhook"

// WebhookPostPath is the first part of the webhook payload URL [POST]
const WebhookPostPath = "/webhook/{webhook}"

// ConfigurationPath is the basic URL path for configuring the HPC webhook server
const ConfigurationPath = "/configuration"

// ConfigurationAddPath is the URL path to add a new webhook [PUT]
const ConfigurationAddPath = "/configuration"

// ConfigurationListPath is the URL path to list all webhook for a certain user [GET]
const ConfigurationListPath = "/configuration"

// ConfigurationInfoPath is the URL path to get detailed information about a certain webhook [GET]
const ConfigurationInfoPath = "/configuration/{webhook}"

// ConfigurationDeletePath is the URL path to delete a certain webhook [DELETE]
const ConfigurationDeletePath = "/configuration/{webhook}"

// RunsWithinContainer checks if the program runs in a Docker container or not
func RunsWithinContainer() bool {
	file, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(file), "docker")
}
