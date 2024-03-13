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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "begin transaction")
	}
	defer tx.Rollback()

	booking := &domain.Booking{}
	if err = tx.QueryRowContext(
		ctx, SelectByBookingIdFromBookings, bookingID,
	).Scan(
		&booking.ID, &booking.BookingID, &booking.CampsiteID, &booking.Email,
		&booking.FullName, &booking.StartDate, &booking.EndDate, &booking.Active,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "booking for ID not found: "+bookingID)
		}
		return nil, errors.Wrap(err, "scan booking row")
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return booking, nil
}

func (r BookingRepository) FindForDateRange(
	ctx context.Context, campsiteID string, startDate time.Time, endDate time.Time,
) ([]*domain.Booking, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "begin transaction")
	}
	defer tx.Rollback()

	bookings, err := r.findForDateRange(ctx, tx, campsiteID, startDate, endDate)
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return bookings, nil
}

func (r BookingRepository) findForDateRange(
	ctx context.Context, tx *sql.Tx, campsiteID string, startDate time.Time, endDate time.Time,
) ([]*domain.Booking, error) {
	rows, err := tx.QueryContext(ctx, FindForDateRangeInBookings, campsiteID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "query bookings for date range")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "close booking rows")
		}
	}(rows)

	var bookings []*domain.Booking
	for rows.Next() {
		booking := &domain.Booking{}
		if err = rows.Scan(
			&booking.ID, &booking.BookingID, &booking.CampsiteID, &booking.Email,
			&booking.FullName, &booking.StartDate, &booking.EndDate, &booking.Active,
		); err != nil {
			return nil, errors.Wrap(err, "scan booking row")
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finish booking rows")
	}
	return bookings, nil
}

func (r BookingRepository) Insert(ctx context.Context, booking *domain.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (r BookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	//TODO implement me
	panic("implement me")
}
