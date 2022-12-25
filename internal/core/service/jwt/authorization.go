package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var encryptionKey = []byte("")

const ISSUER = "go-authentication"

type UserClaim struct {
	claims jwt.StandardClaims
}

type RepositoryUser interface {
	GetUser(name string) (string, error)
}

type AuthorizationService struct {
	repo   RepositoryUser
	logger *zap.Logger
}

func NewAuthorizationService(repo RepositoryUser, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{repo: repo, logger: logger}
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

	tokenString, err := a.generateToken(userClaims.claims, encryptionKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthorizationService) Validate(token string) error {
	signedToken := strings.Split(token, " ")[1]
	return a.parseToken(signedToken)
}

func (a *AuthorizationService) parseToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("error in Validate wrong signing algo used")
		}
		return encryptionKey, nil
	})

	if err != nil {
		a.logger.Error("error in parseToken, error parsing token", zap.String("error", err.Error()))
		return fmt.Errorf("error parsing token")
	}

	if !token.Valid {
		a.logger.Error("error in parseToken, invalid token", zap.String("error", err.Error()))
		return fmt.Errorf("token not valid")
	}

	return nil
}

func (a *AuthorizationService) generateToken(claims jwt.StandardClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		a.logger.Error("error in generateToken, could not generate signed token", zap.String("error", err.Error()))
		return "", fmt.Errorf("could not generate signed token")
	}

	return signedToken, nil
}
