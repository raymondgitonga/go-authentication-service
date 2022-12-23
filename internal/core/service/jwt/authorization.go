package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"time"
)

type Authorization struct {
	authRequest dormain.AuthRequest
}

func NewAuthorization(authRequest dormain.AuthRequest) *Authorization {
	return &Authorization{authRequest: authRequest}
}

func (a *Authorization) Authorize() (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   a.authRequest.Email,
	}

	tokenString, err := generateToken(claims, []byte(""))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func generateToken(claims jwt.StandardClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error in generateToken, could not generate signed token: %w", err)
	}

	return signedToken, nil
}
