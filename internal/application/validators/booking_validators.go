package validators

import (
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingValidator interface {
	Validate(b *domain.Booking) error
}

type BookingAllowedStartDateValidator struct{}

type BookingMaximumStayValidator struct{}

type BookingStartDateBeforeEndDateValidator struct{}

type ErrBookingAllowedStartDate struct{}

type ErrBookingMaximumStay struct{}

type ErrBookingStartDateBeforeEndDate struct{}

func (v *BookingAllowedStartDateValidator) Validate(b *domain.Booking) error {
	now := time.Now()
	if b.StartDate.After(now) && b.StartDate.Before(now.AddDate(0, 1, 0)) {
		return nil
	}
	return ErrBookingAllowedStartDate{}
}

func (v *BookingMaximumStayValidator) Validate(b *domain.Booking) error {
	if b.EndDate.Sub(b.StartDate).Hours()/24 <= 3 {
		return nil
	}
	return ErrBookingMaximumStay{}
}

func (v *BookingStartDateBeforeEndDateValidator) Validate(b *domain.Booking) error {
	if b.StartDate.Before(b.EndDate) {
		return nil
	}
	return ErrBookingStartDateBeforeEndDate{}
}

func (e ErrBookingAllowedStartDate) Error() string {
	return "start_date: must be from 1 day to up to 1 month ahead"
}

func (e ErrBookingMaximumStay) Error() string {
	return "maximum stay: must be less or equal to three days"
}

func (e ErrBookingStartDateBeforeEndDate) Error() string {
	return "start_date: must be before end_date"
}

func Apply(validators []BookingValidator, booking *domain.Booking) error {
	var errs []error

	for _, v := range validators {
		err := v.Validate(booking)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		var msg string
		for _, e := range errs {
			msg += "\n - " + e.Error()
		}
		return status.Errorf(codes.InvalidArgument, "validation error: %s", msg)
	}
	return nil
}
