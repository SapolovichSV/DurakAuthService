package handlers

import (
	"context"

	"github.com/SapolovichSV/durak/auth/internal/cookier"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	"github.com/SapolovichSV/durak/auth/internal/storage"
)

type Services struct {
	Ctx       context.Context
	Repo      storage.Repo
	Cookier   cookier.Cookier
	Logger    logger.Logger
	SecretKey string
}

func New(ctx context.Context, logger logger.Logger, repo storage.Repo, cookier cookier.Cookier, secretKey string) *Services {
	return &Services{
		ctx,
		repo,
		cookier,
		logger,
		secretKey,
	}
}
