//go:build integration || database

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/scaffold"
	"github.com/stackus/errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	deleteBookings = "DELETE FROM bookings"
)

type bookingSuite struct {
	container *pg.PostgresContainer
	db        *sql.DB
	repo      postgres.BookingRepository
	suite.Suite
}

func TestBookingRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode: skipping")
	}
	suite.Run(t, &bookingSuite{})
}

func (s *bookingSuite) SetupSuite() {
	var err error
	s.container, err = scaffold.NewPostgresContainer()
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = scaffold.NewDB(s.container)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *bookingSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal(err)
	}
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal("failed to terminate postgres container", err)
	}
}

func (s *bookingSuite) SetupTest() {
	s.repo = postgres.NewBookingRepository(s.db)
}
func (s *bookingSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), deleteBookings)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.db.ExecContext(context.Background(), deleteCampsites)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *bookingSuite) TestBookingRepository_Find_Success() {
	// given
	campsite, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking, err := scaffold.FakeBooking(campsite.CampsiteID)
	s.NoError(err)

	err = scaffold.InsertBooking(s.db, booking)
	s.NoError(err)
	// when
	result, err := s.repo.Find(context.Background(), booking.BookingID)
	// then
	if s.NoError(err) {
		s.NotNil(result)
		s.NotEqual(booking.ID, result.ID)
		booking.ID = result.ID
		s.Equal(booking, result)
	}
}

func (s *bookingSuite) TestBookingRepository_Find_ErrNotFound() {
	// given
	booking := &domain.Booking{
		BookingID: "non-existing-booking-id",
	}
	// when
	result, err := s.repo.Find(context.Background(), booking.BookingID)
	// then
	if s.Error(err) {
		s.Nil(result)
		s.True(errors.Is(err, domain.ErrBookingNotFound{BookingID: booking.BookingID}))
		s.Equal("booking not found for BookingID non-existing-booking-id", err.Error())
	}
}

func (s *bookingSuite) TestBookingRepository_FindForDateRange_Success() {
	tests := map[string]struct {
		s, e   string // existing booking start/end dates in ISO 8601 format, denoted as S and E
		rs, re string // given date range start/end in ISO 8601 format, denoted as |-|
		len    int    // len of returned bookings slice
	}{
		"SE|-|----|-|--": {s: "2006-01-02", e: "2006-01-03", rs: "2006-01-04", re: "2006-01-07", len: 0},
		"-S|E|----|-|--": {s: "2006-01-02", e: "2006-01-03", rs: "2006-01-03", re: "2006-01-07", len: 0},
		"-S|-|E---|-|--": {s: "2006-01-02", e: "2006-01-04", rs: "2006-01-03", re: "2006-01-07", len: 1},
		"--|S|E---|-|--": {s: "2006-01-02", e: "2006-01-03", rs: "2006-01-02", re: "2006-01-07", len: 1},
		"--|-|S--E|-|--": {s: "2006-01-03", e: "2006-01-04", rs: "2006-01-02", re: "2006-01-07", len: 1},
		"--|-|---S|E|--": {s: "2006-01-06", e: "2006-01-07", rs: "2006-01-02", re: "2006-01-07", len: 1},
		"--|-|---S|-|E-": {s: "2006-01-06", e: "2006-01-08", rs: "2006-01-02", re: "2006-01-07", len: 1},
		"--|-|----|S|E-": {s: "2006-01-07", e: "2006-01-08", rs: "2006-01-02", re: "2006-01-07", len: 1},
		"--|-|----|-|SE": {s: "2006-01-08", e: "2006-01-09", rs: "2006-01-02", re: "2006-01-07", len: 0},
		"-S|-|----|-|E-": {s: "2006-01-02", e: "2006-01-08", rs: "2006-01-03", re: "2006-01-07", len: 1},
	}

	for name, test := range tests {
		s.T().Run(name, func(t *testing.T) {
			campsite, err := scaffold.FakeCampsite()
			s.NoError(err)

			err = scaffold.InsertCampsite(s.db, campsite)
			s.NoError(err)

			booking, err := scaffold.FakeBooking(campsite.CampsiteID)
			s.NoError(err)
			booking.StartDate, _ = time.Parse(time.DateOnly, test.s)
			booking.EndDate, _ = time.Parse(time.DateOnly, test.e)

			err = scaffold.InsertBooking(s.db, booking)
			s.NoError(err)
			start, _ := time.Parse(time.DateOnly, test.rs)
			end, _ := time.Parse(time.DateOnly, test.re)
			// when
			result, err := s.repo.FindForDateRange(
				context.Background(), campsite.CampsiteID, start, end)
			// then
			if s.NoError(err) {
				s.Equal(test.len, len(result))
			}
		})
	}
}

