// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/theme"
)

var todoCounter int

func newTestTodo(title, description string) model.Todo {
	todoCounter++
	return model.Todo{
		ID:          fmt.Sprintf("test-%d", todoCounter),
		Title:       title,
		Description: description,
		CreatedAt:   time.Unix(int64(todoCounter), 0),
	}
}

// newTestModel creates a new model for testing with sample data.
func newTestModel() Model {
	todayList := &model.TodoList{
		Name: "Today",
		Todos: []model.Todo{
			newTestTodo("Test todo 1", "Test description 1"),
			newTestTodo("Test todo 2", "Test description 2"),
			newTestTodo("Test todo 3", "Test description 3"),
		},
	}
	tomorrowList := &model.TodoList{
		Name:  "Tomorrow",
		Todos: []model.Todo{newTestTodo("Tomorrow task", "")},
	}
	todoList := &model.TodoList{
		Name:  "Todos",
		Todos: []model.Todo{newTestTodo("General task", "")},
	}
	return New(theme.Default(), todayList, tomorrowList, todoList)
}

func TestNew(t *testing.T) {
	m := newTestModel()

	if m.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today, got %v", m.activeTab)
	}

	if m.todayList == nil {
		t.Error("Expected todayList to be initialised")
	}

	if m.tomorrowList == nil {
		t.Error("Expected tomorrowList to be initialised")
	}

	if m.todoList == nil {
		t.Error("Expected todoList to be initialised")
	}

	if len(m.todayList.Todos) == 0 {
		t.Error("Expected todayList to have sample todos")
	}
}

func TestTabNavigation(t *testing.T) {
	m := newTestModel()

	m.nextTab()
	if m.activeTab != TabTomorrow {
		t.Errorf("Expected active tab to be Tomorrow, got %v", m.activeTab)
	}

	m.previousTab()
	if m.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today, got %v", m.activeTab)
	}

	m.previousTab()
	if m.activeTab != TabTodo {
		t.Errorf("Expected active tab to wrap to Todo, got %v", m.activeTab)
	}
}

func TestCursorMovement(t *testing.T) {
	m := newTestModel()

	if m.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", m.cursor)
	}

	m.cursorDown()
	if m.cursor != 1 {
		t.Errorf("Expected cursor to be at 1, got %d", m.cursor)
	}

	m.cursorUp()
	if m.cursor != 0 {
		t.Errorf("Expected cursor to be at 0, got %d", m.cursor)
	}

	m.cursorUp()
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", m.cursor)
	}
}

func TestToggleTodo(t *testing.T) {
	m := newTestModel()

	todo := &m.todayList.Todos[0]
	initialCompletionState := todo.Completed

	m.toggleCurrent()

	if todo.Completed == initialCompletionState {
		t.Error("Expected todo completion state to change")
	}

	m.toggleCurrent()

	if todo.Completed != initialCompletionState {
		t.Error("Expected todo completion state to return to initial state")
	}
}

func TestKeyboardInput(t *testing.T) {
	m := newTestModel()
	ptr := &m

	updated, _ := ptr.Update(tea.KeyMsg{Type: tea.KeyRight})
	ptr = updated.(*Model)
	if ptr.activeTab != TabTomorrow {
		t.Errorf("Expected active tab to be Tomorrow after right key, got %v", ptr.activeTab)
	}

	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyLeft})
	ptr = updated.(*Model)
	if ptr.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today after left key, got %v", ptr.activeTab)
	}

	initialCursor := ptr.cursor
	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyDown})
	ptr = updated.(*Model)
	if ptr.cursor != initialCursor+1 {
		t.Errorf("Expected cursor to move down, got %d", ptr.cursor)
	}

	todo := &ptr.todayList.Todos[ptr.cursor]
	initialState := todo.Completed
	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = updated.(*Model)
	if todo.Completed == initialState {
		t.Error("Expected todo to be toggled after Enter key")
	}
}

func TestEditKeyOpensForm(t *testing.T) {
	m := newTestModel()
	ptr := &m

	if len(ptr.todayList.Todos) == 0 {
		t.Fatal("Expected today list to contain todos")
	}

	updated, _ := ptr.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	ptr = updated.(*Model)
	if ptr.formMode != FormModeEdit {
		t.Fatalf("Expected form mode to be edit, got %v", ptr.formMode)
	}

	ptr.closeForm()

	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyRight})
	ptr = updated.(*Model)
	if ptr.activeTab != TabTomorrow {
		t.Fatalf("Expected to switch to tomorrow tab, got %v", ptr.activeTab)
	}

	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	ptr = updated.(*Model)
	if ptr.formMode != FormModeEdit {
		t.Fatalf("Expected form mode to be edit on tomorrow tab, got %v", ptr.formMode)
	}
}

func TestSubmitFormAddsTodoToSelectedList(t *testing.T) {
	m := newTestModel()
	initialCount := len(m.todoList.Todos)

	m.openForm()
	m.titleInput.SetValue("Write docs")
	m.descriptionInput.SetValue("Document the add/edit flow")
	m.formTargetList = TabTodo

	m.submitForm()

	if len(m.todoList.Todos) != initialCount+1 {
		t.Fatalf("expected todo list count to increase, got %d", len(m.todoList.Todos))
	}

	added := m.todoList.Todos[len(m.todoList.Todos)-1]
	if added.Title != "Write docs" {
		t.Fatalf("expected title to be updated, got %q", added.Title)
	}
	if added.Description != "Document the add/edit flow" {
		t.Fatalf("expected description to be updated, got %q", added.Description)
	}
	if m.formMode != FormModeNone {
		t.Fatalf("expected form to close after submit, got mode %v", m.formMode)
	}
}

func TestEditFormMoveKeepsCorrectTodo(t *testing.T) {
	m := newTestModel()
	m.cursor = 1

	target := m.todayList.Todos[m.cursor]

	m.openEditForm()
	m.titleInput.SetValue("Updated Title")
	m.formTargetList = TabTomorrow
	m.submitForm()

	for _, todo := range m.todayList.Todos {
		if todo.ID == target.ID {
			t.Fatalf("expected todo with id %s to be removed from today list", target.ID)
		}
	}

	found := false
	for _, todo := range m.tomorrowList.Todos {
		if todo.ID == target.ID {
			found = true
			if todo.Title != "Updated Title" {
				t.Fatalf("expected moved todo title to be updated, got %q", todo.Title)
			}
		}
	}

	if !found {
		t.Fatalf("expected todo with id %s to exist in tomorrow list after move", target.ID)
	}
}

func TestView(t *testing.T) {
	m := newTestModel()
	view := m.View()

	if view == "" {
		t.Error("Expected view to return non-empty string")
	}

	if !contains(view, "Today") {
		t.Error("Expected view to contain 'Today'")
	}
}

func TestRenderHelpShowsEditWhenTodosPresent(t *testing.T) {
	m := newTestModel()
	help := stripANSI(m.renderHelp())

	if !contains(help, "E to edit") {
		t.Fatalf("Expected help to mention edit when todos exist, got %q", help)
	}

	m.todayList.Todos = nil
	help = stripANSI(m.renderHelp())
	if contains(help, "E to edit") {
		t.Fatalf("Expected help to omit edit when no todos exist, got %q", help)
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b { // ESC
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') {
				inEscape = false
			}
			continue
		}
		b.WriteByte(s[i])
	}

	return b.String()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
