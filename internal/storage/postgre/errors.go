package postgre

import (
	"errors"

	"github.com/SapolovichSV/durak/auth/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

// need to implement Unwrap() method for errors.Is() to work
func remakeError(err error, at string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return storage.ErrSuchUserExists{
				StorErr: storage.StorErr{At: at, Err: err},
				Email:   pgErr.Detail,
			}
		default:
			return storage.StorErr{At: at, Err: err}
		}
	}
	return storage.StorErr{At: at, Err: err}
}
