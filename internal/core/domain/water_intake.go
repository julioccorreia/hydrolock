package domain

import "time"

type WaterIntake struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id" gorm:"index"`
	ImageURL      string    `json:"image_url"`
	IsWater       bool      `json:"is_water"`
	Confidence    string    `json:"confidence"`
	AIExplanation string    `json:"ai_explanation"`
	CreatedAt     time.Time `json:"created_at"`
}

func (WaterIntake) TableName() string {
	return "water_intakes"
}
