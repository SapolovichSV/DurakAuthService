package postgre

import (
	"context"
	"errors"

	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type hasher interface {
	hash(string) (string, error)
	unhash(string) (string, error)
}
type RepoPostgre struct {
	pgpool *pgxpool.Pool
	hasher hasher
	logger logger.Logger
}

func New(pgpool *pgxpool.Pool, hasher hasher, logger logger.Logger) *RepoPostgre {
	return &RepoPostgre{
		pgpool: pgpool,
		hasher: hasher,
		logger: logger.WithGroup("Postgre"),
	}
}

// TODO implement error types
func (r *RepoPostgre) AddUser(ctx context.Context, email, username, password string) error {
	var ErrLogTopicName = "Error at AddUser"
	r.logger.Logger.Info("Starts transaction")

	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.TxIsoLevel(pgx.ReadUncommitted),
		AccessMode: pgx.TxAccessMode(pgx.ReadWrite),
	})
	defer tx.Rollback(ctx)

	if err != nil {
		r.logger.Logger.Warn(
			ErrLogTopicName,
			"Can't start transaction", err,
		)
		return errors.New("can't start transaction")
	}

	hashedPass, err := r.hasher.hash(password)
	if err != nil {
		return err
	}

	sqlExec := `INSERT INTO users (email,username,passwordHASH)
	VALUES ($1,$2,$3);`

	res, err := tx.Exec(ctx, sqlExec, email, username, hashedPass)
	if err != nil {
		r.logger.Logger.Error(
			ErrLogTopicName,
			"Can't execute transaction", err)
		return errors.New("can't exec transaction")
	}

	if !(res.RowsAffected() == 1 && res.Insert()) {
		r.logger.Logger.Warn(
			ErrLogTopicName,
			"something goes wrong", "user not added or added incorectly")
		return errors.New("something bad at adduser")
	}

	r.logger.Logger.Info("Succesful ended transaction")
	return nil
}
