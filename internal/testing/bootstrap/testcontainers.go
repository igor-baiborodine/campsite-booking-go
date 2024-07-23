package bootstrap

import (
	"context"
	"database/sql"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/db/migrations"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer() (*postgres.PostgresContainer, error) {
	ctx := context.Background()
	dbName := "test_campgrounds"
	dbUser := "test_campgrounds_user"
	dbPassword := "test_campgrounds_pass"

	return postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithInitScripts("../testing/bootstrap/db/test_create_campgrounds_db.sh"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func NewDB(c *postgres.PostgresContainer) (*sql.DB, error) {
	ctx := context.Background()
	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	goose.SetLogger(&logger.SilentLogger{})
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
