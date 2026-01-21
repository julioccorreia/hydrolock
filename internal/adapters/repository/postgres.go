package repository

import (
	"fmt"

	"github.com/julioccorreia/hydrolock/config"
	"github.com/julioccorreia/hydrolock/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = db.AutoMigrate(&domain.WaterIntake{})
	if err != nil {
		return nil, fmt.Errorf("failed to run auto-migration: %w", err)
	}

	return db, nil
}
