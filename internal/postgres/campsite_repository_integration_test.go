//go:build integration || database

package postgres_test

import (
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger/log"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"math"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/db/migrations"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
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
	ctx := context.Background()
	dbName := "test_campgrounds"
	dbUser := "test_campgrounds_user"
	dbPassword := "test_campgrounds_pass"

	s.container, err = pg.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		pg.WithDatabase(dbName),
		pg.WithUsername(dbUser),
		pg.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	s.checkError(err)

	connStr, err := s.container.ConnectionString(ctx, "sslmode=disable")
	s.checkError(err)

	s.db, err = sql.Open("pgx", connStr)
	s.checkError(err)

	goose.SetLogger(&log.SilentLogger{})
	goose.SetBaseFS(migrations.FS)
	err = goose.SetDialect("postgres")
	s.checkError(err)

	err = goose.Up(s.db, ".")
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
	s.repo = postgres.NewCampsiteRepository("campsites", s.db)
}
func (s *campsiteSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), "TRUNCATE campsites")
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
	const query = "INSERT INTO campsites " +
		"(campsite_id, campsite_code, capacity, restrooms, drinking_water, picnic_table, fire_pit, active, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	createdAt := time.Now()
	_, err := s.db.ExecContext(context.Background(), query, campsite.CampsiteID, campsite.CampsiteCode,
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
	row := s.db.QueryRow("SELECT campsite_code FROM campsites WHERE campsite_id = $1", campsite.CampsiteID)
	if s.NoError(row.Err()) {
		var campsiteCode string
		s.NoError(row.Scan(&campsiteCode))
		s.Equal(campsite.CampsiteCode, campsiteCode)
	}
}
