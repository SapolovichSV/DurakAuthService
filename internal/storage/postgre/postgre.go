package postgre

import (
	"context"
	"errors"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/jackc/pgx/v5"
)

type Hasher interface {
	Hash(string) (string, error)
}
type pgPool interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}
type RepoPostgre struct {
	pgpool pgPool
	hasher Hasher
	logger logger.Logger
}

func New(pgpool pgPool, hasher Hasher, logger logger.Logger) *RepoPostgre {
	return &RepoPostgre{
		pgpool: pgpool,
		hasher: hasher,
		logger: logger.WithGroup("Postgre"),
	}
}

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
func (r *RepoPostgre) UserByEmailAndPassword(ctx context.Context, email string, password string) (user.User, error) {
	var ErrLogTopicName = "Error at UserByEmailAndPassword"
	r.logger.Logger.Info("starts transaction")

	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		cause := "can't start transaction"
		r.logger.Logger.Error(ErrLogTopicName, cause, err)
		return user.User{}, remakeError(err, cause)
	}

	sqlQuery := `SELECT id,email,username,status,user_role,passwordHASH FROM users WHERE email=$1;`
	row := tx.QueryRow(ctx, sqlQuery, email)

	var (
		id             int
		username       string
		status         string
		user_role      string
		passHashFromDB string
	)
	//Really vpadlu pisat source ili che to tam
	if err := row.Scan(&id, &email, &username, &status, &user_role, &passHashFromDB); err != nil {
		cause := "error at getting user Info"
		r.logger.Logger.Debug(ErrLogTopicName, cause, err)
		return user.User{}, remakeError(err, "can't get user info")
	}

	passwordHASH, err := r.hasher.Hash(password)
	if err != nil {
		cause := "can't hash password"
		r.logger.Logger.Error(ErrLogTopicName, cause, err)
	}
	if passHashFromDB != passwordHASH {
		cause := "incorrect password"
		r.logger.Logger.Debug(ErrLogTopicName, cause, errors.New("bad input pass have hash"+passHashFromDB+"want hash"+passwordHASH))
		return user.User{}, remakeError(errors.New("incorrect pass"), cause)
	}

	userInfo := user.User{
		ID:       id,
		Email:    email,
		Username: username,
		Password: password,
		Status: func() user.Status {
			userStatus, err := user.BuildStatus(status)
			if err != nil {
				r.logger.Logger.Error(ErrLogTopicName, "impossible status", errors.New("get status"+status+"???"))
				userStatus, _ = user.BuildStatus("offline")
			}
			return userStatus
		}(),
	}
	return userInfo, nil
}
