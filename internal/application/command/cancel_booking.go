package command

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/decorator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/handler"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CancelBooking struct {
		BookingID string
	}

	// CancelBookingHandler is a logging decorator for the cancelBookingHandler struct.
	CancelBookingHandler handler.Command[CancelBooking]

	cancelBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewCancelBookingHandler(bookings domain.BookingRepository) CancelBookingHandler {
	return decorator.ApplyCommandDecorator[CancelBooking](cancelBookingHandler{bookings: bookings})
}

func (h cancelBookingHandler) Handle(ctx context.Context, cmd CancelBooking) error {
	booking, err := h.bookings.Find(ctx, cmd.BookingID)
	if err != nil {
		return err
	}
	if !booking.Active {
		return domain.ErrBookingAlreadyCancelled{BookingID: cmd.BookingID}
	}
	booking.Active = false

	return h.bookings.Update(ctx, booking)
}
