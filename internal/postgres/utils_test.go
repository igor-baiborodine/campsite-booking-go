package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
)

func TestRollbackTx(t *testing.T) {

	tests := map[string]struct {
		err error
	}{
		"Success": {
			err: nil,
		},
		"Error_ErTxDone": {
			err: sql.ErrTxDone,
		},
		"Error_Unexpected": {
			err: errors.Wrap(errors.ErrUnknown, "unexpected error during rollback"),
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
			l := logger.NewStdout(nil)
			// when
			rollbackTx(tx, l)
			// then
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
