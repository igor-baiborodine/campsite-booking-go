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
	rows, err := r.db.QueryContext(ctx, FindAllInCampsites)
	if err != nil {
		return nil, errors.Wrap(err, "querying campsites")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing campsite rows")
		}
	}(rows)

	for rows.Next() {
		campsite := &domain.Campsite{}
		err := rows.Scan(&campsite.CampsiteID, &campsite.CampsiteCode, &campsite.Capacity,
			&campsite.Restrooms, &campsite.DrinkingWater, &campsite.PicnicTable, &campsite.FirePit,
			&campsite.Active)
		if err != nil {
			return nil, errors.Wrap(err, "scanning campsite")
		}
		campsites = append(campsites, campsite)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing campsite rows")
	}
	return campsites, nil
}

func (r CampsiteRepository) Insert(ctx context.Context, campsite *domain.Campsite) error {
	createdAt := time.Now()
	_, err := r.db.ExecContext(ctx, InsertIntoCampsites, campsite.CampsiteID, campsite.CampsiteCode,
		campsite.Capacity, campsite.Restrooms, campsite.DrinkingWater, campsite.PicnicTable,
		campsite.FirePit, campsite.Active, createdAt, createdAt)

	return errors.Wrap(err, "inserting campsite")
}
