package query

import (
	"context"
	"log/slog"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/decorator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/handler"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	GetBooking struct {
		BookingID string
	}

	// GetBookingHandler is a logging decorator for the getBookingHandler struct.
	GetBookingHandler handler.Query[GetBooking, *domain.Booking]

	getBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewGetBookingHandler(bookings domain.BookingRepository, logger *slog.Logger) GetBookingHandler {
	return decorator.ApplyQueryDecorator[GetBooking, *domain.Booking](
		getBookingHandler{bookings: bookings},
		logger,
	)
}

func (h getBookingHandler) Handle(ctx context.Context, qry GetBooking) (*domain.Booking, error) {
	return h.bookings.Find(ctx, qry.BookingID)
}
