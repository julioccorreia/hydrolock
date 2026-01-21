package repository

import (
	"context"

	"github.com/julioccorreia/hydrolock/internal/core/domain"
	"gorm.io/gorm"
)

type WaterIntakeRepository struct {
	db *gorm.DB
}

func NewWaterIntakeRepository(db *gorm.DB) *WaterIntakeRepository {
	return &WaterIntakeRepository{
		db: db,
	}
}

func (r *WaterIntakeRepository) Save(ctx context.Context, intake *domain.WaterIntake) error {
	result := r.db.WithContext(ctx).Create(intake)
	return result.Error
}

func (r *WaterIntakeRepository) GetByID(ctx context.Context, id uint) (*domain.WaterIntake, error) {
	var intake domain.WaterIntake

	result := r.db.WithContext(ctx).First(&intake, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &intake, nil
}
