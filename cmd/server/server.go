package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Donders-Institute/hpc-qaas/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	// Set Qaas server variables
	qaasHost := os.Getenv("QAAS_HOST")
	qaasPort := os.Getenv("QAAS_PORT")
	address := fmt.Sprintf("%s:%s", qaasHost, qaasPort)

	// Set target computer variables
	relayNode := os.Getenv("RELAY_NODE")
	relayNodeTestUser := os.Getenv("RELAY_NODE_TEST_USER")
	relayNodeTestUserPassword := os.Getenv("RELAY_NODE_TEST_USER_PASSWORD")

	// Set the database variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")

	// Override settings if we run the server in a Docker container
	if server.RunsWithinContainer() {
		host = "db"
		address = fmt.Sprintf("0.0.0.0:%s", qaasPort)
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
		DataDir:                   "",
		VaultDir:                  "",
		HomeDir:                   "",
		RelayNode:                 relayNode,
		RelayNodeTestUser:         relayNodeTestUser,
		RelayNodeTestUserPassword: relayNodeTestUserPassword,
		QaasHost:                  qaasHost,
		QaasPort:                  qaasPort,
	}

	// Set the data dir and create it
	api.SetDataDir("/data")
	err = os.MkdirAll(api.DataDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Set the vault dir and create it
	api.SetVaultDir("/vault")
	err = os.MkdirAll(api.VaultDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Set the /home dir
	api.SetHomeDir("/home")

	app := &api

	// Handle external webhook payloads
	http.HandleFunc(server.WebhookPath, app.WebhookHandler)

	// Handle internal webhook configuration payloads
	http.HandleFunc(server.ConfigurationPath, app.ConfigurationHandler)

	log.Fatal(http.ListenAndServe(address, nil))
}
