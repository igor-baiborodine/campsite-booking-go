package commands

import (
	"context"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"time"
)

type (
	UpdateBooking struct {
		BookingID  string
		CampsiteID string
		Email      string
		FullName   string
		StartDate  time.Time
		EndDate    time.Time
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
	bookingBuilder := domain.NewBookingBuilder().
		BookingID(cmd.BookingID).
		CampsiteID(cmd.CampsiteID).
		Email(cmd.Email).
		FullName(cmd.FullName).
		StartDate(cmd.StartDate).
		EndDate(cmd.EndDate)

	booking, err := bookingBuilder.Build()
	if err != nil {
		return err
	}
	return h.bookings.Update(ctx, booking)
}
