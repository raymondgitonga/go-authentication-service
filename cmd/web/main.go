package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading configs: %s", err)
		return
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error setting up logger: %s", err)
		return
	}
	defer func() {
		err = logger.Sync()
	}()

	if err != nil {
		log.Fatalf("error setting up logger: %s", err)
		return
	}

	config, err := NewAppConfigs(
		os.Getenv("DB_CONNECTION_URL"),
		os.Getenv("DB_NAME"),
		os.Getenv("BASE_URL"),
		logger,
	)

	if err != nil {
		config.logger.Error("error initialising configs", zap.String("error", err.Error()))
		return
	}

	router, err := config.StartApp()
	if err != nil {
		config.logger.Error("error starting app", zap.String("error", err.Error()))
		return
	}

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router,
	}

	config.logger.Info("app started", zap.String("port", port))

	err = server.ListenAndServe()
	if err != nil {
		config.logger.Error("error starting server", zap.String("error", err.Error()))
	}
}
