// cmd/server/main.go
package main

import (
	"context"
	"hc-task-service/config"
	"hc-task-service/internal/handler"
	"hc-task-service/internal/middleware"
	"hc-task-service/internal/model"
	"hc-task-service/internal/repository"
	"hc-task-service/internal/service"
	"hc-task-service/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Миграция с проверкой ошибки
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	r := gin.Default()

	// Публичные роуты (без авторизации)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Защищённые роуты (с авторизацией)
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg))
	{
		v1.POST("/tasks", h.Create)
		v1.GET("/tasks", h.GetAll)
		v1.GET("/tasks/:id", h.GetByID)
		v1.PUT("/tasks/:id/assign", h.Assign)
		v1.DELETE("/tasks/:id", h.Delete)
		v1.PUT("/tasks/:id/status", h.UpdateStatus)
		v1.PUT("/tasks/:id/assignee", h.ChangeAssignee)
	}

	// Graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Task Service starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
