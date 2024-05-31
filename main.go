package main

import (
	"atalkowski/go-rpc/api"
	"database/sql"
	"log"
	"testing"

	// db "github.com/atalkowski/go-rpc/db/sqlc"

	db "github.com/atalkowski/go-rpc/db/sqlc"
	_ "github.com/lib/pq"
)

// Above ^^^ lib/pq functions not called directly here so a save will remove it; use _ to prevent this.

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:9090"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}
