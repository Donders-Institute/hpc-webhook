package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Donders-Institute/hpc-qaas/internal/server"
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
	db, err := server.InitDB(psqlInfo)
	if err != nil {
		panic(err)
	}

	api := server.API{DB: db}
	app := &api

	http.HandleFunc(server.WebhookPath, app.WebhookHandler)
	http.HandleFunc(server.ConfigurationPath, app.ConfigurationHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:4444", nil))
}
