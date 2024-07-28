package validator

import (
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type BookingAllowedStartDateValidator struct{}

type BookingMaximumStayValidator struct{}

type BookingStartDateBeforeEndDateValidator struct{}

type ErrBookingAllowedStartDate struct{}

type ErrBookingMaximumStay struct{}

type ErrBookingStartDateBeforeEndDate struct{}

func (v BookingAllowedStartDateValidator) Validate(b *domain.Booking) error {
	now := time.Now()
	if b.StartDate.After(now) && b.StartDate.Before(now.AddDate(0, 1, 0)) {
		return nil
	}
	return ErrBookingAllowedStartDate{}
}

func (v BookingMaximumStayValidator) Validate(b *domain.Booking) error {
	if b.EndDate.Sub(b.StartDate).Hours()/24 <= 3 {
		return nil
	}
	return ErrBookingMaximumStay{}
}

func (v BookingStartDateBeforeEndDateValidator) Validate(b *domain.Booking) error {
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

func Apply(validators []domain.BookingValidator, booking *domain.Booking) error {
	merr := domain.ErrBookingValidation{}

	for _, v := range validators {
		if err := v.Validate(booking); err != nil {
			merr.Append(err)
		}
	}
	if merr.MultiErr.ErrorOrNil() != nil {
		return merr
	}
	return nil
}
