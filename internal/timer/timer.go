package timer

import (
	"fmt"
	"task-manager/internal/service"
	"time"

	"github.com/sirupsen/logrus"
)

// TaskTimer управляет периодической проверкой задач.
// Используется для поиска задач со статусом "Выполнена" каждые 5 минут, их вывода в консоль и удаления.
type TaskTimer struct {
	service *service.TaskService
	logger  *logrus.Logger
}

// NewTaskTimer создаёт новый экземпляр TaskTimer.
// Принимает сервис и логгер, возвращает объект для управления таймером.
func NewTaskTimer(service *service.TaskService, logger *logrus.Logger) *TaskTimer {
	return &TaskTimer{service: service, logger: logger}
}

// Start запускает таймер для проверки задач.
// Создаёт тикер, который каждые 5 минут вызывает функцию checkCompletedTasks в отдельной горутине.
func (t *TaskTimer) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			t.checkCompletedTasks()
		}
	}()
}

// checkCompletedTasks проверяет наличие задач со статусом "Выполнена".
// Получает список завершенных задач, выводит их в консоль и удаляет из базы данных.
func (t *TaskTimer) checkCompletedTasks() {
	tasks, err := t.service.GetCompletedTasks()
	if err != nil {
		t.logger.Error("Failed to retrieve completed tasks:", err)
		return
	}

	if len(tasks) == 0 {
		t.logger.Info("No completed tasks found")
		return
	}

	fmt.Println("Completed tasks found:")
	for _, task := range tasks {
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %d, Created At: %s\n",
			task.ID, task.Title, task.Description, task.Status, task.CreatedAt)
		if err := t.service.DeleteTask(task.ID); err != nil {
			t.logger.Error("Failed to delete task:", err)
		} else {
			t.logger.Info("Deleted completed task:", task.ID)
		}
	}
}
