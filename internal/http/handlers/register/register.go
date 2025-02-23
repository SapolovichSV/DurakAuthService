package register

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/response"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/SapolovichSV/durak/auth/internal/storage"
	"github.com/go-playground/validator"
)

var ErrRegisterLogTopicName = "Register failed"

type userRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type strge interface {
	AddUser(ctx context.Context, email, username, password string) error
}
type Handler struct {
	log  logger.Logger
	repo strge
}

func New(services handlers.Services) Handler {
	return Handler{
		log:  services.Logger.WithGroup("register"),
		repo: services.Repo,
	}
}

// Register godoc
//
//	@Tags			Auth
//	@Summary		Register user
//	@Description	Registering user by email,username,password
//	@Accept			json
//	@Produce		json
//
//	@Param			RegisterData	body		userRegister	true	"need Email Username Password"
//
//	@Success		201				{object}	response.OkResp
//	@Failure		400				{object}	response.ErrorResp
//
//	@Failure		500				{object}	response.ErrorResp
//
//	@Router			/auth/register [POST]
func (c Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.log.Logger.Info(
		"Register",
		"with URI", r.Pattern,
	)
	user, err := userRegisterData(r)
	if err != nil {
		c.log.Logger.Error(ErrRegisterLogTopicName, "parse error:", err)

		responseError(map[string]error{"can't parse json data": err}, w, http.StatusBadRequest)
		return
	}
	if err := validator.New().Struct(user); err != nil {
		c.log.Logger.Warn(ErrRegisterLogTopicName, "validation error", err)

		errs := response.BeatifyValidationErrors(err.(validator.ValidationErrors))
		responseError(errs, w, http.StatusBadRequest)
		return
	}
	if err := c.repo.AddUser(ctx, user.Email, user.Username, user.Password); err != nil {
		if errors.Is(err, storage.ErrSuchUserExists{}) {
			c.log.Logger.Info(ErrRegisterLogTopicName, "user already exists", user.Email)

			responseError(map[string]error{"user already exists": err}, w, http.StatusBadRequest)
			return

		} else {
			c.log.Logger.Error(ErrRegisterLogTopicName, "add user error: ", err)

			responseError(map[string]error{"can't add user": err}, w, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(
		response.NewOkResp("created").JsonString(),
	)); err != nil {
		c.log.Logger.Error(ErrRegisterLogTopicName, "write error: ", err)

		responseError(map[string]error{"can't write payload to response's body": err}, w, http.StatusInternalServerError)
		return
	}
}
func userRegisterData(r *http.Request) (userRegister, error) {
	var user userRegister
	userData, err := io.ReadAll(r.Body)
	if err != nil {
		return user, err
	}
	if err := json.Unmarshal(userData, &user); err != nil {
		return user, err
	}
	return user, nil
}
func responseError(errs map[string]error, w http.ResponseWriter, statusCode int) {
	http.Error(w, response.NewErrorResp(errs).JsonString(), statusCode)
}
