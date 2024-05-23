package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	api "github.com/igor-baiborodine/campsite-booking-go/campsitespb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/commands"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/queries"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	api.UnimplementedCampsitesServiceServer
}

var _ api.CampsitesServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	api.RegisterCampsitesServiceServer(registrar, server{app: app})
	return nil
}

func (s server) GetCampsites(ctx context.Context, _ *api.GetCampsitesRequest) (*api.GetCampsitesResponse, error) {
	campsites, err := s.app.GetCampsites(ctx, queries.GetCampsites{})
	if err != nil {
		return nil, err
	}

	var protoCampsites []*api.Campsite
	for _, campsite := range campsites {
		protoCampsites = append(protoCampsites, s.campsiteFromDomain(campsite))
	}

	return &api.GetCampsitesResponse{
		Campsites: protoCampsites,
	}, nil
}

func (s server) CreateCampsite(ctx context.Context, req *api.CreateCampsiteRequest) (*api.CreateCampsiteResponse, error) {
	campsiteID := uuid.New().String()

	err := s.app.CreateCampsite(ctx, commands.CreateCampsite{
		CampsiteId:    campsiteID,
		CampsiteCode:  req.CampsiteCode,
		Capacity:      req.Capacity,
		DrinkingWater: req.DrinkingWater,
		Restrooms:     req.Restrooms,
		PicnicTable:   req.PicnicTable,
		FirePit:       req.FirePit,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateCampsiteResponse{
		CampsiteId: campsiteID,
	}, nil
}

func (s server) GetBooking(ctx context.Context, req *api.GetBookingRequest) (*api.GetBookingResponse, error) {
	booking, err := s.app.GetBooking(ctx, queries.GetBooking{BookingID: req.BookingId})
	if err != nil {
		return nil, err
	}

	return &api.GetBookingResponse{
		Booking: s.bookingFromDomain(booking),
	}, nil
}

func (s server) CreateBooking(ctx context.Context, req *api.CreateBookingRequest) (*api.CreateBookingResponse, error) {
	bookingID := uuid.New().String()

	err := s.app.CreateBooking(ctx, commands.CreateBooking{
		BookingID:  bookingID,
		CampsiteID: req.CampsiteId,
		Email:      req.Email,
		FullName:   req.FullName,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateBookingResponse{
		BookingId: bookingID,
	}, nil
}

func (s server) UpdateBooking(ctx context.Context, req *api.UpdateBookingRequest) (*api.UpdateBookingResponse, error) {
	err := s.app.UpdateBooking(ctx, commands.UpdateBooking{
		BookingID:  req.Booking.BookingId,
		CampsiteID: req.Booking.CampsiteId,
		Email:      req.Booking.Email,
		FullName:   req.Booking.FullName,
		StartDate:  req.Booking.StartDate,
		EndDate:    req.Booking.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return &api.UpdateBookingResponse{}, nil
}

func (s server) CancelBooking(ctx context.Context, req *api.CancelBookingRequest) (*api.CancelBookingResponse, error) {
	err := s.app.CancelBooking(ctx, commands.CancelBooking{
		BookingID: req.GetBookingId(),
	})
	if err != nil {
		return nil, err
	}
	return &api.CancelBookingResponse{}, nil
}

func (s server) GetVacantDates(ctx context.Context, req *api.GetVacantDatesRequest) (*api.GetVacantDatesResponse, error) {
	vacantDates, err := s.app.GetVacantDates(ctx, queries.GetVacantDates{
		CampsiteID: req.CampsiteId,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	return &api.GetVacantDatesResponse{
		VacantDates: vacantDates,
	}, nil
}

func (s server) campsiteFromDomain(campsite *domain.Campsite) *api.Campsite {
	return &api.Campsite{
		CampsiteId:    campsite.CampsiteID,
		CampsiteCode:  campsite.CampsiteCode,
		Capacity:      campsite.Capacity,
		DrinkingWater: campsite.DrinkingWater,
		Restrooms:     campsite.Restrooms,
		PicnicTable:   campsite.PicnicTable,
		FirePit:       campsite.FirePit,
		Active:        campsite.Active,
	}
}

func (s server) bookingFromDomain(booking *domain.Booking) *api.Booking {
	return &api.Booking{
		BookingId:  booking.BookingID,
		CampsiteId: booking.CampsiteID,
		Email:      booking.Email,
		FullName:   booking.FullName,
		StartDate:  booking.StartDate.Format(time.DateOnly),
		EndDate:    booking.EndDate.Format(time.DateOnly),
	}
}
