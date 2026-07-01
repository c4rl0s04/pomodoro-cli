package cmd

import (
	"time"

	"github.com/carlosandreshuete/pomodoro-cli/pomodoro"
	"github.com/spf13/cobra"
)

var (
	workDuration       int
	shortBreakDuration int
	longBreakDuration  int
	cycles             int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Pomodoro session",
	Long:  `Start a sequence of work and break timers according to the Pomodoro technique.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 1; i <= cycles; i++ {
			// Work session
			workSession := pomodoro.Session{
				Type:     pomodoro.Work,
				Duration: time.Duration(workDuration) * time.Minute,
			}
			workSession.Run()

			// Decide break type
			if i%4 == 0 {
				longBreakSession := pomodoro.Session{
					Type:     pomodoro.LongBreak,
					Duration: time.Duration(longBreakDuration) * time.Minute,
				}
				longBreakSession.Run()
			} else if i < cycles {
				shortBreakSession := pomodoro.Session{
					Type:     pomodoro.ShortBreak,
					Duration: time.Duration(shortBreakDuration) * time.Minute,
				}
				shortBreakSession.Run()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&workDuration, "work", "w", 25, "Work duration in minutes")
	startCmd.Flags().IntVarP(&shortBreakDuration, "short-break", "s", 5, "Short break duration in minutes")
	startCmd.Flags().IntVarP(&longBreakDuration, "long-break", "l", 15, "Long break duration in minutes")
	startCmd.Flags().IntVarP(&cycles, "cycles", "c", 4, "Number of pomodoro cycles")
}
