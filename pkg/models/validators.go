package models

import (
	"calc_service/pkg/errors"
	"strings"
)

// ValidateExpression проверяет корректность структуры выражения
func (e *Expression) Validate() error {
	if e.ID == "" {
		return errors.ErrInvalidExpression
	}

	allowedStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"done":       true,
		"error":      true,
	}

	if !allowedStatuses[e.Status] {
		return errors.ErrInvalidExpression
	}

	return nil
}

// ValidateTask проверяет корректность структуры задачи
func (t *Task) Validate() error {
	if t.ID == "" || t.ExpressionID == "" {
		return errors.ErrTaskNotFound
	}

	allowedOperations := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
	}

	if !allowedOperations[t.Operation] {
		return errors.ErrInvalidOperation
	}

	if t.Operation == "/" && t.Arg2 == 0 {
		return errors.ErrDivisionByZero
	}

	return nil
}

// SanitizeExpression очищает ввод выражения
func SanitizeExpression(expr string) string {
	return strings.TrimSpace(expr)
}
