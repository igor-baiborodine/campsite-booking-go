package commands

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CancelBooking struct {
		BookingID string
	}

	CancelBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewCancelBookingHandler(bookings domain.BookingRepository) CancelBookingHandler {
	return CancelBookingHandler{bookings: bookings}
}

func (h CancelBookingHandler) CancelBooking(ctx context.Context, cmd CancelBooking) error {
	booking, err := h.bookings.Find(ctx, cmd.BookingID)
	if err != nil {
		return err
	}
	if err := booking.Cancel(); err != nil {
		return err
	}
	return h.bookings.Update(ctx, booking)
}
