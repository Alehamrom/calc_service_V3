package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"calc_service/pkg/errors"
	"calc_service/pkg/models"
)

const (
	defaultTimeout = 10 * time.Second
	maxRetries     = 3
	retryDelay     = 1 * time.Second
)

type OrchestratorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *OrchestratorClient {
	return &OrchestratorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Получение задачи от оркестратора
func (c *OrchestratorClient) FetchTask() (*models.Task, error) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := c.httpClient.Get(c.baseURL + "/internal/task")
		if err != nil {
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			var task models.Task
			if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
				return nil, fmt.Errorf("ошибка декодирования задачи: %w", err)
			}
			return &task, nil
		case http.StatusNotFound:
			return nil, errors.ErrTaskNotFound
		default:
			return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
		}
	}
	return nil, errors.ErrConnectionFailed
}

// Отправка результата выполнения задачи
func (c *OrchestratorClient) SubmitResult(taskID string, result float64) error {
	payload := struct {
		TaskID string  `json:"task_id"`
		Result float64 `json:"result"`
	}{
		TaskID: taskID,
		Result: result,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка кодирования результата: %w", err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := c.httpClient.Post(
			c.baseURL+"/internal/task",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("ошибка отправки результата: %s", string(body))
		}
		return nil
	}
	return errors.ErrConnectionFailed
}
