package task

import (
	"strings"
)

type Status string

const (
	StatusTodo Status = "todo"
	StatusDone Status = "done"
)

type Task struct {
	Title  string
	Status Status
}

func FormatTitle(s string) string {
	var cadenaBuena = s
	cadenaBuena = strings.Trim(cadenaBuena, " ")
	cadenaBuena = strings.Title(cadenaBuena)
	return cadenaBuena
}

func isValid(s string) bool {
	s = strings.Trim(s, " ")
	if s == "" {
		return false
	}
	return true
}

func FilterByStatus(tasks []Task, status Status) []Task {
	var filtered []Task
	for _, task := range tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func CountByStatus(tasks []Task) map[Status]int {
	counts := make(map[Status]int)
	for _, task := range tasks {
		counts[task.Status]++
	}
	return counts
}
