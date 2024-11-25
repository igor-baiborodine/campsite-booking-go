//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stackus/errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
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
	s.container, err = bootstrap.NewPostgresContainer()
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = bootstrap.NewDB(s.container)
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
		s.T().Fatal("terminate postgres container", err)
	}
}

func (s *bookingSuite) SetupTest() {
	s.repo = postgres.NewBookingRepository(s.db)
}

func (s *bookingSuite) TearDownTest() {
	err := bootstrap.DeleteBookings(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	err = bootstrap.DeleteCampsites(s.db)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *bookingSuite) TestBookingRepository_Find_Success() {
	// given
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking, err := bootstrap.NewBooking(campsite.CampsiteID)
	s.NoError(err)

	err = bootstrap.InsertBooking(s.db, booking)
	s.NoError(err)
	// when
	got, err := s.repo.Find(context.Background(), booking.BookingID)
	// then
	if s.NoError(err) {
		s.NotNil(got)
		s.NotEqual(booking.ID, got.ID)
		booking.ID = got.ID
		s.Equal(booking, got)
	}
}

func (s *bookingSuite) TestBookingRepository_Find_ErrNotFound() {
	// given
	booking := &domain.Booking{
		BookingID: "non-existing-booking-id",
	}
	// when
	got, err := s.repo.Find(context.Background(), booking.BookingID)
	// then
	if s.Error(err) {
		s.Nil(got)
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
		"SE|-|----|-|--": {
			s:   "2006-01-02",
			e:   "2006-01-03",
			rs:  "2006-01-04",
			re:  "2006-01-07",
			len: 0,
		},
		"-S|E|----|-|--": {
			s:   "2006-01-02",
			e:   "2006-01-03",
			rs:  "2006-01-03",
			re:  "2006-01-07",
			len: 0,
		},
		"-S|-|E---|-|--": {
			s:   "2006-01-02",
			e:   "2006-01-04",
			rs:  "2006-01-03",
			re:  "2006-01-07",
			len: 1,
		},
		"--|S|E---|-|--": {
			s:   "2006-01-02",
			e:   "2006-01-03",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 1,
		},
		"--|-|S--E|-|--": {
			s:   "2006-01-03",
			e:   "2006-01-04",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 1,
		},
		"--|-|---S|E|--": {
			s:   "2006-01-06",
			e:   "2006-01-07",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 1,
		},
		"--|-|---S|-|E-": {
			s:   "2006-01-06",
			e:   "2006-01-08",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 1,
		},
		"--|-|----|S|E-": {
			s:   "2006-01-07",
			e:   "2006-01-08",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 1,
		},
		"--|-|----|-|SE": {
			s:   "2006-01-08",
			e:   "2006-01-09",
			rs:  "2006-01-02",
			re:  "2006-01-07",
			len: 0,
		},
		"-S|-|----|-|E-": {
			s:   "2006-01-02",
			e:   "2006-01-08",
			rs:  "2006-01-03",
			re:  "2006-01-07",
			len: 1,
		},
	}

	for name, test := range tests {
		s.T().Run(name, func(t *testing.T) {
			campsite, err := bootstrap.NewCampsite()
			s.NoError(err)

			err = bootstrap.InsertCampsite(s.db, campsite)
			s.NoError(err)

			booking, err := bootstrap.NewBooking(campsite.CampsiteID)
			s.NoError(err)
			booking.StartDate, _ = time.Parse(time.DateOnly, test.s)
			booking.EndDate, _ = time.Parse(time.DateOnly, test.e)

			err = bootstrap.InsertBooking(s.db, booking)
			s.NoError(err)
			start, _ := time.Parse(time.DateOnly, test.rs)
			end, _ := time.Parse(time.DateOnly, test.re)
			// when
			got, err := s.repo.FindForDateRange(
				context.Background(), campsite.CampsiteID, start, end)
			// then
			if s.NoError(err) {
				s.Equal(test.len, len(got))
			}
		})
	}
}

func (s *bookingSuite) TestBookingRepository_Insert_Success() {
	// given
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking, err := bootstrap.NewBooking(campsite.CampsiteID)
	s.NoError(err)
	// when
	s.NoError(s.repo.Insert(context.Background(), booking))
	// then
	query := "SELECT campsite_id, created_at, updated_at, version FROM bookings WHERE campsite_id = $1"
	row := s.db.QueryRow(query, campsite.CampsiteID)

	if s.NoError(row.Err()) {
		var campsiteID string
		var createdAt, updatedAt time.Time
		var version int
		s.NoError(row.Scan(&campsiteID, &createdAt, &updatedAt, &version))
		s.Equal(booking.CampsiteID, campsiteID)
		s.NotNil(createdAt)
		s.Equal(createdAt, updatedAt)
		s.Equal(1, version)
	}
}

func (s *bookingSuite) TestBookingRepository_Insert_ErrBookingDatesNotAvailable() {
	// given
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking1, err := bootstrap.NewBooking(campsite.CampsiteID)
	s.NoError(err)

	err = bootstrap.InsertBooking(s.db, booking1)
	s.NoError(err)

	booking2, err := bootstrap.NewBooking(campsite.CampsiteID)
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
	campsite1, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite1)
	s.NoError(err)

	booking, err := bootstrap.NewBookingWithAddDays(campsite1.CampsiteID, 1, 2)
	s.NoError(err)

	err = bootstrap.InsertBooking(s.db, booking)
	s.NoError(err)

	campsite2, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite2)
	s.NoError(err)

	existingBooking, err := bootstrap.FindBooking(s.db, booking.BookingID)
	s.NoError(err)

	bookingToUpdate, err := bootstrap.NewBookingWithAddDays(campsite2.CampsiteID, 2, 3)
	s.NoError(err)

	bookingToUpdate.BookingID = existingBooking.BookingID
	bookingToUpdate.Active = !existingBooking.Active
	// when
	err = s.repo.Update(context.Background(), bookingToUpdate)
	// then
	if s.NoError(err) {
		updatedBooking, err := bootstrap.FindBooking(s.db, bookingToUpdate.BookingID)
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
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking1, err := bootstrap.NewBookingWithAddDays(campsite.CampsiteID, 1, 2)
	s.NoError(err)

	err = bootstrap.InsertBooking(s.db, booking1)
	s.NoError(err)

	booking2, err := bootstrap.NewBookingWithAddDays(campsite.CampsiteID, 2, 3)
	s.NoError(err)

	err = bootstrap.InsertBooking(s.db, booking2)
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
