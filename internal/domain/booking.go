package domain

import "time"

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

func (b *Booking) Cancel() (err error) {
	// TODO: add error handling for already canceled booking
	b.Active = false
	return nil
}
