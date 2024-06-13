package bootstrap

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	"github.com/stackus/errors"
)

func NewCampsite() (*domain.Campsite, error) {
	campsite := domain.Campsite{}
	err := faker.FakeData(&campsite)

	if err != nil {
		return nil, err
	}
	campsite.ID = math.MaxInt64
	campsite.CampsiteID = uuid.New().String()
	campsite.Active = true

	return &campsite, nil
}

func NewBooking(campsiteID string) (*domain.Booking, error) {
	return NewBookingWithAddDays(campsiteID, 1, 2)
}

func NewBookingWithAddDays(campsiteID string, startAddDays int, endAddDays int) (*domain.Booking, error) {
	booking := domain.Booking{}
	err := faker.FakeData(&booking)

	if err != nil {
		return nil, err
	}
	now := AsStartOfDayUTC(time.Now())

	booking.ID = math.MaxInt64
	booking.BookingID = uuid.New().String()
	booking.CampsiteID = campsiteID
	booking.StartDate = now.AddDate(0, 0, startAddDays)
	booking.EndDate = now.AddDate(0, 0, endAddDays)
	booking.Active = true

	return &booking, nil
}

func InsertCampsite(db *sql.DB, c *domain.Campsite) error {
	createdAt := time.Now()
	_, err := db.ExecContext(context.Background(), postgres.InsertCampsiteQuery,
		c.CampsiteID, c.CampsiteCode, c.Capacity, c.Restrooms, c.DrinkingWater, c.PicnicTable,
		c.FirePit, c.Active, createdAt, createdAt)
	return err
}

func InsertBooking(db *sql.DB, b *domain.Booking) error {
	createdAt := time.Now()
	_, err := db.ExecContext(context.Background(), postgres.InsertBookingQuery,
		b.BookingID, b.CampsiteID, b.Email, b.FullName, b.StartDate, b.EndDate, b.Active, createdAt,
		createdAt)
	return err
}

func FindBooking(db *sql.DB, bookingID string) (*domain.Booking, error) {
	booking := &domain.Booking{}
	if err := db.QueryRowContext(
		context.Background(), postgres.FindBookingByBookingIDQuery, bookingID,
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

func AsStartOfDayUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
