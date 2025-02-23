package postgre

import (
	"context"
	"errors"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Hasher interface {
	Hash(string) (string, error)
	Unhash(string) (string, error)
}
type RepoPostgre struct {
	pgpool *pgxpool.Pool
	hasher Hasher
	logger logger.Logger
}

func New(pgpool *pgxpool.Pool, hasher Hasher, logger logger.Logger) *RepoPostgre {
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
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			cause := "can't rollback transaction at defer"
			r.logger.Logger.Error(
				ErrLogTopicName,
				cause, err,
			)
		}
	}()

	if err != nil {
		cause := "Can't start transaction"
		r.logger.Logger.Warn(
			ErrLogTopicName,
			cause, err,
		)
		return remakeError(err, cause)
	}

	hashedPass, err := r.hasher.Hash(password)

	if err != nil {
		cause := "Can't hash password"
		r.logger.Logger.Error(
			ErrLogTopicName,
			cause, err,
		)
		return err
	}

	sqlExec := `INSERT INTO users (email,username,passwordHASH)
	VALUES ($1,$2,$3);`

	res, err := tx.Exec(ctx, sqlExec, email, username, hashedPass)
	if err != nil {
		cause := "Can't execute transaction"
		r.logger.Logger.Error(
			ErrLogTopicName,
			cause, err,
		)
		return remakeError(err, cause)
	}

	if !(res.RowsAffected() == 1 && res.Insert()) {
		cause := "user not added or added incorectly"
		err := errors.New("somenthing bad at adduser")
		r.logger.Logger.Warn(
			ErrLogTopicName,
			cause, err,
		)
		return remakeError(err, cause)
	}

	r.logger.Logger.Info("Succesful ended transaction")
	return tx.Commit(ctx)
}
func (r *RepoPostgre) GetUser(username string) {
	panic("implement me")
}
func (r *RepoPostgre) DeleteUser() {
	panic("implement me")
}
func (r *RepoPostgre) UpdateUser() {
	panic("implement me")
}
func (r *RepoPostgre) UserByEmailAndPassword(email string, password string) (user.User, error) {
	panic("implement me")
}
