package main

import (
	"database/sql"
	"log"

	api "github.com/atalkowski/go-rpc/api"
	db "github.com/atalkowski/go-rpc/db/sqlc"
	_ "github.com/lib/pq"
)

// Above ^^^ lib/pq functions not called directly here so a save will remove it; use _ to prevent this.

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:9090"
)

func main() {
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
