package repository

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func isConflict(err error) bool {
	perr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}

	return perr.Code == pgerrcode.UniqueViolation
}
