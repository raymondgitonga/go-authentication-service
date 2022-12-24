package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	claims := jwt.StandardClaims{
		IssuedAt:  1671807731,
		ExpiresAt: 1671894131,
		Id:        "id",
		Subject:   "subject",
		Issuer:    "auth",
	}

	key := []byte("")
	token, err := generateToken(claims, key)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzE4OTQxMzEsImp0aSI6ImlkIiwiaWF0IjoxNjcxODA3NzMxLCJpc3MiOiJhdXRoIiwic3ViIjoic3ViamVjdCJ9.nVPXUXj7YMP7C7rQXOhde6S10Te6yRjT7jCWGB5sgUY", token)
}

func TestParseToken(t *testing.T) {
	var claims = jwt.StandardClaims{}
	key := []byte("")

	t.Run("valid token", func(t *testing.T) {
		claims = jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
			Id:        "id",
			Subject:   "subject",
			Issuer:    "auth",
		}

		token, err := generateToken(claims, key)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		err = parseToken(token)
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

		err := parseToken(wrongToken)

		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("expired token", func(t *testing.T) {
		claims = jwt.StandardClaims{
			IssuedAt:  1671807731,
			ExpiresAt: 1671894131,
			Id:        "id",
			Subject:   "subject",
			Issuer:    "auth",
		}

		token, err := generateToken(claims, key)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		err = parseToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})
}
