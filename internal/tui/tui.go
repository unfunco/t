// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/theme"
)

// Tab represents a tab in the UI, which corresponds to a todo list.
type Tab int

const (
	TabToday Tab = iota
	TabTomorrow
	TabTodo
	TabCount
)

// String implements the fmt.Stringer interface and returns the title of a Tab.
func (t Tab) String() string {
	switch t {
	case TabToday:
		return list.Today().Name
	case TabTomorrow:
		return list.Tomorrow().Name
	case TabTodo:
		return list.Todos().Name
	default:
		return "Unknown"
	}
}

// FormMode represents the current form state.
type FormMode int

const (
	FormModeNone FormMode = iota
	FormModeAdd
	FormModeEdit
)

// FormField represents which field is currently focused in the form.
type FormField int

const (
	FormFieldTitle FormField = iota
	FormFieldDescription
	FormFieldList
	formFieldCount
)

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
	Add      key.Binding
	Edit     key.Binding
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
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add todo"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit todo"),
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

	// Form state
	formMode         FormMode
	formField        FormField
	titleInput       textinput.Model
	descriptionInput textarea.Model
	formTargetList   Tab
	editingIndex     int
}

// New creates a new TUI model with the provided todo lists and theme.
func New(th theme.Theme, todayList, tomorrowList, todoList *model.TodoList) Model {
	ti := textinput.New()
	ti.Placeholder = "Todo title"
	ti.CharLimit = 100
	ti.Width = 50

	ta := textarea.New()
	ta.Placeholder = "Description (optional)"
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(3)

	return Model{
		keys:             DefaultKeyMap(),
		activeTab:        TabToday,
		cursor:           0,
		theme:            th,
		todayList:        todayList,
		tomorrowList:     tomorrowList,
		todoList:         todoList,
		formMode:         FormModeNone,
		formField:        FormFieldTitle,
		titleInput:       ti,
		descriptionInput: ta,
	}
}

// Init initialises the model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.formMode == FormModeAdd || m.formMode == FormModeEdit {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.closeForm()
				return m, nil
			case "ctrl+s":
				m.submitForm()
				return m, nil
			case "tab", "down":
				cmd = m.nextFormField()
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case "shift+tab", "up":
				cmd = m.previousFormField()
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case "left":
				if m.formField == FormFieldList {
					m.previousFormList()
					return m, nil
				}
			case "right":
				if m.formField == FormFieldList {
					m.nextFormList()
					return m, nil
				}
			}
		}

		// Update the focused input.
		switch m.formField {
		case FormFieldTitle:
			m.titleInput, cmd = m.titleInput.Update(msg)
			cmds = append(cmds, cmd)
		case FormFieldDescription:
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
			cmds = append(cmds, cmd)
		case FormFieldList:
		case formFieldCount:
		}

		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.exited = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Submit):
			m.submitted = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Add):
			cmd = m.openForm()
			return m, cmd
		case key.Matches(msg, m.keys.Edit):
			cmd = m.openEditForm()
			return m, cmd
		case key.Matches(msg, m.keys.Left), key.Matches(msg, m.keys.ShiftTab):
			// Only allow tab navigation if there are todos.
			if m.hasAnyTodos() {
				m.previousTab()
				m.cursor = 0
			}
		case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Tab):
			// Only allow tab navigation if there are todos.
			if m.hasAnyTodos() {
				m.nextTab()
				m.cursor = 0
			}
		case key.Matches(msg, m.keys.Up):
			m.cursorUp()
		case key.Matches(msg, m.keys.Down):
			m.cursorDown()
		case key.Matches(msg, m.keys.Enter), key.Matches(msg, m.keys.Space):
			m.toggleCurrent()
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

	// Render form if in add or edit mode
	if m.formMode == FormModeAdd || m.formMode == FormModeEdit {
		return m.renderForm()
	}

	var b strings.Builder

	if m.hasAnyTodos() {
		b.WriteString(m.renderTabs())
		b.WriteString("\n\n")
	}

	b.WriteString(m.renderList())
	b.WriteString("\n\n")
	b.WriteString(m.renderHelp())

	return m.theme.ContainerStyle().Render(b.String())
}

