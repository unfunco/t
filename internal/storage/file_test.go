// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
)

func TestFileSaveListPrivateAtomic(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "store")

	store, err := NewFileStorageWithDir(dataDir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	def := list.Definition{
		Name:     "Inbox",
		Filename: "inbox.json",
	}
	todoList := &model.TodoList{
		Name: def.Name,
		Todos: []model.Todo{
			{Title: "secret"},
		},
	}

	if err := store.SaveList(def, todoList); err != nil {
		t.Fatalf("SaveList() returned error: %v", err)
	}

	stored, err := store.LoadList(def)
	if err != nil {
		t.Fatalf("LoadList() returned error: %v", err)
	}

	if len(stored.Todos) != len(todoList.Todos) || stored.Todos[0].Title != todoList.Todos[0].Title {
		t.Fatalf("stored todos do not match original: %+v vs %+v", stored.Todos, todoList.Todos)
	}

	if runtime.GOOS != "windows" {
		dirInfo, err := os.Stat(dataDir)
		if err != nil {
			t.Fatalf("failed to stat data dir: %v", err)
		}

		if perm := dirInfo.Mode().Perm(); perm != 0o700 {
			t.Fatalf("unexpected dir permissions, want 0700 got %o", perm)
		}

		fileInfo, err := os.Stat(filepath.Join(dataDir, def.Filename))
		if err != nil {
			t.Fatalf("failed to stat data file: %v", err)
		}

		if perm := fileInfo.Mode().Perm(); perm != 0o600 {
			t.Fatalf("unexpected file permissions, want 0600 got %o", perm)
		}
	}

	tmpFiles, err := filepath.Glob(filepath.Join(dataDir, def.Filename+".tmp-*"))
	if err != nil {
		t.Fatalf("glob failed: %v", err)
	}

	if len(tmpFiles) != 0 {
		t.Fatalf("temporary files leaked: %v", tmpFiles)
	}
}
