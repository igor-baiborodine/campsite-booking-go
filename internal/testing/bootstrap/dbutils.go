package bootstrap

import (
	"context"
	"database/sql"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/stackus/errors"
)

const (
	deleteBookingsQuery = `
		DELETE FROM bookings
	`
	deleteCampsitesQuery = `
		DELETE FROM campsites
	`
)

func InsertCampsite(db *sql.DB, c *domain.Campsite) error {
	_, err := db.ExecContext(context.Background(), queries.InsertCampsite,
		c.CampsiteID, c.CampsiteCode, c.Capacity, c.Restrooms, c.DrinkingWater, c.PicnicTable,
		c.FirePit, c.Active,
	)
	return err
}

func InsertBooking(db *sql.DB, b *domain.Booking) error {
	_, err := db.ExecContext(context.Background(), queries.InsertBooking,
		b.BookingID, b.CampsiteID, b.Email, b.FullName, b.StartDate, b.EndDate, b.Active,
	)
	return err
}

func FindBooking(db *sql.DB, bookingID string) (*domain.Booking, error) {
	booking := &domain.Booking{}
	if err := db.QueryRowContext(
		context.Background(), queries.FindBookingByBookingID, bookingID,
	).Scan(
		&booking.ID, &booking.BookingID, &booking.CampsiteID, &booking.Email,
		&booking.FullName, &booking.StartDate, &booking.EndDate, &booking.Active,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookingNotFound{BookingID: bookingID}
		}
		return nil, errors.Wrap(err, "scan booking row")
	}
	return booking, nil
}

func DeleteBookings(db *sql.DB) error {
	_, err := db.ExecContext(context.Background(), deleteBookingsQuery)
	return err
}

func DeleteCampsites(db *sql.DB) error {
	_, err := db.ExecContext(context.Background(), deleteCampsitesQuery)
	return err
}
