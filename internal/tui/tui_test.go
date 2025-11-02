package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	m := New()

	if m.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today, got %v", m.activeTab)
	}

	if m.todayList == nil {
		t.Error("Expected todayList to be initialized")
	}

	if m.tomorrowList == nil {
		t.Error("Expected tomorrowList to be initialized")
	}

	if m.todoList == nil {
		t.Error("Expected todoList to be initialized")
	}

	if len(m.todayList.Todos) == 0 {
		t.Error("Expected todayList to have sample todos")
	}
}

func TestTabNavigation(t *testing.T) {
	m := New()

	// Test next tab
	m.nextTab()
	if m.activeTab != TabTomorrow {
		t.Errorf("Expected active tab to be Tomorrow, got %v", m.activeTab)
	}

	// Test previous tab
	m.previousTab()
	if m.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today, got %v", m.activeTab)
	}

	// Test wrap around
	m.previousTab()
	if m.activeTab != TabTodo {
		t.Errorf("Expected active tab to wrap to Todo, got %v", m.activeTab)
	}
}

func TestCursorMovement(t *testing.T) {
	m := New()

	// Start at 0
	if m.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", m.cursor)
	}

	// Move down
	m.cursorDown()
	if m.cursor != 1 {
		t.Errorf("Expected cursor to be at 1, got %d", m.cursor)
	}

	// Move up
	m.cursorUp()
	if m.cursor != 0 {
		t.Errorf("Expected cursor to be at 0, got %d", m.cursor)
	}

	// Can't go below 0
	m.cursorUp()
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", m.cursor)
	}
}

func TestToggleTodo(t *testing.T) {
	m := New()

	// Get the first todo
	todo := &m.todayList.Todos[0]
	initialState := todo.Completed

	// Toggle it
	m.toggleCurrent()

	if todo.Completed == initialState {
		t.Error("Expected todo completion state to change")
	}

	// Toggle it back
	m.toggleCurrent()

	if todo.Completed != initialState {
		t.Error("Expected todo completion state to return to initial state")
	}
}

func TestKeyboardInput(t *testing.T) {
	m := New()
	ptr := &m

	// Test right arrow
	updated, _ := ptr.Update(tea.KeyMsg{Type: tea.KeyRight})
	ptr = updated.(*Model)
	if ptr.activeTab != TabTomorrow {
		t.Errorf("Expected active tab to be Tomorrow after right key, got %v", ptr.activeTab)
	}

	// Test left arrow
	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyLeft})
	ptr = updated.(*Model)
	if ptr.activeTab != TabToday {
		t.Errorf("Expected active tab to be Today after left key, got %v", ptr.activeTab)
	}

	// Test down arrow
	initialCursor := ptr.cursor
	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyDown})
	ptr = updated.(*Model)
	if ptr.cursor != initialCursor+1 {
		t.Errorf("Expected cursor to move down, got %d", ptr.cursor)
	}

	// Test Enter to toggle
	todo := &ptr.todayList.Todos[ptr.cursor]
	initialState := todo.Completed
	updated, _ = ptr.Update(tea.KeyMsg{Type: tea.KeyEnter})
	ptr = updated.(*Model)
	if todo.Completed == initialState {
		t.Error("Expected todo to be toggled after Enter key")
	}
}

func TestView(t *testing.T) {
	m := New()

	view := m.View()

	if view == "" {
		t.Error("Expected view to return non-empty string")
	}

	// Should contain tab names
	if !contains(view, "Today") {
		t.Error("Expected view to contain 'Today'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
