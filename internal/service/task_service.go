// internal/service/task_service.go
package service

import (
	"errors"
	"hc-task-service/internal/dto"
	"hc-task-service/internal/model"
	"hc-task-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("task not found")
	ErrForbidden = errors.New("forbidden")
)

type TaskService interface {
	Create(req dto.CreateTaskRequest, creatorID uuid.UUID) (*dto.TaskResponse, error)
	GetByID(id uuid.UUID) (*dto.TaskResponse, error)
	Assign(taskID, assigneeID, requesterID uuid.UUID) error
	ChangeAssignee(taskID, newAssigneeID, requesterID uuid.UUID) error
	Delete(taskID, requesterID uuid.UUID) error
	UpdateStatus(taskID uuid.UUID, status model.TaskStatus, requesterID uuid.UUID) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(req dto.CreateTaskRequest, creatorID uuid.UUID) (*dto.TaskResponse, error) {
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      model.StatusCreated,
		CreatorID:   creatorID,
		TeamID:      req.TeamID,
		AssigneeID:  req.AssigneeID,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return &dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		AssigneeID:  task.AssigneeID,
		CreatorID:   task.CreatorID,
		TeamID:      task.TeamID,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *taskService) GetByID(id uuid.UUID) (*dto.TaskResponse, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	response := dto.ToTaskResponse(*task)
	return &response, nil
}

func (s *taskService) Assign(taskID, assigneeID, requesterID uuid.UUID) error {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return ErrNotFound
	}

	// Только создатель или участник команды может назначать
	if task.CreatorID != requesterID {
		return ErrForbidden
	}

	return s.repo.Assign(taskID, assigneeID)
}

func (s *taskService) ChangeAssignee(taskID, newAssigneeID, requesterID uuid.UUID) error {
	return s.Assign(taskID, newAssigneeID, requesterID) // та же логика
}

func (s *taskService) Delete(taskID, requesterID uuid.UUID) error {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return ErrNotFound
	}
	if task.CreatorID != requesterID {
		return ErrForbidden
	}
	return s.repo.Delete(taskID)
}

func (s *taskService) UpdateStatus(taskID uuid.UUID, status model.TaskStatus, requesterID uuid.UUID) error {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return ErrNotFound
	}

	// Только исполнитель или создатель
	if task.CreatorID != requesterID && (task.AssigneeID == nil || *task.AssigneeID != requesterID) {
		return ErrForbidden
	}

	return s.repo.UpdateStatus(taskID, status)
}
