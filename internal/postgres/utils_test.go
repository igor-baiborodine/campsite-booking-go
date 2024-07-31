//go:build !integration

package postgres

import (
	"bytes"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
)

func TestRollbackTx(t *testing.T) {
	tests := map[string]struct {
		mockErr error
		want    string
	}{
		"Success": {
			mockErr: nil,
			want:    "",
		},
		"Error_ErTxDone": {
			mockErr: sql.ErrTxDone,
			want:    "",
		},
		"Error_Unexpected": {
			mockErr: errors.Wrap(errors.ErrUnknown, "unexpected error during rollback"),
			want:    "ERROR rollback transaction error=unexpected error during rollback",
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

			mock.ExpectBegin().WillReturnError(nil)
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf("begin transaction error: %v", err)
			}
			mock.ExpectRollback().WillReturnError(tc.mockErr)

			var buf bytes.Buffer
			slog.SetDefault(logger.NewDefault(&buf, nil))
			// when
			rollbackTx(tx)
			// then
			if tc.want != "" {
				got := buf.String()
				assert.Contains(t, got, tc.want,
					"rollbackTx() got = %s, want %s", got, tc.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCloseRows(t *testing.T) {
	tests := map[string]struct {
		mockErr error
		want    string
	}{
		"Success": {
			mockErr: nil,
			want:    "",
		},
		"Error_Unexpected": {
			mockErr: errors.Wrap(errors.ErrUnknown, "unexpected error during close rows"),
			want:    "ERROR close rows error=unexpected error during close rows",
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

			mockRows := sqlmock.NewRows([]string{"column"}).CloseError(tc.mockErr)
			mock.ExpectQuery("^SELECT (.+)$").WillReturnRows(mockRows)

			rows, err := db.Query("SELECT 1")
			if err != nil {
				t.Fatalf("execute query error: %v", err)
			}

			var buf bytes.Buffer
			slog.SetDefault(logger.NewDefault(&buf, nil))
			// when
			closeRows(rows)
			// then
			if tc.want != "" {
				got := buf.String()
				assert.Contains(t, got, tc.want,
					"closeRows() got = %s, want %s", got, tc.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
