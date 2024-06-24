package gapi

import (
	db "github.com/atalkowski/go-rpc/db/sqlc"
	"github.com/atalkowski/go-rpc/pb"
)

// Server serves gRPC requests
type Server struct {
	store *db.Store
	pb.UnimplementedSimpleBankServer
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) (*Server, error) {
	server := &Server{store: store}
	return server, nil
}
