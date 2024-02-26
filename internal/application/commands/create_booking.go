package commands

import (
	"context"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"time"
)

type (
	CreateBooking struct {
		CampsiteID string
		Email      string
		FullName   string
		StartDate  time.Time
		EndDate    time.Time
	}

	CreateBookingHandler struct {
		bookings domain.BookingRepository
	}
)

func NewCreateBookingHandler(bookings domain.BookingRepository) CreateBookingHandler {
	return CreateBookingHandler{bookings: bookings}
}

func (h CreateBookingHandler) CreateBooking(ctx context.Context, cmd CreateBooking) error {
	bookingBuilder := domain.NewBookingBuilder().
		CampsiteID(cmd.CampsiteID).
		Email(cmd.Email).
		FullName(cmd.FullName).
		StartDate(cmd.StartDate).
		EndDate(cmd.EndDate).
		Active(true)

	booking, err := bookingBuilder.Build()
	if err != nil {
		return err
	}
	return h.bookings.Insert(ctx, booking)
}
