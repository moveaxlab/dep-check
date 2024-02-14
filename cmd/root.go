package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dep-check",
	Short: "project structure and dependencies",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(os.Stderr)
		value, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return fmt.Errorf("failed to check flags: %w", err)
		}
		if value {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	viper.SetConfigName("dep-check")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfig()
		} else {
			panic(err)
		}
	}
	rootCmd.PersistentFlags().Bool("debug", false, "print debug messages")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
