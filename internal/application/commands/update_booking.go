package commands

import (
	"context"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/validators"
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
		bookings   domain.BookingRepository
		validators []validators.BookingValidator
	}
)

func NewUpdateBookingHandler(bookings domain.BookingRepository) UpdateBookingHandler {
	return UpdateBookingHandler{
		bookings: bookings,
		validators: []validators.BookingValidator{
			&validators.BookingAllowedStartDateValidator{},
		},
	}
}

func (h UpdateBookingHandler) UpdateBooking(ctx context.Context, cmd UpdateBooking) error {
	booking, err := h.bookings.Find(ctx, cmd.BookingID)
	if err != nil {
		return err
	}
	if !booking.Active {
		return domain.ErrBookingAlreadyCancelled{BookingID: cmd.BookingID}
	}

	if cmd.CampsiteID != "" {
		booking.CampsiteID = cmd.CampsiteID
	}
	if cmd.Email != "" {
		booking.Email = cmd.Email
	}
	if cmd.FullName != "" {
		booking.FullName = cmd.FullName
	}

	if cmd.StartDate != "" && cmd.EndDate != "" {
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
	}

	err = validators.ApplyValidators(booking, h.validators)
	if err != nil {
		return err
	}
	return h.bookings.Update(ctx, booking)
}
