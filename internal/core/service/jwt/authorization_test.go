package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.NoError(t, err)
	defer func() {
		err = logger.Sync()
	}()
	assert.NoError(t, err)

	claims := jwt.StandardClaims{
		IssuedAt:  1671807731,
		ExpiresAt: 1671894131,
		Id:        "id",
		Subject:   "subject",
		Issuer:    "auth",
	}

	key := []byte("")
	service := NewAuthorizationService(nil, logger)
	token, err := service.generateToken(claims, key)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzE4OTQxMzEsImp0aSI6ImlkIiwiaWF0IjoxNjcxODA3NzMxLCJpc3MiOiJhdXRoIiwic3ViIjoic3ViamVjdCJ9.nVPXUXj7YMP7C7rQXOhde6S10Te6yRjT7jCWGB5sgUY", token)
}

func TestParseToken(t *testing.T) {
	var claims = jwt.StandardClaims{}
	key := []byte("")

	logger, err := zap.NewProduction()
	assert.NoError(t, err)
	defer func() {
		err = logger.Sync()
	}()
	assert.NoError(t, err)

	service := NewAuthorizationService(nil, logger)

	t.Run("valid token", func(t *testing.T) {
		claims = jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
			Id:        "id",
			Subject:   "subject",
			Issuer:    "auth",
		}

		token, err := service.generateToken(claims, key)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		err = service.parseToken(token)
		assert.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		claims = jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
			Id:        "id",
			Subject:   "subject",
			Issuer:    "auth",
		}

		wrongToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzE5MTM3OTAsImlhdCI6MTY3MTkxMzY3MCwic3ViIjoiZW1pbEBnbWFpbC5jb20ifQ.fDE0KdLE8bU2TxY5cDTkFihtCCRUKPHJDS30UQd-zy0"

		err = service.parseToken(wrongToken)
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		claims = jwt.StandardClaims{
			IssuedAt:  1671807731,
			ExpiresAt: 1671894131,
			Id:        "id",
			Subject:   "subject",
			Issuer:    "auth",
		}

		token, err := service.generateToken(claims, key)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		err = service.parseToken(token)
		assert.Error(t, err)
	})
}
