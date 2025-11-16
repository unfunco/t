// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package automation

import (
	"testing"
	"time"

	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
)

func TestSyncMovesDueTomorrowTodos(t *testing.T) {
	now := time.Date(2025, time.January, 2, 9, 0, 0, 0, time.UTC)
	yesterday := now.Add(-day)

	due := list.DefaultDueDate(list.TomorrowID, yesterday)
	if due == nil {
		t.Fatalf("expected tomorrow list to have a due date")
	}

	store := newMemoryStorage(map[list.ID]*model.TodoList{
		list.TodayID: {
			Name: list.Today().Name,
		},
		list.TomorrowID: {
			Name: list.Tomorrow().Name,
			Todos: []model.Todo{
				{
					ID:        "1",
					Title:     "Move me",
					CreatedAt: yesterday,
					DueDate:   due,
				},
			},
		},
	})

	lists, err := Sync(store, now)
	if err != nil {
		t.Fatalf("Sync returned error: %v", err)
	}

	if got := len(lists[list.TomorrowID].Todos); got != 0 {
		t.Fatalf("expected tomorrow list to be empty, got %d todos", got)
	}

	if got := len(lists[list.TodayID].Todos); got != 1 {
		t.Fatalf("expected today list to have 1 todo, got %d", got)
	}

	todo := lists[list.TodayID].Todos[0]
	wantDue := list.DefaultDueDate(list.TodayID, now)
	if todo.DueDate == nil || wantDue == nil || !todo.DueDate.Equal(*wantDue) {
		t.Fatalf("expected due date %v, got %v", wantDue, todo.DueDate)
	}
}

func TestSyncLeavesFutureTomorrowTodos(t *testing.T) {
	now := time.Date(2025, time.January, 1, 9, 0, 0, 0, time.UTC)

	due := list.DefaultDueDate(list.TomorrowID, now)
	store := newMemoryStorage(map[list.ID]*model.TodoList{
		list.TomorrowID: {
			Name: list.Tomorrow().Name,
			Todos: []model.Todo{
				{
					ID:        "1",
					Title:     "Too early",
					CreatedAt: now,
					DueDate:   due,
				},
			},
		},
	})

	lists, err := Sync(store, now)
	if err != nil {
		t.Fatalf("Sync returned error: %v", err)
	}

	if got := len(lists[list.TomorrowID].Todos); got != 1 {
		t.Fatalf("expected tomorrow list to keep todo, got %d", got)
	}
}

func TestSyncBackfillsMissingDueDates(t *testing.T) {
	now := time.Date(2025, time.January, 3, 9, 0, 0, 0, time.UTC)
	created := now.Add(-2 * day)

	store := newMemoryStorage(map[list.ID]*model.TodoList{
		list.TodayID: {
			Name: list.Today().Name,
			Todos: []model.Todo{
				{
					ID:        "a",
					Title:     "Missing due date",
					CreatedAt: created,
				},
			},
		},
	})

	lists, err := Sync(store, now)
	if err != nil {
		t.Fatalf("Sync returned error: %v", err)
	}

	todos := lists[list.TodayID].Todos
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo after sync, got %d", len(todos))
	}

	want := list.DefaultDueDate(list.TodayID, created)
	if want == nil {
		t.Fatalf("expected today list to produce a due date")
	}

	if todos[0].DueDate == nil || !todos[0].DueDate.Equal(*want) {
		t.Fatalf("expected due date %v, got %v", want, todos[0].DueDate)
	}
}

type memoryStorage struct {
	lists map[list.ID]*model.TodoList
}

func newMemoryStorage(initial map[list.ID]*model.TodoList) *memoryStorage {
	lists := make(map[list.ID]*model.TodoList, len(initial))
	for id, todoList := range initial {
		lists[id] = todoList
	}

	return &memoryStorage{lists: lists}
}

func (m *memoryStorage) LoadList(def list.Definition) (*model.TodoList, error) {
	if l, ok := m.lists[def.ID]; ok {
		return l, nil
	}

	l := &model.TodoList{Name: def.Name}
	m.lists[def.ID] = l

	return l, nil
}

func (m *memoryStorage) SaveList(def list.Definition, todoList *model.TodoList) error {
	m.lists[def.ID] = todoList
	return nil
}
