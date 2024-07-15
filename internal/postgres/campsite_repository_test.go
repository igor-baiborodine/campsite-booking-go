//go:build !integration

package postgres

import (
	"context"
	"database/sql/driver"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var campsites []*domain.Campsite
	for i := 1; i < 4; i++ {
		campsite, err := bootstrap.NewCampsite()
		if err != nil {
			t.Fatalf("create campsite[%d] error: %v", i, err)
		}
		campsite.ID = int64(i)
		campsites = append(campsites, campsite)
	}

	columnsRow := []string{
		"id",
		"campsite_id",
		"campsite_code",
		"capacity",
		"restrooms",
		"drinking_water",
		"picnic_table",
		"fire_pit",
		"active",
	}
	dbErr := errors.Wrap(errors.ErrUnknown, "unexpected error during db query")

	tests := map[string]struct {
		mockQuery     func(ctx context.Context, mock sqlmock.Sqlmock)
		wantCampsites []*domain.Campsite
		wantErr       error
	}{
		"Success": {
			mockQuery: func(ctx context.Context, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(rowValues(campsites[0])...).
					AddRow(rowValues(campsites[1])...).
					AddRow(rowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			wantCampsites: campsites,
			wantErr:       nil,
		},
		"NoCampsitesFound": {
			mockQuery: func(ctx context.Context, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			wantCampsites: nil,
		},
		"Error_Unexpected": {
			mockQuery: func(ctx context.Context, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			wantErr: dbErr,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("open stub database connection error: %v", err)
			}
			defer db.Close()

			ctx := context.TODO()
			tc.mockQuery(ctx, mock)
			repo := NewCampsiteRepository(db, logger.NewDefault(os.Stdout, nil))
			// when
			campsites, err := repo.FindAll(ctx)
			// then
			assert.Equal(t, tc.wantCampsites, campsites)
			assert.ErrorIs(t, err, tc.wantErr, "FindAll() error = %v, wantErr %v",
				err, tc.wantErr)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func rowValues(c *domain.Campsite) []driver.Value {
	return []driver.Value{
		c.ID,
		c.CampsiteID,
		c.CampsiteCode,
		c.Capacity,
		c.Restrooms,
		c.DrinkingWater,
		c.PicnicTable,
		c.FirePit,
		c.Active,
	}
}
