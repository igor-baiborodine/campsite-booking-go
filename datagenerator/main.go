package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-faker/faker/v4"
	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: datagenerator [SERVER_ADDR] [CAMPSITES_COUNT]")
	}
	addr := args[0]
	campsitesCount := 100
	if len(args) > 1 {
		campsitesCount, _ = strconv.Atoi(args[1])
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to grpc server at %s: %v", addr, err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)

	client := api.NewCampgroundsServiceClient(conn)
	campsitesIDs := createCampsites(client, campsitesCount)
	log.Printf("created %d campsites: %v", len(campsitesIDs), campsitesIDs)
}

func createCampsites(c api.CampgroundsServiceClient, count int) (campsiteIDs []string) {
	for i := 0; i < count; i++ {
		response, err := c.CreateCampsite(context.Background(), newCreateCampsiteRequest())
		if err != nil {
			log.Fatalf("failed to create campsite: %v", err)
		}
		campsiteIDs = append(campsiteIDs, response.CampsiteId)
	}
	return campsiteIDs
}

func newCreateCampsiteRequest() *api.CreateCampsiteRequest {
	req := api.CreateCampsiteRequest{}
	err := faker.FakeData(&req)
	if err != nil {
		log.Fatalf("failed to fake CreateCampsiteRequest: %v", err)
		return nil
	}
	return &req
}
