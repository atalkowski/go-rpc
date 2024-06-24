package gapi

import (
	db "github.com/atalkowski/go-rpc/db/sqlc"
	"github.com/atalkowski/go-rpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertAccount(acc db.Account) *pb.Account {
	return &pb.Account{
		Id:        acc.ID, // Warning please consider removing this
		Owner:     acc.Owner,
		Currency:  acc.Currency,
		Balance:   acc.Balance,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}
}
