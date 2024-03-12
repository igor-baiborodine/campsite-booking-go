//go:build integration || database

package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	ct "github.com/igor-baiborodine/campsite-booking-go/internal/common_testing"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	deleteBookings = "DELETE FROM campgrounds.bookings"
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
	s.container, err = ct.NewPostgresContainer()
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = ct.NewDB(s.container)
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
	campsite, err := ct.FakeCampsite()
	s.NoError(err)

	err = ct.InsertCampsite(s.db, campsite)
	s.NoError(err)

	booking, err := ct.FakeBooking(campsite.CampsiteID)
	s.NoError(err)

	err = ct.InsertBooking(s.db, booking)
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

func (s *bookingSuite) TestBookingRepository_Find_NotFound() {
	// TODO: implement me
}
