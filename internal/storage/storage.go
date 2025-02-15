package storage

import (
	"context"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
)

type Repo interface {
	AddUser(ctx context.Context, user user.User) error
	GetUser(username string)
	DeleteUser()
	UpdateUser()
}
