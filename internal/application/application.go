package application

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/commands"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/queries"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CancelBooking(ctx context.Context, cmd commands.CancelBooking) error
		CreateBooking(ctx context.Context, cmd commands.CreateBooking) (*domain.Booking, error)
		CreateCampsite(ctx context.Context, cmd commands.CreateCampsite) error
		UpdateBooking(ctx context.Context, cmd commands.UpdateBooking) error
	}

	Queries interface {
		GetBooking(ctx context.Context, qry queries.GetBooking) (*domain.Booking, error)
		GetCampsites(ctx context.Context, _ queries.GetCampsites) ([]*domain.Campsite, error)
		GetVacantDates(ctx context.Context, qry queries.GetVacantDates) ([]string, error)
	}

	CampsitesApp struct {
		CampsitesCommands
		CampsitesQueries
	}

	CampsitesCommands struct {
		commands.CreateCampsiteHandler
		commands.CreateBookingHandler
		commands.UpdateBookingHandler
		commands.CancelBookingHandler
	}

	CampsitesQueries struct {
		queries.GetCampsitesHandler
		queries.GetBookingHandler
		queries.GetVacantDatesHandler
	}
)

var _ App = (*CampsitesApp)(nil)

func New(campsites domain.CampsiteRepository, bookings domain.BookingRepository) *CampsitesApp {
	return &CampsitesApp{
		CampsitesCommands: CampsitesCommands{
			CreateCampsiteHandler: commands.NewCreateCampsiteHandler(campsites),
			CreateBookingHandler:  commands.NewCreateBookingHandler(bookings),
			UpdateBookingHandler:  commands.NewUpdateBookingHandler(bookings),
			CancelBookingHandler:  commands.NewCancelBookingHandler(bookings),
		},
		CampsitesQueries: CampsitesQueries{
			GetCampsitesHandler:   queries.NewGetCampsitesHandler(campsites),
			GetBookingHandler:     queries.NewGetBookingHandler(bookings),
			GetVacantDatesHandler: queries.NewGetVacantDatesHandler(bookings),
		},
	}
}
