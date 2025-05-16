package rpc

import (
	authv1 "github.com/co1seam/ember-backend-api-contracts/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	authConn *grpc.ClientConn
)

var AuthClient authv1.AuthClient

func NewAuthClient() error {
	var err error

	authConn, err = grpc.NewClient("auth:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	AuthClient = authv1.NewAuthClient(authConn)

	return nil
}
