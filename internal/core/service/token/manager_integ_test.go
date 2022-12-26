package manager_test

import (
	"context"
	"fmt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/token"
	"github.com/raymondgitonga/go-authentication/internal/core/service/token/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestTokenService(t *testing.T) {
	t.Run("successfully rotate encryption keys", func(t *testing.T) {
		tokenRepo := &manager_mocks.TokenRepositoryMock{
			AddTokenFunc: func(ctx context.Context, encryptionKey string, tokenID int64) error {
				return nil
			},
			ClearExpiredTokensFunc: func(ctx context.Context) error {
				return nil
			},
		}

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)

		tokenService := manager.NewTokenService(tokenRepo, observedLogger)
		tokenService.RotateEncryptionTokens(context.Background())

		assert.Equal(t, "successfully added token", observedLogs.All()[0].Message)
	})

	t.Run("failed to rotate encryption keys, new token not added", func(t *testing.T) {
		tokenRepo := &manager_mocks.TokenRepositoryMock{
			AddTokenFunc: func(ctx context.Context, encryptionKey string, tokenID int64) error {
				return fmt.Errorf("could not save token")
			},
			ClearExpiredTokensFunc: func(ctx context.Context) error {
				return nil
			},
		}

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)

		tokenService := manager.NewTokenService(tokenRepo, observedLogger)
		tokenService.RotateEncryptionTokens(context.Background())

		assert.Equal(t, "error adding encryption token", observedLogs.All()[0].Message)
	})

	t.Run("failed to rotate encryption keys, expired tokens not cleared", func(t *testing.T) {
		tokenRepo := &manager_mocks.TokenRepositoryMock{
			AddTokenFunc: func(ctx context.Context, encryptionKey string, tokenID int64) error {
				return nil
			},
			ClearExpiredTokensFunc: func(ctx context.Context) error {
				return fmt.Errorf("could not clear tokens")
			},
		}

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)

		tokenService := manager.NewTokenService(tokenRepo, observedLogger)
		tokenService.RotateEncryptionTokens(context.Background())

		assert.Equal(t, "error clearing expired tokens", observedLogs.All()[0].Message)
	})
}
