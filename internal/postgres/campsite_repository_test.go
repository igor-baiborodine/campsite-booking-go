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
	beginTxErr := errors.Wrap(errors.ErrUnknown, "unexpected begin transaction error")
	queryErr := errors.Wrap(errors.ErrUnknown, "unexpected query error")
	rowErr := errors.Wrap(errors.ErrUnknown, "unexpected rows error")
	commitErr := errors.Wrap(errors.ErrUnknown, "unexpected commit error")

	tests := map[string]struct {
		mockQuery func(mock sqlmock.Sqlmock)
		want      []*domain.Campsite
		wantErr   error
	}{
		"Success": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(rowValues(campsites[0])...).
					AddRow(rowValues(campsites[1])...).
					AddRow(rowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want:    campsites,
			wantErr: nil,
		},
		"NoCampsitesFound": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want:    nil,
			wantErr: nil,
		},
		"Error_BeginTx": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(beginTxErr)
			},
			want:    nil,
			wantErr: beginTxErr,
		},
		"Error_Query": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnError(queryErr)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: queryErr,
		},
		"Error_Rows": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(rowValues(campsites[0])...).
					AddRow(rowValues(campsites[1])...).
					AddRow(rowValues(campsites[2])...)
				rows.RowError(2, rowErr)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: rowErr,
		},
		"Error_Commit": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(rowValues(campsites[0])...).
					AddRow(rowValues(campsites[1])...).
					AddRow(rowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit().WillReturnError(commitErr)
			},
			want:    nil,
			wantErr: commitErr,
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

			tc.mockQuery(mock)
			repo := NewCampsiteRepository(db, logger.NewDefault(os.Stdout, nil))
			// when
			campsites, err := repo.FindAll(context.TODO())
			// then
			assert.Equal(t, tc.want, campsites)
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
