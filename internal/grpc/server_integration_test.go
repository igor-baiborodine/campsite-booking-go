//go:build integration

package grpc_test

import (
	"context"
	"log/slog"
	"net"
	"testing"
	"time"

	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	rpc "github.com/igor-baiborodine/campsite-booking-go/internal/grpc"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type mocks struct {
	campsites *domain.MockCampsiteRepository
	bookings  *domain.MockBookingRepository
}

type serverSuite struct {
	server *grpc.Server
	client api.CampgroundsServiceClient
	mocks
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, &serverSuite{})
}

func (s *serverSuite) SetupSuite() {}

func (s *serverSuite) TearDownSuite() {}

func (s *serverSuite) SetupTest() {
	l := logger.New(logger.LogConfig{
		Environment: "integration",
		LogLevel:    "DEBUG",
	})
	slog.SetDefault(l)

	var err error
	s.server, err = rpc.NewServer()
	if err != nil {
		s.T().Fatal(err)
	}

	var listener net.Listener
	const grpcTestPort = ":10912"
	listener, err = net.Listen("tcp", grpcTestPort)
	if err != nil {
		s.T().Fatal(err)
	}

	s.mocks = mocks{
		campsites: domain.NewMockCampsiteRepository(s.T()),
		bookings:  domain.NewMockBookingRepository(s.T()),
	}
	app := application.New(s.mocks.campsites, s.mocks.bookings)

	if err = rpc.RegisterServer(app, s.server); err != nil {
		s.T().Fatal(err)
	}
	go func(listener net.Listener) {
		err := s.server.Serve(listener)
		if err != nil {
			s.T().Fatal(err)
		}
	}(listener)

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(grpcTestPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = api.NewCampgroundsServiceClient(conn)
}

func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}

func (s *serverSuite) TestCampgroundsService_CreateCampsite() {
	tests := map[string]struct {
		req     *api.CreateCampsiteRequest
		on      func(f mocks)
		want    *api.CreateCampsiteResponse
		wantErr string
	}{
		"Success": {
			req: &api.CreateCampsiteRequest{
				CampsiteCode:  "campsite-code",
				Capacity:      1,
				DrinkingWater: true,
				Restrooms:     true,
				PicnicTable:   true,
				FirePit:       true,
			},
			on: func(f mocks) {
				s.mocks.campsites.On(
					"Insert", mock.Anything, mock.AnythingOfType("*domain.Campsite"),
				).Return(nil)
			},
			want:    nil,
			wantErr: "",
		},
		"InvalidArgument_CampsiteCode": {
			req: &api.CreateCampsiteRequest{
				CampsiteCode: "",
				Capacity:     1,
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_Capacity": {
			req: &api.CreateCampsiteRequest{
				CampsiteCode: "campsite-code",
				Capacity:     0,
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
	}
	for name, tc := range tests {
		s.T().Run(name, func(t *testing.T) {
			// given
			if tc.on != nil {
				tc.on(s.mocks)
			}
			// when
			resp, err := s.client.CreateCampsite(context.Background(), tc.req)
			// then
			if tc.wantErr != "" {
				s.Empty(resp)
				assert.Containsf(
					t,
					err.Error(),
					tc.wantErr,
					"CreateCampsite() error = %v, wantErr %v",
					err,
					tc.wantErr,
				)
				return
			}
			s.NotEmpty(resp.CampsiteId)
		})
	}
}

func (s *serverSuite) TestCampgroundsService_GetBooking() {
	booking, err := bootstrap.NewBooking("campsite-id")
	s.NoError(err)

	tests := map[string]struct {
		req     *api.GetBookingRequest
		on      func(f mocks)
		want    *api.GetCampsitesResponse
		wantErr string
	}{
		"Success": {
			req: &api.GetBookingRequest{BookingId: "ddd1b17a-9baf-403a-98ac-ef36b643c04d"},
			on: func(f mocks) {
				s.mocks.bookings.On(
					"Find", mock.Anything, mock.AnythingOfType("string"),
				).Return(booking, nil)
			},
			want:    nil,
			wantErr: "",
		},
		"InvalidArgument_BookingId": {
			req:     &api.GetBookingRequest{BookingId: "invalid-uuid-booking-id"},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
	}
	for name, tc := range tests {
		s.T().Run(name, func(t *testing.T) {
			// given
			if tc.on != nil {
				tc.on(s.mocks)
			}
			// when
			resp, err := s.client.GetBooking(context.Background(), tc.req)
			// then
			if tc.wantErr != "" {
				s.Empty(resp)
				assert.Containsf(
					t,
					err.Error(),
					tc.wantErr,
					"GetBooking() error = %v, wantErr %v",
					err,
					tc.wantErr,
				)
				return
			}
			s.Assert().Equal(rpc.BookingFromDomain(booking), resp.Booking)
		})
	}
}

func (s *serverSuite) TestCampgroundsService_CreateBooking() {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	tests := map[string]struct {
		req     *api.CreateBookingRequest
		on      func(f mocks)
		want    *api.CreateBookingResponse
		wantErr string
	}{
		"Success": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "john.smith@example.com",
				FullName:   "John Smith",
				StartDate:  now.AddDate(0, 0, 1).Format(time.DateOnly),
				EndDate:    now.AddDate(0, 0, 2).Format(time.DateOnly),
			},
			on: func(f mocks) {
				s.mocks.bookings.On(
					"Insert", mock.Anything, mock.AnythingOfType("*domain.Booking"),
				).Return(nil)
			},
			want:    nil,
			wantErr: "",
		},
		"InvalidArgument_CampsiteId": {
			req: &api.CreateBookingRequest{
				CampsiteId: "invalid-uuid-campsite-id",
				Email:      "john.smith@example.com",
				FullName:   "John Smith",
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_Email": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "invalid-email",
				FullName:   "John Smith",
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_FullName": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "john.smith@example.com",
				FullName:   "",
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_StartDate": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "john.smith@example.com",
				FullName:   "John Smith",
				StartDate:  "9999-99-99",
				EndDate:    "2006-01-03",
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_EndDate": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "john.smith@example.com",
				FullName:   "",
				StartDate:  "2006-01-02",
				EndDate:    "9999-99-99",
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
		"InvalidArgument_BookingStartDateBeforeEndDateValidator": {
			req: &api.CreateBookingRequest{
				CampsiteId: "b5839e4a-1dab-4c0a-8aa5-6a4e6910ce46",
				Email:      "john.smith@example.com",
				FullName:   "John Smith",
				StartDate:  now.AddDate(0, 0, 2).Format(time.DateOnly),
				EndDate:    now.AddDate(0, 0, 1).Format(time.DateOnly),
			},
			on:      nil,
			want:    nil,
			wantErr: codes.InvalidArgument.String(),
		},
	}
	for name, tc := range tests {
		s.T().Run(name, func(t *testing.T) {
			// given
			if tc.on != nil {
				tc.on(s.mocks)
			}
			// when
			resp, err := s.client.CreateBooking(context.Background(), tc.req)
			// then
			if tc.wantErr != "" {
				s.Empty(resp)
				assert.Containsf(
					t,
					err.Error(),
					tc.wantErr,
					"CreateBooking() error = %v, wantErr %v",
					err,
					tc.wantErr,
				)
				return
			}
			s.NotEmpty(resp.BookingId)
		})
	}
}
