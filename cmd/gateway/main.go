package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"git.netflux.io/rob/solar-toolkit/gateway/handler"
	"git.netflux.io/rob/solar-toolkit/gateway/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const defaultBindAddr = ":8888"

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("missing configuration DATABASE_URL")
	}

	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = defaultBindAddr
	}

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	store := store.NewSQL(db)
	handler := handler.New(store)
	srv := http.Server{
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		Handler:      handler,
		Addr:         bindAddr,
	}

	log.Printf("Listening on %s...", bindAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
