package httpserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/raymondgitonga/go-authentication/internal/core/repository"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt"
	"github.com/raymondgitonga/go-authentication/internal/core/service/user"
	"net/http"
)

type Handler struct {
	DB *sql.DB
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
	repo := repository.NewUserRepository(h.DB)
	service := jwt.NewAuthorizationService(repo)

	token, err := service.Authorize(key, secret)
	if err != nil {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error writing httpserver response: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(token))
	if err != nil {
		fmt.Printf("error writing httpserver response: %s", err)
	}
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	repo := repository.NewUserRepository(h.DB)
	service := jwt.NewAuthorizationService(repo)

	err := service.Validate(token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("could not authorize"))
		if err != nil {
			fmt.Printf("error writing httpserver response: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("authorized"))
	if err != nil {
		fmt.Printf("error writing httpserver response: %s", err)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	repo := repository.NewUserRepository(h.DB)
	userService := user.NewRegistrationService(repo)

	authReg, err := userService.RegisterUser(name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("could not register"))
		if err != nil {
			fmt.Printf("error writing httpserver response: %s", err)
		}
		return
	}

	resp, err := json.Marshal(authReg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not register"))
		if err != nil {
			fmt.Printf("error marshaling response: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Printf("error writing httpserver response: %s", err)
	}
}
