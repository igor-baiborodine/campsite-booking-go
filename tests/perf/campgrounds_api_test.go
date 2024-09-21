package perf

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkGetCampsites(b *testing.B) {
	conn := newClientConn(b, getServerAddr(b))
	defer closeConn(conn)

	client := api.NewCampgroundsServiceClient(conn)
	for x := 0; x < b.N; x++ {
		campsites, err := client.GetCampsites(context.Background(), &api.GetCampsitesRequest{})
		assert.Nil(b, err)
		assert.NotNil(b, campsites)
	}
}

func BenchmarkGetVacantDates(b *testing.B) {
	conn := newClientConn(b, getServerAddr(b))
	defer closeConn(conn)

	client := api.NewCampgroundsServiceClient(conn)
	resp, err := client.GetCampsites(context.Background(), &api.GetCampsitesRequest{})
	assert.Nil(b, err)
	campsites := resp.Campsites
	assert.True(b, len(campsites) > 0, "no campsites found")

	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, 1).Format(time.DateOnly)
	endDate := now.AddDate(0, 1, 0).Format(time.DateOnly)
	req := &api.GetVacantDatesRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)
	minIndex := 0
	maxIndex := len(campsites)

	for x := 0; x < b.N; x++ {
		b.StopTimer()
		i := r.Intn(maxIndex - minIndex)
		req.CampsiteId = campsites[i].CampsiteId
		b.StartTimer()

		resp, err := client.GetVacantDates(context.Background(), req)
		assert.Nil(b, err)
		assert.NotNil(b, resp.VacantDates)
	}
}

func newClientConn(b *testing.B, addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(b, err)
	return conn
}

func closeConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Fatalf("failed to close connection: %v", err)
	}
}

func getServerAddr(b *testing.B) string {
	addr, ok := os.LookupEnv("SERVER_ADDR")
	assert.True(b, ok)
	return addr
}
