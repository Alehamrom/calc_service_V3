package parser_test

import (
	"testing"

	"calc_service/internal/orchestrator/parser"
	"calc_service/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // Количество задач
		err      error
	}{
		{
			name:     "Простое сложение",
			input:    "2+3",
			expected: 1,
		},
		{
			name:     "Приоритет операций",
			input:    "2+3*4",
			expected: 2,
		},
		{
			name:     "Скобки",
			input:    "(2+3)*4",
			expected: 2,
		},
		{
			name:     "Деление",
			input:    "10/2",
			expected: 1,
		},
		{
			name:     "Неверное выражение",
			input:    "2++3",
			expected: 0,
			err:      errors.ErrInvalidExpression,
		},
		{
			name:     "Несбалансированные скобки",
			input:    "(2+3",
			expected: 0,
			err:      errors.ErrInvalidParentheses,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := parser.Parse(tt.input)

			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, len(tasks))
		})
	}
}
