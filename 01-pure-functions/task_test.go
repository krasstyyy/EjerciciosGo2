package task

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase title",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "title with spaces",
			input:    "  hello  ",
			expected: "Hello",
		},
		{
			name:     "already capitalized",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTitle(tt.input)

			if result != tt.expected {
				t.Errorf("FormatTitle(%q) = %q; want %q",
					tt.input,
					result,
					tt.expected)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid normal title",
			input:    "Hello",
			expected: true,
		},
		{
			name:     "title with spaces before",
			input:    "  Hello",
			expected: true,
		},
		{
			name:     "Lots of spaces and letter",
			input:    "   h    ",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValid(tt.input)

			if result != tt.expected {
				t.Errorf("isValid(%q) = %v; want %v",
					tt.input,
					result,
					tt.expected)
			}
		})
	}
}

func TestFilterByStatus(t *testing.T) {
	task1 := Task{Title: "Tarea 1", Status: StatusDone}
	task2 := Task{Title: "Tarea 2", Status: StatusTodo}
	task3 := Task{Title: "Tarea 3", Status: StatusTodo}

	tests := []struct {
		name     string
		tasks    []Task
		status   Status
		expected []Task
	}{
		{
			name:     "filtra correctamente tareas TODO",
			tasks:    []Task{task1, task2, task3},
			status:   StatusTodo,
			expected: []Task{task2, task3},
		},
		{
			name:     "filtra correctamente tareas DONE",
			tasks:    []Task{task1, task2, task3},
			status:   StatusDone,
			expected: []Task{task1},
		},
		{
			name:     "ninguna tarea coincide con el estado",
			tasks:    []Task{task2, task3},
			status:   StatusDone,
			expected: nil,
		},
		{
			name:     "slice de entrada vacío",
			tasks:    []Task{},
			status:   StatusTodo,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByStatus(tt.tasks, tt.status)

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("FilterByStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCountByStatus(t *testing.T) {
	task1 := Task{Title: "Tarea 1", Status: StatusDone}
	task2 := Task{Title: "Tarea 2", Status: StatusTodo}
	task3 := Task{Title: "Tarea 3", Status: StatusTodo}

	tests := []struct {
		name     string
		tasks    []Task
		expected map[Status]int
	}{
		{
			name:  "cuenta múltiples estados mezclados",
			tasks: []Task{task1, task2, task3},

			expected: map[Status]int{
				StatusTodo: 2,
				StatusDone: 1,
			},
		},
		{
			name:  "cuenta solo cuando hay un tipo de estado",
			tasks: []Task{task2, task3},
			expected: map[Status]int{
				StatusTodo: 2,
			},
		},
		{
			name:     "slice de entrada vacío",
			tasks:    []Task{},
			expected: map[Status]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountByStatus(tt.tasks)

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("CountByStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
