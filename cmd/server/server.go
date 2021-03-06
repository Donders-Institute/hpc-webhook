package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Donders-Institute/hpc-webhook/internal/server"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Set HPC webhook server variables
	hpcWebhookHost := os.Getenv("HPC_WEBHOOK_HOST")
	hpcWebhookInternalPort := os.Getenv("HPC_WEBHOOK_INTERNAL_PORT")
	hpcWebhookExternalPort := os.Getenv("HPC_WEBHOOK_EXTERNAL_PORT")
	address := fmt.Sprintf("%s:%s", hpcWebhookHost, hpcWebhookInternalPort)
	homeDir := os.Getenv("HOME_DIR")
	dataDir := os.Getenv("DATA_DIR")
	privateKeyFilename := os.Getenv("PRIVATE_KEY_FILE")
	publicKeyFilename := os.Getenv("PUBLIC_KEY_FILE")

	// Set target computer variables
	relayNode := os.Getenv("RELAY_NODE")
	relayNodeTestUser := os.Getenv("RELAY_NODE_TEST_USER")
	relayNodeTestUserPassword := os.Getenv("RELAY_NODE_TEST_USER_PASSWORD")
	connectionTimeoutSeconds, err := strconv.Atoi(os.Getenv("CONNECTION_TIMEOUT_SECONDS"))
	if err != nil {
		panic(err)
	}

	// Set the database variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")

	// Override settings if we run the server in a Docker container
	if server.RunsWithinContainer() {
		host = "db"
		address = fmt.Sprintf("0.0.0.0:%s", hpcWebhookInternalPort)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := server.InitDB(psqlInfo)
	if err != nil {
		panic(err)
	}

	// Setup the app
	api := server.API{
		DB: db,
		Connector: server.SSHConnector{
			Description: "SSH connection to relay node",
		},
		DataDir:                   dataDir,
		HomeDir:                   homeDir,
		RelayNode:                 relayNode,
		RelayNodeTestUser:         relayNodeTestUser,
		RelayNodeTestUserPassword: relayNodeTestUserPassword,
		ConnectionTimeoutSeconds:  connectionTimeoutSeconds,
		HPCWebhookHost:            hpcWebhookHost,
		HPCWebhookInternalPort:    hpcWebhookInternalPort,
		HPCWebhookExternalPort:    hpcWebhookExternalPort,
		PrivateKeyFilename:        privateKeyFilename,
		PublicKeyFilename:         publicKeyFilename,
	}

	// Set the data dir and create it
	err = os.MkdirAll(api.DataDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	app := &api

	r := mux.NewRouter()

	// Handle external webhook payloads
	r.HandleFunc(server.WebhookPostPath, app.WebhookHandler).Methods("POST")

	// Handle internal webhook configuration payloads
	r.HandleFunc(server.ConfigurationAddPath, app.ConfigurationAddHandler).Methods("PUT")
	r.HandleFunc(server.ConfigurationInfoPath, app.ConfigurationInfoHandler).Methods("GET")
	r.HandleFunc(server.ConfigurationListPath, app.ConfigurationListHandler).Methods("GET")
	r.HandleFunc(server.ConfigurationDeletePath, app.ConfigurationDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(address, r))
}
