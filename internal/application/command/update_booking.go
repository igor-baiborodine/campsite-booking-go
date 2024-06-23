package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/decorator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/handler"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/validator"
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

	// UpdateBookingHandler is a logging decorator for the updateBookingHandler struct.
	UpdateBookingHandler handler.Command[UpdateBooking]

	updateBookingHandler struct {
		bookings   domain.BookingRepository
		validators []validator.BookingValidator
	}
)

func NewUpdateBookingHandler(bookings domain.BookingRepository, logger *slog.Logger) UpdateBookingHandler {
	return decorator.ApplyCommandDecorator[UpdateBooking](
		updateBookingHandler{
			bookings: bookings,
			validators: []validator.BookingValidator{
				&validator.BookingStartDateBeforeEndDateValidator{},
				&validator.BookingAllowedStartDateValidator{},
				&validator.BookingMaximumStayValidator{},
			},
		},
		logger,
	)
}

func (h updateBookingHandler) Handle(ctx context.Context, cmd UpdateBooking) error {
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

	err = validator.Apply(h.validators, booking)
	if err != nil {
		return err
	}
	return h.bookings.Update(ctx, booking)
}
