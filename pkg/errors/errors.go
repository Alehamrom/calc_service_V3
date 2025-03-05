package errors

import "fmt"

var (
	ErrInvalidExpression   = fmt.Errorf("некорректное выражение")
	ErrDivisionByZero      = fmt.Errorf("деление на ноль")
	ErrTaskNotFound        = fmt.Errorf("задача не найдена")
	ErrExpressionNotFound  = fmt.Errorf("выражение не найдено")
	ErrInvalidOperation    = fmt.Errorf("неподдерживаемая операция")
	ErrExpressionExists    = fmt.Errorf("выражение уже существует")
	ErrTaskExists          = fmt.Errorf("задача уже существует")
	ErrConnectionFailed    = fmt.Errorf("ошибка соединения")
	ErrTimeout             = fmt.Errorf("превышено время выполнения")
	ErrInvalidParentheses  = fmt.Errorf("несбалансированные скобки")
	ErrInvalidJSON         = fmt.Errorf("некорректный JSON")
	ErrEmptyExpression     = fmt.Errorf("пустое выражение")
	ErrInternalServerError = fmt.Errorf("внутренняя ошибка сервера")
)
