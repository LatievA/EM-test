package main

import (
	"context"
	"em-test/internal/adapter/db"
	httpHandler "em-test/internal/adapter/http"
	"em-test/internal/config"
	"em-test/internal/logger"
	"em-test/internal/service"
	"net/http"
	"os"
)

func main() {
	log := logger.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("Error reading config", "err", err)
	}

	dbRepo, err := db.NewSubscriptionRepo(context.Background(), cfg.ConnectionString)
	if err != nil {
		log.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer dbRepo.Close()

	subscriptionService := service.NewSubscriptionService(dbRepo)

	handler := httpHandler.NewHandler(subscriptionService, log)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Info("Server started on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Error("Server failed", "err", err)
		os.Exit(1)
	}
}
