package pomodoro

import (
	"time"

	"github.com/carlosandreshuete/pomodoro-cli/core/notifier"
	"github.com/carlosandreshuete/pomodoro-cli/core/store"
)

// SessionType represents the type of a pomodoro session
type SessionType string

const (
	Work       SessionType = "Focus"
	ShortBreak SessionType = "Break"
	LongBreak  SessionType = "Break"
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
	IsPaused      bool
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
	config   Config
	notifier notifier.Notifier
	store    store.Store
}

// NewEngine creates a new Engine
func NewEngine(cfg Config) *Engine {
	return &Engine{config: cfg}
}

// SetNotifier attaches a desktop notifier to the engine
func (e *Engine) SetNotifier(n notifier.Notifier) {
	e.notifier = n
}

// SetStore attaches a persistent storage backend to the engine
func (e *Engine) SetStore(s store.Store) {
	e.store = s
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
		IsPaused:      paused,
	}

	if remaining <= 0 {
		e.notifySessionEnd(sessionType)
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
				tickChan <- Tick{
					Type:          sessionType,
					TimeRemaining: time.Duration(remaining) * time.Second,
					IsPaused:      paused,
				}
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
					IsPaused:      paused,
				}
				if remaining <= 0 {
					e.notifySessionEnd(sessionType)
					return true
				}
			}
		}
	}
}

func (e *Engine) notifySessionEnd(sessionType SessionType) {
	// Log the work session explicitly
	if sessionType == Work && e.store != nil {
		_ = e.store.SaveSession(store.SessionRecord{
			Duration:    int(e.config.WorkDuration.Minutes()),
			CompletedAt: time.Now(),
		})
	}

	if e.notifier != nil {
		var title, msg string
		if sessionType == Work {
			title = "Work Session Complete!"
			msg = "Time to take a break."
		} else {
			title = "Break Complete!"
			msg = "Time to get back to work."
		}
		_ = e.notifier.Notify(title, msg)
	}
}
