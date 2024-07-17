//go:build !integration

package postgres

import (
	"context"
	"database/sql/driver"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	queries "github.com/igor-baiborodine/campsite-booking-go/internal/postgres/sql"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	campsiteID := uuid.New().String()
	booking, err := bootstrap.NewBooking(campsiteID)
	if err != nil {
		t.Fatalf("create booking error: %v", err)
	}
	booking.ID = 1

	columnsRow := []string{
		"id",
		"booking_id",
		"campsite_id",
		"email",
		"full_name",
		"start_date",
		"end_date",
		"active",
	}
	errBookingNotFound := domain.ErrBookingNotFound{BookingID: booking.BookingID}

	tests := map[string]struct {
		mockTxPhases func(mock sqlmock.Sqlmock)
		want         *domain.Booking
		wantErr      error
	}{
		"Success": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(bookingRowValues(booking)...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindBookingByBookingIDQuery).
					WithArgs(booking.BookingID).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want:    booking,
			wantErr: nil,
		},
		"NoBookingFound": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindBookingByBookingIDQuery).
					WithArgs(booking.BookingID).
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: errBookingNotFound,
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
				mock.ExpectQuery(queries.FindBookingByBookingIDQuery).
					WithArgs(booking.BookingID).
					WillReturnError(bootstrap.ErrQuery)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrQuery,
		},
		"Error_Rows": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(bookingRowValues(booking)...)
				rows.RowError(0, bootstrap.ErrRow)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindBookingByBookingIDQuery).
					WithArgs(booking.BookingID).
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: bootstrap.ErrRow,
		},
		"Error_Commit": {
			mockTxPhases: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columnsRow).
					AddRow(bookingRowValues(booking)...)
				mock.ExpectBegin()
				mock.ExpectQuery(queries.FindBookingByBookingIDQuery).
					WithArgs(booking.BookingID).
					WillReturnRows(rows)
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

			tc.mockTxPhases(mock)
			repo := NewBookingRepository(db, logger.NewDefault(os.Stdout, nil))
			// when
			got, err := repo.Find(context.TODO(), booking.BookingID)
			// then
			assert.Equal(t, tc.want, got, "Find() got = %v, want %v",
				got, tc.want)
			assert.ErrorIs(t, err, tc.wantErr, "Find() error = %v, wantErr %v",
				err, tc.wantErr)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func bookingRowValues(b *domain.Booking) []driver.Value {
	return []driver.Value{
		b.ID,
		b.BookingID,
		b.CampsiteID,
		b.Email,
		b.FullName,
		b.StartDate,
		b.EndDate,
		b.Active,
	}
}
