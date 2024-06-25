package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/command"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/query"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	app application.App
	api.UnimplementedCampgroundsServiceServer
}

var _ api.CampgroundsServiceServer = (*server)(nil)

func NewServer(l *slog.Logger) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logServiceCalls(l)),
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
	}
	return grpc.NewServer(opts...), nil
}

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	api.RegisterCampgroundsServiceServer(registrar, server{app: app})
	return nil
}

func (s server) GetCampsites(
	ctx context.Context,
	_ *api.GetCampsitesRequest,
) (*api.GetCampsitesResponse, error) {
	campsites, err := s.app.GetCampsites(ctx, query.GetCampsites{})
	if err != nil {
		return nil, err
	}

	var protoCampsites []*api.Campsite
	for _, campsite := range campsites {
		protoCampsites = append(protoCampsites, CampsiteFromDomain(campsite))
	}

	return &api.GetCampsitesResponse{
		Campsites: protoCampsites,
	}, nil
}

func (s server) CreateCampsite(
	ctx context.Context,
	req *api.CreateCampsiteRequest,
) (*api.CreateCampsiteResponse, error) {
	campsite := command.CreateCampsite{
		CampsiteID:    uuid.New().String(),
		CampsiteCode:  req.CampsiteCode,
		Capacity:      req.Capacity,
		DrinkingWater: req.DrinkingWater,
		Restrooms:     req.Restrooms,
		PicnicTable:   req.PicnicTable,
		FirePit:       req.FirePit,
	}
	err := s.app.CreateCampsite(ctx, campsite)
	if err != nil {
		return nil, handleDomainError(err)
	}

	return &api.CreateCampsiteResponse{
		CampsiteId: campsite.CampsiteID,
	}, nil
}

func (s server) GetBooking(
	ctx context.Context,
	req *api.GetBookingRequest,
) (*api.GetBookingResponse, error) {
	booking, err := s.app.GetBooking(ctx, query.GetBooking{BookingID: req.BookingId})
	if err != nil {
		return nil, handleDomainError(err)
	}

	return &api.GetBookingResponse{
		Booking: BookingFromDomain(booking),
	}, nil
}

func (s server) CreateBooking(
	ctx context.Context,
	req *api.CreateBookingRequest,
) (*api.CreateBookingResponse, error) {
	booking := command.CreateBooking{
		BookingID:  uuid.New().String(),
		CampsiteID: req.CampsiteId,
		Email:      req.Email,
		FullName:   req.FullName,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}
	err := s.app.CreateBooking(ctx, booking)
	if err != nil {
		return nil, handleDomainError(err)
	}

	return &api.CreateBookingResponse{
		BookingId: booking.BookingID,
	}, nil
}

func (s server) UpdateBooking(
	ctx context.Context,
	req *api.UpdateBookingRequest,
) (*api.UpdateBookingResponse, error) {
	booking := command.UpdateBooking{
		BookingID:  req.Booking.BookingId,
		CampsiteID: req.Booking.CampsiteId,
		Email:      req.Booking.Email,
		FullName:   req.Booking.FullName,
		StartDate:  req.Booking.StartDate,
		EndDate:    req.Booking.EndDate,
	}
	err := s.app.UpdateBooking(ctx, booking)
	if err != nil {
		return nil, handleDomainError(err)
	}
	return &api.UpdateBookingResponse{}, nil
}

func (s server) CancelBooking(
	ctx context.Context,
	req *api.CancelBookingRequest,
) (*api.CancelBookingResponse, error) {
	booking := command.CancelBooking{
		BookingID: req.GetBookingId(),
	}
	err := s.app.CancelBooking(ctx, booking)
	if err != nil {
		return nil, handleDomainError(err)
	}
	return &api.CancelBookingResponse{}, nil
}

func (s server) GetVacantDates(
	ctx context.Context,
	req *api.GetVacantDatesRequest,
) (*api.GetVacantDatesResponse, error) {
	vacantDates, err := s.app.GetVacantDates(ctx, query.GetVacantDates{
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

func CampsiteFromDomain(campsite *domain.Campsite) *api.Campsite {
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

func BookingFromDomain(booking *domain.Booking) *api.Booking {
	return &api.Booking{
		BookingId:  booking.BookingID,
		CampsiteId: booking.CampsiteID,
		Email:      booking.Email,
		FullName:   booking.FullName,
		StartDate:  booking.StartDate.Format(time.DateOnly),
		EndDate:    booking.EndDate.Format(time.DateOnly),
		Active:     booking.Active,
	}
}

func handleDomainError(e error) error {
	switch e.(type) {
	case domain.ErrBookingNotFound:
		return status.Error(codes.NotFound, e.Error())
	case domain.ErrBookingAlreadyCancelled, domain.ErrBookingDatesNotAvailable:
		return status.Error(codes.FailedPrecondition, e.Error())
	default:
		return e
	}
}
