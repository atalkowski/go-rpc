package main

import (
	"database/sql"
	"log"
	"net"

	api "github.com/atalkowski/go-rpc/api"
	db "github.com/atalkowski/go-rpc/db/sqlc"
	gapi "github.com/atalkowski/go-rpc/gapi"
	"github.com/atalkowski/go-rpc/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Above ^^^ lib/pq functions not called directly here so a save will remove it; use _ to prevent this.

const (
	dbDriver          = "postgres"
	dbSource          = "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable"
	HTTPServerAddress = "0.0.0.0:9091"
	GRPCServerAddress = "0.0.0.0:9090"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	storeDb := db.NewStore(conn)
	runGrpcServer(storeDb)
	runGinServer(storeDb)
}

func runGrpcServer(storeDb *db.Store) {
	server, err := gapi.NewServer(storeDb)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener for gRPC at "+GRPCServerAddress, err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start gRPC server at "+GRPCServerAddress, err)
	}
	log.Printf("gRPC server now listening on %s", GRPCServerAddress)
}

func runGinServer(storeDb *db.Store) {
	server := api.NewServer(storeDb)
	err := server.Start(HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start HTTP server at "+HTTPServerAddress, err)
	}
}
