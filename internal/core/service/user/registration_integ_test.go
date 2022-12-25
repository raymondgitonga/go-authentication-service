package user_test

import (
	"fmt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/user"
	"github.com/raymondgitonga/go-authentication/internal/core/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestRegistrationService_RegisterUser(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.NoError(t, err)
	defer func() {
		err = logger.Sync()
	}()
	assert.NoError(t, err)

	t.Run("successfully add user", func(t *testing.T) {
		repo := &user_mocks.RepositoryMock{AddUserFunc: func(name string, secret []byte) error {
			return nil
		}}

		service := user.NewRegistrationService(repo, logger)

		auth, err := service.RegisterUser("example-service")
		assert.NoError(t, err)
		assert.Equal(t, "example-service", auth.Name)
		assert.Equal(t, 36, len(auth.Secret))
	})

	t.Run("failed to add user", func(t *testing.T) {
		repo := &user_mocks.RepositoryMock{AddUserFunc: func(name string, secret []byte) error {
			return fmt.Errorf("error adding user")
		}}

		service := user.NewRegistrationService(repo, logger)

		auth, err := service.RegisterUser("example-service")
		assert.Error(t, err)
		assert.Nil(t, auth)
	})
}
