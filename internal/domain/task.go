package domain

import (
	"errors"
	"time"
)

// Task representa una tarea en el sistema.
type Task struct {
	ID          string
	ProjectID   string
	Title       string
	Description string
	Status      TaskStatus
	AssigneeID  *string // ID del usuario asignado a la tarea (opcional).
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask crea una nueva instancia de Task con los datos proporcionados.
func NewTask(projectID, title, description, createdBy string) *Task {
	now := time.Now()

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      TaskStatusTodo,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

var (
	// ErrTaskAlreadyDone se devuelve cuando se intenta realizar una acción en una tarea que ya está completada.
	ErrTaskAlreadyDone = errors.New("la tarea ya está completada")
	// ErrTaskNotAssignable se devuelve cuando se intenta asignar una tarea que no puede ser asignada en su estado actual.
	ErrTaskNotAssignable = errors.New("la tarea no puede asignarse en este estado")
)

// AssignTo asigna la tarea a un usuario específico.
func (t *Task) AssignTo(userID string) error {
	if t.Status == TaskStatusDone {
		return ErrTaskNotAssignable
	}
	t.AssigneeID = &userID
	t.UpdatedAt = time.Now()
	return nil
}

// MarkAsDone marca la tarea como completada.
func (t *Task) MarkAsDone() error {
	if t.Status == TaskStatusDone {
		return ErrTaskAlreadyDone
	}
	t.Status = TaskStatusDone
	t.UpdatedAt = time.Now()
	return nil
}

// MoveToStatus cambia el estado de la tarea a un nuevo estado válido.
func (t *Task) MoveToStatus(newStatus TaskStatus) error {
	if !newStatus.IsValid() {
		return errors.New("estado de tarea inválido")
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}
