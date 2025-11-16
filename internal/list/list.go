// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package list

import "time"

// ID identifies a todo list.
type ID string

const (
	// TodayID represents the "Today" list.
	TodayID ID = "today"
	// TomorrowID represents the "Tomorrow" list.
	TomorrowID ID = "tomorrow"
	// TodosID represents the general "Todos" list.
	TodosID ID = "todos"
)

// Definition contains the metadata needed to load/store a list.
type Definition struct {
	ID       ID
	Name     string
	Filename string
}

// day is the number of hours in a full calendar day.
const day = 24 * time.Hour

var definitions = map[ID]Definition{
	TodayID: {
		ID:       TodayID,
		Name:     "Today",
		Filename: "today.json",
	},
	TomorrowID: {
		ID:       TomorrowID,
		Name:     "Tomorrow",
		Filename: "tomorrow.json",
	},
	TodosID: {
		ID:       TodosID,
		Name:     "Todos",
		Filename: "todo.json",
	},
}

var orderedDefinitions = []Definition{
	definitions[TodayID],
	definitions[TomorrowID],
	definitions[TodosID],
}

// Default returns a copy of the default list definitions in UI order.
func Default() []Definition {
	out := make([]Definition, len(orderedDefinitions))
	copy(out, orderedDefinitions)
	return out
}

// Today returns the default Today list definition.
func Today() Definition {
	return definitions[TodayID]
}

// Tomorrow returns the default Tomorrow list definition.
func Tomorrow() Definition {
	return definitions[TomorrowID]
}

// Todos returns the default Todos list definition.
func Todos() Definition {
	return definitions[TodosID]
}

// DefaultDueDate returns the default due date for items added to the provided
// list ID. Lists that do not have a due date return nil.
func DefaultDueDate(id ID, now time.Time) *time.Time {
	switch id {
	case TodayID:
		t := startOfDay(now)
		return &t
	case TomorrowID:
		t := startOfDay(now).Add(day)
		return &t
	default:
		return nil
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
