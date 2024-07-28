package application

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/command"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/query"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/validator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	App interface {
		CreateCampsite(ctx context.Context, cmd command.CreateCampsite) error
		CreateBooking(ctx context.Context, cmd command.CreateBooking) error
		UpdateBooking(ctx context.Context, cmd command.UpdateBooking) error
		CancelBooking(ctx context.Context, cmd command.CancelBooking) error
		GetCampsites(ctx context.Context, qry query.GetCampsites) ([]*domain.Campsite, error)
		GetBooking(ctx context.Context, qry query.GetBooking) (*domain.Booking, error)
		GetVacantDates(ctx context.Context, qry query.GetVacantDates) ([]string, error)
	}

	commands struct {
		command.CreateCampsiteHandler
		command.CreateBookingHandler
		command.UpdateBookingHandler
		command.CancelBookingHandler
	}

	queries struct {
		query.GetCampsitesHandler
		query.GetBookingHandler
		query.GetVacantDatesHandler
	}

	CampgroundsApp struct {
		commands
		queries
	}
)

var bookingValidators = []domain.BookingValidator{
	validator.BookingStartDateBeforeEndDate{},
	validator.BookingAllowedStartDate{},
	validator.BookingMaximumStay{},
}

func (a CampgroundsApp) CreateCampsite(ctx context.Context, cmd command.CreateCampsite) error {
	return a.commands.CreateCampsiteHandler.Handle(ctx, cmd)
}

func (a CampgroundsApp) CreateBooking(ctx context.Context, cmd command.CreateBooking) error {
	return a.commands.CreateBookingHandler.Handle(ctx, cmd)
}

func (a CampgroundsApp) UpdateBooking(ctx context.Context, cmd command.UpdateBooking) error {
	return a.commands.UpdateBookingHandler.Handle(ctx, cmd)
}

func (a CampgroundsApp) CancelBooking(ctx context.Context, cmd command.CancelBooking) error {
	return a.commands.CancelBookingHandler.Handle(ctx, cmd)
}

func (a CampgroundsApp) GetCampsites(
	ctx context.Context,
	qry query.GetCampsites,
) ([]*domain.Campsite, error) {
	return a.queries.GetCampsitesHandler.Handle(ctx, qry)
}

func (a CampgroundsApp) GetBooking(
	ctx context.Context,
	qry query.GetBooking,
) (*domain.Booking, error) {
	return a.queries.GetBookingHandler.Handle(ctx, qry)
}

func (a CampgroundsApp) GetVacantDates(
	ctx context.Context,
	qry query.GetVacantDates,
) ([]string, error) {
	return a.queries.GetVacantDatesHandler.Handle(ctx, qry)
}

var _ App = (*CampgroundsApp)(nil)

func New(campsites domain.CampsiteRepository, bookings domain.BookingRepository) *CampgroundsApp {
	return &CampgroundsApp{
		commands: commands{
			CreateCampsiteHandler: command.NewCreateCampsiteHandler(campsites),
			CreateBookingHandler:  command.NewCreateBookingHandler(bookings, bookingValidators),
			UpdateBookingHandler:  command.NewUpdateBookingHandler(bookings, bookingValidators),
			CancelBookingHandler:  command.NewCancelBookingHandler(bookings),
		},
		queries: queries{
			GetCampsitesHandler:   query.NewGetCampsitesHandler(campsites),
			GetBookingHandler:     query.NewGetBookingHandler(bookings),
			GetVacantDatesHandler: query.NewGetVacantDatesHandler(bookings),
		},
	}
}
