package storage

import (
	"calc_service/pkg/errors"
	"calc_service/pkg/models"
	"sync"
)

type Storage interface {
	AddExpression(*models.Expression) error
	GetExpression(string) (*models.Expression, bool)
	GetAllExpressions() ([]*models.Expression, error)
	UpdateExpression(*models.Expression) error
	AddTask(*models.Task) error
	GetNextTask() (*models.Task, error)
	CompleteTask(string, float64) error
	GetTask(string) (*models.Task, bool)
	UpdateTask(*models.Task) error
}

// Реализация MemoryStorage
type MemoryStorage struct {
	expressions     map[string]*models.Expression
	tasks           map[string]*models.Task
	pendingTasks    []string
	processingTasks map[string]struct{}
	mu              sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		expressions:     make(map[string]*models.Expression),
		tasks:           make(map[string]*models.Task),
		pendingTasks:    make([]string, 0),
		processingTasks: make(map[string]struct{}),
	}
}

// Реализация методов интерфейса Storage
func (s *MemoryStorage) AddExpression(expr *models.Expression) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.expressions[expr.ID]; exists {
		return errors.ErrExpressionExists
	}

	s.expressions[expr.ID] = expr
	return nil
}

func (s *MemoryStorage) GetExpression(id string) (*models.Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, exists := s.expressions[id]
	return expr, exists
}

func (s *MemoryStorage) GetAllExpressions() ([]*models.Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		result = append(result, expr)
	}
	return result, nil
}

func (s *MemoryStorage) UpdateExpression(expr *models.Expression) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.expressions[expr.ID]; !exists {
		return errors.ErrExpressionNotFound
	}

	s.expressions[expr.ID] = expr
	return nil
}

// Методы для работы с задачами

func (s *MemoryStorage) AddTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return errors.ErrTaskExists
	}

	s.tasks[task.ID] = task
	s.pendingTasks = append(s.pendingTasks, task.ID)
	return nil
}

func (s *MemoryStorage) GetNextTask() (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pendingTasks) == 0 {
		return nil, errors.ErrTaskNotFound
	}

	taskID := s.pendingTasks[0]
	s.pendingTasks = s.pendingTasks[1:]

	task := s.tasks[taskID]
	task.Status = "processing"
	s.processingTasks[taskID] = struct{}{}

	return task, nil
}

func (s *MemoryStorage) UpdateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return errors.ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStorage) CompleteTask(taskID string, result float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return errors.ErrTaskNotFound
	}

	task.Status = "done"
	task.Result = result
	delete(s.processingTasks, taskID)

	// Обновляем статус выражения
	expr, exists := s.expressions[task.ExpressionID]
	if !exists {
		return errors.ErrExpressionNotFound
	}

	// Логика обновления статуса выражения
	// (должна проверять завершение всех задач выражения)
	expr.Status = "processing"

	return nil
}

func (s *MemoryStorage) GetTask(taskID string) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	return task, exists
}

// Вспомогательные методы

func (s *MemoryStorage) GetPendingTasksCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.pendingTasks)
}

func (s *MemoryStorage) GetProcessingTasksCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.processingTasks)
}
