package bootstrap

import (
	"math"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
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

func NewBookingWithAddDays(
	campsiteID string,
	startAddDays int,
	endAddDays int,
) (*domain.Booking, error) {
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

func AsStartOfDayUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
