package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

const ISSUER = "go-authentication-service"

type UserClaim struct {
	claims jwt.StandardClaims
}

type UserRepository interface {
	GetUser(name string) (string, error)
}

type TokenRepository interface {
	GetToken(ctx context.Context, tokenID int64) (*redis.Z, error)
	GetLatestToken(ctx context.Context) (*redis.Z, error)
}

type AuthorizationService struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
	logger    *zap.Logger
}

func NewAuthorizationService(userRepo UserRepository, tokenRepo TokenRepository, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{userRepo: userRepo, tokenRepo: tokenRepo, logger: logger}
}

func (a *AuthorizationService) Authorize(key, secret string) (string, error) {
	encryptionToken, err := a.tokenRepo.GetLatestToken(context.Background())
	if err != nil {
		return "", err
	}

	userClaims := UserClaim{claims: jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        strconv.Itoa(int(encryptionToken.Score)),
		Subject:   secret,
		Issuer:    ISSUER,
	}}

	encryptedSecret, err := a.userRepo.GetUser(key)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(encryptedSecret), []byte(secret))
	if err != nil {
		return "", err
	}

	tokenString, err := a.generateToken(userClaims.claims, []byte(encryptionToken.Member.(string)))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthorizationService) Validate(token string) error {
	userClaims := &jwt.StandardClaims{}
	signedToken := strings.Split(token, " ")[1]
	payload := strings.Split(signedToken, ".")[1]

	decodedPayload, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodedPayload, &userClaims)
	if err != nil {
		return err
	}

	tokenID, err := strconv.Atoi(userClaims.Id)
	if err != nil {
		return err
	}

	tokenResp, err := a.tokenRepo.GetToken(context.Background(), int64(tokenID))
	if err != nil {
		return err
	}

	return a.parseToken(signedToken, []byte(tokenResp.Member.(string)))
}

func (a *AuthorizationService) parseToken(signedToken string, encryptionKey []byte) error {
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
