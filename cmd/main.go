package main

import (
	"fmt"
	"net/http"
	"os"
	_ "task-manager/docs"
	"task-manager/internal/config"
	"task-manager/internal/handler"
	"task-manager/internal/repository"
	"task-manager/internal/service"
	"task-manager/internal/timer"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API менеджера задач
// @version 1.0
// @description Это REST API для управления задачами с использованием PostgreSQL
// @host localhost:8080
// @BasePath /
func main() {
	// Инициализация логгера для записи событий и ошибок.
	// Используется JSON-формат для структурированных логов.
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Логирование в текстовой файл, если отсутсвует то создается автоматически.
	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	logger.SetOutput(logFile)

	// Инициализация конфигурации приложения, включая подключение к базе данных.
	cfg := config.NewConfig()

	// Инициализация слоёв приложения: репозиторий, сервис и обработчик.
	// Репозиторий отвечает за работу с базой данных, сервис — за бизнес-логику,
	// обработчик — за обработку HTTP-запросов.
	repo := repository.NewTaskRepository(cfg.DB)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc, logger)

	// Инициализация и запуск таймера.
	// Таймер каждые 5 минут проверяет завершенные задачи и удаляет их.
	taskTimer := timer.NewTaskTimer(svc, logger)
	taskTimer.Start()

	// Настройка HTTP-сервера с маршрутами для API.
	// Используется gorilla/mux для маршрутизации запросов.
	r := mux.NewRouter()
	r.HandleFunc("/tasks", h.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", h.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", h.DeleteTask).Methods("DELETE")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Запуск сервера на порту 8080.
	logger.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
