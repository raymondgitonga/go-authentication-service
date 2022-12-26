package repository_test

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/raymondgitonga/go-authentication-service/internal/adapters/cache"
	"github.com/raymondgitonga/go-authentication-service/internal/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"strings"
	"testing"
)

func TestTokenRepository(t *testing.T) {
	redisClient, container := setupTestCache(t)
	defer destroyCache(container)

	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	tokenRepo := repository.NewTokenRepository(redisClient, observedLogger)

	t.Run("successfully add and retrieve latest token", func(t *testing.T) {
		err := tokenRepo.AddToken(context.Background(), "key", int64(123456))
		assert.NoError(t, err)

		token, err := tokenRepo.GetLatestToken(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "key", token.Member)
		assert.Equal(t, 123456, int(token.Score))
	})

	t.Run("successfully add and retrieve specific token", func(t *testing.T) {
		err := tokenRepo.AddToken(context.Background(), "key1", int64(1234567))
		assert.NoError(t, err)

		token, err := tokenRepo.GetToken(context.Background(), int64(1234567))
		assert.NoError(t, err)
		assert.Equal(t, "key1", token.Member)
		assert.Equal(t, 1234567, int(token.Score))
	})

	t.Run("fail retrieve token", func(t *testing.T) {
		token, err := tokenRepo.GetToken(context.Background(), int64(12345678))
		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("fail retrieve specific token", func(t *testing.T) {
		token, err := tokenRepo.GetToken(context.Background(), int64(12345678))
		assert.Error(t, err)
		assert.Nil(t, token)
	})
}

func setupTestCache(t *testing.T) (*redis.Client, *testcontainers.LocalDockerCompose) {
	redisContainer := testcontainers.NewLocalDockerCompose([]string{"../../../test_docker_files/redis-compose.yml"},
		strings.ToLower(uuid.New().String()))
	redisContainer.WithCommand([]string{"up", "-d"}).Invoke()

	conn, err := cache.NewClient(context.Background(), "localhost:6378", "redis")
	assert.NoError(t, err)

	return conn, redisContainer
}

func destroyCache(compose *testcontainers.LocalDockerCompose) {
	compose.WithCommand([]string{"down"}).Invoke()
}
