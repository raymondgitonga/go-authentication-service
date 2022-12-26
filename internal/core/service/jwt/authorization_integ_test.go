package jwt_test

import (
	"fmt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthorizationService(t *testing.T) {
	t.Run("successfully get auth token", func(t *testing.T) {
		userRepo := &jwt_mocks.RepositoryUserMock{
			GetUserFunc: func(name string) (string, error) {
				return "$2a$05$gKFDk14pCAIpk8uBiBRu2euxDt97BFSSABCZ2OPYCHkX0ZSRCT8aq", nil
			},
		}

		service := jwt.NewAuthorizationService(userRepo, nil, nil)

		token, err := service.Authorize("key", "sample")
		assert.NoError(t, err)
		assert.NotNil(t, token)
	})

	t.Run("successfully get auth token", func(t *testing.T) {
		userRepo := &jwt_mocks.RepositoryUserMock{
			GetUserFunc: func(name string) (string, error) {
				return "", fmt.Errorf("error getting user")
			},
		}

		service := jwt.NewAuthorizationService(userRepo, nil, nil)

		token, err := service.Authorize("key", "secret")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
