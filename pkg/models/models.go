package models

import (
	"time"
)

// Expression представляет арифметическое выражение для вычисления
type Expression struct {
	ID        string    `json:"id"`         // Уникальный идентификатор
	Status    string    `json:"status"`     // Статус: pending/processing/done/error
	Result    float64   `json:"result"`     // Результат вычисления
	CreatedAt time.Time `json:"created_at"` // Время создания
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления
}

// Task представляет отдельную вычислительную операцию
type Task struct {
	ID            string  `json:"id"`             // Уникальный идентификатор
	ExpressionID  string  `json:"expression_id"`  // Связь с выражением
	Arg1          float64 `json:"arg1"`           // Первый операнд
	Arg2          float64 `json:"arg2"`           // Второй операнд
	Operation     string  `json:"operation"`      // Операция: +, -, *, /
	OperationTime int     `json:"operation_time"` // Время выполнения в мс
	Status        string  `json:"status"`         // Статус: pending/processing/done
	Result        float64 `json:"result"`         // Результат вычисления
}
