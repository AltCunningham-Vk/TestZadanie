package service

import (
	"task-manager/internal/model"
	"task-manager/internal/repository"
	"time"
)

// TaskService обрабатывает бизнес-логику для задач.
// Служит промежуточным слоем между HTTP-обработчиками и репозиторием, упрощая управление задачами.
type TaskService struct {
	repo *repository.TaskRepository // Репозиторий для доступа к базе данных
}

// NewTaskService создаёт новый экземпляр TaskService.
// Принимает репозиторий и возвращает объект для работы с задачами.
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask создаёт новую задачу.
// Принимает название, описание и статус, создаёт структуру Task и сохраняет её в базе данных.
func (s *TaskService) CreateTask(title, description string, status int) (*model.Task, error) {
	task := &model.Task{
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now(), // Устанавливает текущую дату и время
	}
	return task, s.repo.Create(task)
}

// GetAllTasks возвращает все задачи из базы данных.
// Вызывает соответствующий метод репозитория для получения списка задач.
func (s *TaskService) GetAllTasks() ([]model.Task, error) {
	return s.repo.GetAll()
}

// GetTaskByID возвращает задачу по её ID.
// Вызывает метод репозитория для получения задачи по указанному идентификатору.
func (s *TaskService) GetTaskByID(id int) (*model.Task, error) {
	return s.repo.GetByID(id)
}

// UpdateTask обновляет задачу.
// Принимает ID и новые данные задачи, обновляет запись в базе данных через репозиторий.
func (s *TaskService) UpdateTask(id int, title, description string, status int) error {
	task := &model.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	}
	return s.repo.Update(task)
}

// DeleteTask удаляет задачу по её ID.
// Вызывает метод репозитория для удаления задачи из базы данных.
func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}

// GetCompletedTasks возвращает список задач со статусом "Выполнена".
// Используется для периодической проверки завершенных задач таймером.
func (s *TaskService) GetCompletedTasks() ([]model.Task, error) {
	return s.repo.GetCompletedTasks()
}
