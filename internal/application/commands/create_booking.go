package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
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
	booking := domain.Booking{}
	booking.BookingID = uuid.New().String()
	booking.CampsiteID = cmd.CampsiteID
	booking.Email = cmd.Email
	booking.FullName = cmd.FullName
	booking.StartDate = cmd.StartDate
	booking.EndDate = cmd.EndDate
	booking.Active = true

	return h.bookings.Insert(ctx, &booking)
}
