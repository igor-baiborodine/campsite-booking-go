package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/igor-baiborodine/campsite-booking-go/db/migrations"
	"github.com/igor-baiborodine/campsite-booking-go/internal/config"
	"github.com/igor-baiborodine/campsite-booking-go/internal/service"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		slog.Error("campgrounds exited abnormally:", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	s, err := service.New(cfg)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())
	if err = s.MigrateDB(migrations.FS); err != nil {
		return err
	}

	if err = s.Startup(); err != nil {
		return err
	}

	s.Logger().Info("✅ campgrounds app starting...")
	defer s.Logger().Info("🚫 stopped campgrounds app")

	s.Waiter().Add(s.WaitForRPC)

	return s.Waiter().Wait()
}
