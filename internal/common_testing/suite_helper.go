package common_testing

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"math"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"

	"github.com/go-faker/faker/v4"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"

	"github.com/igor-baiborodine/campsite-booking-go/db/migrations"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger/log"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer() (*pg.PostgresContainer, error) {
	ctx := context.Background()
	dbName := "test_campgrounds"
	dbUser := "test_campgrounds_user"
	dbPassword := "test_campgrounds_pass"

	return pg.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		pg.WithDatabase(dbName),
		pg.WithUsername(dbUser),
		pg.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func NewDB(c *pg.PostgresContainer) (*sql.DB, error) {
	ctx := context.Background()
	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	goose.SetLogger(&log.SilentLogger{})
	goose.SetBaseFS(migrations.FS)

	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FakeCampsite() (*domain.Campsite, error) {
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

func FakeBooking(campsiteId string) (*domain.Booking, error) {
	booking := domain.Booking{}
	err := faker.FakeData(&booking)

	if err != nil {
		return nil, err
	}
	now := truncateToStartOfDayInUTC(time.Now())

	booking.ID = math.MaxInt64
	booking.BookingID = uuid.New().String()
	booking.CampsiteID = campsiteId
	booking.StartDate = now.AddDate(0, 0, 1)
	booking.EndDate = now.AddDate(0, 0, 2)
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
		context.Background(), postgres.FindBookingByBookingIdQuery, bookingID,
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

func truncateToStartOfDayInUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
