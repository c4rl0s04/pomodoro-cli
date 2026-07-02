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

// Engine manages the pomodoro timer state
type Engine struct {
	config Config
}

// NewEngine creates a new Engine
func NewEngine(cfg Config) *Engine {
	return &Engine{config: cfg}
}

// Run starts the engine and emits ticks via the provided channel
func (e *Engine) Run(tickChan chan<- Tick) {
	defer close(tickChan)

	for i := 1; i <= e.config.Cycles; i++ {
		// Work session
		e.runSession(Work, e.config.WorkDuration, tickChan)

		// Decide break type
		if i%4 == 0 {
			e.runSession(LongBreak, e.config.LongBreakDuration, tickChan)
		} else if i < e.config.Cycles {
			e.runSession(ShortBreak, e.config.ShortBreakDuration, tickChan)
		}
	}
}

func (e *Engine) runSession(sessionType SessionType, duration time.Duration, tickChan chan<- Tick) {
	totalSeconds := int(duration.Seconds())

	for i := totalSeconds; i >= 0; i-- {
		tickChan <- Tick{
			Type:          sessionType,
			TimeRemaining: time.Duration(i) * time.Second,
		}

		if i > 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
