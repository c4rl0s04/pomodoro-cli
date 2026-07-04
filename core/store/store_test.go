package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestJSONStore_SaveAndGetSessions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pomodoro-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	filePath := filepath.Join(tempDir, "stats.json")
	store := &JSONStore{filePath: filePath}

	// 1. Initial state should be empty
	sessions, err := store.GetSessions()
	if err != nil {
		t.Fatalf("expected no error getting empty sessions, got %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions initially, got %d", len(sessions))
	}

	// 2. Save a session
	now := time.Now().Round(time.Second) // round for JSON comparison safety
	record1 := SessionRecord{Duration: 25, CompletedAt: now}
	if err := store.SaveSession(record1); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// 3. Verify it was saved
	sessions, err = store.GetSessions()
	if err != nil {
		t.Fatalf("failed to get sessions after save: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Duration != record1.Duration {
		t.Errorf("expected duration %d, got %d", record1.Duration, sessions[0].Duration)
	}

	// 4. Save another session
	record2 := SessionRecord{Duration: 5, CompletedAt: now.Add(5 * time.Minute)}
	if err := store.SaveSession(record2); err != nil {
		t.Fatalf("failed to save second session: %v", err)
	}

	// 5. Verify both
	sessions, err = store.GetSessions()
	if err != nil {
		t.Fatalf("failed to get sessions: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestJSONStore_CorruptData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pomodoro-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "corrupt.json")
	// Write corrupt JSON
	os.WriteFile(filePath, []byte("{ invalid json ]"), 0644)

	store := &JSONStore{filePath: filePath}
	_, err = store.GetSessions()
	if err == nil {
		t.Errorf("expected error when reading corrupt JSON, got nil")
	}
}

func TestJSONStore_EmptyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pomodoro-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "empty.json")
	os.WriteFile(filePath, []byte(""), 0644)

	store := &JSONStore{filePath: filePath}
	sessions, err := store.GetSessions()
	if err != nil {
		t.Errorf("expected no error for empty file, got %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions for empty file, got %d", len(sessions))
	}
}
