package agent_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"calc_service/internal/agent"
	"calc_service/pkg/models"
)

func TestWorker(t *testing.T) {
	// Мок сервера оркестратора
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/internal/task":
			task := models.Task{
				ID:            "test_task",
				Operation:     "+",
				Arg1:          2,
				Arg2:          3,
				OperationTime: 100,
			}
			json.NewEncoder(w).Encode(task)
		}
	}))
	defer ts.Close()

	client := agent.NewClient(ts.URL)
	worker := agent.NewWorker(client)

	go worker.Start()

	// Даем время на выполнение
	time.Sleep(200 * time.Millisecond)
}
