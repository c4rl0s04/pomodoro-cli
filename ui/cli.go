package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/carlosandreshuete/pomodoro-cli/core/pomodoro"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
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
		_ = c.area.Stop()
	}
	// Leave Alternate Screen Buffer and show cursor
	fmt.Print("\033[?1049l\033[?25h")
}

// ListenKeyboard captures raw keystrokes and sends control messages
func (c *CLI) ListenKeyboard(controlChan chan<- pomodoro.ControlMsg) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			return
		}

		switch b[0] {
		case ' ':
			controlChan <- pomodoro.PauseResume
		case 's', 'S', 'n', 'N':
			controlChan <- pomodoro.Skip
		case 'q', 'Q', 3: // 3 is Ctrl+C
			controlChan <- pomodoro.Quit
			return
		}
	}
}

// RenderTick renders a single tick event to the screen
func (c *CLI) RenderTick(tick pomodoro.Tick) {
	minutes := int(tick.TimeRemaining.Minutes())
	seconds := int(tick.TimeRemaining.Seconds()) % 60
	timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

	// Create big text for the phase
	phaseRaw, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromString(strings.ToUpper(string(tick.Type))),
	).Srender()
	phaseStr := pterm.RemoveColorFromString(phaseRaw)

	// Create big text for the time
	timeRaw, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromString(timeString),
	).Srender()
	timeStr := pterm.RemoveColorFromString(timeRaw)

	phaseStr, timeStr = c.centerBlocks(phaseStr, timeStr)

	combined := phaseStr + "\n" + timeStr
	shadowedText := c.addShadow(combined, tick.IsPaused)

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

	// Because we put the terminal in Raw mode to capture keystrokes,
	// the terminal no longer automatically translates \n to \r\n.
	// We must manually add carriage returns to prevent diagonal text scattering.
	finalOutput = strings.ReplaceAll(finalOutput, "\r\n", "\n")
	finalOutput = strings.ReplaceAll(finalOutput, "\n", "\r\n")

	c.area.Update(finalOutput)
}

func (c *CLI) addShadow(text string, isPaused bool) string {
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
				if isPaused {
					sb.WriteString(pterm.FgRed.Sprintf("%c", char))
				} else {
					sb.WriteString(pterm.FgYellow.Sprintf("%c", char))
				}
			} else {
				sb.WriteRune(char)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c *CLI) centerBlocks(block1, block2 string) (string, string) {
	w1 := c.getBlockWidth(block1)
	w2 := c.getBlockWidth(block2)

	if w1 < w2 {
		pad := strings.Repeat(" ", (w2-w1)/2)
		block1 = c.padBlockLeft(block1, pad)
	} else if w2 < w1 {
		pad := strings.Repeat(" ", (w1-w2)/2)
		block2 = c.padBlockLeft(block2, pad)
	}
	return block1, block2
}

func (c *CLI) getBlockWidth(block string) int {
	max := 0
	for _, line := range strings.Split(block, "\n") {
		count := len([]rune(line))
		if count > max {
			max = count
		}
	}
	return max
}

func (c *CLI) padBlockLeft(block, padding string) string {
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if len(line) > 0 { // Only pad lines that actually have content
			lines[i] = padding + line
		}
	}
	return strings.Join(lines, "\n")
}
