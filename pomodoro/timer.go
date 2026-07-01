package pomodoro

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Timer type for Pomodoro
type TimerType string

const (
	Work       TimerType = "Work"
	ShortBreak TimerType = "Short Break"
	LongBreak  TimerType = "Long Break"
)

// Session represents a single pomodoro session
type Session struct {
	Type     TimerType
	Duration time.Duration
}

// Run executes the pomodoro session
func (s *Session) Run() {
	fmt.Printf("\nStarting %s for %v\n", s.Type, s.Duration)

	totalSeconds := int(s.Duration.Seconds())

	bar := progressbar.NewOptions(totalSeconds,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s[reset]", s.Type)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i := 0; i < totalSeconds; i++ {
		time.Sleep(1 * time.Second)
		bar.Add(1)
	}
	fmt.Printf("\n%s completed!\n", s.Type)
}
