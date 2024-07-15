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
		err     error
		wantLog string
	}{
		"Success": {
			err:     nil,
			wantLog: "",
		},
		"Error_ErTxDone": {
			err:     sql.ErrTxDone,
			wantLog: "",
		},
		"Error_Unexpected": {
			err:     errors.Wrap(errors.ErrUnknown, "unexpected error during rollback"),
			wantLog: "ERROR rollback transaction error=unexpected error during rollback",
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

			gotLog := buf.String()
			if tc.wantLog != "" {
				assert.Containsf(
					t,
					gotLog,
					tc.wantLog,
					"rollbackTx() gotLog = %s, wantLog %s",
					gotLog,
					tc.wantLog,
				)
			}
		})
	}
}

func TestCloseRows(t *testing.T) {
	tests := map[string]struct {
		err     error
		wantLog string
	}{
		"Success": {
			err:     nil,
			wantLog: "",
		},
		"Error_Unexpected": {
			err:     errors.Wrap(errors.ErrUnknown, "unexpected error during close rows"),
			wantLog: "ERROR close rows error=unexpected error during close rows",
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

			gotLog := buf.String()
			if tc.wantLog != "" {
				assert.Containsf(
					t,
					gotLog,
					tc.wantLog,
					"closeRows() gotLog = %s, wantLog %s",
					gotLog,
					tc.wantLog,
				)
			}
		})
	}
}
