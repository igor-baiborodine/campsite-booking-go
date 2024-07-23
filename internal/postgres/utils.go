package postgres

import (
	"database/sql"
	"log/slog"

	"github.com/stackus/errors"
)

func rollbackTx(tx *sql.Tx) {
	// Rollback returns sql.ErrTxDone if the transaction was already closed.
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		slog.Error("rollback transaction", slog.Any("error", err))
	}
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		slog.Error("close rows", slog.Any("error", err))
	}
}
