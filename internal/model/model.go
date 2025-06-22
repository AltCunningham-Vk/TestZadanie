package model

import "time"

// Task представляет задачу в базе данных.
// Используется для хранения информации о задаче, включая её идентификатор, название, описание, статус и дату создания.
// Поля структуры соответствуют столбцам таблицы tasks в PostgreSQL.
type Task struct {
	ID          int       `json:"id" example:"1"`                                  // Уникальный идентификатор задачи, автоматически генерируется базой данных
	Title       string    `json:"title" example:"Test Task"`                       // Название задачи, строка до 255 символов
	Description string    `json:"description" example:"Test description"`          // Описание задачи, может быть пустым
	Status      int       `json:"status" example:"0"`                              // Статус задачи: 0 (Новая), 1 (В процессе), 2 (Выполнена)
	CreatedAt   time.Time `json:"created_at" example:"2025-06-22T13:12:30.37482Z"` // Дата и время создания задачи, устанавливается автоматически
}
