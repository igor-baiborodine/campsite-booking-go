package client

import (
	"log"

	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCampgroundsServiceClient(addr string) api.CampgroundsServiceClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to grpc server at %s: %v", addr, err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)
	return api.NewCampgroundsServiceClient(conn)
}
