package domain

import (
	"context"
	"time"
)

type BookingRepository interface {
	Find(ctx context.Context, bookingID string) (*Booking, error)
	FindForDateRange(
		ctx context.Context,
		campsiteID string,
		startDate time.Time,
		endDate time.Time,
	) ([]*Booking, error)
	Insert(ctx context.Context, booking *Booking) error
	Update(ctx context.Context, booking *Booking) error
}
