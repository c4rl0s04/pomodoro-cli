package pomodoro

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
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

// Run executes the pomodoro session with a full screen BigText clock
func (s *Session) Run() {
	// Enter Alternate Screen Buffer and hide cursor
	fmt.Print("\033[?1049h\033[?25l")
	// Ensure we leave Alternate Screen Buffer and show cursor when done
	defer fmt.Print("\033[?1049l\033[?25h")

	// Clear the screen explicitly
	fmt.Print("\033[2J\033[H")

	// Enable pterm output
	pterm.EnableOutput()

	// Create an area for updating the screen
	area, _ := pterm.DefaultArea.WithFullscreen().Start()
	defer area.Stop()

	totalSeconds := int(s.Duration.Seconds())

	for i := totalSeconds; i >= 0; i-- {
		// Calculate minutes and seconds
		minutes := i / 60
		seconds := i % 60

		// Format the time string e.g., "25:00"
		timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

		// Create big text for the time
		bigTextStr, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromString(timeString),
		).Srender()

		// Center the text in the terminal horizontally and vertically
		centeredText := pterm.DefaultCenter.Sprint(fmt.Sprintf("%s\n\n%s", s.Type, bigTextStr))

		// Update the area
		area.Update(centeredText)

		// Sleep for 1 second if not done
		if i > 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
