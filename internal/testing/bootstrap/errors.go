package bootstrap

import (
	"github.com/jackc/pgconn"
	"github.com/stackus/errors"
)

var (
	ErrBeginTx         = errors.Wrap(errors.ErrUnknown, "unexpected begin transaction error")
	ErrQuery           = errors.Wrap(errors.ErrUnknown, "unexpected query error")
	ErrExec            = errors.Wrap(errors.ErrUnknown, "unexpected exec error")
	ErrRow             = errors.Wrap(errors.ErrUnknown, "unexpected rows error")
	ErrCommitTx        = errors.Wrap(errors.ErrUnknown, "unexpected commit transaction error")
	ErrSerializationTx = pgconn.PgError{
		Severity: "ERROR",
		Code:     "40001",
		Detail:   "transaction serialization failure",
	}
)
