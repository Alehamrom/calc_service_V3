package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"calc_service/pkg/errors"
	"calc_service/pkg/models"
)

// Parse разбивает выражение на задачи
func Parse(expr string) ([]*models.Task, error) {
	expr = strings.ReplaceAll(expr, " ", "") // Удаляем пробелы

	// Проверка на пустое выражение
	if len(expr) == 0 {
		return nil, errors.ErrInvalidExpression
	}

	// Проверка баланса скобок
	if !checkParentheses(expr) {
		return nil, errors.ErrInvalidParentheses
	}

	// Преобразуем выражение в обратную польскую запись (RPN)
	rpn, err := toRPN(expr)
	if err != nil {
		return nil, err
	}

	// Преобразуем RPN в задачи
	return rpnToTasks(rpn)
}

// toRPN преобразует выражение в обратную польскую запись
func toRPN(expr string) ([]string, error) {
	var output []string
	var operators []string

	for i := 0; i < len(expr); i++ {
		char := rune(expr[i])

		if unicode.IsDigit(char) || char == '.' {
			// Считываем число целиком
			numStr := readNumber(expr, &i)
			output = append(output, numStr)
			continue
		}

		switch char {
		case '(':
			operators = append(operators, string(char))
		case ')':
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, errors.ErrInvalidParentheses
			}
			operators = operators[:len(operators)-1] // Убираем "("
		case '+', '-', '*', '/':
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(string(char)) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, string(char))
		default:
			return nil, fmt.Errorf("неподдерживаемый символ: %c", char)
		}
	}

	// Добавляем оставшиеся операторы
	for len(operators) > 0 {
		op := operators[len(operators)-1]
		if op == "(" {
			return nil, errors.ErrInvalidParentheses
		}
		output = append(output, op)
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// rpnToTasks преобразует RPN в задачи
func rpnToTasks(rpn []string) ([]*models.Task, error) {
	var stack []string
	var tasks []*models.Task

	for _, token := range rpn {
		if isNumber(token) {
			stack = append(stack, token)
			continue
		}

		// Для оператора нужны два операнда
		if len(stack) < 2 {
			return nil, errors.ErrInvalidExpression
		}

		// Создаем задачу
		arg2 := stack[len(stack)-1]
		arg1 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		task := &models.Task{
			ID:            generateTaskID(),
			Operation:     token,
			Arg1:          parseNumber(arg1),
			Arg2:          parseNumber(arg2),
			OperationTime: getOperationTime(token),
			Status:        "pending",
		}
		tasks = append(tasks, task)

		// Результат задачи становится новым операндом
		stack = append(stack, fmt.Sprintf("task_%s", task.ID))
	}

	if len(stack) != 1 {
		return nil, errors.ErrInvalidExpression
	}

	return tasks, nil
}

// Вспомогательные функции

func checkParentheses(expr string) bool {
	var balance int
	for _, char := range expr {
		switch char {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return false
			}
		}
	}
	return balance == 0
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func readNumber(expr string, i *int) string {
	var numStr strings.Builder
	for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
		numStr.WriteByte(expr[*i])
		*i++
	}
	*i-- // Возвращаем индекс на последний символ числа
	return numStr.String()
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func parseNumber(s string) float64 {
	if strings.HasPrefix(s, "task_") {
		return 0 // Временное значение, будет заменено при вычислении
	}
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

func getOperationTime(op string) int {
	switch op {
	case "+":
		return 1000 // Время в миллисекундах
	case "-":
		return 1000
	case "*":
		return 2000
	case "/":
		return 2000
	}
	return 0
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
