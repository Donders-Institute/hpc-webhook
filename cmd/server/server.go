package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Donders-Institute/hpc-webhook/internal/server"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	homeDir            = "/home"
	dataDir            = "/data"
	privateKeyFilename = "/run/secrets/hpc_webhook_private_key"
	publicKeyFilename  = "/run/secrets/hpc_webhook_public_key"
)

func main() {
	// Set HPC webhook server variables
	webhookBaseURL := strings.TrimRight(os.Getenv("WEBHOOK_BASEURL"), "/")

	// Set target computer variables
	relayNode := os.Getenv("RELAY_ACCESS_NODE")
	connectionTimeoutSeconds, err := strconv.Atoi(os.Getenv("RELAY_CONNECTION_TIMEOUT_SECONDS"))
	if err != nil {
		panic(err)
	}

	// Set the database variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")

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
		DataDir:                  dataDir,
		HomeDir:                  homeDir,
		RelayNode:                relayNode,
		ConnectionTimeoutSeconds: connectionTimeoutSeconds,
		WebhookBaseURL:           webhookBaseURL,
		PrivateKeyFilename:       privateKeyFilename,
		PublicKeyFilename:        publicKeyFilename,
		Scheduler:                server.MyScheduler(os.Getenv("RELAY_JOB_SCHEDULER")),
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

	fmt.Printf("starting server on port 5111 ...")

	log.Fatal(http.ListenAndServe("0.0.0.0:5111", r))
}
