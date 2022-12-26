package repository_test

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/raymondgitonga/go-authentication-service/internal/adapters/db"
	"github.com/raymondgitonga/go-authentication-service/internal/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"strings"
	"testing"
)

func TestUserRepository(t *testing.T) {
	dbCLient, postgres := setupTestDatabase(t)
	defer destroyDB(postgres)

	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	userRepo := repository.NewUserRepository(dbCLient, observedLogger)

	t.Run("successfully add and retrieve user", func(t *testing.T) {
		err := userRepo.AddUser("test", []byte("test"))
		assert.NoError(t, err)

		secret, err := userRepo.GetUser("test")
		assert.NoError(t, err)
		assert.Equal(t, "test", secret)
	})

	t.Run("fail to get user", func(t *testing.T) {
		secret, err := userRepo.GetUser("test1")
		assert.Error(t, err)
		assert.Empty(t, secret)
	})
}

func setupTestDatabase(t *testing.T) (*sql.DB, *testcontainers.LocalDockerCompose) {
	postgres := testcontainers.NewLocalDockerCompose([]string{"../../../test_docker_files/docker-compose.yml"},
		strings.ToLower(uuid.New().String()))
	postgres.WithCommand([]string{"up", "-d"}).Invoke()
	postgresURL := "postgres://postgres:postgres@localhost:9876/postgres?sslmode=disable"

	conn, err := db.NewClient(context.Background(), postgresURL)
	assert.NoError(t, err)

	err = db.RunMigrations(conn, "postgres")
	assert.NoError(t, err)

	return conn, postgres
}

func destroyDB(compose *testcontainers.LocalDockerCompose) {
	compose.WithCommand([]string{"down"}).Invoke()
}
