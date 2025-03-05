package api_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"calc_service/internal/orchestrator/api"
	"calc_service/internal/orchestrator/storage"
	"calc_service/pkg/models"
)

func TestCalculateHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	t.Run("Успешное создание выражения", func(t *testing.T) {
		body := bytes.NewBufferString(`{"expression": "2+2"}`)
		req := httptest.NewRequest("POST", "/api/v1/calculate", body)
		w := httptest.NewRecorder()

		handler.CalculateHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Ожидался статус 201, получен %d", w.Code)
		}
	})

	t.Run("Пустое выражение", func(t *testing.T) {
		body := bytes.NewBufferString(`{"expression": ""}`)
		req := httptest.NewRequest("POST", "/api/v1/calculate", body)
		w := httptest.NewRecorder()

		handler.CalculateHandler(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Ожидался статус 422, получен %d", w.Code)
		}
	})
}

func TestGetExpressionHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	// Добавляем тестовое выражение
	expr := &models.Expression{
		ID:        "test123",
		Status:    "done",
		Result:    42,
		CreatedAt: time.Now(),
	}
	store.AddExpression(expr)

	t.Run("Получение существующего выражения", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/expressions/test123", nil)
		w := httptest.NewRecorder()

		handler.GetExpressionHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}
	})

	t.Run("Получение несуществующего выражения", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/expressions/invalid", nil)
		w := httptest.NewRecorder()

		handler.GetExpressionHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен %d", w.Code)
		}
	})
}
