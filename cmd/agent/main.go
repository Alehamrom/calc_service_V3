package main

import (
	"os"
	"strconv"

	"calc_service/internal/agent"
)

func main() {
	// Конфигурация
	orchURL := getEnv("ORCHESTRATOR_URL", "http://localhost:8080")
	workersNum := getEnvAsInt("COMPUTING_POWER", 1)

	// Инициализация клиента
	client := agent.NewClient(orchURL)

	// Создание и запуск агента
	a := agent.NewAgent(client, workersNum)
	a.Run()
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil || value < 1 {
		return defaultValue
	}
	return value
}
