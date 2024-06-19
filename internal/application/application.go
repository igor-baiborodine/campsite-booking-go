package application

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/command"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/query"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CancelBooking(ctx context.Context, cmd command.CancelBooking) error
		CreateBooking(ctx context.Context, cmd command.CreateBooking) error
		CreateCampsite(ctx context.Context, cmd command.CreateCampsite) error
		UpdateBooking(ctx context.Context, cmd command.UpdateBooking) error
	}

	Queries interface {
		GetBooking(ctx context.Context, qry query.GetBooking) (*domain.Booking, error)
		GetCampsites(ctx context.Context, _ query.GetCampsites) ([]*domain.Campsite, error)
		GetVacantDates(ctx context.Context, qry query.GetVacantDates) ([]string, error)
	}

	CampsitesApp struct {
		CampsitesCommands
		CampsitesQueries
	}

	CampsitesCommands struct {
		command.CreateCampsiteHandler
		command.CreateBookingHandler
		command.UpdateBookingHandler
		command.CancelBookingHandler
	}

	CampsitesQueries struct {
		query.GetCampsitesHandler
		query.GetBookingHandler
		query.GetVacantDatesHandler
	}
)

var _ App = (*CampsitesApp)(nil)

func New(campsites domain.CampsiteRepository, bookings domain.BookingRepository) *CampsitesApp {
	return &CampsitesApp{
		CampsitesCommands: CampsitesCommands{
			CreateCampsiteHandler: command.NewCreateCampsiteHandler(campsites),
			CreateBookingHandler:  command.NewCreateBookingHandler(bookings),
			UpdateBookingHandler:  command.NewUpdateBookingHandler(bookings),
			CancelBookingHandler:  command.NewCancelBookingHandler(bookings),
		},
		CampsitesQueries: CampsitesQueries{
			GetCampsitesHandler:   query.NewGetCampsitesHandler(campsites),
			GetBookingHandler:     query.NewGetBookingHandler(bookings),
			GetVacantDatesHandler: query.NewGetVacantDatesHandler(bookings),
		},
	}
}
