package rpc

import (
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mediaConn *grpc.ClientConn
)

var MediaClient mediav1.MediaServiceClient

func NewMediaClient() error {
	var err error

	mediaConn, err = grpc.NewClient("media:50052", grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(100*1024*1024),
		grpc.MaxCallSendMsgSize(100*1024*1024),
	), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	MediaClient = mediav1.NewMediaServiceClient(mediaConn)

	return nil
}
