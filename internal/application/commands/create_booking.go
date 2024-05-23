package commands

import (
	"context"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CreateBooking struct {
		BookingID  string
		CampsiteID string
		Email      string
		FullName   string
		StartDate  string
		EndDate    string
	}

	CreateBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewCreateBookingHandler(bookings domain.BookingRepository) CreateBookingHandler {
	return CreateBookingHandler{bookings: bookings}
}

func (h CreateBookingHandler) CreateBooking(ctx context.Context, cmd CreateBooking) error {
	booking := domain.Booking{}
	booking.BookingID = cmd.BookingID
	booking.CampsiteID = cmd.CampsiteID
	booking.Email = cmd.Email
	booking.FullName = cmd.FullName

	startDate, err := time.Parse(time.DateOnly, cmd.StartDate)
	if err != nil {
		return err
	}
	booking.StartDate = startDate

	endDate, err := time.Parse(time.DateOnly, cmd.EndDate)
	if err != nil {
		return err
	}
	booking.EndDate = endDate
	booking.Active = true

	return h.bookings.Insert(ctx, &booking)
}
