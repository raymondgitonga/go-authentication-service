package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"time"
)

var secret = []byte("")

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

	tokenString, err := generateToken(claims, secret)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func parseToken(signedToken string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("error in Validate wrong signing algo used")
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in Validate, error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("error in Validate: token not valid")
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func generateToken(claims jwt.StandardClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error in generateToken, could not generate signed token: %w", err)
	}

	return signedToken, nil
}