// renderTabs renders the tab navigation.
func (m *Model) renderTabs() string {
	var tabs []string

	for i := TabToday; i < TabCount; i++ {
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
	l := m.getCurrentList()
	if l == nil || len(l.Todos) == 0 {
		// If no todos in any l, show a simpler message
		if !m.hasAnyTodos() {
			return "No todos"
		}
		// If there are todos in other lists, keep the original message
		return "No todos yet."
	}

	now := time.Now()
	var items []string
	for i, todo := range l.Todos {
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
				titleStyle = m.theme.HighlightedItemStyle().Foreground(m.theme.Muted.LipGloss()).Strikethrough(true)
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
			descStyle = m.theme.HighlightedItemStyle().Foreground(m.theme.Muted.LipGloss())
		} else {
			descStyle = m.theme.DescriptionStyle()
		}

		item := fmt.Sprintf("%s%s %s",
			cursor,
			checkbox,
			titleStyle.Render(todo.Title),
		)

		if todo.IsOverdue(now) {
			overdueLabel := m.theme.WorryStyle().Render("! Overdue")
			item += " " + overdueLabel
		}

		if todo.Description != "" {
			item += "\n      " + descStyle.Render(todo.Description)
		}

		items = append(items, item)
	}

	return strings.Join(items, "\n\n")
}

