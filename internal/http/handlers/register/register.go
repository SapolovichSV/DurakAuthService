package register

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/response"
	handler "github.com/SapolovichSV/durak/auth/internal/http/handlers"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/go-playground/validator"
)

var ErrRegisterLogTopicName = "Register failed"

type userRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=storage
type storage interface {
	AddUser(ctx context.Context, email, username, password string) error
}
type Handler struct {
	log  logger.Logger
	repo storage
	ctx  context.Context
}

func New(services handler.Services) Handler {
	return Handler{
		log:  services.Logger.WithGroup("register"),
		repo: services.Repo,
	}
}

// TODO beautify validator errors
// TODO if map[string]error{cause:err} work refactor code on down
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
	c.log.Logger.Info(
		"Register",
		"with URI", r.Pattern,
	)
	userData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		c.log.Logger.Error(
			ErrRegisterLogTopicName,
			"read error:", err,
		)
		errs := make(map[string]error, 1)
		errs["can't read request data"] = err
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}

	var user userRegister
	if err := json.Unmarshal(userData, &user); err != nil {
		c.log.Logger.Error(
			ErrRegisterLogTopicName,
			"parse erorr:", err,
		)
		errs := make(map[string]error, 1)
		errs["can't parse json data"] = err
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	//TODO beatify validator errors
	if err := validator.New().Struct(user); err != nil {
		c.log.Logger.Warn(
			ErrRegisterLogTopicName,
			"validation error", err,
		)

		errs := make(map[string]error, 1)
		errs["invalid user data"] = err
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	if err := c.repo.AddUser(c.ctx, user.Email, user.Username, user.Password); err != nil {
		c.log.Logger.Error(
			ErrRegisterLogTopicName,
			"add user error: ", err,
		)
		errs := make(map[string]error, 1)
		errs["can't add user to repo"] = err
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(
		response.NewOkResp("created").JsonString(),
	)); err != nil {
		errs := make(map[string]error, 1)
		errs["can't write payload to response's body"] = err
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusInternalServerError)
		return
	}
}
