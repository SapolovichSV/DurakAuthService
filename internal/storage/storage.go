package storage

import (
	"context"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
)

type Repo interface {
	AddUser(ctx context.Context, email, username, password string) error
	GetUser(username string)
	DeleteUser()
	UpdateUser()
	UserByEmailAndPassword(email string, password string) (user.User, error)
}
