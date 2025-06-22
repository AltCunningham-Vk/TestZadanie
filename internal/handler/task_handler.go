package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager/internal/model"
	"task-manager/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// TaskHandler обрабатывает HTTP-запросы для работы с задачами.
// Отвечает за взаимодействие с клиентом через REST API, передавая запросы в сервисный слой.
type TaskHandler struct {
	service *service.TaskService
	logger  *logrus.Logger
}

// NewTaskHandler создаёт новый экземпляр TaskHandler.
// Принимает сервис и логгер, возвращает объект для обработки HTTP-запросов.
func NewTaskHandler(service *service.TaskService, logger *logrus.Logger) *TaskHandler {
	return &TaskHandler{service: service, logger: logger}
}

// CreateTask создаёт новую задачу через HTTP-запрос.
// Обрабатывает POST-запрос, принимает JSON с данными задачи, сохраняет её в базе и возвращает созданную задачу.
// CreateTask creates a new task
// @Summary Create a new task
// @Description Create a new task with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Task details"
// @Success 201 {object} model.Task
// @Failure 400 {string} string "Invalid request"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.logger.Error("Invalid request body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	createdTask, err := h.service.CreateTask(task.Title, task.Description, task.Status)
	if err != nil {
		h.logger.Error("Failed to create task:", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Task created:", createdTask.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, err := json.MarshalIndent(createdTask, "", "  ")
	if err != nil {
		h.logger.Error("Failed to marshal task:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// GetAllTasks возвращает список всех задач.
// Обрабатывает GET-запрос, получает все задачи из базы данных и возвращает их в формате JSON.
// GetAllTasks retrieves all tasks
// @Summary Get all tasks
// @Description Retrieve a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} model.Task
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		h.logger.Error("Failed to retrieve tasks:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		h.logger.Error("Failed to marshal tasks:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Retrieved all tasks")
	w.Write(response)
}

// GetTaskByID возвращает задачу по её ID.
// Обрабатывает GET-запрос, получает задачу по указанному ID и возвращает её в формате JSON.
// @Summary Get a task by ID
// @Description Retrieve a task by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} model.Task
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		h.logger.Error("Failed to retrieve task:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if task == nil {
		h.logger.Warn("Task not found:", id)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		h.logger.Error("Failed to marshal task:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Retrieved task:", id)
	w.Write(response)
}

// UpdateTask обновляет задачу.
// Обрабатывает PUT-запрос, принимает JSON с новыми данными задачи и обновляет запись в базе.
// @Summary Update a task
// @Description Update a task with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body model.Task true "Task details"
// @Success 200 {string} string "Task updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.logger.Error("Invalid request body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(id, task.Title, task.Description, task.Status); err != nil {
		h.logger.Error("Failed to update task:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Task updated:", id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`"Task updated"`))
}

// DeleteTask удаляет задачу.
// Обрабатывает DELETE-запрос, удаляет задачу по указанному ID из базы данных.
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "Task deleted"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.DeleteTask(id); err != nil {
		h.logger.Error("Failed to delete task:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Task deleted:", id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`"Task deleted"`))
}
