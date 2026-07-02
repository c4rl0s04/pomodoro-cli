package pomodoro

import (
	"testing"
)

func TestEngine_Run(t *testing.T) {
	// Configure for a very short session (0 duration so it doesn't sleep but runs the loop once)
	cfg := Config{
		WorkDuration:       0,
		ShortBreakDuration: 0,
		LongBreakDuration:  0,
		Cycles:             1,
	}

	engine := NewEngine(cfg)
	tickChan := make(chan Tick, 10)

	engine.Run(tickChan)

	// Since work duration is 0, we expect 1 tick (at 0 seconds) for Work
	// And since it's cycle 1 (out of 1), it's the last cycle, so it stops after short break?
	// Wait, the logic for short break says "if i < e.config.Cycles", so it will NOT run a break on the last cycle.
	// Therefore, we only expect 1 tick for Work.

	ticks := []Tick{}
	for tick := range tickChan {
		ticks = append(ticks, tick)
	}

	if len(ticks) != 1 {
		t.Fatalf("Expected 1 tick, got %d", len(ticks))
	}

	if ticks[0].Type != Work {
		t.Errorf("Expected session type Work, got %s", ticks[0].Type)
	}

	if ticks[0].TimeRemaining != 0 {
		t.Errorf("Expected 0 time remaining, got %v", ticks[0].TimeRemaining)
	}
}
