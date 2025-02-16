package controller

import (
	"context"

	"github.com/SapolovichSV/durak/auth/internal/cookier"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/SapolovichSV/durak/auth/internal/storage"
)

type Controller struct {
	ctx       context.Context
	repo      storage.Repo
	cookier   cookier.Cookier
	logger    logger.Logger
	secretKey string
}

func New(ctx context.Context, logger logger.Logger, repo storage.Repo, cookier cookier.Cookier, secretKey string) *Controller {
	return &Controller{
		ctx,
		repo,
		cookier,
		logger,
		secretKey,
	}
}
