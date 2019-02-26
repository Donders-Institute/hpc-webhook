package main

import (
	"log"
	"net/http"

	"github.com/Donders-Institute/hpc-qaas/pkg/server"
)

func main() {
	http.HandleFunc(server.WebhookPath, server.WebhookHandler)
	http.HandleFunc(server.ConfigurationPath, server.ConfigurationHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:4444", nil))
}
