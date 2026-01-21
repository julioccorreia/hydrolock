package ports

import (
	"context"
	"mime/multipart"

	"github.com/julioccorreia/hydrolock/internal/core/domain"
)

type WaterIntakeRepository interface {
	Save(ctx context.Context, intake *domain.WaterIntake) error
	GetByID(ctx context.Context, id uint) (*domain.WaterIntake, error)
}

type AIService interface {
	AnalyzeImage(ctx context.Context, file multipart.File) (isWater bool, explanation string, err error)
}

type WaterIntakeService interface {
	RegisterIntake(ctx context.Context, file multipart.File, userID string) (*domain.WaterIntake, error)
}
