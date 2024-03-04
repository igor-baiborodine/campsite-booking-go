package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"time"

	"github.com/stackus/errors"
)

type CampsiteRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.CampsiteRepository = (*CampsiteRepository)(nil)

func NewCampsiteRepository(tableName string, db *sql.DB) CampsiteRepository {
	return CampsiteRepository{tableName: tableName, db: db}
}

func (r CampsiteRepository) FindAll(ctx context.Context) (campsites []*domain.Campsite, err error) {
	const query = "SELECT " +
		"campsite_id, campsite_code, capacity, restrooms, drinking_water, picnic_table, fire_pit, active " +
		"FROM %s"

	rows, err := r.db.QueryContext(ctx, r.table(query))
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
	const query = "INSERT INTO %s " +
		"(campsite_id, campsite_code, capacity, restrooms, drinking_water, picnic_table, fire_pit, active, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	createdAt := time.Now()
	_, err := r.db.ExecContext(ctx, r.table(query), campsite.CampsiteID, campsite.CampsiteCode,
		campsite.Capacity, campsite.Restrooms, campsite.DrinkingWater, campsite.PicnicTable,
		campsite.FirePit, campsite.Active, createdAt, createdAt)

	return errors.Wrap(err, "inserting campsite")
}

func (r CampsiteRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
