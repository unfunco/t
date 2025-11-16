// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package list

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
