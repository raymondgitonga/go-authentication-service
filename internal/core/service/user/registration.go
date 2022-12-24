package user

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	AddUser(name string, secret []byte) error
}

type RegistrationService struct {
	repo Repository
}

func NewRegistrationService(repo Repository) *RegistrationService {
	return &RegistrationService{repo: repo}
}

func (r *RegistrationService) RegisterUser(name string) (*dormain.AuthRequest, error) {
	secret, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating user secret")
	}

	encryptedKey, err := bcrypt.GenerateFromPassword([]byte(secret.String()), 5)
	if err != nil {
		return nil, fmt.Errorf("error generating user secret")
	}

	err = r.repo.AddUser(name, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("error generating user secret")
	}

	authReq := &dormain.AuthRequest{
		Name:   name,
		Secret: secret.String(),
	}

	return authReq, nil
}
