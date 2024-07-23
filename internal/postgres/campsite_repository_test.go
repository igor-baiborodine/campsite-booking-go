//go:build !integration

package postgres

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
)

func TestCampsiteRepository_FindAll(t *testing.T) {
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

	tests := map[string]struct {
		mockTxPhases func(mock sqlmock.Sqlmock)
		want         []*domain.Campsite
		wantErr      error
	}{
		"Success": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsites).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want:    campsites,
			wantErr: nil,
		},
		"NoCampsitesFound": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsites).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want:    nil,
			wantErr: nil,
		},
		"Error_BeginTx": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(bootstrap.ErrBeginTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrBeginTx,
		},
		"Error_Query": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsites).WillReturnError(bootstrap.ErrQuery)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrQuery,
		},
		"Error_Rows": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
				rows.RowError(2, bootstrap.ErrRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsites).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrRow,
		},
		"Error_CommitTx": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsites).WillReturnRows(rows)
				mock.ExpectCommit().WillReturnError(bootstrap.ErrCommitTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("open stub database connection error: %v", err)
			}
			defer db.Close()

			tc.mockTxPhases(mock)
			repo := NewCampsiteRepository(db)
			// when
			got, err := repo.FindAll(context.TODO())
			// then
			assert.Equal(t, tc.want, got, "FindAll() got = %v, want %v",
				got, tc.want)
			assert.ErrorIs(t, err, tc.wantErr, "FindAll() error = %v, wantErr %v",
				err, tc.wantErr)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestCampsiteRepository_Insert(t *testing.T) {
	campsite, err := bootstrap.NewCampsite()
	if err != nil {
		t.Fatalf("create campsite error: %v", err)
	}

	tests := map[string]struct {
		mockTxPhases func(mock sqlmock.Sqlmock)
		wantErr      error
	}{
		"Success": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(queries.InsertCampsite).
					WithArgs(campsiteArgs(campsite)...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		"Error_BeginTx": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(bootstrap.ErrBeginTx)
			},
			wantErr: bootstrap.ErrBeginTx,
		},
		"Error_Exec": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(queries.InsertCampsite).
					WithArgs(campsiteArgs(campsite)...).
					WillReturnError(bootstrap.ErrExec)
				mock.ExpectRollback()
			},
			wantErr: bootstrap.ErrExec,
		},
		"Error_CommitTx": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(queries.InsertCampsite).
					WithArgs(campsiteArgs(campsite)...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(bootstrap.ErrCommitTx)
			},
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("open stub database connection error: %v", err)
			}
			defer db.Close()

			tc.mockTxPhases(mock)
			repo := NewCampsiteRepository(db)
			// when
			err = repo.Insert(context.TODO(), campsite)
			// then
			assert.ErrorIs(t, err, tc.wantErr, "Insert() error = %v, wantErr %v",
				err, tc.wantErr)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func campsiteArgs(c *domain.Campsite) []driver.Value {
	return campsiteRowValues(c)[1:] // remove ID
}

func campsiteRowValues(c *domain.Campsite) []driver.Value {
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
