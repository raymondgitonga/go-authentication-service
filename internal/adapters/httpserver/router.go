package httpserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/raymondgitonga/go-authentication/internal/core/repository"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt"
	token_manager "github.com/raymondgitonga/go-authentication/internal/core/service/token"
	"github.com/raymondgitonga/go-authentication/internal/core/service/user"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	db     *sql.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewHandler(db *sql.DB, redis *redis.Client, logger *zap.Logger) *Handler {
	return &Handler{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal("Healthy")
	if err != nil {
		fmt.Printf("error writing marshalling response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("error writing httpserver response: %s", err)
	}
}

func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	key, secret, _ := r.BasicAuth()

	repo := repository.NewUserRepository(h.db, h.logger)
	service := jwt.NewAuthorizationService(repo, h.logger)

	token, err := service.Authorize(key, secret)
	if err != nil {
		h.logger.Error("error at Authorize, token generation failed", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(token))
	if err != nil {
		h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
	}
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	repo := repository.NewUserRepository(h.db, h.logger)
	service := jwt.NewAuthorizationService(repo, h.logger)

	err := service.Validate(token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("could not authorize"))
		if err != nil {
			h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("authorized"))
	if err != nil {
		h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	repo := repository.NewUserRepository(h.db, h.logger)
	userService := user.NewRegistrationService(repo, h.logger)

	authReg, err := userService.RegisterUser(name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("could not register"))
		if err != nil {
			h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
		}
		return
	}

	resp, err := json.Marshal(authReg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not register"))
		if err != nil {
			h.logger.Error("error marshaling response", zap.String("error", err.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Error("error writing httpserver response", zap.String("error", err.Error()))
	}
}

func (h *Handler) RotateTokens(w http.ResponseWriter, _ *http.Request) {
	repo := repository.NewTokenRepository(h.redis, h.logger)
	tokenService := token_manager.NewTokenService(repo, h.logger)
	go tokenService.RotateEncryptionTokens(context.Background())

	response, err := json.Marshal("success")
	if err != nil {
		fmt.Printf("error writing marshalling response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("error writing httpserver response: %s", err)
	}
}
