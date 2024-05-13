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
	deleteCampsites    = "DELETE FROM campgrounds.campsites"
	selectByCampsiteId = "SELECT campsite_code FROM campgrounds.campsites WHERE campsite_id = $1"
)

type campsiteSuite struct {
	container *pg.PostgresContainer
	db        *sql.DB
	repo      postgres.CampsiteRepository
	suite.Suite
}

func TestCampsiteRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode: skipping")
	}
	suite.Run(t, &campsiteSuite{})
}

func (s *campsiteSuite) SetupSuite() {
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

func (s *campsiteSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal(err)
	}
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal("failed to terminate postgres container", err)
	}
}

func (s *campsiteSuite) SetupTest() {
	s.repo = postgres.NewCampsiteRepository(s.db)
}
func (s *campsiteSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), deleteCampsites)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *campsiteSuite) TestCampsiteRepository_FindAll() {
	// given
	campsite, err := ct.FakeCampsite()
	s.NoError(err)

	err = ct.InsertCampsite(s.db, campsite)
	s.NoError(err)
	// when
	result, err := s.repo.FindAll(context.Background())
	// then
	if s.NoError(err) {
		s.Equal(1, len(result))
		s.NotEqual(campsite.ID, result[0].ID)
		campsite.ID = result[0].ID
		s.Equal(campsite, result[0])
	}
}

func (s *campsiteSuite) TestCampsiteRepository_Insert() {
	// given
	campsite, err := ct.FakeCampsite()
	s.NoError(err)
	// when
	s.NoError(s.repo.Insert(context.Background(), campsite))
	// then
	row := s.db.QueryRow(selectByCampsiteId, campsite.CampsiteID)
	if s.NoError(row.Err()) {
		var campsiteCode string
		s.NoError(row.Scan(&campsiteCode))
		s.Equal(campsite.CampsiteCode, campsiteCode)
	}
}
