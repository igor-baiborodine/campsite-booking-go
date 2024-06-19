package command

import (
	"context"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/validator"
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
		bookings   domain.BookingRepository
		validators []validator.BookingValidator
	}
)

func NewCreateBookingHandler(bookings domain.BookingRepository) CreateBookingHandler {
	return CreateBookingHandler{
		bookings: bookings,
		validators: []validator.BookingValidator{
			&validator.BookingStartDateBeforeEndDateValidator{},
			&validator.BookingAllowedStartDateValidator{},
			&validator.BookingMaximumStayValidator{},
		},
	}
}

func (h CreateBookingHandler) CreateBooking(ctx context.Context, cmd CreateBooking) error {
	booking := &domain.Booking{
		BookingID:  cmd.BookingID,
		CampsiteID: cmd.CampsiteID,
		Email:      cmd.Email,
		FullName:   cmd.FullName,
	}
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

	err = validator.Apply(h.validators, booking)
	if err != nil {
		return err
	}
	return h.bookings.Insert(ctx, booking)
}
