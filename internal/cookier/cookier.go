package cookier

import (
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
)

type Cookier interface {
	Auth()
	Login(user user.User, w http.ResponseWriter) error
	Logout()
}
