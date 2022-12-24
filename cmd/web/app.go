package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/raymondgitonga/go-authentication/internal/adapters/db"
	"github.com/raymondgitonga/go-authentication/internal/adapters/httpserver"
	"log"
)

type AppConfigs struct {
	baseURL string
	dbURL   string
	dbName  string
	logger  *log.Logger
}

func NewAppConfigs(dbURL, dbName, baseURL string, logger *log.Logger) (*AppConfigs, error) {
	if len(baseURL) < 1 {
		return nil, fmt.Errorf("error in NewAppConfigs, incorrect baseURl")
	}
	if dbURL == "" {
		return nil, fmt.Errorf("kindly provide dbURL")
	}
	if dbName == "" {
		return nil, fmt.Errorf("kindly provide dbName")
	}
	return &AppConfigs{baseURL: baseURL, dbURL: dbURL, dbName: dbName, logger: logger}, nil
}

func (c *AppConfigs) StartApp() (*mux.Router, error) {
	r := mux.NewRouter()

	dbClient, err := db.NewClient(context.Background(), c.dbURL)
	if err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}

	err = db.RunMigrations(dbClient, c.dbName)
	if err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}

	handler := httpserver.Handler{DB: dbClient}
	r.HandleFunc(fmt.Sprintf("%s/health-check", c.baseURL), handler.HealthCheck)
	r.HandleFunc(fmt.Sprintf("%s/authorize", c.baseURL), handler.Authorize)
	r.HandleFunc(fmt.Sprintf("%s/validate", c.baseURL), handler.Validate)
	r.HandleFunc(fmt.Sprintf("%s/register", c.baseURL), handler.Register)
	return r, nil
}
