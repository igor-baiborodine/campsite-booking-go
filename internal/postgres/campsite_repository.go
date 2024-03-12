package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"

	"github.com/stackus/errors"
)

type CampsiteRepository struct {
	db *sql.DB
}

var _ domain.CampsiteRepository = (*CampsiteRepository)(nil)

func NewCampsiteRepository(db *sql.DB) CampsiteRepository {
	return CampsiteRepository{db: db}
}

func (r CampsiteRepository) FindAll(ctx context.Context) (campsites []*domain.Campsite, err error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "begin transaction")
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, FindAllInCampsites)
	if err != nil {
		return nil, errors.Wrap(err, "query campsites")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrapf(err, "close campsite rows")
		}
	}(rows)

	for rows.Next() {
		campsite := &domain.Campsite{}
		err := rows.Scan(&campsite.ID, &campsite.CampsiteID, &campsite.CampsiteCode,
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
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit transaction")
	}
	return campsites, nil
}

func (r CampsiteRepository) Insert(ctx context.Context, campsite *domain.Campsite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "begin transaction")
	}
	defer tx.Rollback()

	createdAt := time.Now()
	_, err = tx.ExecContext(ctx, InsertIntoCampsites,
		campsite.CampsiteID, campsite.CampsiteCode, campsite.Capacity, campsite.Restrooms,
		campsite.DrinkingWater, campsite.PicnicTable, campsite.FirePit, campsite.Active,
		createdAt, createdAt)
	if err != nil {
		return errors.Wrapf(err, "insert campsite")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "commit transaction")
	}
	return nil
}
