package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/raymondgitonga/go-authentication/internal/adapters/httpserver"
	"log"
)

type AppConfigs struct {
	baseURL string
	logger  *log.Logger
}

func NewAppConfigs(baseURL string, logger *log.Logger) (*AppConfigs, error) {
	if len(baseURL) < 1 {
		return nil, fmt.Errorf("error in NewAppConfigs, incorrect baseURl")
	}
	return &AppConfigs{baseURL: baseURL, logger: logger}, nil
}

func (c *AppConfigs) StartApp() (*mux.Router, error) {
	r := mux.NewRouter()
	handler := httpserver.Handler{}

	fmt.Println(fmt.Sprintf("%s/health-check", c.baseURL))
	r.HandleFunc(fmt.Sprintf("%s/health-check", c.baseURL), handler.HealthCheck)
	return r, nil
}
