package postgres

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
)

func TestRollbackTx(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"Success": {
			err:  nil,
			want: "",
		},
		"Error_ErTxDone": {
			err:  sql.ErrTxDone,
			want: "",
		},
		"Error_Unexpected": {
			err:  errors.Wrap(errors.ErrUnknown, "unexpected error during rollback"),
			want: "ERROR rollback transaction error=unexpected error during rollback",
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
			mock.ExpectRollback().WillReturnError(tc.err)

			var buf bytes.Buffer
			dl := logger.NewDefault(&buf, nil)
			// when
			rollbackTx(tx, dl)
			// then
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)

			got := buf.String()
			if tc.want != "" {
				assert.Containsf(t, got, tc.want, "rollbackTx() got = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestCloseRows(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"Success": {
			err:  nil,
			want: "",
		},
		"Error_Unexpected": {
			err:  errors.Wrap(errors.ErrUnknown, "unexpected error during close rows"),
			want: "ERROR close rows error=unexpected error during close rows",
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

			mockRows := sqlmock.NewRows([]string{"column"}).CloseError(tc.err)
			mock.ExpectQuery("^SELECT (.+)$").WillReturnRows(mockRows)

			rows, err := db.Query("SELECT 1")
			if err != nil {
				t.Fatalf("execute query error: %v", err)
			}

			var buf bytes.Buffer
			dl := logger.NewDefault(&buf, nil)
			// when
			closeRows(rows, dl)
			// then
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)

			got := buf.String()
			if tc.want != "" {
				assert.Containsf(t, got, tc.want, "closeRows() got = %s, want %s", got, tc.want)
			}
		})
	}
}
