package commands

import (
	"context"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	UpdateBooking struct {
		BookingID  string
		CampsiteID string
		Email      string
		FullName   string
		StartDate  string
		EndDate    string
	}

	UpdateBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewUpdateBookingHandler(bookings domain.BookingRepository) UpdateBookingHandler {
	return UpdateBookingHandler{bookings: bookings}
}

func (h UpdateBookingHandler) UpdateBooking(ctx context.Context, cmd UpdateBooking) error {
	_, err := h.bookings.Find(ctx, cmd.BookingID)
	if err != nil {
		return err
	}
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

	return h.bookings.Update(ctx, &booking)
}
