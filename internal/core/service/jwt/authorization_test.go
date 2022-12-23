package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	claims := jwt.StandardClaims{
		IssuedAt:  1671807731,
		ExpiresAt: 1671894131,
		Subject:   "email@gmail.com",
	}

	secret := []byte("")
	token, err := generateToken(claims, secret)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzE4OTQxMzEsImlhdCI6MTY3MTgwNzczMSwic3ViIjoiZW1haWxAZ21haWwuY29tIn0.pPGxkoXX8AE861CVvjTvf5lhLGJwRCRyhWZLDeQauGY", token)
}

func TestParseToken(t *testing.T) {
	claims := jwt.StandardClaims{
		IssuedAt:  1671807731,
		ExpiresAt: 1671894131,
		Subject:   "email@gmail.com",
	}

	secret := []byte("")
	token, err := generateToken(claims, secret)
	claim, err := parseToken(token)

	assert.Equal(t, "email@gmail.com", claim.Subject)
	assert.NoError(t, err)
}