// renderHelp renders the help text.
func (m *Model) renderHelp() string {
	var helpItems []string

	helpItems = append(helpItems, "A to add")

	l := m.getCurrentList()
	if l != nil && len(l.Todos) > 0 {
		helpItems = append(helpItems, "E to edit")
		helpItems = append(helpItems, "Enter to select")
	}

	if m.hasAnyTodos() {
		helpItems = append(helpItems, "Tab/Arrow keys to navigate")
		helpItems = append(helpItems, "Ctrl+S to submit")
	}

	helpItems = append(helpItems, "Esc to cancel")

	return m.theme.HelpStyle().Render(strings.Join(helpItems, " · "))
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

// hasAnyTodos returns true if any list has at least one todo.
func (m *Model) hasAnyTodos() bool {
	return (m.todayList != nil && len(m.todayList.Todos) > 0) ||
		(m.tomorrowList != nil && len(m.tomorrowList.Todos) > 0) ||
		(m.todoList != nil && len(m.todoList.Todos) > 0)
}

// cursorUp moves the cursor up.
func (m *Model) cursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// cursorDown moves the cursor down.
func (m *Model) cursorDown() {
	l := m.getCurrentList()
	if l != nil && m.cursor < len(l.Todos)-1 {
		m.cursor++
	}
}

// nextTab moves to the next tab.
func (m *Model) nextTab() {
	m.activeTab = (m.activeTab + 1) % TabCount
}

// previousTab moves to the previous tab.
func (m *Model) previousTab() {
	m.activeTab = (m.activeTab + TabCount - 1) % TabCount
}

// toggleCurrent toggles the completion status of the current todo.
func (m *Model) toggleCurrent() {
	l := m.getCurrentList()
	if l != nil && m.cursor < len(l.Todos) {
		l.Todos[m.cursor].ToggleCompleted()
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

// ListByID returns the list matching the provided list ID.
func (m *Model) ListByID(id list.ID) *model.TodoList {
	switch id {
	case list.TodayID:
		return m.todayList
	case list.TomorrowID:
		return m.tomorrowList
	case list.TodosID:
		return m.todoList
	default:
		return nil
	}
}

// WasSubmitted returns true if the user submitted the form.
func (m *Model) WasSubmitted() bool {
	return m.submitted
}

// openForm opens the add todo form.
func (m *Model) openForm() tea.Cmd {
	m.formMode = FormModeAdd
	m.formField = FormFieldTitle
	m.formTargetList = m.activeTab

	// Reset and focus title input
	m.titleInput.SetValue("")
	m.descriptionInput.SetValue("")
	m.descriptionInput.Blur()

	// Return the focus command for the title input
	return m.titleInput.Focus()
}

// openEditForm opens the edit form for the currently selected todo.
func (m *Model) openEditForm() tea.Cmd {
	l := m.getCurrentList()
	if l == nil {
		return nil
	}

	if len(l.Todos) == 0 {
		return nil
	}

	if m.cursor < 0 || m.cursor >= len(l.Todos) {
		return nil
	}

	todo := l.Todos[m.cursor]
	m.formMode = FormModeEdit
	m.formField = FormFieldTitle
	m.formTargetList = m.activeTab
	m.editingIndex = m.cursor

	m.titleInput.SetValue(todo.Title)
	m.descriptionInput.SetValue(todo.Description)
	m.descriptionInput.Blur()

	return m.titleInput.Focus()
}

// closeForm closes the form without saving.
func (m *Model) closeForm() {
	m.formMode = FormModeNone
	m.titleInput.Blur()
	m.descriptionInput.Blur()
}

// submitForm saves the new or edited todo and closes the form.
func (m *Model) submitForm() {
	title := strings.TrimSpace(m.titleInput.Value())
	if title == "" {
		// Don't save empty todos
		m.closeForm()
		return
	}

	description := strings.TrimSpace(m.descriptionInput.Value())

	if m.formMode == FormModeEdit {
		currentList := m.getCurrentList()
		if currentList != nil && m.editingIndex < len(currentList.Todos) {
			todo := currentList.Todos[m.editingIndex]
			todo.Title = title
			todo.Description = description

			if m.formTargetList != m.activeTab {
				todo.SetDueDate(m.dueDateForTab(m.formTargetList))
				currentList.Todos = append(currentList.Todos[:m.editingIndex], currentList.Todos[m.editingIndex+1:]...)

				targetList := m.getListByTab(m.formTargetList)
				if targetList != nil {
					targetList.Todos = append(targetList.Todos, todo)
				}

				if m.cursor >= len(currentList.Todos) && m.cursor > 0 {
					m.cursor--
				}
			} else {
				currentList.Todos[m.editingIndex] = todo
			}
		}
	} else {
		newTodo := model.NewTodo(title, description, m.dueDateForTab(m.formTargetList))
		targetList := m.getListByTab(m.formTargetList)
		if targetList != nil {
			targetList.Todos = append(targetList.Todos, newTodo)
		}
	}

	m.closeForm()
}

// nextFormField moves to the next form field.
func (m *Model) nextFormField() tea.Cmd {
	m.formField = (m.formField + 1) % formFieldCount
	return m.updateFormFocus()
}

// previousFormField moves to the previous form field.
func (m *Model) previousFormField() tea.Cmd {
	m.formField = (m.formField + formFieldCount - 1) % formFieldCount
	return m.updateFormFocus()
}

// updateFormFocus updates which input is focused based on current field.
func (m *Model) updateFormFocus() tea.Cmd {
	switch m.formField {
	case FormFieldTitle:
		m.descriptionInput.Blur()
		return m.titleInput.Focus()
	case FormFieldDescription:
		m.titleInput.Blur()
		return m.descriptionInput.Focus()
	case FormFieldList:
		m.titleInput.Blur()
		m.descriptionInput.Blur()
		return nil
	default:
		return nil
	}
}

// nextFormList cycles to the next list option.
func (m *Model) nextFormList() {
	m.formTargetList = (m.formTargetList + 1) % TabCount
}

// previousFormList cycles to the previous list option.
func (m *Model) previousFormList() {
	m.formTargetList = (m.formTargetList + TabCount - 1) % TabCount
}

// getListByTab returns the todo list for the given tab.
func (m *Model) getListByTab(tab Tab) *model.TodoList {
	switch tab {
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

func (m *Model) dueDateForTab(tab Tab) *time.Time {
	now := time.Now()
	switch tab {
	case TabToday:
		return list.DefaultDueDate(list.TodayID, now)
	case TabTomorrow:
		return list.DefaultDueDate(list.TomorrowID, now)
	case TabTodo:
		return list.DefaultDueDate(list.TodosID, now)
	default:
		return nil
	}
}

// renderForm renders the add or edit todo form.
func (m *Model) renderForm() string {
	var b strings.Builder

	titleStyle := m.theme.ActiveTabStyle()
	formTitle := "Add Todo"
	if m.formMode == FormModeEdit {
		formTitle = "Edit Todo"
	}
	b.WriteString(titleStyle.Render(formTitle))
	b.WriteString("\n\n")

	titleLabel := "Title:"
	if m.formField == FormFieldTitle {
		titleLabel = m.theme.HighlightedItemStyle().Render("❯ Title:")
	} else {
		titleLabel = "  " + titleLabel
	}
	b.WriteString(titleLabel + "\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	descLabel := "Description:"
	if m.formField == FormFieldDescription {
		descLabel = m.theme.HighlightedItemStyle().Render("❯ Description:")
	} else {
		descLabel = "  " + descLabel
	}
	b.WriteString(descLabel + "\n")
	b.WriteString(m.descriptionInput.View())
	b.WriteString("\n\n")

	listLabel := "Add to list:"
	if m.formMode == FormModeEdit {
		listLabel = "Move to list:"
	}
	if m.formField == FormFieldList {
		listLabel = m.theme.HighlightedItemStyle().Render("❯ " + listLabel)
	} else {
		listLabel = "  " + listLabel
	}

	b.WriteString(listLabel + "\n")

	for i := TabToday; i < TabCount; i++ {
		var listStyle lipgloss.Style
		if i == m.formTargetList {
			if m.formField == FormFieldList {
				listStyle = m.theme.ActiveTabStyle()
			} else {
				listStyle = m.theme.HighlightedItemStyle()
			}
		} else {
			listStyle = m.theme.DescriptionStyle()
		}

		indicator := "  "
		if i == m.formTargetList {
			indicator = "▸ "
		}

		b.WriteString("  " + indicator + listStyle.Render(i.String()) + "  ")
	}

	b.WriteString("\n\n")

	helpText := "Tab/Arrows to navigate · ←/→ to select list · Ctrl+S to save · Esc to cancel"
	b.WriteString(m.theme.HelpStyle().Render(helpText))

	return m.theme.ContainerStyle().Render(b.String())
}
