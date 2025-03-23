package models

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Description string    `gorm:"not null"`
	isCorrect   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
