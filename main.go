package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arxxm/CRUD_test/api"
	"github.com/arxxm/CRUD_test/handler"
	_ "github.com/lib/pq"
)

const (
	// Initialize connection constants.
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "2245core"
	dbname   = "crud_pg"
)

func main() {

	psqlconn := fmt.Sprintf("host= %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	// psqlconn := fmt.Sprintf("user = %s password = %s dbname = %s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	repo := api.NewRepository(db)

	h, err := handler.NewAPIHandler(repo)
	if err != nil {
		log.Fatal(err)
	}

	m := http.NewServeMux()
	rtr := h.InitRoutes()

	m.Handle("/", rtr)
	var srv = &http.Server{
		Addr:    ":8080",
		Handler: m,
		// TLSConfig: nil,

		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal()
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
