package login

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/response"
	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/go-playground/validator"
)

var ErrLoginLogTopicName = "Login failed"

type userLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type storage interface {
	UserByEmailAndPassword(email, password string) (user.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=cookier
type cookier interface {
	Login(user user.User, w http.ResponseWriter) error
}
type Handler struct {
	log     logger.Logger
	repo    storage
	cookier cookier
}

func New(svc handlers.Services) *Handler {
	return &Handler{
		log:     svc.Logger.WithGroup("login"),
		repo:    svc.Repo,
		cookier: svc.Cookier,
	}
}

// Login godoc
//
//	@Tags			Auth
//	@Summary		login user
//	@Description	Login user by email,password
//	@Accept			json
//	@Produce		json
//	@Param			LoginData	body		userLogin	true	"Need only email and password"
//
//	@Success		202			{object}	user.User
//	@Failure		400			{object}	response.ErrorResp
//	@Failure		401			{object}	response.ErrorResp
//	@Failure		500			{object}	response.ErrorResp
//	@Router			/auth/login [POST]
func (c *Handler) Login(w http.ResponseWriter, r *http.Request) {
	c.log.Logger.Info(
		"Login",
		"with URI", r.Pattern,
	)
	userData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		cause := "can't parse user data"
		c.log.Logger.Error(
			ErrLoginLogTopicName,
			cause, err,
		)
		errs := map[string]error{cause: err}
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	var userLoginInfo userLogin
	if err := json.Unmarshal(userData, &userLoginInfo); err != nil {
		cause := "can't parse json user data"
		c.log.Logger.Error(
			ErrLoginLogTopicName,
			cause, err,
		)
		errs := map[string]error{cause: err}
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	if err := validator.New().Struct(userLoginInfo); err != nil {
		cause := "invalid user data"
		c.log.Logger.Warn(
			ErrLoginLogTopicName,
			cause, err,
		)
		errs := map[string]error{cause: err}
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	//TODO make error intospection
	user, err := c.repo.UserByEmailAndPassword(userLoginInfo.Email, userLoginInfo.Password)
	if err != nil {
		cause := "can't get user"
		abstractError := errors.New("incorrect userInfo,or not exist such user")
		c.log.Logger.Warn(
			ErrLoginLogTopicName,
			cause, abstractError,
		)
		errs := map[string]error{cause: abstractError}
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusUnauthorized)
		return
	}
	if err := c.cookier.Login(user, w); err != nil {
		cause := "can't set auth-cookie to user"
		c.log.Logger.Warn(
			ErrLoginLogTopicName,
			cause, err,
		)
		errs := map[string]error{cause: err}
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusInternalServerError)
		return
	}

}
