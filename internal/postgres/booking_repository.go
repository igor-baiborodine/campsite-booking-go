package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/jackc/pgconn"
	"github.com/stackus/errors"
)

type BookingRepository struct {
	db *sql.DB
}

var _ domain.BookingRepository = (*BookingRepository)(nil)

func NewBookingRepository(db *sql.DB) BookingRepository {
	return BookingRepository{db}
}

func (r BookingRepository) Find(ctx context.Context, bookingID string) (*domain.Booking, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "begin transaction")
	}
	defer rollbackTx(tx)

	booking := &domain.Booking{}
	if err = tx.QueryRowContext(
		ctx, queries.FindBookingByBookingID, bookingID,
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
	defer rollbackTx(tx)

	bookings, err := r.findForDateRangeWithTx(
		ctx, tx, queries.FindAllBookingsForDateRange, campsiteID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return bookings, nil
}

func (r BookingRepository) Insert(ctx context.Context, booking *domain.Booking) error {
	const (
		maxAttempts = 2
		backoffBase = 500 * time.Millisecond
	)
	txName := "insert booking"

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := r.insertTx(ctx, booking)
		if err == nil {
			return nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "40001" { // serialization failure
				backoff := backoffBase * time.Duration(attempt)
				slog.Warn(
					"failed to execute transaction due to serialization error",
					"tx_name",
					txName,
					"attempt",
					attempt,
					"retry_in_ms",
					backoff.Milliseconds(),
				)
				time.Sleep(backoff)
				continue
			}
		}
		return err
	}
	return fmt.Errorf("%s: exhaust retries after %d attempts", txName, maxAttempts)
}

func (r BookingRepository) insertTx(ctx context.Context, booking *domain.Booking) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer rollbackTx(tx)

	query := queries.FindAllBookingsForDateRange + " FOR UPDATE"
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

	_, err = tx.ExecContext(
		ctx, queries.InsertBooking, booking.BookingID, booking.CampsiteID, booking.Email,
		booking.FullName, booking.StartDate, booking.EndDate, booking.Active,
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
	defer rollbackTx(tx)

	query := queries.FindAllBookingsForDateRange + "FOR UPDATE"
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
	_, err = tx.ExecContext(
		ctx, queries.UpdateBooking, booking.BookingID, booking.CampsiteID, booking.Email,
		booking.FullName, booking.StartDate, booking.EndDate, booking.Active,
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
) (bookings []*domain.Booking, err error) {
	rows, err := tx.QueryContext(ctx, query, campsiteID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "query bookings for date range")
	}
	defer closeRows(rows)

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
