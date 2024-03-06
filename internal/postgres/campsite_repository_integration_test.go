//go:build integration || database

package postgres_test

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	ct "github.com/igor-baiborodine/campsite-booking-go/internal/common_testing"

	"github.com/go-faker/faker/v4"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	truncateCampsites  = "TRUNCATE campsites.campsites"
	selectByCampsiteId = "SELECT campsite_code FROM campsites.campsites WHERE campsite_id = $1"
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
	s.checkError(err)

	s.db, err = ct.NewDB(s.container)
	s.checkError(err)
}

func (s *campsiteSuite) TearDownSuite() {
	err := s.db.Close()
	s.checkError(err)
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal("failed to terminate postgres container", err)
	}
}

func (s *campsiteSuite) SetupTest() {
	s.repo = postgres.NewCampsiteRepository(s.db)
}
func (s *campsiteSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), truncateCampsites)
	s.checkError(err)
}

func (s *campsiteSuite) createCampsite() domain.Campsite {
	campsite := domain.Campsite{}
	err := faker.FakeData(&campsite)
	s.checkError(err)
	return campsite
}

func (s *campsiteSuite) checkError(err error) {
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *campsiteSuite) TestCampsiteRepository_FindAll() {
	// given
	campsite := s.createCampsite()
	campsite.ID = math.MaxInt64
	createdAt := time.Now()
	_, err := s.db.ExecContext(context.Background(), postgres.InsertIntoCampsites, campsite.CampsiteID, campsite.CampsiteCode,
		campsite.Capacity, campsite.Restrooms, campsite.DrinkingWater, campsite.PicnicTable,
		campsite.FirePit, campsite.Active, createdAt, createdAt)
	s.NoError(err)
	// when
	campsites, err := s.repo.FindAll(context.Background())
	// then
	if s.NoError(err) {
		s.Equal(1, len(campsites))
		s.NotEqual(campsite.ID, campsites[0].ID)
		campsite.ID = campsites[0].ID
		s.Equal(&campsite, campsites[0])
	}
}

func (s *campsiteSuite) TestCampsiteRepository_Insert() {
	// given
	campsite := s.createCampsite()
	// when
	s.NoError(s.repo.Insert(context.Background(), &campsite))
	// then
	row := s.db.QueryRow(selectByCampsiteId, campsite.CampsiteID)
	if s.NoError(row.Err()) {
		var campsiteCode string
		s.NoError(row.Scan(&campsiteCode))
		s.Equal(campsite.CampsiteCode, campsiteCode)
	}
}
