package main

import (
	"fmt"
	"github.com/joho/godotenv"
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

	logger := log.New(os.Stdout, os.Getenv("APP_NAME"), 5)

	appConfig, err := NewAppConfigs(os.Getenv("BASE_URL"), logger)
	if err != nil {
		appConfig.logger.Println(err)
		return
	}

	router, err := appConfig.StartApp()
	if err != nil {
		appConfig.logger.Println(err)
		return
	}

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router,
	}

	appConfig.logger.Printf("starting server on %s", port)

	err = server.ListenAndServe()
	if err != nil {
		appConfig.logger.Println("error starting server: %s", err)
	}
}
