package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var encryptionKey = []byte("")

const ISSUER = "go-authentication"

type RepositoryUser interface {
	GetUser(name string) (string, error)
}

type AuthorizationService struct {
	repo RepositoryUser
}

type UserClaim struct {
	claims jwt.StandardClaims
}

func NewAuthorizationService(repo RepositoryUser) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

func (a *AuthorizationService) Authorize(key, secret string) (string, error) {
	userClaims := UserClaim{claims: jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        key,
		Subject:   secret,
		Issuer:    ISSUER,
	}}

	encryptedSecret, err := a.repo.GetUser(key)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(encryptedSecret), []byte(secret))
	if err != nil {
		return "", err
	}

	tokenString, err := generateToken(userClaims.claims, encryptionKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthorizationService) Validate(token string) error {
	signedToken := strings.Split(token, " ")[1]
	return parseToken(signedToken)
}

func parseToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("error in Validate wrong signing algo used")
		}
		return encryptionKey, nil
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
