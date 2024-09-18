package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	os "os"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/config"
	rpc "github.com/igor-baiborodine/campsite-booking-go/internal/grpc"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/igor-baiborodine/campsite-booking-go/internal/postgres"
	"github.com/igor-baiborodine/campsite-booking-go/internal/waiter"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	cfg    config.AppConfig
	db     *sql.DB
	rpc    *grpc.Server
	waiter waiter.Waiter
}

func New(cfg config.AppConfig) (*Service, error) {
	s := &Service{cfg: cfg}
	s.initLogger()

	if err := s.initDB(); err != nil {
		return nil, err
	}
	if err := s.initRPC(); err != nil {
		return nil, err
	}
	s.initWaiter()

	return s, nil
}

func (s *Service) Config() config.AppConfig {
	return s.cfg
}

func (s *Service) DB() *sql.DB {
	return s.db
}

func (s *Service) RPC() *grpc.Server {
	return s.rpc
}

func (s *Service) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *Service) initDB() (err error) {
	s.db, err = sql.Open("pgx", config.ReplaceEnvPlaceholders(s.cfg.PG.Conn))
	return err
}

func (s *Service) initRPC() (err error) {
	srv, err := rpc.NewServer()
	if err != nil {
		return err
	}
	s.rpc = srv
	reflection.Register(s.rpc)

	return nil
}

func (s *Service) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *Service) initLogger() {
	l := logger.New(logger.LogConfig{
		Environment: s.cfg.Environment,
		LogLevel:    logger.Level(s.cfg.LogLevel),
	})
	slog.SetDefault(l)

	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		slog.Default().Info("slog debug level enabled")
	}
}

func (s *Service) MigrateDB(fs fs.FS) error {
	goose.SetLogger(&logger.SilentLogger{})
	goose.SetBaseFS(fs)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}
	return nil
}

func (s *Service) Startup() error {
	// setup driven adapters
	campsites := postgres.NewCampsiteRepository(s.db)
	bookings := postgres.NewBookingRepository(s.db)
	// setup application
	app := application.New(campsites, bookings)
	// setup driver adapters
	if err := rpc.RegisterServer(app, s.rpc); err != nil {
		return err
	}
	return nil
}

func (s *Service) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.RPC.Address())
	if err != nil {
		return err
	}
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		slog.Info("✅ rpc server started")
		defer slog.Info("🚫 rpc server shut down")

		if err := s.RPC().Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})

	var pprofServer *http.Server
	if os.Getenv("ENABLE_PPROF") == "true" {
		group.Go(func() error {
			slog.Info("✅ pprof server started")
			defer slog.Info("🚫 pprof server shut down")

			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

			pprofServer = &http.Server{
				Addr:    ":6060",
				Handler: mux,
			}
			if err := pprofServer.ListenAndServe(); err != nil {
				return err
			}
			return nil
		})
	}

	group.Go(func() error {
		<-gCtx.Done()
		slog.Info("rpc server to be shut down")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			if pprofServer != nil {
				_ = pprofServer.Shutdown(gCtx)
			}
			close(stopped)
		}()
		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// force it to stop
			s.RPC().Stop()
			if pprofServer != nil {
				_ = pprofServer.Shutdown(gCtx)
			}
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}
