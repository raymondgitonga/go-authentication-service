package jwt_test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthorizationService(t *testing.T) {
	t.Run("successfully get auth token", func(t *testing.T) {
		tokenRepo := &jwt_mocks.TokenRepositoryMock{
			GetLatestTokenFunc: func(ctx context.Context) (*redis.Z, error) {
				return &redis.Z{
					Score:  1672008811,
					Member: "f46f9f82-84a6-11ed-b27d-0a8d12596476",
				}, nil
			},
		}
		userRepo := &jwt_mocks.RepositoryUserMock{
			GetUserFunc: func(name string) (string, error) {
				return "$2a$05$gKFDk14pCAIpk8uBiBRu2euxDt97BFSSABCZ2OPYCHkX0ZSRCT8aq", nil
			},
		}

		service := jwt.NewAuthorizationService(userRepo, tokenRepo, nil)

		token, err := service.Authorize("key", "sample")
		assert.NoError(t, err)
		assert.NotNil(t, token)
	})

	t.Run("fail get auth token, user not found", func(t *testing.T) {
		tokenRepo := &jwt_mocks.TokenRepositoryMock{
			GetLatestTokenFunc: func(ctx context.Context) (*redis.Z, error) {
				return &redis.Z{
					Score:  1672008811,
					Member: "f46f9f82-84a6-11ed-b27d-0a8d12596476",
				}, nil
			},
		}

		userRepo := &jwt_mocks.RepositoryUserMock{
			GetUserFunc: func(name string) (string, error) {
				return "", fmt.Errorf("error getting user")
			},
		}

		service := jwt.NewAuthorizationService(userRepo, tokenRepo, nil)

		token, err := service.Authorize("key", "secret")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("fail get auth token, token not found", func(t *testing.T) {
		tokenRepo := &jwt_mocks.TokenRepositoryMock{
			GetLatestTokenFunc: func(ctx context.Context) (*redis.Z, error) {
				return nil, fmt.Errorf("token not found")
			},
		}

		userRepo := &jwt_mocks.RepositoryUserMock{
			GetUserFunc: func(name string) (string, error) {
				return "$2a$05$gKFDk14pCAIpk8uBiBRu2euxDt97BFSSABCZ2OPYCHkX0ZSRCT8aq", nil
			},
		}

		service := jwt.NewAuthorizationService(userRepo, tokenRepo, nil)

		token, err := service.Authorize("key", "secret")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
