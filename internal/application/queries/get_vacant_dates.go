package queries

import (
	"context"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"time"
)

type (
	GetVacantDates struct {
		CampsiteID string
		StartDate  time.Time
		EndDate    time.Time
	}

	GetVacantDatesHandler struct {
		bookings domain.BookingRepository
	}
)

func NewGetVacantDatesHandler(bookings domain.BookingRepository) GetVacantDatesHandler {
	return GetVacantDatesHandler{bookings: bookings}
}

func (h GetVacantDatesHandler) GetVacantDates(ctx context.Context, qry GetVacantDates) ([]*domain.Booking, error) {
	return h.bookings.FindForDateRange(ctx, qry.CampsiteID, qry.StartDate, qry.EndDate)
}
