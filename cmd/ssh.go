/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wish/recover"
	"github.com/franciscolkdo/breach-protocol/config"
	"github.com/franciscolkdo/breach-protocol/game"
	"github.com/franciscolkdo/breach-protocol/game/style"

	"github.com/spf13/cobra"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

var sshKeyPath string

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Start the game in ssh port!",
	Long: `Start the breach-protocol game, it will look into /config/game.json by default
If you want to provide a specific path for the config, use the -c option.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options := []ssh.Option{
			wish.WithAddress(net.JoinHostPort(host, port)),

			// Allocate a pty.
			// This creates a pseudoconsole on windows, compatibility is limited in
			// that case, see the open issues for more details.
			ssh.AllocatePty(),
			wish.WithMiddleware(
				// ensure the user has requested a tty
				activeterm.Middleware(),
				logging.Middleware(),
				recover.Middleware(bubbletea.Middleware(teaHandler)),
				// run our Bubble Tea handler
			),
		}
		if sshKeyPath != "" {
			options = append(options, wish.WithHostKeyPath(sshKeyPath))
		}

		s, err := wish.NewServer(options...)
		if err != nil {
			log.Error("Could not start server", "error", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Info("Starting SSH server", "host", host, "port", port)
		go func() {
			if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				log.Error("Could not start server", "error", err)
				done <- nil
			}
		}()

		<-done
		log.Info("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not stop server", "error", err)
		}
		return nil
	},
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	renderer := bubbletea.MakeRenderer(s)
	style.SetRenderer(renderer)

	cfg, err := config.GetConfig(configPath)

	if err != nil {
		panic(err)
	}
	m := game.NewGame(cfg, s.User())
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func init() {
	sshCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file to use")
	sshCmd.Flags().StringVarP(&sshKeyPath, "sshKey", "s", "", "ssh key file to use")
	rootCmd.AddCommand(sshCmd)
}
