// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package automation

import (
	"fmt"
	"time"

	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/storage"
)

const day = 24 * time.Hour

// Sync loads all default lists, applies scheduled automations, persists any
// changes, and returns the resulting lists keyed by their ID.
func Sync(store storage.Storage, now time.Time) (map[list.ID]*model.TodoList, error) {
	defs := list.Default()
	lists := make(map[list.ID]*model.TodoList, len(defs))

	for _, def := range defs {
		l, err := store.LoadList(def)
		if err != nil {
			return nil, fmt.Errorf("load %s list: %w", def.Name, err)
		}

		lists[def.ID] = l
	}

	todayStart := startOfDay(now)
	changed := false

	if ensureDueDates(lists[list.TodayID], list.TodayID) {
		changed = true
	}

	if ensureDueDates(lists[list.TomorrowID], list.TomorrowID) {
		changed = true
	}

	if moveTomorrowTodos(lists[list.TomorrowID], lists[list.TodayID], todayStart) {
		changed = true
	}

	if changed {
		for _, def := range defs {
			l := lists[def.ID]
			if l == nil {
				continue
			}
			if err := store.SaveList(def, l); err != nil {
				return nil, fmt.Errorf("save %s list: %w", def.Name, err)
			}
		}
	}

	return lists, nil
}

func ensureDueDates(todoList *model.TodoList, id list.ID) bool {
	if todoList == nil {
		return false
	}

	changed := false
	for i := range todoList.Todos {
		todo := &todoList.Todos[i]

		switch id {
		case list.TodayID:
			if applyDueDate(todo, list.DefaultDueDate(list.TodayID, todo.CreatedAt)) {
				changed = true
			}
		case list.TomorrowID:
			if applyDueDate(todo, list.DefaultDueDate(list.TomorrowID, todo.CreatedAt)) {
				changed = true
			}
		default:
			if todo.DueDate != nil {
				normalized := startOfDay(*todo.DueDate)
				if !todo.DueDate.Equal(normalized) {
					todo.SetDueDate(&normalized)
					changed = true
				}
			}
		}
	}

	return changed
}

func moveTomorrowTodos(tomorrowList, todayList *model.TodoList, todayStart time.Time) bool {
	if tomorrowList == nil || todayList == nil {
		return false
	}

	var remaining []model.Todo
	changed := false

	for _, todo := range tomorrowList.Todos {
		if todo.Completed {
			remaining = append(remaining, todo)
			continue
		}

		if todo.DueDate == nil {
			remaining = append(remaining, todo)
			continue
		}

		due := startOfDay(*todo.DueDate)
		if due.After(todayStart) {
			remaining = append(remaining, todo)
			continue
		}

		todayList.Todos = append(todayList.Todos, todo)
		changed = true
	}

	if len(remaining) != len(tomorrowList.Todos) {
		tomorrowList.Todos = remaining
	}

	return changed
}

func applyDueDate(todo *model.Todo, due *time.Time) bool {
	if due == nil {
		if todo.DueDate == nil {
			return false
		}
		todo.SetDueDate(nil)
		return true
	}

	normalizedTarget := startOfDay(*due)

	if todo.DueDate == nil {
		todo.SetDueDate(&normalizedTarget)
		return true
	}

	current := startOfDay(*todo.DueDate)
	if current.Equal(normalizedTarget) {
		return false
	}

	todo.SetDueDate(&normalizedTarget)
	return true
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
