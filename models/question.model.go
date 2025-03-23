package models

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Description string    `gorm:"not null"`
	Answers     []Answer  `Allow multiple answers`
	TimerOpt    bool      `Specific question timer option`
	Timer       int       `Only used if sp_timer_opt enabled`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
