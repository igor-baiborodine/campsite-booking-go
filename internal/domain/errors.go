package domain

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

type (
	ErrBookingNotFound struct {
		BookingID string
	}

	ErrBookingDatesNotAvailable struct {
		StartDate time.Time
		EndDate   time.Time
	}

	ErrBookingAlreadyCancelled struct {
		BookingID string
	}

	ErrBookingValidation struct {
		MultiErr *multierror.Error
	}
)

func (e ErrBookingNotFound) Error() string {
	return fmt.Sprintf("booking not found for BookingID %s", e.BookingID)
}

func (e ErrBookingDatesNotAvailable) Error() string {
	return fmt.Sprintf("booking dates not available from %s to %s",
		e.StartDate.Format(time.DateOnly), e.EndDate.Format(time.DateOnly))
}

func (e ErrBookingAlreadyCancelled) Error() string {
	return fmt.Sprintf("booking already cancelled for BookingID %s", e.BookingID)
}

func (e ErrBookingValidation) Error() string {
	if e.MultiErr != nil {
		return fmt.Sprintf("booking validation: %s", e.MultiErr.Error())
	}
	return ""
}

func (e *ErrBookingValidation) Append(err error) {
	if err != nil {
		e.MultiErr = multierror.Append(e.MultiErr, err)
	}
}
