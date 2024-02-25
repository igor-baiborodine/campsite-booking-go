package domain

import (
	"context"
)

type BookingRepository interface {
	Find(ctx context.Context, bookingID string) (*Booking, error)
	FindForDateRange(ctx context.Context, startDate string, endDate string) ([]*Booking, error)
	Save(ctx context.Context, order *Booking) error
	Update(ctx context.Context, order *Booking) error
}
