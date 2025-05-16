package rpc

import (
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mediaConn *grpc.ClientConn
)

var MediaClient mediav1.MediaClient

func NewMediaClient() error {
	var err error

	mediaConn, err = grpc.NewClient("media:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	MediaClient = mediav1.NewMediaClient(mediaConn)

	return nil
}
