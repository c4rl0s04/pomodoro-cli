package pomodoro

import (
	"time"
)

// SessionType represents the type of a pomodoro session
type SessionType string

const (
	Work       SessionType = "Work"
	ShortBreak SessionType = "Short Break"
	LongBreak  SessionType = "Long Break"
)

// Config holds the configuration for a Pomodoro session
type Config struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	Cycles             int
}

// Tick represents a single second tick in the timer
type Tick struct {
	Type          SessionType
	TimeRemaining time.Duration
}

// ControlMsg represents an interactive command from the user
type ControlMsg int

const (
	PauseResume ControlMsg = iota
	Skip
	Quit
)

// Engine manages the pomodoro timer state
type Engine struct {
	config Config
}

// NewEngine creates a new Engine
func NewEngine(cfg Config) *Engine {
	return &Engine{config: cfg}
}

// Run starts the engine and emits ticks via the provided channel
func (e *Engine) Run(tickChan chan<- Tick, controlChan <-chan ControlMsg) {
	defer close(tickChan)

	for i := 1; i <= e.config.Cycles; i++ {
		// Work session
		if !e.runSession(Work, e.config.WorkDuration, tickChan, controlChan) {
			return // Quit requested
		}

		// Decide break type
		if i%4 == 0 {
			if !e.runSession(LongBreak, e.config.LongBreakDuration, tickChan, controlChan) {
				return
			}
		} else if i < e.config.Cycles {
			if !e.runSession(ShortBreak, e.config.ShortBreakDuration, tickChan, controlChan) {
				return
			}
		}
	}
}

// runSession returns true if completed or skipped, false if quit requested
func (e *Engine) runSession(sessionType SessionType, duration time.Duration, tickChan chan<- Tick, controlChan <-chan ControlMsg) bool {
	remaining := int(duration.Seconds())
	paused := false

	// Emit initial tick immediately
	tickChan <- Tick{
		Type:          sessionType,
		TimeRemaining: time.Duration(remaining) * time.Second,
	}

	if remaining <= 0 {
		return true
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-controlChan:
			switch msg {
			case PauseResume:
				paused = !paused
			case Skip:
				return true
			case Quit:
				return false
			}
		case <-ticker.C:
			if !paused {
				remaining--
				tickChan <- Tick{
					Type:          sessionType,
					TimeRemaining: time.Duration(remaining) * time.Second,
				}
				if remaining <= 0 {
					return true
				}
			}
		}
	}
}
