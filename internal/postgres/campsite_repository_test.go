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

	tests := map[string]struct {
		mockQuery func(mock sqlmock.Sqlmock)
		want      []*domain.Campsite
		wantErr   error
	}{
		"Success": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
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
				mock.ExpectBegin().WillReturnError(bootstrap.ErrBeginTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrBeginTx,
		},
		"Error_Query": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnError(bootstrap.ErrQuery)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrQuery,
		},
		"Error_Rows": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
				rows.RowError(2, bootstrap.ErrRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrRow,
		},
		"Error_Commit": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(campsiteRowValues(campsites[0])...).
					AddRow(campsiteRowValues(campsites[1])...).
					AddRow(campsiteRowValues(campsites[2])...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindAllCampsitesQuery).WillReturnRows(rows)
				mock.ExpectCommit().WillReturnError(bootstrap.ErrCommit)
			},
			want:    nil,
			wantErr: bootstrap.ErrCommit,
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

			tc.mockQuery(mock)
			repo := NewCampsiteRepository(db, logger.NewDefault(os.Stdout, nil))
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

func TestInsert(t *testing.T) {
	campsite, err := bootstrap.NewCampsite()
	if err != nil {
		t.Fatalf("create campsite error: %v", err)
	}

	tests := map[string]struct {
		mockQuery func(mock sqlmock.Sqlmock)
		wantErr   error
	}{
		// TODO: implement after fixing issue with created_at and updated_at
		//"Success": {
		//	mockQuery: func(mock sqlmock.Sqlmock) {
		//		mock.ExpectBegin()
		//		mock.ExpectExec(queries.InsertCampsiteQuery).
		//			WithArgs(campsiteArgs(campsite)...).
		//			WillReturnResult(sqlmock.NewResult(1, 1))
		//		mock.ExpectCommit()
		//	},
		//	wantErr: nil,
		//},
		"Error_BeginTx": {
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(bootstrap.ErrBeginTx)
			},
			wantErr: bootstrap.ErrBeginTx,
		},
		// TODO: implement after fixing issue with created_at and updated_at
		//"Error_Query": {
		//	mockQuery: func(mock sqlmock.Sqlmock) {
		//		mock.ExpectBegin()
		//		mock.ExpectExec(queries.InsertCampsiteQuery).
		//			WithArgs(campsiteArgs(campsite)...).
		//			WillReturnError(bootstrap.ErrQuery)
		//		mock.ExpectRollback()
		//	},
		//	wantErr: bootstrap.ErrQuery,
		//},
		// TODO: implement after fixing issue with created_at and updated_at
		//"Error_Commit": {
		//	mockQuery: func(mock sqlmock.Sqlmock) {
		//		mock.ExpectBegin()
		//		mock.ExpectExec(queries.InsertCampsiteQuery).
		//			WithArgs(campsiteArgs(campsite)...).
		//			WillReturnResult(sqlmock.NewResult(1, 1))
		//		mock.ExpectCommit().WillReturnError(bootstrap.ErrCommit)
		//	},
		//	wantErr: bootstrap.ErrCommit,
		//},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("open stub database connection error: %v", err)
			}
			defer db.Close()

			tc.mockQuery(mock)
			repo := NewCampsiteRepository(db, logger.NewDefault(os.Stdout, nil))
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

// TODO: implement after fixing issue with created_at and updated_at
//func campsiteArgs(c *domain.Campsite) []driver.Value {
//	args := campsiteRowValues(c)[1:] // remove ID
//
//	createdAt := time.Now()
//	args = append(args, createdAt)
//	args = append(args, createdAt) // updated_at
//	return args
//}

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
