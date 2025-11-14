// internal/dto/task_dto.go
package dto

import (
	"hc-task-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=100"`
	Description string     `json:"description" validate:"max=1000"`
	TeamID      uuid.UUID  `json:"team_id" validate:"required"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
}

type UpdateAssigneeRequest struct {
	AssigneeID uuid.UUID `json:"assignee_id" validate:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=created in_process done"`
}

type TaskResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	CreatorID   uuid.UUID  `json:"creator_id"`
	TeamID      uuid.UUID  `json:"team_id"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

func ToTaskResponse(t model.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID,
		CreatorID:   t.CreatorID,
		TeamID:      t.TeamID,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}
