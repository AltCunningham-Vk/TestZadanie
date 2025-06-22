package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Config содержит конфигурацию приложения, включая подключение к базе данных.
// Используется для передачи объекта базы данных в другие компоненты приложения.
type Config struct {
	DB *sql.DB
}

// NewConfig создаёт и инициализирует конфигурацию приложения.
// Устанавливает соединение с базой данных PostgreSQL, используя строку подключения.
// Если подключение не удалось, программа завершится с ошибкой.
func NewConfig() *Config {
	connStr := "user=postgres password=12345 dbname=task_manager sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	return &Config{DB: db}
}
