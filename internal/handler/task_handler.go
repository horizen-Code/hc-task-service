// internal/handler/task_handler.go
package handler

import (
	"hc-task-service/internal/dto"
	"hc-task-service/internal/model"
	"hc-task-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	creatorID, _ := uuid.Parse(userIDStr)

	task, err := h.service.Create(req, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	task, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Assign(c *gin.Context) {
	taskID, _ := uuid.Parse(c.Param("id"))
	var req dto.UpdateAssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil || validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
		return
	}

	userIDStr := c.GetString("user_id")
	requesterID, _ := uuid.Parse(userIDStr)

	if err := h.service.Assign(taskID, req.AssigneeID, requesterID); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// ChangeAssignee — смена исполнителя (только создатель задачи)
func (h *TaskHandler) ChangeAssignee(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateAssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	requesterID, _ := uuid.Parse(userIDStr)

	if err := h.service.ChangeAssignee(taskID, req.AssigneeID, requesterID); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator can change assignee"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change assignee"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assignee changed"})
}

// Delete — удаление задачи (только создатель)
func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userIDStr := c.GetString("user_id")
	requesterID, _ := uuid.Parse(userIDStr)

	if err := h.service.Delete(taskID, requesterID); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator can delete task"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// UpdateStatus — смена статуса (создатель или исполнитель)
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := model.TaskStatus(req.Status)
	if status != model.StatusCreated && status != model.StatusInProcess && status != model.StatusDone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	userIDStr := c.GetString("user_id")
	requesterID, _ := uuid.Parse(userIDStr)

	if err := h.service.UpdateStatus(taskID, status, requesterID); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator or assignee can update status"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated", "status": req.Status})
}
