package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"strings"
	"time"
)

var key = []byte("")

const ISSUER = "go-authentication"

type Authorization struct {
	authRequest dormain.AuthRequest
}

type UserClaim struct {
	claims jwt.StandardClaims
}

func NewAuthorization(authRequest dormain.AuthRequest) *Authorization {
	return &Authorization{authRequest: authRequest}
}

func (a *Authorization) Authorize() (string, error) {
	userClaims := UserClaim{claims: jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        a.authRequest.Key,
		Subject:   a.authRequest.Secret,
		Issuer:    ISSUER,
	}}

	tokenString, err := generateToken(userClaims.claims, key)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (a *Authorization) Validate() error {
	signedToken := strings.Split(a.authRequest.Token, " ")[1]
	return parseToken(signedToken)
}

func parseToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("error in Validate wrong signing algo used")
		}
		return key, nil
	})

	if err != nil {
		return fmt.Errorf("error in Validate, error parsing token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("error in Validate: token not valid")
	}

	return nil
}

func generateToken(claims jwt.StandardClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error in generateToken, could not generate signed token: %w", err)
	}

	return signedToken, nil
}
