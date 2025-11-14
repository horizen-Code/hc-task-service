// cmd/server/main.go
package main

import (
	"hc-task-service/config"
	"hc-task-service/internal/handler"
	"hc-task-service/internal/middleware"
	"hc-task-service/internal/model"
	"hc-task-service/internal/repository"
	"hc-task-service/internal/service"
	"hc-task-service/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Task{})

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	r := gin.Default()
	r.Use(middleware.AuthMiddleware(cfg))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", h.Create)
		v1.GET("/tasks/:id", h.GetByID)
		v1.PUT("/tasks/:id/assign", h.Assign)
		v1.PUT("/tasks/:id/assignee", h.Assign) // или отдельно
		v1.DELETE("/tasks/:id", h.Delete)
		v1.PUT("/tasks/:id/status", h.UpdateStatus)
		v1.PUT("/tasks/:id/assignee", h.ChangeAssignee) // смена исполнителя
	}

	log.Printf("Task Service on :%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
