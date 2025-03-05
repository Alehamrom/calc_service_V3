package agent

import (
	"log"
	"time"

	"calc_service/pkg/errors"
	"calc_service/pkg/models"
)

type Worker struct {
	client *OrchestratorClient
}

func NewWorker(client *OrchestratorClient) *Worker {
	return &Worker{client: client}
}

func (w *Worker) Start() {
	for {
		task, err := w.client.FetchTask()
		if err != nil {
			if err == errors.ErrTaskNotFound {
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("Ошибка получения задачи: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		result, err := w.executeTask(task)
		if err != nil {
			log.Printf("Ошибка выполнения задачи %s: %v", task.ID, err)
			continue
		}

		if err := w.client.SubmitResult(task.ID, result); err != nil {
			log.Printf("Ошибка отправки результата для задачи %s: %v", task.ID, err)
		}
	}
}

func (w *Worker) executeTask(task *models.Task) (float64, error) {
	// Имитация долгого выполнения операции
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0, errors.ErrDivisionByZero
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, errors.ErrInvalidOperation
	}
}
