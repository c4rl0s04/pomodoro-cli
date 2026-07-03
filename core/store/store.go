package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionRecord holds data about a single completed Pomodoro session
type SessionRecord struct {
	Duration    int       `json:"duration_minutes"`
	CompletedAt time.Time `json:"completed_at"`
}

// Store defines the interface for persisting completed sessions
type Store interface {
	SaveSession(record SessionRecord) error
	GetSessions() ([]SessionRecord, error)
}

// JSONStore implements Store by saving records to a local JSON file
type JSONStore struct {
	filePath string
	mu       sync.Mutex
}

// NewJSONStore creates a new JSONStore. It uses ~/.pomodoro/stats.json by default.
func NewJSONStore() (*JSONStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	dir := filepath.Join(home, ".pomodoro")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create directory %s: %w", dir, err)
	}

	filePath := filepath.Join(dir, "stats.json")
	return &JSONStore{filePath: filePath}, nil
}

// SaveSession appends a new session record to the JSON file
func (s *JSONStore) SaveSession(record SessionRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.readSessions()
	if err != nil {
		return err
	}

	sessions = append(sessions, record)
	return s.writeSessions(sessions)
}

// GetSessions retrieves all session records from the JSON file
func (s *JSONStore) GetSessions() ([]SessionRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readSessions()
}

func (s *JSONStore) readSessions() ([]SessionRecord, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionRecord{}, nil
		}
		return nil, fmt.Errorf("failed to read stats file: %w", err)
	}

	if len(data) == 0 {
		return []SessionRecord{}, nil
	}

	var sessions []SessionRecord
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("failed to parse stats file: %w", err)
	}

	return sessions, nil
}

func (s *JSONStore) writeSessions(sessions []SessionRecord) error {
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode stats data: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}

	return nil
}
