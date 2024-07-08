package domain

import (
	"encoding/json"
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

func (b *Booking) String() string {
	result, _ := json.Marshal(b)
	return string(result)
}
