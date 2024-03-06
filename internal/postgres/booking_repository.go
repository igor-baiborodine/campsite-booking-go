package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/stackus/errors"
)

type BookingRepository struct {
	db *sql.DB
}

var _ domain.BookingRepository = (*BookingRepository)(nil)

func NewBookingRepository(db *sql.DB) BookingRepository {
	return BookingRepository{db: db}
}

func (r BookingRepository) Find(ctx context.Context, bookingID string) (*domain.Booking, error) {
	booking := &domain.Booking{}
	err := r.db.QueryRowContext(ctx, SelectByBookingIdFromBookings, bookingID).Scan(
		&booking.ID, &booking.BookingID, &booking.CampsiteID, &booking.Email,
		&booking.FullName, &booking.StartDate, &booking.EndDate, &booking.Active)
	if err != nil {
		return nil, errors.Wrap(err, "scanning booking")
	}
	return booking, nil
}

func (r BookingRepository) FindForDateRange(
	ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time) (
	[]*domain.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (r BookingRepository) Insert(ctx context.Context, booking *domain.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (r BookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	//TODO implement me
	panic("implement me")
}
