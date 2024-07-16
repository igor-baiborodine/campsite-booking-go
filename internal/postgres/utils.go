package postgres

import (
	"database/sql"
	"log/slog"

	"github.com/stackus/errors"
)

func rollbackTx(tx *sql.Tx, logger *slog.Logger) {
	// Rollback returns sql.ErrTxDone if the transaction was already closed.
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Error("rollback transaction", slog.Any("error", err))
	}
}

func closeRows(rows *sql.Rows, logger *slog.Logger) {
	if err := rows.Close(); err != nil {
		logger.Error("close rows", slog.Any("error", err))
	}
}
