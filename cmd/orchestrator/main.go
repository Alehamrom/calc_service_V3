package main

import (
	"log"
	"net/http"
	"os"

	"calc_service/internal/orchestrator/api"
	"calc_service/internal/orchestrator/storage"
)

func main() {
	//Инициализация хранилища
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	http.HandleFunc("/api/v1/calculate", handler.CalculateHandler)
	http.HandleFunc("/api/v1/expressions", handler.GetExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", handler.GetExpressionHandler)
	http.HandleFunc("/internal/task", handler.TaskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Оркестратор запущен на порту :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
