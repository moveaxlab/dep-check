package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

const (
	debugFlag = "debug"
)

var (
	DebugMode = false
	rootCmd   = &cobra.Command{
		Use:   "dep-check",
		Short: "Compute the dependency graph of a monorepo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debugFlagValue, err := cmd.Flags().GetBool(debugFlag)
			if err != nil {
				slog.Error(fmt.Sprint("Error reading flag", debugFlag))
				panic(err)
			}
			if debugFlagValue {
				DebugMode = debugFlagValue
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().Bool(debugFlag, false, "Add aditional information for debuing purpose")
}

func Execute() {
	rootCmd.Execute()
}
