package query

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	GetBooking struct {
		BookingID string
	}

	GetBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewGetBookingHandler(bookings domain.BookingRepository) GetBookingHandler {
	return GetBookingHandler{bookings: bookings}
}

func (h GetBookingHandler) GetBooking(ctx context.Context, qry GetBooking) (*domain.Booking, error) {
	return h.bookings.Find(ctx, qry.BookingID)
}
