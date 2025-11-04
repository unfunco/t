package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unfunco/t/internal/model"
)

const (
	todayFile    = "today.json"
	tomorrowFile = "tomorrow.json"
	todoFile     = "todo.json"
)

// Storage handles TODO persistence.
type Storage struct {
	dataDir string
}

// New creates a new Storage instance.
func New() (*Storage, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	return &Storage{dataDir}, nil
}

// getDataDir returns the path to the data directory.
func getDataDir() (string, error) {
	var baseDir string

	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		baseDir = xdgDataHome
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".local", "share")
	}

	return filepath.Join(baseDir, "t"), nil
}

// ensureDataDir creates the data directory if it doesn't exist.
func (s *Storage) ensureDataDir() error {
	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	return nil
}

// LoadTodayList loads the today todo list.
func (s *Storage) LoadTodayList() (*model.TodoList, error) {
	return s.loadList(todayFile, "Today")
}

// LoadTomorrowList loads the tomorrow todo list.
func (s *Storage) LoadTomorrowList() (*model.TodoList, error) {
	return s.loadList(tomorrowFile, "Tomorrow")
}

// LoadTodoList loads the general todo list.
func (s *Storage) LoadTodoList() (*model.TodoList, error) {
	return s.loadList(todoFile, "Todos")
}

// loadList loads a todo list from the specified file.
// If the file doesn't exist, it returns an empty list with the given name.
func (s *Storage) loadList(filename, listName string) (*model.TodoList, error) {
	filePath := filepath.Join(s.dataDir, filename)

	// If the file doesn't exist, return an empty list
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &model.TodoList{
			Name:  listName,
			Todos: []model.Todo{},
		}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var todos []model.Todo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
		}
	}

	return &model.TodoList{
		Name:  listName,
		Todos: todos,
	}, nil
}

// SaveToday saves the today todo list to disk.
func (s *Storage) SaveToday(list *model.TodoList) error {
	return s.saveList(todayFile, list)
}

// SaveTomorrow saves the tomorrow todo list to disk.
func (s *Storage) SaveTomorrow(list *model.TodoList) error {
	return s.saveList(tomorrowFile, list)
}

// SaveTodo saves the general todo list to disk.
func (s *Storage) SaveTodo(list *model.TodoList) error {
	return s.saveList(todoFile, list)
}

// saveList saves a todo list to the specified file.
func (s *Storage) saveList(filename string, list *model.TodoList) error {
	if err := s.ensureDataDir(); err != nil {
		return err
	}

	filePath := filepath.Join(s.dataDir, filename)

	data, err := json.MarshalIndent(list.Todos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}

// SaveAll saves all three todo lists to disk.
func (s *Storage) SaveAll(today, tomorrow, todo *model.TodoList) error {
	if err := s.SaveToday(today); err != nil {
		return err
	}
	if err := s.SaveTomorrow(tomorrow); err != nil {
		return err
	}
	if err := s.SaveTodo(todo); err != nil {
		return err
	}
	return nil
}

// DataDir returns the data directory path.
func (s *Storage) DataDir() string {
	return s.dataDir
}
