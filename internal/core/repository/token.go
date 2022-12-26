package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type TokenRepository struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewTokenRepository(redis *redis.Client, logger *zap.Logger) *TokenRepository {
	return &TokenRepository{
		redis:  redis,
		logger: logger,
	}
}

func (r *TokenRepository) AddToken(ctx context.Context, encryptionKey string, tokenID int64) error {
	status := r.redis.ZAdd(ctx, "encryption_keys", redis.Z{
		Score:  float64(tokenID),
		Member: encryptionKey,
	})

	if status.Err() != nil {
		r.logger.Error("error at AddToken, could not save token", zap.String("error", status.Err().Error()))
		return fmt.Errorf("could not save token")
	}
	return nil
}

func (r *TokenRepository) GetLatestToken(ctx context.Context) (*redis.Z, error) {
	result := r.redis.ZRevRangeWithScores(ctx, "encryption_keys", 0, 0)

	if result.Err() != nil {
		r.logger.Error("error at GetToken", zap.String("error", result.Err().Error()))
		return nil, fmt.Errorf("could not save token")
	}

	token, err := result.Result()
	if err != nil {
		r.logger.Info("token not found", zap.String("error", err.Error()))
		return nil, fmt.Errorf("token not found")
	}

	if len(token) < 1 {
		return nil, fmt.Errorf("token not found")
	}

	return &token[0], nil
}

func (r *TokenRepository) GetToken(ctx context.Context, tokenID int64) (*redis.Z, error) {
	result := r.redis.ZRangeByScoreWithScores(ctx, "encryption_keys", &redis.ZRangeBy{
		Min:    strconv.Itoa(int(tokenID)),
		Max:    strconv.Itoa(int(tokenID)),
		Offset: 0,
		Count:  1,
	})

	if result.Err() != nil {
		r.logger.Error("error at GetToken", zap.String("error", result.Err().Error()))
		return nil, fmt.Errorf("could not save token")
	}

	token, err := result.Result()
	if err != nil {
		r.logger.Info("token not found", zap.String("error", err.Error()))
		return nil, fmt.Errorf("token not found")
	}

	if len(token) < 1 {
		return nil, fmt.Errorf("token not found")
	}

	return &token[0], nil
}

func (r *TokenRepository) ClearExpiredTokens(ctx context.Context) error {
	end := strconv.Itoa(int(time.Now().Add(-time.Hour * 48).Unix()))
	result := r.redis.ZRemRangeByScore(ctx, "encryption_keys", "0", end)

	if result.Err() != nil {
		return result.Err()
	}

	_, err := result.Result()
	if err != nil {
		return err
	}
	return nil
}