func (s *bookingSuite) TestBookingRepository_Insert_Success() {
	// given
	campsite, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking, err := scaffold.FakeBooking(campsite.CampsiteID)
	s.NoError(err)
	// when
	err = s.repo.Insert(context.Background(), booking)
	// then
	if s.NoError(err) {
		found, err := scaffold.FindBooking(s.db, booking.BookingID)
		s.NoError(err)
		s.NotNil(found)
		s.NotEqual(booking.ID, found.ID)
		booking.ID = found.ID
		s.Equal(booking, found)
	}
}

func (s *bookingSuite) TestBookingRepository_Insert_ErrBookingDatesNotAvailable() {
	// given
	campsite, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking1, err := scaffold.FakeBooking(campsite.CampsiteID)
	s.NoError(err)

	err = scaffold.InsertBooking(s.db, booking1)
	s.NoError(err)

	booking2, err := scaffold.FakeBooking(campsite.CampsiteID)
	s.NoError(err)
	booking2.StartDate = booking1.StartDate
	booking2.EndDate = booking1.EndDate
	// when
	err = s.repo.Insert(context.Background(), booking2)
	// then
	if s.Error(err) {
		s.True(errors.Is(err, domain.ErrBookingDatesNotAvailable{
			StartDate: booking2.StartDate,
			EndDate:   booking2.EndDate,
		}))
		errMsg := fmt.Sprintf("booking dates not available from %s to %s",
			booking2.StartDate.Format(time.DateOnly), booking2.EndDate.Format(time.DateOnly))
		s.Equal(errMsg, err.Error())
	}
}

func (s *bookingSuite) TestBookingRepository_Update_Success() {
	// given
	campsite1, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite1)
	s.NoError(err)

	booking, err := scaffold.FakeBookingWithAddDays(campsite1.CampsiteID, 1, 2)
	s.NoError(err)

	err = scaffold.InsertBooking(s.db, booking)
	s.NoError(err)

	campsite2, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite2)
	s.NoError(err)

	existingBooking, err := scaffold.FindBooking(s.db, booking.BookingID)
	s.NoError(err)

	bookingToUpdate, err := scaffold.FakeBookingWithAddDays(campsite2.CampsiteID, 2, 3)
	s.NoError(err)

	bookingToUpdate.BookingID = existingBooking.BookingID
	bookingToUpdate.Active = !existingBooking.Active
	// when
	err = s.repo.Update(context.Background(), bookingToUpdate)
	// then
	if s.NoError(err) {
		updatedBooking, err := scaffold.FindBooking(s.db, bookingToUpdate.BookingID)
		s.NoError(err)
		s.NotNil(updatedBooking)
		s.Equal(existingBooking.ID, updatedBooking.ID)

		s.NotEqual(existingBooking.CampsiteID, updatedBooking.CampsiteID)
		s.NotEqual(existingBooking.Email, updatedBooking.Email)
		s.NotEqual(existingBooking.FullName, updatedBooking.FullName)
		s.NotEqual(existingBooking.StartDate, updatedBooking.StartDate)
		s.NotEqual(existingBooking.EndDate, updatedBooking.EndDate)
		s.NotEqual(existingBooking.Active, updatedBooking.Active)
	}
}

func (s *bookingSuite) TestBookingRepository_Update_ErrBookingDatesNotAvailable() {
	// given
	campsite, err := scaffold.FakeCampsite()
	s.NoError(err)

	err = scaffold.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking1, err := scaffold.FakeBookingWithAddDays(campsite.CampsiteID, 1, 2)
	s.NoError(err)

	err = scaffold.InsertBooking(s.db, booking1)
	s.NoError(err)

	booking2, err := scaffold.FakeBookingWithAddDays(campsite.CampsiteID, 2, 3)
	s.NoError(err)

	err = scaffold.InsertBooking(s.db, booking2)
	s.NoError(err)
	booking2.StartDate = booking1.StartDate
	booking2.EndDate = booking1.EndDate
	// when
	err = s.repo.Update(context.Background(), booking2)
	// then
	if s.Error(err) {
		s.True(errors.Is(err, domain.ErrBookingDatesNotAvailable{
			StartDate: booking2.StartDate,
			EndDate:   booking2.EndDate,
		}))
		errMsg := fmt.Sprintf("booking dates not available from %s to %s",
			booking2.StartDate.Format(time.DateOnly), booking2.EndDate.Format(time.DateOnly))
		s.Equal(errMsg, err.Error())
	}
}
