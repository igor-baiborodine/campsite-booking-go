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

type ErrBookingAllowedStartDate struct{}

func (v *BookingAllowedStartDateValidator) Validate(b *domain.Booking) error {
	now := time.Now()
	if now.Before(b.StartDate) && b.StartDate.Before(now.AddDate(0, 1, 1)) {
		return &ErrBookingAllowedStartDate{}
	}
	return nil
}

func (e ErrBookingAllowedStartDate) Error() string {
	return "start_date: must be from 1 day to up to 1 month ahead"
}

func Apply(booking *domain.Booking, validators []BookingValidator) error {
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
