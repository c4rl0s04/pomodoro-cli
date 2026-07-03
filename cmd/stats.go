package cmd

import (
	"fmt"
	"time"

	"github.com/carlosandreshuete/pomodoro-cli/core/store"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your Pomodoro session statistics",
	Long:  "Displays a beautiful bar chart of your completed focus sessions (in minutes) over the past 7 days.",
	Run: func(cmd *cobra.Command, args []string) {
		jsonStore, err := store.NewJSONStore()
		if err != nil {
			pterm.Error.Println("Failed to initialize session store:", err)
			return
		}

		sessions, err := jsonStore.GetSessions()
		if err != nil {
			pterm.Error.Println("Failed to read session data:", err)
			return
		}

		if len(sessions) == 0 {
			pterm.Info.Println("No sessions completed yet. Start working!")
			return
		}

		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// Group minutes by date string
		data := make(map[string]int)
		for _, s := range sessions {
			dateStr := s.CompletedAt.Local().Format("2006-01-02")
			data[dateStr] += s.Duration
		}

		bars := make(pterm.Bars, 7)
		for i := 6; i >= 0; i-- {
			date := midnight.AddDate(0, 0, -i)
			dateStr := date.Format("2006-01-02")

			bars[6-i] = pterm.Bar{
				Label: date.Format("Mon"),
				Value: data[dateStr],
			}
		}

		pterm.DefaultHeader.WithFullWidth().
			WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
			WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
			Println("Focus Time (Minutes) - Last 7 Days")
		fmt.Println()

		chart := pterm.DefaultBarChart.WithBars(bars).WithShowValue()
		rendered, _ := chart.Srender()
		fmt.Println(rendered)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
