package manager

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type TokenRepository interface {
	AddToken(ctx context.Context, encryptionKey string, tokenID int64) error
	ClearExpiredTokens(ctx context.Context) error
}

type TokenService struct {
	repo   TokenRepository
	logger *zap.Logger
}

func NewTokenService(repo TokenRepository, logger *zap.Logger) *TokenService {
	return &TokenService{repo: repo, logger: logger}
}

func (s *TokenService) RotateEncryptionTokens(ctx context.Context) {
	encryptionToken, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error("error generating encryption token", zap.String("error", err.Error()))
		return
	}

	tokenID := time.Now().Unix()
	err = s.repo.AddToken(ctx, encryptionToken.String(), tokenID)
	if err != nil {
		s.logger.Error("error adding encryption token", zap.String("error", err.Error()))
		return
	}

	err = s.repo.ClearExpiredTokens(ctx)
	if err != nil {
		s.logger.Error("error clearing expired tokens", zap.String("error", err.Error()))
		return
	}

	s.logger.Info("successfully added token", zap.Int("tokenID", int(tokenID)))
}
