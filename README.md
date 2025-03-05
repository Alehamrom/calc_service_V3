# 🧮 Распределенный вычислитель арифметических выражений

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/calc_service)](https://goreportcard.com/report/github.com/yourusername/calc_service)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Система для параллельного вычисления математических выражений с поддержкой распределенных вычислений.

## Особенности

- Параллельная обработка операций (`+`, `-`, `*`, `/`)
- Таймауты выполнения операций
- Отслеживание статуса выражений в реальном времени
- Готовые Docker-образы
- Поддержка масштабирования агентов

## Быстрый старт

### Запуск через Docker Compose

```bash
git clone https://github.com/yourusername/calc_service.git
cd calc_service
docker-compose up --build
```

## Ручная установка

```bash
# Сборка бинарников
make build

# Запуск оркестратора
export PORT=8080 TIME_ADDITION_MS=1000 TIME_SUBTRACTION_MS=1000
./bin/orchestrator

# Запуск агента (в отдельном терминале)
export ORCHESTRATOR_URL=http://localhost:8080 COMPUTING_POWER=4
./bin/agent
```

# 📡 API Endpoints

## ➕ Добавление выражения
```json
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "(2+3)*4"}'
```

### Пример ответа:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "created_at": "2024-03-20T12:00:00Z"
}
```

### 📊 Получение статуса выражения
```bash
curl http://localhost:8080/api/v1/expressions/550e8400-e29b-41d4-a716-446655440000
```

### Пример ответа:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "done",
  "result": 20.0,
  "created_at": "2024-03-20T12:00:00Z"
}
```

## 🏗️ Архитектура системы
A[Пользователь] --> B[Оркестратор]
B --> C[Парсер]
C --> D[Очередь задач]
D --> E[Агент 1]
D --> F[Агент 2]
E -->|Результат| B
F -->|Результат| B
B --> G[(Хранилище)]

# ⚙️ Конфигурация

### Оркестратор
| Переменная               | По умолчанию | Описание                     |
|--------------------------|--------------|------------------------------|
| `PORT`                   | 8080         | Порт HTTP-сервера            |
| `TIME_ADDITION_MS`       | 1000         | Время выполнения сложения    |
| `TIME_SUBTRACTION_MS`    | 1000         | Время выполнения вычитания   |
| `TIME_MULTIPLICATION_MS` | 2000         | Время выполнения умножения   |
| `TIME_DIVISIONS_MS`      | 2000         | Время выполнения деления     |

### Агент
| Переменная             | Обязательно | Описание                          |
|------------------------|-------------|-----------------------------------|
| `ORCHESTRATOR_URL`     | Да          | URL оркестратора (например: `http://localhost:8080`) |
| `COMPUTING_POWER`      | Да          | Количество параллельных воркеров  |



# 🧪 Тестирование

### Все тесты:
```bash
make test
```

### Тесты с проверкой гонок:
```bash
make race-test
```

### Покрытие кода:
```bash
make coverage
```
