package gapi

import (
	"context"
	"log"

	db "github.com/atalkowski/go-rpc/db/sqlc"
	"github.com/atalkowski/go-rpc/pb"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {

	log.Printf("RPC CreateAccount(%s, %s)", req.GetOwner(), req.GetCurrency())
	arg := db.CreateAccountParams{
		Owner:    req.GetOwner(),
		Currency: req.GetCurrency(),
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		// Ig this is the postgress error for unique constraint violation...
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "account already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create account: %s", err)
	}

	response := &pb.CreateAccountResponse{
		Account: convertAccount(account),
	}
	log.Printf("RPC Created Account ID: %v %s %s", response.Account.GetId(),
		response.Account.GetOwner(), response.Account.GetCurrency())
	return response, nil
}
