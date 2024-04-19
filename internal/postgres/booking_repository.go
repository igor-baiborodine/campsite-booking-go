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
		return nil, errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	booking := &domain.Booking{}
	if err = tx.QueryRowContext(
		ctx, FindBookingByBookingIdQuery, bookingID,
	).Scan(
		&booking.ID, &booking.BookingID, &booking.CampsiteID, &booking.Email,
		&booking.FullName, &booking.StartDate, &booking.EndDate, &booking.Active,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookingNotFound{BookingID: bookingID}
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
		return nil, errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	bookings, err := r.findForDateRangeWithTx(
		ctx, tx, FindAllBookingsForDateRangeQuery, campsiteID, startDate, endDate)
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return bookings, nil
}

func (r BookingRepository) Insert(ctx context.Context, booking *domain.Booking) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	query := FindAllBookingsForDateRangeQuery + " FOR UPDATE"
	bookings, err := r.findForDateRangeWithTx(
		ctx, tx, query, booking.CampsiteID, booking.StartDate, booking.EndDate)
	if err != nil {
		return errors.Wrap(err, "query bookings for date range")
	}
	if len(bookings) > 0 {
		return domain.ErrBookingDatesNotAvailable{
			StartDate: booking.StartDate,
			EndDate:   booking.EndDate,
		}
	}

	createdAt := time.Now()
	_, err = tx.ExecContext(
		ctx, InsertBookingQuery, booking.BookingID, booking.CampsiteID, booking.Email,
		booking.FullName, booking.StartDate, booking.EndDate, booking.Active, createdAt, createdAt,
	)
	if err != nil {
		return errors.Wrap(err, "insert booking")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}

func (r BookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	query := FindAllBookingsForDateRangeQuery + "FOR UPDATE"
	bookings, err := r.findForDateRangeWithTx(
		ctx, tx, query, booking.CampsiteID, booking.StartDate, booking.EndDate)
	if err != nil {
		return errors.Wrap(err, "query bookings for date range")
	}

	for _, b := range bookings {
		if b.BookingID != booking.BookingID {
			return domain.ErrBookingDatesNotAvailable{
				StartDate: booking.StartDate,
				EndDate:   booking.EndDate,
			}
		}
	}
	updatedAt := time.Now()
	_, err = tx.ExecContext(
		ctx, UpdateBookingQuery, booking.BookingID, booking.CampsiteID, booking.Email,
		booking.FullName, booking.StartDate, booking.EndDate, booking.Active, updatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "update booking")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}

func (r BookingRepository) findForDateRangeWithTx(
	ctx context.Context, tx *sql.Tx, query string, campsiteID string, startDate time.Time,
	endDate time.Time,
) ([]*domain.Booking, error) {
	rows, err := tx.QueryContext(ctx, query, campsiteID, startDate, endDate)
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
