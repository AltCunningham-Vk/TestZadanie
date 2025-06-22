package repository

import (
	"database/sql"
	"task-manager/internal/model"
)

// TaskRepository управляет операциями с базой данных для задач.
// Используется для выполнения CRUD-операций (создание, чтение, обновление, удаление) с таблицей tasks.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository создаёт новый экземпляр TaskRepository.
// Принимает подключение к базе данных и возвращает объект для работы с задачами.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create добавляет новую задачу в базу данных.
// Принимает указатель на структуру Task, вставляет её в таблицу tasks и возвращает ID созданной задачи.
func (r *TaskRepository) Create(task *model.Task) error {
	query := `INSERT INTO tasks (title, description, status, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, task.Title, task.Description, task.Status, task.CreatedAt).
		Scan(&task.ID)
}

// GetAll возвращает список всех задач из базы данных.
// Выполняет SQL-запрос для получения всех записей из таблицы tasks и возвращает их как массив структур Task.// GetAll retrieves all tasks from the database
func (r *TaskRepository) GetAll() ([]model.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID возвращает задачу по её ID.
// Выполняет SQL-запрос для получения одной записи из таблицы tasks по указанному ID.
// Если задача не найдена, возвращает nil без ошибки.
func (r *TaskRepository) GetByID(id int) (*model.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks WHERE id = $1`
	var task model.Task
	err := r.db.QueryRow(query, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

// Update обновляет задачу в базе данных.
// Принимает указатель на структуру Task и обновляет соответствующие поля в таблице tasks по ID.
func (r *TaskRepository) Update(task *model.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4`
	_, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.ID)
	return err
}

// Delete удаляет задачу по её ID.
// Выполняет SQL-запрос для удаления записи из таблицы tasks по указанному ID.
func (r *TaskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetCompletedTasks возвращает список задач со статусом "Выполнена" (status = 2).
// Используется таймером для периодической проверки и удаления завершенных задач.
func (r *TaskRepository) GetCompletedTasks() ([]model.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks WHERE status = 2`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
