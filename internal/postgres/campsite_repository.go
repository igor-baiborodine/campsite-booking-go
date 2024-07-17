package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/stackus/errors"
)

type CampsiteRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

var _ domain.CampsiteRepository = (*CampsiteRepository)(nil)

func NewCampsiteRepository(db *sql.DB, logger *slog.Logger) CampsiteRepository {
	return CampsiteRepository{db, logger}
}

func (r CampsiteRepository) FindAll(ctx context.Context) (campsites []*domain.Campsite, err error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "begin transaction")
	}
	defer rollbackTx(tx, r.logger)

	rows, err := tx.QueryContext(ctx, queries.FindAllCampsitesQuery)
	if err != nil {
		return nil, errors.Wrap(err, "query campsites")
	}
	defer closeRows(rows, r.logger)

	for rows.Next() {
		campsite := &domain.Campsite{}
		err = rows.Scan(&campsite.ID, &campsite.CampsiteID, &campsite.CampsiteCode,
			&campsite.Capacity, &campsite.Restrooms, &campsite.DrinkingWater, &campsite.PicnicTable,
			&campsite.FirePit, &campsite.Active)
		if err != nil {
			return nil, errors.Wrap(err, "scan campsite row")
		}
		campsites = append(campsites, campsite)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finish campsite rows")
	}
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return campsites, nil
}

func (r CampsiteRepository) Insert(ctx context.Context, campsite *domain.Campsite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer rollbackTx(tx, r.logger)

	_, err = tx.ExecContext(ctx, queries.InsertCampsiteQuery,
		campsite.CampsiteID, campsite.CampsiteCode, campsite.Capacity, campsite.Restrooms,
		campsite.DrinkingWater, campsite.PicnicTable, campsite.FirePit, campsite.Active)
	if err != nil {
		return errors.Wrap(err, "insert campsite")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}
