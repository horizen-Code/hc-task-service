// internal/repository/task_repository.go
package repository

import (
	"hc-task-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *model.Task) error
	FindByID(id uuid.UUID) (*model.Task, error)
	FindByTeamID(teamID uuid.UUID) ([]model.Task, error)
	FindByCreatorID(creatorID uuid.UUID) ([]model.Task, error)
	Update(task *model.Task) error
	Delete(id uuid.UUID) error
	Assign(taskID, userID uuid.UUID) error
	UpdateStatus(id uuid.UUID, status model.TaskStatus) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByTeamID(teamID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Where("team_id = ?", teamID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) FindByCreatorID(creatorID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Where("creator_id = ?", creatorID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}

func (r *taskRepository) Assign(taskID, userID uuid.UUID) error {
	return r.db.Model(&model.Task{}).Where("id = ?", taskID).
		Update("assignee_id", userID).Error
}

func (r *taskRepository) UpdateStatus(id uuid.UUID, status model.TaskStatus) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).
		Update("status", status).Error
}
