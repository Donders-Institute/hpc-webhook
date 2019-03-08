package server

import (
	"database/sql"
	"io/ioutil"
	"strings"
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
	QaasHost                  string
	QaasPort                  string
	PrivateKeyFilename        string
	PublicKeyFilename         string
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
