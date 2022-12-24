package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/raymondgitonga/go-authentication/internal/core/dormain"
	"github.com/raymondgitonga/go-authentication/internal/core/service/jwt"
	"net/http"
)

type Handler struct{}

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

	authRequest := dormain.AuthRequest{Key: key, Secret: secret}

	token, err := jwt.NewAuthorization(authRequest).Authorize()
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
	authRequest := dormain.AuthRequest{Token: token}
	err := jwt.NewAuthorization(authRequest).Validate()

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
