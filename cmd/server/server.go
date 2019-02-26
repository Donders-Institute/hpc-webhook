package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Donders-Institute/hpc-qaas/internal/pkg/server"
	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	server.InitDB(psqlInfo)

	http.HandleFunc(server.WebhookPath, server.WebhookHandler)
	http.HandleFunc(server.ConfigurationPath, server.ConfigurationHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:4444", nil))
}
