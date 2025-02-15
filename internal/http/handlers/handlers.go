package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/SapolovichSV/durak/auth/internal/storage"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	repo      storage.Repo
	logger    logger.Logger
	ctx       context.Context
	secretKey string
}

func New(ctx context.Context, logger logger.Logger, repo storage.Repo, secretKey string) *Controller {
	return &Controller{
		repo,
		logger,
		ctx,
		secretKey,
	}
}

var ErrorLogTopicName = "Register failed"

// TODO add jwt
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	c.logger.Logger.Info(
		"Register",
		"with URI", r.Pattern,
	)
	userData := []byte{}
	if _, err := r.Body.Read(userData); err != nil {
		c.logger.Logger.Error(
			ErrorLogTopicName,
			"read error:", err,
		)
		http.Error(w, fmt.Errorf("can't read form data: %w", err).Error(), http.StatusBadRequest)
		return
	}
	var user user.User
	if err := json.Unmarshal(userData, &user); err != nil {
		c.logger.Logger.Error(
			ErrorLogTopicName,
			"parse erorr:", err,
		)
		http.Error(w, fmt.Errorf("can't parse form data %w", err).Error(), http.StatusBadRequest)
		return
	}
	if err := validator.New(nil).Struct(user); err != nil {
		c.logger.Logger.Warn(ErrorLogTopicName,
			"validation error", err,
		)
		http.Error(w, fmt.Errorf("invalid user data: %w", err).Error(), http.StatusBadRequest)
		return
	}
	if err := c.repo.AddUser(c.ctx, user); err != nil {
		c.logger.Logger.Error(ErrorLogTopicName,
			"add user error: ", err,
		)
		http.Error(w, fmt.Errorf("can't add user to repo: %w", err).Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
