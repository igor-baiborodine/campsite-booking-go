package bootstrap

import "github.com/stackus/errors"

var (
	ErrBeginTx = errors.Wrap(errors.ErrUnknown, "unexpected begin transaction error")
	ErrQuery   = errors.Wrap(errors.ErrUnknown, "unexpected query error")
	ErrRow     = errors.Wrap(errors.ErrUnknown, "unexpected rows error")
	ErrCommit  = errors.Wrap(errors.ErrUnknown, "unexpected commit error")
)
