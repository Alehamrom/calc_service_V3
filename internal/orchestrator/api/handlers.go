package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"calc_service/internal/orchestrator/storage"
	"calc_service/pkg/errors"
	"calc_service/pkg/models"

	"github.com/google/uuid"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{storage: store}
}

// Обработчик добавления выражения
func (h *Handler) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if request.Expression == "" {
		http.Error(w, "Пустое выражение", http.StatusUnprocessableEntity)
		return
	}

	// Генерация ID выражения
	exprID := uuid.New().String()

	// Создаем новое выражение
	newExpr := &models.Expression{
		ID:        exprID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Сохраняем в хранилище
	if err := h.storage.AddExpression(newExpr); err != nil {
		http.Error(w, "Ошибка сохранения выражения", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})

	// Запускаем обработку выражения в фоне
	go h.processExpression(newExpr, request.Expression)
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTaskHandler(w, r)
	case http.MethodPost:
		h.SubmitTaskResultHandler(w, r)
	default:
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

// Обработчик получения списка выражений
func (h *Handler) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions, err := h.storage.GetAllExpressions()
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressions,
	})
}

// Обработчик получения выражения по ID
func (h *Handler) GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Некорректный URL", http.StatusBadRequest)
		return
	}

	exprID := pathParts[3]
	expr, exists := h.storage.GetExpression(exprID)
	if !exists {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expr)
}

// Обработчик получения задачи для агента
func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := h.storage.GetNextTask()
	if err != nil {
		if err == errors.ErrTaskNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Обновляем статус задачи
	task.Status = "processing"
	h.storage.UpdateTask(task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Обработчик отправки результата задачи
func (h *Handler) SubmitTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	var result struct {
		TaskID string  `json:"task_id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	// Обновляем задачу и выражение
	if err := h.storage.CompleteTask(result.TaskID, result.Result); err != nil {
		if err == errors.ErrTaskNotFound {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Внутренняя логика обработки выражения
func (h *Handler) processExpression(expr *models.Expression, rawExpr string) {
	// Заглушка для парсинга выражения
	// Реальная реализация будет разбивать выражение на задачи
	// Пример временной реализации:

	// Имитация долгой обработки
	time.Sleep(2 * time.Second)

	// Обновляем статус
	expr.Status = "done"
	expr.Result = 42 // Пример результата

	h.storage.UpdateExpression(expr)
}
