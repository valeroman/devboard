// Package domain task status
package domain

import (
	"encoding/json"
	"fmt"
)

// TaskStatus represents the lifecycle status of a task.
type TaskStatus string

// TaskStatus values represent the lifecycle of a task.
const (
	// TaskStatusTodo indicates the task has not started yet.
	TaskStatusTodo TaskStatus = "todo"
	// TaskStatusInProgress indicates the task is currently being worked on.
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusDone indicates the task has been completed.
	TaskStatusDone TaskStatus = "done"
)

// IsValid reports whether status is a defined task status.
func (status TaskStatus) IsValid() bool {
	switch status {
	case TaskStatusTodo, TaskStatusDone, TaskStatusInProgress:
		return true
	default:
		return false
	}
}

// UnmarshalJSON decodes and validates a task status from JSON.
func (status *TaskStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	statusTask := TaskStatus(raw)
	if !statusTask.IsValid() {
		return fmt.Errorf("estado de tarea inválido: %q", raw)
	}

	*status = statusTask
	return nil
}
