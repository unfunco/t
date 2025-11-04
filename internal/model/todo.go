package model

import "time"

// Todo represents a single todo item.
type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// TodoList represents a collection of todos with a name.
type TodoList struct {
	Name  string
	Todos []Todo
}

// NewTodo creates a new todo item with the given title and description.
func NewTodo(title, description string) Todo {
	return Todo{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}

// Toggle toggles the completion status of a Todo.
func (t *Todo) Toggle() {
	t.Completed = !t.Completed
	if t.Completed {
		now := time.Now()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}

// generateID generates a unique ID for a todo item.
func generateID() string {
	return time.Now().Format("20060102150405.000000")
}
