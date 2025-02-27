/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/config"
	"github.com/franciscolkdo/breach-protocol/game"
	"github.com/spf13/cobra"
)

var configPath string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the game!",
	Long: `Start the breach-protocol game, it will look into /config/game.json by default
If you want to provide a specific path for the config, use the -c option.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig(configPath)
		if err != nil {
			return fmt.Errorf("error on reading config file: %s", err)
		}
		g := game.NewGame(cfg)

		_, err = tea.NewProgram(g, tea.WithMouseCellMotion()).Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	startCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file to use")
	rootCmd.AddCommand(startCmd)
}
