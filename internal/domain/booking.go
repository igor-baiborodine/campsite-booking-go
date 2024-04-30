package domain

import (
	"fmt"
	"time"
)

type Booking struct {
	// Persistence ID
	ID int64
	// Business ID
	BookingID  string
	CampsiteID string
	Email      string
	FullName   string
	StartDate  time.Time
	EndDate    time.Time
	Active     bool
}

func (b *Booking) BookingDates() []time.Time {
	var dates []time.Time
	for d := b.StartDate; d.Before(b.EndDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

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

type ErrBookingAlreadyCancelled struct {
	BookingID string
}

func (e ErrBookingAlreadyCancelled) Error() string {
	return fmt.Sprintf("booking already cancelled for BookingID %s", e.BookingID)
}
