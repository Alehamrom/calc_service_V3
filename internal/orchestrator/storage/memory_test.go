package storage_test

import (
	"testing"
	"time"

	"calc_service/internal/orchestrator/storage"
	"calc_service/pkg/errors"
	"calc_service/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	// Тестовые данные
	expr := &models.Expression{
		ID:        "test1",
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	task := &models.Task{
		ID:            "task1",
		ExpressionID:  "test1",
		Operation:     "+",
		Arg1:          2,
		Arg2:          3,
		OperationTime: 1000,
		Status:        "pending",
	}

	t.Run("Добавление и получение выражения", func(t *testing.T) {
		err := store.AddExpression(expr)
		assert.NoError(t, err)

		retrieved, exists := store.GetExpression("test1")
		assert.True(t, exists)
		assert.Equal(t, expr.ID, retrieved.ID)
	})

	t.Run("Дублирование выражения", func(t *testing.T) {
		err := store.AddExpression(expr)
		assert.ErrorIs(t, err, errors.ErrExpressionExists)
	})

	t.Run("Получение несуществующего выражения", func(t *testing.T) {
		_, exists := store.GetExpression("invalid")
		assert.False(t, exists)
	})

	t.Run("Добавление и выполнение задачи", func(t *testing.T) {
		err := store.AddTask(task)
		assert.NoError(t, err)

		// Получаем задачу
		retrievedTask, err := store.GetNextTask()
		assert.NoError(t, err)
		assert.Equal(t, "task1", retrievedTask.ID)
		assert.Equal(t, "processing", retrievedTask.Status)

		// Отмечаем как выполненную
		err = store.CompleteTask("task1", 5)
		assert.NoError(t, err)

		// Проверяем статус
		completedTask, exists := store.GetTask("task1")
		assert.True(t, exists)
		assert.Equal(t, "done", completedTask.Status)
		assert.Equal(t, 5.0, completedTask.Result)
	})

	t.Run("Попытка завершения несуществующей задачи", func(t *testing.T) {
		err := store.CompleteTask("invalid", 0)
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
	})
}
