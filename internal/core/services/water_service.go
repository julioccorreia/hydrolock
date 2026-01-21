package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/julioccorreia/hydrolock/internal/core/domain"
	"github.com/julioccorreia/hydrolock/internal/core/ports"
	"go.uber.org/zap"
)

type WaterIntakeService struct {
	repo   ports.WaterIntakeRepository
	ai     ports.AIService
	logger *zap.Logger
}

func NewWaterIntakeService(
	repo ports.WaterIntakeRepository,
	ai ports.AIService,
	logger *zap.Logger,
) *WaterIntakeService {
	return &WaterIntakeService{
		repo:   repo,
		ai:     ai,
		logger: logger,
	}
}

func (s *WaterIntakeService) RegisterIntake(ctx context.Context, file multipart.File, userID string) (*domain.WaterIntake, error) {
	s.logger.Info("Starting water intake registration", zap.String("user_id", userID))

	isWater, explanation, err := s.ai.AnalyzeImage(ctx, file)
	if err != nil {
		s.logger.Error("Failed to analyze image with AI", zap.Error(err))
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	confidence := "LOW"
	if isWater {
		confidence = "HIGH"
	}

	intake := &domain.WaterIntake{
		UserID:        userID,
		IsWater:       isWater,
		Confidence:    confidence,
		AIExplanation: explanation,
		CreatedAt:     time.Now(),
		ImageURL:      "temp/local/storage",
	}

	if err := s.repo.Save(ctx, intake); err != nil {
		s.logger.Error("Failed to persist intake to repository", zap.Error(err))
		return nil, errors.New("failed to save data to database")
	}

	s.logger.Info(
		"Water intake registered successfully",
		zap.Bool("is_water", isWater),
		zap.Uint("id", intake.ID),
	)

	return intake, nil
}
