package domain

import (
	"context"
	"fmt"
	"time"
)

type ErrBookingNotFound struct {
	BookingID string
}

func (e ErrBookingNotFound) Error() string {
	return fmt.Sprintf("booking not found for BookingID %s", e.BookingID)
}

type ErrBookingDatesNotAvailable struct {
	StartDate time.Time
	EndDate   time.Time
}

func (e ErrBookingDatesNotAvailable) Error() string {
	return fmt.Sprintf("booking dates not available from %s to %s",
		e.StartDate.Format(time.DateOnly), e.EndDate.Format(time.DateOnly))
}

type BookingRepository interface {
	Find(ctx context.Context, bookingID string) (*Booking, error)
	FindForDateRange(ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time) ([]*Booking, error)
	Insert(ctx context.Context, booking *Booking) error
	Update(ctx context.Context, booking *Booking) error
}
