package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/carlosandreshuete/pomodoro-cli/core/pomodoro"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// CLI handles the terminal user interface rendering
type CLI struct {
	area *pterm.AreaPrinter
}

// NewCLI creates a new CLI UI manager
func NewCLI() *CLI {
	return &CLI{}
}

// Start prepares the terminal for rendering
func (c *CLI) Start() error {
	// Enter Alternate Screen Buffer and hide cursor
	fmt.Print("\033[?1049h\033[?25l")
	// Clear the screen explicitly
	fmt.Print("\033[2J\033[H")

	pterm.EnableOutput()
	area, err := pterm.DefaultArea.WithFullscreen().Start()
	if err != nil {
		return fmt.Errorf("failed to start pterm area: %w", err)
	}
	c.area = area
	return nil
}

// Stop restores the terminal
func (c *CLI) Stop() {
	if c.area != nil {
		c.area.Stop()
	}
	// Leave Alternate Screen Buffer and show cursor
	fmt.Print("\033[?1049l\033[?25h")
}

// RenderTick renders a single tick event to the screen
func (c *CLI) RenderTick(tick pomodoro.Tick) {
	minutes := int(tick.TimeRemaining.Minutes())
	seconds := int(tick.TimeRemaining.Seconds()) % 60
	timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

	// Create big text for the phase
	phaseStr, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString(strings.ToUpper(string(tick.Type))),
	).Srender()

	// Create big text for the time
	timeStr, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString(timeString),
	).Srender()

	combined := phaseStr + "\n" + timeStr
	cleanText := pterm.RemoveColorFromString(combined)
	shadowedText := c.addShadow(cleanText)

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

	c.area.Update(finalOutput)
}

func (c *CLI) addShadow(text string) string {
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
