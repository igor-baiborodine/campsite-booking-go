//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
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
	s.container, err = bootstrap.NewPostgresContainer()
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = bootstrap.NewDB(s.container)
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
		s.T().Fatal("terminate postgres container", err)
	}
}

func (s *campsiteSuite) SetupTest() {
	s.repo = postgres.NewCampsiteRepository(s.db)
}

func (s *campsiteSuite) TearDownTest() {
	err := bootstrap.DeleteCampsites(s.db)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *campsiteSuite) TestCampsiteRepository_FindAll() {
	// given
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)

	err = bootstrap.InsertCampsite(s.db, campsite)
	s.NoError(err)
	// when
	got, err := s.repo.FindAll(context.Background())
	// then
	if s.NoError(err) {
		s.Equal(1, len(got))
		s.NotEqual(campsite.ID, got[0].ID)
		campsite.ID = got[0].ID
		s.Equal(campsite, got[0])
	}
}

func (s *campsiteSuite) TestCampsiteRepository_Insert() {
	// given
	campsite, err := bootstrap.NewCampsite()
	s.NoError(err)
	// when
	s.NoError(s.repo.Insert(context.Background(), campsite))
	// then
	query := "SELECT campsite_code, created_at, updated_at FROM campsites WHERE campsite_id = $1"
	row := s.db.QueryRow(query, campsite.CampsiteID)

	if s.NoError(row.Err()) {
		var campsiteCode string
		var createdAt, updatedAt time.Time
		s.NoError(row.Scan(&campsiteCode, &createdAt, &updatedAt))
		s.Equal(campsite.CampsiteCode, campsiteCode)
		s.NotNil(createdAt)
		s.Equal(createdAt, updatedAt)
	}
}
