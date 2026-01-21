package main

import (
	"context"
	"fmt"
	"log"

	"github.com/julioccorreia/hydrolock/config"
	"github.com/julioccorreia/hydrolock/internal/adapters/ai"
	"github.com/julioccorreia/hydrolock/internal/adapters/http/handlers"
	"github.com/julioccorreia/hydrolock/internal/adapters/http/router"
	"github.com/julioccorreia/hydrolock/internal/adapters/repository"
	"github.com/julioccorreia/hydrolock/internal/core/services"
	"github.com/julioccorreia/hydrolock/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logInstance, err := logger.NewLogger(cfg.GoEnv)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logInstance.Sync()

	logInstance.Info("Starting HydroLock API...", zap.String("port", cfg.HTTPPort))

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logInstance.Fatal("Failed to connect to database", zap.Error(err))
	}

	ctx := context.Background()
	aiService, err := ai.NewGeminiService(ctx, cfg.GeminiAPIKey)
	if err != nil {
		logInstance.Fatal("Failed to initialize AI service", zap.Error(err))
	}

	waterRepo := repository.NewWaterIntakeRepository(db)
	waterService := services.NewWaterIntakeService(waterRepo, aiService, logInstance)
	waterHandler := handlers.NewWaterHandler(waterService)

	r := router.NewRouter(waterHandler)

	serverAddr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logInstance.Info("Servier is ready to handle requests", zap.String("address", serverAddr))

	if err := r.Run(serverAddr); err != nil {
		logInstance.Fatal("Failed to start server", zap.Error(err))
	}
}
