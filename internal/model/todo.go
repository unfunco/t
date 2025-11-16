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
}

// TodoList represents a collection of todos with a name.
type TodoList struct {
	Name  string
	Todos []Todo
}

// NewTodo creates a new todo item with the given title and description.
func NewTodo(title, description string) Todo {
	return Todo{
		ID:          time.Now().Format("20060102150405.000000"),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
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
