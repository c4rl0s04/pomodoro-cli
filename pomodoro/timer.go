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

// addShadow applies a manual drop shadow to a pterm BigText string
func addShadow(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	height := len(lines)
	width := 0
	for _, l := range lines {
		runeCount := len([]rune(l))
		if runeCount > width {
			width = runeCount
		}
	}

	grid := make([][]rune, height+1)
	for i := range grid {
		grid[i] = make([]rune, width+2)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for y, line := range lines {
		runes := []rune(line)
		for x, char := range runes {
			if char != ' ' {
				grid[y+1][x+2] = '░'
			}
		}
	}

	for y, line := range lines {
		runes := []rune(line)
		for x, char := range runes {
			if char != ' ' {
				grid[y][x] = char
			}
		}
	}

	var sb strings.Builder
	for _, row := range grid {
		for _, char := range row {
			if char == '░' {
				sb.WriteString(pterm.FgDarkGray.Sprintf("%c", char))
			} else if char != ' ' {
				sb.WriteString(pterm.FgCyan.Sprintf("%c", char))
			} else {
				sb.WriteRune(char)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Run executes the pomodoro session with a full screen BigText clock
func (s *Session) Run() {
	// Enter Alternate Screen Buffer and hide cursor
	fmt.Print("\033[?1049h\033[?25l")
	defer fmt.Print("\033[?1049l\033[?25h")
	fmt.Print("\033[2J\033[H")

	pterm.EnableOutput()
	area, _ := pterm.DefaultArea.WithFullscreen().Start()
	defer area.Stop()

	totalSeconds := int(s.Duration.Seconds())

	for i := totalSeconds; i >= 0; i-- {
		minutes := i / 60
		seconds := i % 60
		timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

		// Create big text for the phase
		phaseStr, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromString(strings.ToUpper(string(s.Type))),
		).Srender()

		// Create big text for the time
		timeStr, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromString(timeString),
		).Srender()

		// Combine, strip ANSI codes, and add shadow
		combined := phaseStr + "\n" + timeStr
		cleanText := pterm.RemoveColorFromString(combined)
		shadowedText := addShadow(cleanText)

		_, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			height = 24
		}

		lines := strings.Split(shadowedText, "\n")
		textHeight := len(lines)
		paddingTop := (height - textHeight) / 2
		if paddingTop < 0 {
			paddingTop = 0
		}

		verticalPadding := strings.Repeat("\n", paddingTop)
		centeredText := pterm.DefaultCenter.Sprint(shadowedText)
		finalOutput := verticalPadding + centeredText

		area.Update(finalOutput)

		if i > 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
