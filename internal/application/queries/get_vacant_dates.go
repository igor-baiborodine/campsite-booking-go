package queries

import (
	"context"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/stackus/errors"
)

type (
	GetVacantDates struct {
		CampsiteID string
		StartDate  string
		EndDate    string
	}

	GetVacantDatesHandler struct {
		bookings domain.BookingRepository
	}
)

func NewGetVacantDatesHandler(bookings domain.BookingRepository) GetVacantDatesHandler {
	return GetVacantDatesHandler{bookings: bookings}
}

func (h GetVacantDatesHandler) GetVacantDates(ctx context.Context, qry GetVacantDates) ([]string, error) {
	startDate, err := time.Parse(time.DateOnly, qry.StartDate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse start date %s", qry.StartDate)
	}

	endDate, err := time.Parse(time.DateOnly, qry.EndDate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse end date %s", qry.EndDate)
	}

	datesForRange := make(map[time.Time]bool)
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		datesForRange[date] = true
	}

	bookings, err := h.bookings.FindForDateRange(ctx, qry.CampsiteID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, booking := range bookings {
		for _, bookingDate := range booking.BookingDates() {
			delete(datesForRange, bookingDate)
		}
	}

	var vacantDates []string
	for d := range datesForRange {
		vacantDates = append(vacantDates, d.Format(time.DateOnly))
	}
	return vacantDates, nil
}
