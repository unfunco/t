// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/paths"
)

// File persists todos on disk.
type File struct {
	dataDir string
}

var _ Storage = (*File)(nil)

// NewFileStorage creates file-backed storage rooted in the default data
// directory, typically ~/.local/share/t.
func NewFileStorage() (*File, error) {
	dataDir, err := paths.DefaultDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	return NewFileStorageWithDir(dataDir)
}

// NewFileStorageWithDir creates a File rooted at the provided directory.
func NewFileStorageWithDir(dataDir string) (*File, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("data directory cannot be empty")
	}

	return &File{dataDir}, nil
}

// ensureDataDir creates the data directory if it doesn't exist.
func (s *File) ensureDataDir() error {
	if err := os.MkdirAll(s.dataDir, 0o700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	if err := os.Chmod(s.dataDir, 0o700); err != nil {
		return fmt.Errorf("failed to set permissions on data directory: %w", err)
	}

	return nil
}

// LoadList loads the todo list represented by the provided definition.
func (s *File) LoadList(def list.Definition) (*model.TodoList, error) {
	filePath := filepath.Join(s.dataDir, def.Filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &model.TodoList{
			Name:  def.Name,
			Todos: []model.Todo{},
		}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", def.Filename, err)
	}

	var todos []model.Todo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", def.Filename, err)
		}
	}

	return &model.TodoList{
		Name:  def.Name,
		Todos: todos,
	}, nil
}

// SaveList saves the provided todo list using the supplied definition.
func (s *File) SaveList(def list.Definition, list *model.TodoList) error {
	if err := s.ensureDataDir(); err != nil {
		return err
	}

	filePath := filepath.Join(s.dataDir, def.Filename)

	data, err := json.MarshalIndent(list.Todos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}

	tmpFile, err := os.CreateTemp(s.dataDir, def.Filename+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file for %s: %w", def.Filename, err)
	}

	tmpPath := tmpFile.Name()
	var writeErr error
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
		}
		if writeErr != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmpFile.Chmod(0o600); err != nil {
		writeErr = fmt.Errorf("failed to set permissions on %s: %w", tmpPath, err)
		return writeErr
	}

	if _, err := tmpFile.Write(data); err != nil {
		writeErr = fmt.Errorf("failed to write %s: %w", tmpPath, err)
		return writeErr
	}

	if err := tmpFile.Sync(); err != nil {
		writeErr = fmt.Errorf("failed to sync %s: %w", tmpPath, err)
		return writeErr
	}

	if err := tmpFile.Close(); err != nil {
		writeErr = fmt.Errorf("failed to close %s: %w", tmpPath, err)
		return writeErr
	}

	tmpFile = nil

	if err := os.Rename(tmpPath, filePath); err != nil {
		writeErr = fmt.Errorf("failed to replace %s: %w", def.Filename, err)
		return writeErr
	}

	return nil
}

// DataDir returns the data directory path.
func (s *File) DataDir() string {
	return s.dataDir
}
