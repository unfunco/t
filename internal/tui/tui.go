package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/theme"
)

// Tab represents a tab in the UI.
type Tab int

const (
	TabToday Tab = iota
	TabTomorrow
	TabTodo
	tabCount
)

// String returns the title of a Tab.
func (t Tab) String() string {
	switch t {
	case TabToday:
		return "Today"
	case TabTomorrow:
		return "Tomorrow"
	case TabTodo:
		return "Todo"
	default:
		return "Unknown"
	}
}

// KeyMap defines the key bindings for the UI.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Space    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Quit     key.Binding
	Submit   key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "previous tab"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next tab"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "cancel"),
		),
		Submit: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "submit"),
		),
	}
}

// Model represents the state of the TUI.
type Model struct {
	keys         KeyMap
	activeTab    Tab
	cursor       int
	todayList    *model.TodoList
	tomorrowList *model.TodoList
	todoList     *model.TodoList
	width        int
	height       int
	submitted    bool
	exited       bool
	theme        theme.Theme
}

// New creates a new TUI model with the provided todo lists.
func New(todayList, tomorrowList, todoList *model.TodoList) Model {
	return Model{
		keys:         DefaultKeyMap(),
		activeTab:    TabToday,
		cursor:       0,
		theme:        theme.Default(),
		todayList:    todayList,
		tomorrowList: tomorrowList,
		todoList:     todoList,
	}
}

// Init initializes the model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.exited = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Submit):
			m.submitted = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Left):
			m.previousTab()
			m.cursor = 0
		case key.Matches(msg, m.keys.Right):
			m.nextTab()
			m.cursor = 0
		case key.Matches(msg, m.keys.Up):
			m.cursorUp()
		case key.Matches(msg, m.keys.Down):
			m.cursorDown()
		case key.Matches(msg, m.keys.Enter), key.Matches(msg, m.keys.Space):
			m.toggleCurrent()
		case key.Matches(msg, m.keys.Tab):
			m.nextTab()
			m.cursor = 0
		case key.Matches(msg, m.keys.ShiftTab):
			m.previousTab()
			m.cursor = 0
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the UI.
func (m *Model) View() string {
	if m.submitted {
		return "✓ Changes saved!\n"
	}

	if m.exited {
		return "Exited without saving changes\n"
	}

	var b strings.Builder
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")
	b.WriteString(m.renderList())
	b.WriteString("\n\n")
	b.WriteString(m.renderHelp())

	return m.theme.ContainerStyle().Render(b.String())
}

// renderTabs renders the tab navigation.
func (m *Model) renderTabs() string {
	var tabs []string

	for i := TabToday; i < tabCount; i++ {
		var style lipgloss.Style
		if i == m.activeTab {
			style = m.theme.ActiveTabStyle()
		} else {
			style = m.theme.TabStyle()
		}
		tabs = append(tabs, style.Render(i.String()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderList renders the current todo list.
func (m *Model) renderList() string {
	list := m.getCurrentList()
	if list == nil {
		return "No todos yet."
	}

	var items []string
	for i, todo := range list.Todos {
		var checkbox string
		if todo.Completed {
			greenCheck := m.theme.SuccessStyle().Render("✓")
			checkbox = "[" + greenCheck + "]"
		} else {
			checkbox = "[ ]"
		}

		cursor := "  "
		if i == m.cursor {
			cursor = m.theme.CursorChar + " "
		}

		var titleStyle lipgloss.Style
		if i == m.cursor {
			if todo.Completed {
				titleStyle = m.theme.HighlightedItemStyle().Foreground(m.theme.MutedText).Strikethrough(true)
			} else {
				titleStyle = m.theme.HighlightedItemStyle()
			}
		} else {
			if todo.Completed {
				titleStyle = m.theme.CompletedTitleStyle()
			} else {
				titleStyle = m.theme.ItemStyle()
			}
		}

		var descStyle lipgloss.Style
		if i == m.cursor {
			descStyle = m.theme.HighlightedItemStyle().Foreground(m.theme.MutedText)
		} else {
			descStyle = m.theme.DescriptionStyle()
		}

		item := fmt.Sprintf("%s%s %s",
			cursor,
			checkbox,
			titleStyle.Render(todo.Title),
		)

		if todo.Description != "" {
			item += "\n      " + descStyle.Render(todo.Description)
		}

		items = append(items, item)
	}

	return strings.Join(items, "\n\n")
}

// renderHelp renders the help text.
func (m *Model) renderHelp() string {
	return m.theme.HelpStyle().Render("Enter to select · Tab/Arrow keys to navigate · Ctrl+S to submit · Esc to cancel")
}

// getCurrentList returns the currently active todo list.
func (m *Model) getCurrentList() *model.TodoList {
	switch m.activeTab {
	case TabToday:
		return m.todayList
	case TabTomorrow:
		return m.tomorrowList
	case TabTodo:
		return m.todoList
	default:
		return nil
	}
}

// cursorUp moves the cursor up.
func (m *Model) cursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// cursorDown moves the cursor down.
func (m *Model) cursorDown() {
	list := m.getCurrentList()
	if list != nil && m.cursor < len(list.Todos)-1 {
		m.cursor++
	}
}

// nextTab moves to the next tab.
func (m *Model) nextTab() {
	m.activeTab = (m.activeTab + 1) % tabCount
}

// previousTab moves to the previous tab.
func (m *Model) previousTab() {
	m.activeTab = (m.activeTab + tabCount - 1) % tabCount
}

// toggleCurrent toggles the completion status of the current todo.
func (m *Model) toggleCurrent() {
	list := m.getCurrentList()
	if list != nil && m.cursor < len(list.Todos) {
		list.Todos[m.cursor].Toggle()
	}
}

// GetTodayList returns the today todo list.
func (m *Model) GetTodayList() *model.TodoList {
	return m.todayList
}

// GetTomorrowList returns the tomorrow todo list.
func (m *Model) GetTomorrowList() *model.TodoList {
	return m.tomorrowList
}

// GetTodosList returns the general todo list.
func (m *Model) GetTodosList() *model.TodoList {
	return m.todoList
}

// WasSubmitted returns true if the user submitted the form.
func (m *Model) WasSubmitted() bool {
	return m.submitted
}
