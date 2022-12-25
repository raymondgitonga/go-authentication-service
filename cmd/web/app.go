package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/raymondgitonga/go-authentication/internal/adapters/cache"
	"github.com/raymondgitonga/go-authentication/internal/adapters/db"
	"github.com/raymondgitonga/go-authentication/internal/adapters/httpserver"
	"go.uber.org/zap"
)

type AppConfigs struct {
	baseURL   string
	dbURL     string
	dbName    string
	redisAddr string
	redisPass string
	logger    *zap.Logger
}

func NewAppConfigs(baseURL, dbURL, dbName, redisAddr, redisPass string, logger *zap.Logger) (*AppConfigs, error) {
	if len(baseURL) < 1 {
		return nil, fmt.Errorf("error in NewAppConfigs, incorrect baseURl")
	}
	if dbURL == "" {
		return nil, fmt.Errorf("kindly provide dbURL")
	}
	if dbName == "" {
		return nil, fmt.Errorf("kindly provide dbName")
	}

	if redisAddr == "" {
		return nil, fmt.Errorf("kindly provide redisAddr")
	}
	if redisPass == "" {
		return nil, fmt.Errorf("kindly provide redisPass")
	}
	return &AppConfigs{baseURL: baseURL, dbURL: dbURL, dbName: dbName, redisAddr: redisAddr, redisPass: redisPass, logger: logger}, nil
}

func (c *AppConfigs) StartApp() (*mux.Router, error) {
	r := mux.NewRouter()

	dbClient, err := db.NewClient(context.Background(), c.dbURL)
	if err != nil {
		return nil, fmt.Errorf("error setting up database: %w", err)
	}

	err = db.RunMigrations(dbClient, c.dbName)
	if err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}

	redisClient, err := cache.NewClient(context.Background(), c.redisAddr, c.redisPass)
	if err != nil {
		return nil, fmt.Errorf("error setting up redis: %w", err)
	}

	handler := httpserver.NewHandler(dbClient, redisClient, c.logger)

	r.HandleFunc(fmt.Sprintf("%s/health-check", c.baseURL), handler.HealthCheck)
	r.HandleFunc(fmt.Sprintf("%s/authorize", c.baseURL), handler.Authorize)
	r.HandleFunc(fmt.Sprintf("%s/validate", c.baseURL), handler.Validate)
	r.HandleFunc(fmt.Sprintf("%s/register", c.baseURL), handler.Register)
	r.HandleFunc(fmt.Sprintf("%s/rotate", c.baseURL), handler.RotateTokens)
	return r, nil
}
