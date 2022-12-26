package user

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	AddUser(name string, secret []byte) error
}

type RegistrationService struct {
	repo   Repository
	logger *zap.Logger
}

func NewRegistrationService(repo Repository, logger *zap.Logger) *RegistrationService {
	return &RegistrationService{repo: repo, logger: logger}
}

func (r *RegistrationService) RegisterUser(name string) (*dormain.AuthResponse, error) {
	secret, err := uuid.NewUUID()
	if err != nil {
		r.logger.Error("error in RegisterUser, error generating user secret", zap.String("error", err.Error()))
		return nil, fmt.Errorf("error generating user secret")
	}

	encryptedKey, err := bcrypt.GenerateFromPassword([]byte(secret.String()), 5)
	if err != nil {
		r.logger.Error("error in RegisterUser, error generating user secret", zap.String("error", err.Error()))
		return nil, fmt.Errorf("error generating user secret")
	}

	err = r.repo.AddUser(name, encryptedKey)
	if err != nil {
		r.logger.Error("error in RegisterUser, error generating user secret", zap.String("error", err.Error()))
		return nil, fmt.Errorf("error generating user secret")
	}

	authResp := &dormain.AuthResponse{
		Name:   name,
		Secret: secret.String(),
	}

	return authResp, nil
}
