package domain

import (
	"fmt"
	"time"
)

type ErrBookingNotFound struct {
	BookingID string
}

type ErrBookingDatesNotAvailable struct {
	StartDate time.Time
	EndDate   time.Time
}

type ErrBookingAlreadyCancelled struct {
	BookingID string
}

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
