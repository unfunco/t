// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package model

import "time"

// Todo represents a single todo item.
type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DueDate     *time.Time `json:"due_date"`
}

// TodoList represents a collection of todos with a name.
type TodoList struct {
	Name  string
	Todos []Todo
}

// NewTodo creates a new todo item with the given title, description, and
// optional due date.
func NewTodo(title, description string, dueDate *time.Time) Todo {
	now := time.Now()
	return Todo{
		ID:          now.Format("20060102150405.000000"),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		DueDate:     cloneTimePtr(dueDate),
	}
}

// ToggleCompleted toggles the completion status of a Todo.
func (t *Todo) ToggleCompleted() {
	t.Completed = !t.Completed
	if t.Completed {
		now := time.Now()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}

// SetDueDate updates the due date for the todo.
func (t *Todo) SetDueDate(dueDate *time.Time) {
	t.DueDate = cloneTimePtr(dueDate)
}

// IsOverdue reports whether the todo is overdue relative to the provided time.
func (t *Todo) IsOverdue(reference time.Time) bool {
	if t.Completed || t.DueDate == nil {
		return false
	}

	due := startOfDay(*t.DueDate)
	ref := startOfDay(reference)

	return due.Before(ref)
}

func cloneTimePtr(in *time.Time) *time.Time {
	if in == nil {
		return nil
	}

	clone := *in
	return &clone
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
