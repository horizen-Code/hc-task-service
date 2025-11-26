// internal/model/task.go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusCreated   TaskStatus = "created"
	StatusInProcess TaskStatus = "in_process"
	StatusDone      TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title       string         `gorm:"not null" json:"title" validate:"required,min=3,max=100"`
	Description string         `json:"description" validate:"max=1000"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'created';index" json:"status"`
	AssigneeID  *uuid.UUID     `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	CreatorID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"creator_id"`
	TeamID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"team_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
