package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/response"
	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/go-playground/validator"
)

var ErrorLogTopicName = "Register failed"

// TODO beautify validator errors

// Register godoc
//
//	@Tags			Auth
//	@Summary		Register user
//	@Description	Registering user by email,username,password
//	@Accept			json
//	@Produce		json
//
//	@Param			RegisterData	body		user.User	true	"need Email Username Password"
//
//	@Success		201				{object}	response.OkResp
//	@Failure		400				{object}	response.ErrorResp
//
//	@Failure		500				{object}	response.ErrorResp
//
//	@Router			/auth/register [POST]
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	c.logger.Logger.Info(
		"Register",
		"with URI", r.Pattern,
	)
	userData, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Logger.Error(
			ErrorLogTopicName,
			"read error:", err,
		)
		errs := make(map[error]string, 1)
		errs[err] = "can't read request data"
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}

	var user user.User
	if err := json.Unmarshal(userData, &user); err != nil {
		c.logger.Logger.Error(
			ErrorLogTopicName,
			"parse erorr:", err,
		)
		errs := make(map[error]string, 1)
		errs[err] = "can't parse json data"
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	//TODO beatify validator errors
	if err := validator.New().Struct(user); err != nil {
		err := fmt.Errorf("%w", err.(validator.ValidationErrors))
		c.logger.Logger.Warn(ErrorLogTopicName,
			"validation error", err,
		)
		//TODO refactor make hashing string, not errors
		errs := make(map[error]string, 1)
		errs[err] = "invalid user data"
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusBadRequest)
		return
	}
	if err := c.repo.AddUser(c.ctx, user); err != nil {
		c.logger.Logger.Error(ErrorLogTopicName,
			"add user error: ", err,
		)
		errs := make(map[error]string, 1)
		errs[err] = "can't add user to repo"
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(
		response.NewOkResp("created").JsonString(),
	)); err != nil {
		errs := make(map[error]string, 1)
		errs[err] = "can't write payload to response's body"
		http.Error(w, response.NewErrorResp(errs).JsonString(), http.StatusInternalServerError)
		return
	}
}
