package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/carlosandreshuete/pomodoro-cli/core/pomodoro"
	"github.com/carlosandreshuete/pomodoro-cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Pomodoro session",
	Long:  `Start a sequence of work and break timers according to the Pomodoro technique.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read config from viper (which includes flags, env vars, and config file)
		cfg := pomodoro.Config{
			WorkDuration:       time.Duration(viper.GetInt("work")) * time.Minute,
			ShortBreakDuration: time.Duration(viper.GetInt("short-break")) * time.Minute,
			LongBreakDuration:  time.Duration(viper.GetInt("long-break")) * time.Minute,
			Cycles:             viper.GetInt("cycles"),
		}

		// Initialize Engine
		engine := pomodoro.NewEngine(cfg)

		// Initialize UI
		termUI := ui.NewCLI()
		if err := termUI.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to initialize UI:", err)
			os.Exit(1)
		}
		defer termUI.Stop()

		// Channel for ticks
		tickChan := make(chan pomodoro.Tick)

		// Start engine in a goroutine
		go engine.Run(tickChan)

		// Listen to ticks and update UI
		for tick := range tickChan {
			termUI.RenderTick(tick)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Define flags
	startCmd.Flags().IntP("work", "w", 25, "Work duration in minutes")
	startCmd.Flags().IntP("short-break", "s", 5, "Short break duration in minutes")
	startCmd.Flags().IntP("long-break", "l", 15, "Long break duration in minutes")
	startCmd.Flags().IntP("cycles", "c", 4, "Number of pomodoro cycles")

	// Bind flags to viper
	viper.BindPFlag("work", startCmd.Flags().Lookup("work"))
	viper.BindPFlag("short-break", startCmd.Flags().Lookup("short-break"))
	viper.BindPFlag("long-break", startCmd.Flags().Lookup("long-break"))
	viper.BindPFlag("cycles", startCmd.Flags().Lookup("cycles"))
}
