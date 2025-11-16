// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
)

// Storage defines the behavior for persisting todo lists.
type Storage interface {
	// LoadList loads a todo list from storage.
	LoadList(list.Definition) (*model.TodoList, error)
	// SaveList saves a todo list to storage.
	SaveList(list.Definition, *model.TodoList) error
}
