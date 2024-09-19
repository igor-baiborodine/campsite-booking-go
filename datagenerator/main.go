package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-faker/faker/v4"
	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	CampsiteFaker struct {
		Capacity int `faker:"boundary_start=1, boundary_end=10"`
	}

	BookingFaker struct {
		Email     string `faker:"email"`
		FirstName string `faker:"first_name"`
		LastName  string `faker:"last_name"`
	}

	BookingStayFaker struct {
		Period int `faker:"boundary_start=1, boundary_end=3"`
	}
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
	log.Printf("created %d campsites", len(campsitesIDs))

	count := createBookings(client, campsitesIDs)
	log.Printf("created total %d bookings", count)
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
	campsite := CampsiteFaker{}
	err := faker.FakeData(&campsite)
	if err != nil {
		log.Fatalf("failed to create CampsiteFaker: %v", err)
		return nil
	}

	req := api.CreateCampsiteRequest{}
	err = faker.FakeData(&req)
	if err != nil {
		log.Fatalf("failed to create CreateCampsiteRequest: %v", err)
		return nil
	}
	req.Capacity = int32(campsite.Capacity)
	return &req
}

func createBookings(c api.CampgroundsServiceClient, campsiteIDs []string) (count int) {
	now := time.Now().UTC()
	maxAllowedEndDate := now.AddDate(0, 1, 0)

	for _, campsiteID := range campsiteIDs {
		countPerCampsite := 0
		startDate := now.AddDate(0, 0, 1)
		for {
			endDate := startDate.AddDate(0, 0, newBookingStayPeriod())
			if !endDate.Before(maxAllowedEndDate) {
				break
			}
			req := newCreateBookingRequest(campsiteID, startDate, endDate)

			_, err := c.CreateBooking(context.Background(), req)
			if err != nil {
				log.Fatalf("failed to create booking for campsite ID %s: %v", campsiteID, err)
			}
			startDate = endDate
			countPerCampsite++
		}
		log.Printf("...created %d bookings for campsite ID %s", countPerCampsite, campsiteID)
		count += countPerCampsite
	}
	return count
}

func newCreateBookingRequest(
	campsiteId string,
	startDate time.Time,
	endDate time.Time,
) *api.CreateBookingRequest {
	booking := BookingFaker{}
	err := faker.FakeData(&booking)
	if err != nil {
		log.Fatalf("failed to create BookingFaker: %v", err)
	}

	return &api.CreateBookingRequest{
		CampsiteId: campsiteId,
		Email:      booking.Email,
		FullName:   booking.FirstName + " " + booking.LastName,
		StartDate:  startDate.Format(time.DateOnly),
		EndDate:    endDate.Format(time.DateOnly),
	}
}

func newBookingStayPeriod() int {
	bookingStay := BookingStayFaker{}
	err := faker.FakeData(&bookingStay)
	if err != nil {
		log.Fatalf("failed to create BookingStayFaker: %v", err)
	}
	return bookingStay.Period
}
