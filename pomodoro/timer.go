package pomodoro

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"
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

		// Get terminal size to center vertically
		_, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			height = 24
		}

		// Split big text into lines to count height
		lines := strings.Split(bigTextStr, "\n")
		textHeight := len(lines) + 2 // +2 for the "Work" type header

		// Calculate vertical padding
		paddingTop := (height - textHeight) / 2
		if paddingTop < 0 {
			paddingTop = 0
		}

		verticalPadding := strings.Repeat("\n", paddingTop)

		// Center the text horizontally
		centeredText := pterm.DefaultCenter.Sprint(fmt.Sprintf("%s\n\n%s", s.Type, bigTextStr))

		// Apply vertical padding and final output
		finalOutput := verticalPadding + centeredText

		// Update the area
		area.Update(finalOutput)

		// Sleep for 1 second if not done
		if i > 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
