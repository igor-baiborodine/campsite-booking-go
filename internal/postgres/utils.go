package postgres

import (
	"database/sql"
	"log/slog"

	"github.com/stackus/errors"
)

func rollbackTx(tx *sql.Tx, logger *slog.Logger) {
	// Rollback returns sql.ErrTxDone if the transaction was already closed.
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Error("rollback transaction", slog.Any("error", err))
	}
}
